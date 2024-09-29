package dbf

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestRow2Json(t *testing.T) {
	// define test cases
	tests := []struct {
		name     string
		columns  []*sql.ColumnType
		values   []sql.RawBytes
		expected string
		isErr    bool
	}{
		{
			name:     "Empty input",
			columns:  []*sql.ColumnType{},
			values:   []sql.RawBytes{},
			expected: "",
			isErr:    true,
		},
		{
			name: "Single column with single value",
			columns: []*sql.ColumnType{
				mockColumnType("id", "INT"),
			},
			values:   []sql.RawBytes{[]byte{'1'}},
			expected: `{"id":1}`,
			isErr:    false,
		},
		{ // assuming error occurs when column and value numbers are mismatched
			name:    "Mismatched column and value count",
			columns: []*sql.ColumnType{},
			values:  []sql.RawBytes{[]byte{'1'}},
			isErr:   true,
		},
		{ // normal behavior
			name: "Mismatched column and value count",
			columns: []*sql.ColumnType{
				mockColumnType("id", "INT"),
				mockColumnType("name", "VARCHAR"),
			},
			values:   []sql.RawBytes{[]byte{'1'}, []byte{'f', 'o', 'o'}},
			expected: `{"id":1,"name":"foo"}`,
			isErr:    false,
		},
		{ // normal behavior with composite column names
			name: "Mismatched column and value count",
			columns: []*sql.ColumnType{
				mockColumnType("id", "INT"),
				mockColumnType("name", "VARCHAR"),
				mockColumnType("foo_bar", "VARCHAR"),
			},
			values:   []sql.RawBytes{[]byte{'1'}, []byte{'f', 'o', 'o'}, []byte{'b', 'a', 'z'}},
			expected: `{"id":1,"name":"foo","foo_bar":"baz"}`,
			isErr:    false,
		},
		{ // normal behavior with difference column types
			name: "Mismatched column and value count",
			columns: []*sql.ColumnType{
				mockColumnType("id", "INT"),
				mockColumnType("name", "VARCHAR"),
				mockColumnType("foo_bar", "VARCHAR"),
				mockColumnType("price", "DOUBLE"),
				mockColumnType("qty", "INT"),
			},
			values: []sql.RawBytes{
				[]byte{'1'},
				[]byte{'f', 'o', 'o'},
				[]byte{'b', 'a', 'z'},
				[]byte{'1', '.', '2', '3'},
				[]byte{'1', '0'}},
			expected: `{"id":1,"name":"foo","foo_bar":"baz","price":1.23,"qty":10}`,
			isErr:    false,
		},
	}

	// test loop
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := Row2Json(test.columns, test.values)
			if (err != nil) != test.isErr {
				t.Errorf("Got unexpected error: %v", err)
			}
			if !test.isErr && !reflect.DeepEqual(output, test.expected) {
				t.Errorf("Expected %s, but got %s", test.expected, output)
			}
		})
	}
}

func BenchmarkRow2Json(b *testing.B) {
	columns := []*sql.ColumnType{
		mockColumnType("id", "INT"),
		mockColumnType("name", "VARCHAR"),
		mockColumnType("foo_bar", "VARCHAR"),
		mockColumnType("price", "DOUBLE"),
		mockColumnType("qty", "INT"),
	}
	values := []sql.RawBytes{
		[]byte{'1'},
		[]byte{'f', 'o', 'o'},
		[]byte{'b', 'a', 'z'},
		[]byte{'1', '.', '2', '3'},
		[]byte{'1', '0'}}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Row2Json(columns, values)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func TestEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "Normal string no escape",
			input:    []byte("Hello World"),
			expected: []byte("Hello World"),
		},
		{
			name:     "String with escape sequences",
			input:    []byte{'T', 'e', 's', 't', '\n', '\t', '\r', '\b', '\f', '"', '\\'},
			expected: []byte{'T', 'e', 's', 't', '\\', '\n', '\\', '\t', '\\', '\r', '\\', '\b', '\\', '\f', '\\', '"', '\\', '\\'},
		},
		{
			name:     "String with Unicode characters",
			input:    []byte("Hello 世界"),
			expected: []byte("Hello 世界"),
		},
		{
			name:     "String with emojis",
			input:    []byte("Hello 👋🏾 🌏"),
			expected: []byte("Hello 👋🏾 🌏"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := escape(test.input)
			if !reflect.DeepEqual(output, test.expected) {
				t.Errorf("Expected %s, but got %s", test.expected, output)
			}
		})
	}
}

