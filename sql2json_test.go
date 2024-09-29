package dbf

import (
	"testing"
)

// Testing AnyToJson function using table driven tests
func TestAnyToJson(t *testing.T) {
	cases := []struct {
		name          string
		data          interface{}
		expected      string
		expectedError error
	}{
		{
			name:          "Valid Struct",
			data:          struct{ Name string }{"John"},
			expected:      `{"Name":"John"}`,
			expectedError: nil,
		},
		{
			name:          "Valid Map",
			data:          map[string]string{"name": "John", "age": "30"},
			expected:      `{"age":"30","name":"John"}`,
			expectedError: nil,
		},
		{
			name:          "Valid String",
			data:          "Hello World",
			expected:      `"Hello World"`,
			expectedError: nil,
		},
		{
			name:          "Valid Int",
			data:          123,
			expected:      `123`,
			expectedError: nil,
		},
		{
			name:          "Valid ArraySlice",
			data:          []string{"John", "Doe"},
			expected:      `["John","Doe"]`,
			expectedError: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBytes, err := AnyToJson(tc.data)
			if err != nil && tc.expectedError == nil {
				t.Errorf("Test %s failed: expected no error, but got: %v", tc.name, err)
				return
			}
			if err == nil && tc.expectedError != nil {
				t.Errorf("Test %s failed: expected error %v, but got nil", tc.name, tc.expectedError)
				return
			}
			if err != nil && tc.expectedError != nil {
				if err.Error() != tc.expectedError.Error() {
					t.Errorf("Test %s failed: expected error %v, but got %v", tc.name, tc.expectedError, err)
				}
				return
			}
			if string(jsonBytes) != tc.expected {
				t.Errorf("Test %s failed:\nexpected:\n%s\n\nbut got:\n%s\n", tc.name, tc.expected, string(jsonBytes))
			}
		})
	}
}

func BenchmarkAnyToJson(b *testing.B) {
	data := struct {
		ID          int64       `json:"id"`
		Product     string      `json:"product"`
		Description interface{} `json:"description"`
		Price       float64     `json:"price"`
		Qty         int64       `json:"qty"`
		Date        string      `json:"date"`
	}{
		ID:          123,
		Product:     "Product test",
		Description: "This is a product",
		Price:       1500.0,
		Qty:         100,
		Date:        "2024-12-25",
	}
	b.Run("Benchmarking AnyToJson", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := AnyToJson(data)
			if err != nil {
				b.Fatalf("Error during benchmarking: %v", err)
			}
		}
	})
}
