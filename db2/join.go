package db2

import "fmt"

//GetColumnsForJoin - adds to all columns prefix
func GetColumnsForJoin(rawColumns []string, tableAlias string) []string {
	result := make([]string, 0, len(rawColumns))
	for _, column := range rawColumns {
		result = append(result, fmt.Sprintf(`%s.%s "%s.%s"`, tableAlias, column, tableAlias, column))
	}

	return result
}
