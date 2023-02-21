package history2

import "encoding/json"

type Data struct {
	ID    int64           `db:"id"`
	Type  int64           `db:"type"`
	Value json.RawMessage `db:"value"`
	Owner string          `db:"owner"`
}
