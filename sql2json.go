package dbf

import (
	"encoding/json"
	"errors"
	"strconv"
)

func Row2Json(columns []string, values []interface{}) ([]byte, error) {
	rowMap := make(map[string]interface{})
	if len(columns) != len(values) {
		return nil, errors.New("columns and values length not equal")
	}
	for i, col := range columns {
		rowMap[col] = assignCellValue(values[i]) // we will help to typify the value
	}
	return json.Marshal(rowMap)
}

func Row2Map(columns []string, values []interface{}) (map[string]interface{}, error) {
	rowMap := make(map[string]interface{})
	if len(columns) != len(values) {
		return nil, errors.New("columns and values length not equal")
	}
	for i, col := range columns {
		rowMap[col] = assignCellValue(values[i]) // we will help to typify the value
	}
	return rowMap, nil
}

func assignCellValue(val interface{}) interface{} {
	if b, ok := val.([]byte); ok {
		if floatValue, err := strconv.ParseFloat(string(b), 64); err == nil {
			return floatValue
		}
		return string(b)
	}
	return val
}
