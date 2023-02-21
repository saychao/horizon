package log

import (
	"github.com/lib/pq"
)

func PQEvent(log *Entry) func(event pq.ListenerEventType, err error) {
	return func(event pq.ListenerEventType, err error) {
		switch event {
		case pq.ListenerEventConnected:
			log.Info("connected to db")
			return
		case pq.ListenerEventDisconnected:
			log.WithError(err).Error("disconnected")
			return
		case pq.ListenerEventReconnected:
			log.Info("reconnected")
			return
		case pq.ListenerEventConnectionAttemptFailed:
			log.WithError(err).Error("failed to connect")
			return
		default:
			log.WithField("event_type", event).Error("unknown event")
			return
		}
	}
}
