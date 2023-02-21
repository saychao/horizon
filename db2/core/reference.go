package core

// Reference is a struct to unmarshal a row from `reference` table of core DB.
type Reference struct {
	Reference string `db:"reference" json:"reference"`
}
