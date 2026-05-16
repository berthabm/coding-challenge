package services

import (
	"fmt"

	"gonum.org/v1/gonum/mat"

	"github.com/interseguros/challenge/api-go/internal/models"
)

// QRService performs QR decomposition on a matrix.
type QRService interface {
	Decompose(matrix [][]float64) (*models.QRResult, error)
}

type qrService struct{}

// NewQRService returns the QR decomposition service implementation.
func NewQRService() QRService {
	return &qrService{}
}

// Decompose computes A = Q·R using gonum/mat QR (LAPACK-backed).
//
// Gonum expects m ≥ n (tall or square A). Then Q is m×m orthogonal and R is m×n upper trapezoidal.
func (s *qrService) Decompose(matrix [][]float64) (*models.QRResult, error) {
	m := len(matrix)
	n := len(matrix[0])

	if m < n {
		return nil, fmt.Errorf("QR decomposition requires rows >= cols (got %d×%d)", m, n)
	}

	a := mat.NewDense(m, n, nil)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			a.Set(i, j, matrix[i][j])
		}
	}

	var fac mat.QR
	fac.Factorize(a)

	var q, r mat.Dense
	fac.QTo(&q)
	fac.RTo(&r)

	return &models.QRResult{
		Q: denseToSlice(&q),
		R: denseToSlice(&r),
	}, nil
}

func denseToSlice(d *mat.Dense) [][]float64 {
	rows, cols := d.Dims()
	out := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		out[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			out[i][j] = d.At(i, j)
		}
	}
	return out
}
