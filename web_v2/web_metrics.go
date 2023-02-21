package web_v2

import (
	"time"

	metrics "github.com/rcrowley/go-metrics"
)

// WebMetrics represents request metrics
type WebMetrics struct {
	requestCounter metrics.Counter
	failureCounter metrics.Counter

	requestTimer metrics.Timer

	failureMeter metrics.Meter
	successMeter metrics.Meter
}

// NewWebMetrics creates new instance of WebMetrics
func NewWebMetrics() *WebMetrics {
	return &WebMetrics{
		requestCounter: metrics.NewCounter(),
		failureCounter: metrics.NewCounter(),

		requestTimer: metrics.NewTimer(),
		failureMeter: metrics.NewMeter(),
		successMeter: metrics.NewMeter(),
	}
}

// Update - updates web metrics with values depending on the state of request
func (m *WebMetrics) Update(requestDuration time.Duration, responseStatus int) {
	m.requestTimer.Update(requestDuration)

	if responseStatus >= 200 && responseStatus < 400 {
		m.successMeter.Mark(1)
	} else if responseStatus >= 400 && responseStatus < 600 {
		m.failureCounter.Inc(1)
		m.failureMeter.Mark(1)
	}
}

// ResetTimer - resets the metrics timer
func (m *WebMetrics) ResetTimer() {
	m.requestTimer = metrics.NewTimer()
}
