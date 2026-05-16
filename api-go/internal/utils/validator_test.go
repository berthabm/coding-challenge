package utils

import (
	"math"
	"testing"
)

func TestValidateMatrix(t *testing.T) {
	tests := []struct {
		name    string
		matrix  [][]float64
		wantErr bool
	}{
		{"valid_rectangular", [][]float64{{1, 2}, {3, 4}}, false},
		{"valid_tall", [][]float64{{1}, {2}, {3}}, false},
		{"valid_single_scalar", [][]float64{{42}}, false},
		{"nil", nil, true},
		{"empty", [][]float64{}, true},
		{"jagged", [][]float64{{1, 2}, {3}}, true},
		{"empty_row", [][]float64{{}}, true},
		{"nan_cell", [][]float64{{math.NaN(), 1}, {2, 3}}, true},
		{"positive_inf_cell", [][]float64{{math.Inf(1), 0}}, true},
		{"negative_inf_cell", [][]float64{{0, math.Inf(-1)}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMatrix(tt.matrix)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateMatrix() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != ErrInvalidMatrix {
				t.Fatalf("ValidateMatrix() err = %v, want %v", err, ErrInvalidMatrix)
			}
		})
	}
}
