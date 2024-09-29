package dbf

import (
	"bytes"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

func isDigit(c *sql.ColumnType) bool {
	switch c.DatabaseTypeName() {
	case "TINYINT":
		return true
	case "SMALLINT":
		return true
	case "MEDIUMINT":
		return true
	case "BIGINT":
		return true
	case "INT":
		return true
	case "INT1":
		return true
	case "INT2":
		return true
	case "INT3":
		return true
	case "INT4":
		return true
	case "INT8":
		return true
	case "BOOL":
		return true
	case "BOOLEAN":
		return true
	case "DECIMAL":
		return true
	case "DEC":
		return true
	case "NUMERIC":
		return true
	case "FIXED":
		return true
	case "NUMBER":
		return true
	case "FLOAT":
		return true
	case "DOUBLE":
		return true
	default:
		return false
	}
}

func isNum(c *sql.ColumnType) bool {
	t := c.ScanType().Kind()
	if t > 1 && t < 15 {
		return true
	}
	return false
}

func escape(in []byte) []byte {
	var out bytes.Buffer
	for _, b := range in {
		switch b {
		case '\n', '\r', '\t', '\b', '\f', '\\', '"':
			out.WriteByte('\\')
			out.WriteByte(b)
		case '/':
			out.WriteByte('\\')
			out.WriteByte(b)
		default:
			if b < 32 || b == 127 {
				out.WriteString(`\u00`)
				out.WriteString(strconv.FormatInt(int64(b), 16))
			} else {
				out.WriteByte(b)
			}
		}
	}
	return out.Bytes()
}

func Row2Json(columns []*sql.ColumnType, values []sql.RawBytes) (string, error) {
	if len(values) == 0 {
		return "", errors.New("no data in values")
	}
	if len(columns) != len(values) {
		return "", errors.New("columns and values length not equal")
	}
	var buff strings.Builder
	buff.WriteByte('{')
	for i, val := range values {
		buff.WriteByte('"')
		buff.WriteString(columns[i].Name())
		buff.WriteByte('"')
		buff.WriteByte(':')
		if len(val) > 0 {
			if !isDigit(columns[i]) {
				buff.WriteByte('"')
			}
			buff.Write(escape(val))
			if !isDigit(columns[i]) {
				buff.WriteByte('"')
			}
		} else {
			buff.WriteString("null")
		}
		if i != len(values)-1 {
			buff.WriteByte(',')
		}
	}
	buff.WriteByte('}')

	return buff.String(), nil
}
