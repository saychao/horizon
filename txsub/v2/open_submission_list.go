package txsub

import (
	"github.com/saychao/horizon/log"
	"sync"
	"time"

	"github.com/go-errors/errors"
)

// NewDefaultSubmissionList returns a list that manages open submissions purely
// in memory.
func NewDefaultSubmissionList(retry time.Duration) openSubmissionList {
	return &SubmissionList{
		submissions:  map[string]*openSubmission{},
		core:         make(map[string]bool),
		history:      make(map[string]bool),
		retryTimeout: retry,
	}
}

// openSubmission tracks a slice of channels that should be emitted to when we
// know the result for the transactions with the provided hash
type openSubmission struct {
	Envelope    *EnvelopeInfo
	SubmittedAt time.Time
	Listeners   []listener
}

type SubmissionList struct {
	sync.RWMutex
	submissions map[string]*openSubmission
	core        map[string]bool
	history     map[string]bool

	retryTimeout time.Duration
}

func (s *SubmissionList) Envelope(hash string) *EnvelopeInfo {
	s.RLock()
	defer s.RUnlock()
	if os, ok := s.submissions[hash]; ok {
		return os.Envelope
	}
	return nil
}

func (s *SubmissionList) Add(envelope *EnvelopeInfo, ingest bool, l listener) error {
	s.Lock()
	defer s.Unlock()

	if cap(l) == 0 {
		panic("Unbuffered listener cannot be added to openSubmissionList")
	}

	if len(envelope.ContentHash) != 64 {
		return errors.New("Unexpected transaction hash length: must be 64 hex characters")
	}
	if ingest {
		s.history[envelope.ContentHash] = true
	} else {
		s.core[envelope.ContentHash] = true
	}

	os, ok := s.submissions[envelope.ContentHash]

	if !ok {
		os = &openSubmission{
			Envelope:    envelope,
			SubmittedAt: time.Now(),
			Listeners:   []listener{},
		}
		s.submissions[envelope.ContentHash] = os
	}

	os.Listeners = append(os.Listeners, l)

	return nil
}

func (s *SubmissionList) Finish(r fullResult) {
	s.Lock()
	defer s.Unlock()

	os, ok := s.submissions[r.Hash]
	if !ok {
		return
	}

	for _, l := range os.Listeners {
		l <- r
		close(l)
	}

	delete(s.submissions, r.Hash)
	delete(s.core, os.Envelope.ContentHash)
	delete(s.history, os.Envelope.ContentHash)
	return
}

func (s *SubmissionList) Clean(maxAge time.Duration) int {
	s.Lock()
	defer s.Unlock()

	for _, os := range s.submissions {
		if time.Since(os.SubmittedAt) > maxAge {
			delete(s.submissions, os.Envelope.ContentHash)
			delete(s.core, os.Envelope.ContentHash)
			delete(s.history, os.Envelope.ContentHash)

			for _, l := range os.Listeners {
				l <- fullResult{Err: cancelError}
				close(l)
			}
		}
	}

	return len(s.submissions)
}

func (s *SubmissionList) PendingCore() []string {
	s.RLock()
	defer s.RUnlock()
	pendingHashes := make([]string, 0, len(s.core))

	for hash := range s.core {
		pendingHashes = append(pendingHashes, hash)
	}

	return pendingHashes
}
func (s *SubmissionList) PendingHistory() []string {
	s.RLock()
	defer s.RUnlock()
	pendingHashes := make([]string, 0, len(s.history))

	for hash := range s.history {
		pendingHashes = append(pendingHashes, hash)
	}

	return pendingHashes
}

func (s *SubmissionList) ShouldRetry(hash string, t time.Time) bool {
	s.RLock()
	defer s.RUnlock()
	submission := s.submissions[hash]
	if submission == nil {
		log.WithField("tx_hash", hash).Error("Submission is nil")
		return false
	}

	return t.After(submission.SubmittedAt.Add(s.retryTimeout))
}
