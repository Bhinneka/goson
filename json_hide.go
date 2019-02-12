package goson

// this utility will hide json response based on struct json field you choose

type fieldSet map[string]bool

func genFieldSet(fields ...string) fieldSet {
	var fieldMap = make(fieldSet, len(fields))
	for _, v := range fields {
		fieldMap[v] = true
	}
	return fieldMap
}

// HideFields will hide json field based on struct json field you choose
// params
// - input (struct only)
// - fields json field that you will be hide
func HideFields(input interface{}, fields ...string) (map[string]interface{}, error) {
	return nil, nil
}
