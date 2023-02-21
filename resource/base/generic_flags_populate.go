package base

import "gitlab.com/tokend/regources"

//go:generate genny -in=$GOFILE -out=asset_policy_type_flags_populate.go gen "flagValueType=xdr.AssetPolicy"
//go:generate genny -in=$GOFILE -out=contract_state_flags_populate.go gen "flagValueType=xdr.ContractState"

func FlagFromflagValueType(mask int32, allFlags []flagValueType) []regources.Flag {
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
