package dbf

import (
	"reflect"
	"testing"
)

func TestRow2Json(t *testing.T) {
	// define test cases
	tests := []struct {
		name     string
		columns  []string
		values   []interface{}
		expected []byte
		isErr    bool
	}{
		{
			name:     "Empty input",
			columns:  []string{},
			values:   []interface{}{},
			expected: []byte("{}"),
			isErr:    false,
		},
		{
			name:     "Single column with single value",
			columns:  []string{"id"},
			values:   []interface{}{1},
			expected: []byte(`{"id":1}`),
			isErr:    false,
		},
		{ // assuming error occurs when column and value numbers are mismatched
			name:    "Mismatched column and value count",
			columns: []string{"id", "name"},
			values:  []interface{}{1},
			isErr:   true,
		},
		{ // normal behavior
			name:     "Mismatched column and value count",
			columns:  []string{"id", "name"},
			values:   []interface{}{1, "foo"},
			expected: []byte(`{"id":1,"name":"foo"}`),
			isErr:    false,
		},
		{ // normal behavior with composite column names
			name:     "Mismatched column and value count",
			columns:  []string{"id", "name", "foo_bar"},
			values:   []interface{}{1, "foo", "baz"},
			expected: []byte(`{"foo_bar":"baz","id":1,"name":"foo"}`),
			isErr:    false,
		},
		{ // normal behavior with difference column types
			name:     "Mismatched column and value count",
			columns:  []string{"id", "name", "foo_bar", "price", "qty"},
			values:   []interface{}{1, "foo", "baz", 1.23, 10},
			expected: []byte(`{"foo_bar":"baz","id":1,"name":"foo","price":1.23,"qty":10}`),
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
				t.Errorf("Expected %s, but got %s", string(test.expected), string(output))
			}
		})
	}
}

func TestRow2Map(t *testing.T) {
	// define test cases
	tests := []struct {
		name     string
		columns  []string
		values   []interface{}
		expected map[string]interface{}
		isErr    bool
	}{
		{
			name:     "Empty input",
			columns:  []string{},
			values:   []interface{}{},
			expected: map[string]interface{}{},
			isErr:    false,
		},
		{
			name:     "Single column with single value",
			columns:  []string{"id"},
			values:   []interface{}{1},
			expected: map[string]interface{}{"id": 1},
			isErr:    false,
		},
		{
			name:    "Mismatched column and value count",
			columns: []string{"id", "name"},
			values:  []interface{}{1},
			isErr:   true,
		},
		{
			name:     "Matched column and values count",
			columns:  []string{"id", "name"},
			values:   []interface{}{1, "foo"},
			expected: map[string]interface{}{"id": 1, "name": "foo"},
			isErr:    false,
		},
	}

	// test loop
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := Row2Map(test.columns, test.values)
			if (err != nil) != test.isErr {
				t.Errorf("Got unexpected error: %v", err)
			}
			if !test.isErr && !reflect.DeepEqual(output, test.expected) {
				t.Errorf("expected %v, but got %v", test.expected, output)
			}
		})
	}
}

func BenchmarkRow2Json(b *testing.B) {
	columns := []string{"id", "name", "foo_bar", "price", "qty"}
	values := []interface{}{1, "foo", "baz", 1.23, 10}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Row2Json(columns, values)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkRow2Map(b *testing.B) {
	columns := []string{"id", "name", "foo_bar", "price", "qty"}
	values := []interface{}{1, "foo", "baz", 1.23, 10}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Row2Map(columns, values)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