func mockSampleValue(t string) interface{} {
	switch strings.ToUpper(t) {
	case "TINYINT":
		return 1
	case "SMALLINT":
		return 0
	case "MEDIUMINT":
		return 0
	case "BIGINT":
		return 0
	case "INT":
		return 0
	case "INT1":
		return 0
	case "INT2":
		return 0
	case "INT3":
		return 0
	case "INT4":
		return 0
	case "INT8":
		return 0
	case "BOOL":
		return 0
	case "BOOLEAN":
		return false
	case "DECIMAL":
		return 1.23
	case "DEC":
		return 1.23
	case "NUMERIC":
		return 1.23
	case "FIXED":
		return 1.23
	case "NUMBER":
		return 1.23
	case "FLOAT":
		return 1.23
	case "DOUBLE":
		return 1.23
	default:
		return ""
	}
}

// mockColumnType mock sql.ColumnType for testing
func mockColumnType(colName string, colType string) *sql.ColumnType {
	db, mock, _ := sqlmock.New()
	column1 := mock.NewColumn(colName).OfType(colType, mockSampleValue(colType)).Nullable(true)
	rows := mock.NewRowsWithColumnDefinition(column1)
	rows.AddRow("foo")
	mock.ExpectQuery("SELECT 1").WillReturnRows(rows)
	query, _ := db.Query("SELECT 1")
	columns, err := query.ColumnTypes()
	if err != nil {
		panic(err)
	}
	if len(columns) != 1 {
		log.Fatalf("unexpected column count: %d", len(columns))
	}
	return columns[0]
}

func BenchmarkEscape(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		escape([]byte("Hello World"))
	}
}

func TestIsDigit(t *testing.T) {
	tests := []struct {
		name     string
		input    *sql.ColumnType
		expected bool
	}{
		{
			name:     "Integer Value",
			input:    mockColumnType("id", "INT"),
			expected: true,
		},
		{
			name:     "Boolean Value",
			input:    mockColumnType("active", "BOOLEAN"),
			expected: true,
		},
		{
			name:     "Float Value",
			input:    mockColumnType("percentage", "FLOAT"),
			expected: true,
		},
		{
			name:     "Numeric Value",
			input:    mockColumnType("points", "NUMERIC"),
			expected: true,
		},
		{
			name:     "Varchar Value",
			input:    mockColumnType("name", "VARCHAR"),
			expected: false,
		},
		{
			name:     "Text Value",
			input:    mockColumnType("description", "TEXT"),
			expected: false,
		},
		{
			name:     "Date Value",
			input:    mockColumnType("created_at", "DATE"),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isDigit(test.input)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func BenchmarkIsDigit(b *testing.B) {
	column := mockColumnType("id", "INT")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		isDigit(column)
	}
}

func TestIsNum(t *testing.T) {
	tests := []struct {
		name     string
		input    *sql.ColumnType
		expected bool
	}{
		{
			name:     "Integer Value",
			input:    mockColumnType("id", "INT"),
			expected: true,
		},
		{
			name:     "Varchar Value",
			input:    mockColumnType("name", "VARCHAR"),
			expected: false,
		},
		{
			name:     "Boolean Value",
			input:    mockColumnType("active", "BOOLEAN"),
			expected: false,
		},
		{
			name:     "Float Value",
			input:    mockColumnType("percentage", "FLOAT"),
			expected: true,
		},
		{
			name:     "Double Value",
			input:    mockColumnType("price", "DOUBLE"),
			expected: true,
		},
		{
			name:     "Date Value",
			input:    mockColumnType("created_at", "DATE"),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isNum(test.input)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func BenchmarkIsNum(b *testing.B) {
	column := mockColumnType("id", "INT")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		isNum(column)
	}
}
