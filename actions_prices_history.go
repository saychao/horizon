package horizon

import (
	"time"

	"github.com/SafeRE-IT/horizon/db2/history"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/render/sse"
	"github.com/SafeRE-IT/horizon/resource"
)

type PricesHistoryAction struct {
	Action

	BaseAsset  string
	QuoteAsset string
	Since      time.Time

	Records  []history.PricePoint
	Resource resource.PriceHistory
}

func (action *PricesHistoryAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		action.addPadding,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *PricesHistoryAction) SSE(stream sse.Stream) {
	action.Setup(
		action.loadParams,
	)
	action.Do(
		func() {
			// we will reuse this variable in sse, so re-initializing is required
			action.Records = []history.PricePoint{}
		},
		action.loadRecords,
		func() {
			records := action.Records[:]

			for _, record := range records {
				stream.Send(sse.Event{
					ID:   record.Timestamp.Format(time.RFC3339),
					Data: record.Price,
				})
				action.Since = record.Timestamp
			}
		})

}

func (action *PricesHistoryAction) loadRecords() {
	points, err := action.HistoryQ().PriceHistory(action.BaseAsset, action.QuoteAsset, action.Since)
	if err != nil {
		action.Log.WithError(err).Error("failed to get price history")
		action.Err = &problem.ServerError
		return
	}

	action.Records = points
}

func (action *PricesHistoryAction) loadParams() {
	var err error
	action.BaseAsset = action.GetNonEmptyString("base_asset")
	action.QuoteAsset = action.GetNonEmptyString("quote_asset")
	since := action.GetString("since")
	if action.Err == nil {
		if since != "" {
			action.Since, err = time.Parse(time.RFC3339, since)
			if err != nil {
				action.SetInvalidField("since", err)
				return
			}
		} else {
			action.Since = time.Unix(0, 0)
		}
	}
}

func (action *PricesHistoryAction) addPadding() {
	nilPointsCount := 360 - len(action.Records)
	if nilPointsCount >= 0 {
		if action.Records == nil || len(action.Records) < 2 {
			action.Err = &problem.BeforeHistory
			return
		}
		diff := action.Records[len(action.Records)-2].Timestamp.Sub(action.Records[len(action.Records)-1].Timestamp)
		for i := 0; i < nilPointsCount; i++ {
			action.Records = append(action.Records, history.PricePoint{Price: 0, Timestamp: action.Records[len(action.Records)-1].Timestamp.Add(diff)})
		}
	}
}

func (action *PricesHistoryAction) loadResource() {
	action.Resource.Prices = action.Records
}
