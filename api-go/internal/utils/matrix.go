package utils

// MatrixDimensions returns row and column counts for a rectangular matrix.
// Business validation (square, non-empty, numeric bounds) will be added later.
func MatrixDimensions(matrix [][]float64) (rows, cols int) {
	if len(matrix) == 0 {
		return 0, 0
	}
	return len(matrix), len(matrix[0])
}
