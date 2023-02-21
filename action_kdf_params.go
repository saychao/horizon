package horizon

import (
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/resource"
)

type KdfParamsAction struct {
	Action
}

func (action *KdfParamsAction) JSON() {
	action.ValidateBodyType()
	action.Do(func() {
		var response resource.KdfParams
		response.Populate()
		hal.Render(action.W, response)
	})
}
