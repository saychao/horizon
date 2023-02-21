package core2

import (
	"encoding/json"
	"fmt"
	"testing"

	"gitlab.com/tokend/go/xdr"
)

func TestTransactionEnvelope(t *testing.T) {
	//This snippet is used to quick unmarsahlling of an XDR string.
	env := xdr.TransactionEnvelope{}
	err := xdr.SafeUnmarshalBase64("AAAAAMqphbZCfO0ZXmfu5DTuJCCxU1QRFH373DSBauCZtroyAAAAAAAAAAAAAAAAAAAAAAAAAABcbqitAAAAAAAAAAEAAAAAAAAAFwAAAABQnSZeABTOqByuDv22XfbBroRigXPBh36gtikIxyA4YgAAAAEAAAAAQMMhkT3vOK3J31Ijrs7u9WeXtHpWCOouivG30mf0oVQAAAAABfXhAAAAAAAAD0JAAAAAAAAPQkAAAAAAAAAAAAAPQkAAAAAAAA9CQAAAAAAAAAABAAAAAAAAAARzdWJqAAAAA3JlZgAAAAAAAAAAAAAAAAGZtroyAAAAQFoAosDzAR36MqVYOUXs54De3fBWQtcE4G0f3WU2XWMEdDPs2YT+eDlCutkN+AJCHq0ZpZUtE/Vh+eCIxwUAJwI=", &env)
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(&env)
	if err != nil {
		panic(err)
	}
	fmt.Printf(string(bytes))
}
