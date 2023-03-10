// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package base

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

func FlagFromXdrContractState(mask int32, allFlags []xdr.ContractState) []regources.Flag {
	result := []regources.Flag{}
	for _, flagValue := range allFlags {
		flagValueAsInt := int32(flagValue)
		if (flagValueAsInt & mask) == flagValueAsInt {
			result = append(result, regources.Flag{
				Value: flagValueAsInt,
				Name:  flagValue.String(),
			})
		}
	}

	return result
}
