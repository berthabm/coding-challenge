package utils

import "math"

// ValidateMatrix checks that matrix exists, is non-empty, rectangular, and numeric.
func ValidateMatrix(matrix [][]float64) error {
	if matrix == nil {
		return ErrInvalidMatrix
	}
	if len(matrix) == 0 {
		return ErrInvalidMatrix
	}

	cols := len(matrix[0])
	if cols == 0 {
		return ErrInvalidMatrix
	}

	for i, row := range matrix {
		if len(row) != cols {
			return ErrInvalidMatrix
		}
		for j, v := range row {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return ErrInvalidMatrix
			}
			_ = i
			_ = j
		}
	}

	return nil
}
