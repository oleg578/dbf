package main

import (
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		name    string
		d       Dummy
		want    []byte
		wantErr bool
	}{
		{
			name: "All fields populated",
			d: Dummy{
				ID:          123,
				Product:     "Product test",
				Description: "This is a product",
				Price:       1500.0,
				Qty:         100,
				Date:        "2024-12-25",
			},
			want:    []byte(`{"id":123,"product":"Product test","description":"This is a product","price":1500,"qty":100,"date":"2024-12-25"}`),
			wantErr: false,
		},
		{
			name:    "No fields populated",
			d:       Dummy{},
			want:    []byte(`{"id":0,"product":"","description":null,"price":0,"qty":0,"date":""}`),
			wantErr: false,
		},
		{
			name: "Description as complex object",
			d: Dummy{
				ID:          456,
				Product:     "Complex product",
				Description: map[string]string{"key": "value"},
				Price:       2000.0,
				Qty:         50,
				Date:        "2025-01-01",
			},
			want:    []byte(`{"id":456,"product":"Complex product","description":{"key":"value"},"price":2000,"qty":50,"date":"2025-01-01"}`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Marshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Dummy.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dummy.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkMarshal(b *testing.B) {
	d := Dummy{
		ID:          123,
		Product:     "Product test",
		Description: "This is a product",
		Price:       1500.0,
		Qty:         100,
		Date:        "2024-12-25",
	}
	b.Run("Benchmarking Marshal", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := d.Marshal()
			if err != nil {
				b.Fatalf("Error during benchmarking: %v", err)
			}
		}
	})
}
