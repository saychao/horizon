package resource

import (
	"gitlab.com/tokend/go/xdr"
)

// AccountThresholds represents an accounts "thresholds", the numerical values
// needed to satisfy the authorization of a given operation.
type AccountThresholds struct {
	LowThreshold  byte `json:"low_threshold"`
	MedThreshold  byte `json:"med_threshold"`
	HighThreshold byte `json:"high_threshold"`
}

func (this *AccountThresholds) Populate(thresholds xdr.Thresholds) {
	this.LowThreshold = thresholds[1]
	this.MedThreshold = thresholds[2]
	this.HighThreshold = thresholds[3]
}
