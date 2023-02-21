package ledger

import (
	"context"

	"strconv"
	"sync"
	"time"

	"github.com/SafeRE-IT/horizon/log"
	"github.com/lib/pq"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

type notificationListener struct {
	state  *State
	chanel string
	log    *log.Entry
	lock   *sync.RWMutex
}

func startNewListener(ctx context.Context, logger *log.Entry, lock *sync.RWMutex, state *State, chanelName,
	dbConnString string) error {
	listener := &notificationListener{
		state:  state,
		log:    logger,
		chanel: chanelName,
		lock:   lock,
	}

	pqListener := pq.NewListener(dbConnString, time.Second, time.Duration(30)*time.Second, log.PQEvent(listener.log))
	err := pqListener.Listen(chanelName)
	if err != nil {
		_ = pqListener.Close()
		return errors.Wrap(err, "failed to subsribe to chanel")
	}

	listener.startHealthChecker(ctx)
	go listener.listen(ctx, pqListener)
	return nil
}

const healthCheckPeriod = time.Duration(30) * time.Second
const healthCheckMinRetryPeriod = time.Duration(1) * time.Minute
const healthCheckMaxRetryPeriod = time.Duration(10) * time.Minute

func (l *notificationListener) startHealthChecker(ctx context.Context) {
	go running.WithBackOff(ctx, l.log, "health_checker", l.checkHealth, healthCheckPeriod,
		healthCheckMinRetryPeriod, healthCheckMaxRetryPeriod)
}

func (l *notificationListener) checkHealth(ctx context.Context) error {
	l.lock.RLock()
	defer l.lock.RUnlock()
	unhealthyPoint := l.state.LastLedgerIncreaseTime.Add(healthCheckPeriod)
	// if we already passed unhealthy point need to notify admin
	if unhealthyPoint.Before(time.Now().UTC()) {
		return errors.From(errors.New("ledger seq is not growing for too long"), logan.F{
			"latest_seq":    l.state.Latest,
			"latest_update": l.state.LastLedgerIncreaseTime,
		})
	}

	return nil
}

func (l *notificationListener) listen(ctx context.Context, listener *pq.Listener) {
	defer func() {
		_ = listener.Close()
	}()

	for notification := range listener.NotificationChannel() {
		if running.IsCancelled(ctx) {
			l.log.Info("Stopped due to ctx been closed")
			return
		}

		err := l.handleNotification(notification)
		// do not exit on error as it was due to issues with particular notification
		if err != nil {
			l.log.WithError(err).Error("failed to handle notification")
			continue
		}

	}

	l.log.Error("Did not expected notification chanel to close")
}

func (l *notificationListener) handleNotification(notification *pq.Notification) error {
	// in some cases chanel might return nil notification (ex. on new connection)
	if notification == nil {
		return nil
	}

	if notification.Channel != l.chanel {
		log.WithField("actual_chanel", notification.Channel).WithField("expected_chanel", l.chanel).
			Error("")
		return errors.From(errors.New("notification from unexpected chanel"), logan.F{
			"actual_chanel":   notification.Channel,
			"expected_chanel": l.chanel,
		})
	}

	ledgerSeq, err := strconv.ParseInt(notification.Extra, 10, 32)
	if err != nil {
		return errors.Wrap(err, "failed to parse payload")
	}

	lock.Lock()
	defer lock.Unlock()
	l.state.LastLedgerIncreaseTime = time.Now().UTC()
	l.state.Latest = int32(ledgerSeq)
	l.log.WithField("new_seq", ledgerSeq).Info("Successfully updated ledger seq")
	return nil
}
