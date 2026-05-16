package services_test

import (
	"math"
	"testing"

	"github.com/interseguros/challenge/api-go/internal/services"
	"gonum.org/v1/gonum/mat"
)

func TestQRService_Decompose_2x2(t *testing.T) {
	svc := services.NewQRService()
	a := [][]float64{{1, 2}, {3, 4}}

	qr, err := svc.Decompose(a)
	if err != nil {
		t.Fatal(err)
	}

	if len(qr.Q) != 2 || len(qr.Q[0]) != 2 {
		t.Fatalf("expected Q 2×2, got %d×%d", len(qr.Q), len(qr.Q[0]))
	}
	if len(qr.R) != 2 || len(qr.R[0]) != 2 {
		t.Fatalf("expected R 2×2, got %d×%d", len(qr.R), len(qr.R[0]))
	}

	assertMatrixProduct(t, a, qr.Q, qr.R)
}

func TestQRService_Decompose_3x2(t *testing.T) {
	svc := services.NewQRService()
	a := [][]float64{{1, 2}, {3, 4}, {5, 6}}

	qr, err := svc.Decompose(a)
	if err != nil {
		t.Fatal(err)
	}

	if len(qr.Q) != 3 || len(qr.Q[0]) != 3 {
		t.Fatalf("expected Q 3×3, got %d×%d", len(qr.Q), len(qr.Q[0]))
	}
	if len(qr.R) != 3 || len(qr.R[0]) != 2 {
		t.Fatalf("expected R 3×2, got %d×%d", len(qr.R), len(qr.R[0]))
	}

	assertMatrixProduct(t, a, qr.Q, qr.R)
}

func TestQRService_Decompose_WideMatrix(t *testing.T) {
	svc := services.NewQRService()
	_, err := svc.Decompose([][]float64{{1, 2, 3}, {4, 5, 6}})
	if err == nil {
		t.Fatal("expected error for wide matrix (rows < cols)")
	}
}

func TestQRService_Decompose_1x1(t *testing.T) {
	svc := services.NewQRService()
	a := [][]float64{{9}}

	qr, err := svc.Decompose(a)
	if err != nil {
		t.Fatal(err)
	}
	if len(qr.Q) != 1 || len(qr.Q[0]) != 1 {
		t.Fatalf("expected Q 1×1, got %d×%d", len(qr.Q), len(qr.Q[0]))
	}
	if len(qr.R) != 1 || len(qr.R[0]) != 1 {
		t.Fatalf("expected R 1×1, got %d×%d", len(qr.R), len(qr.R[0]))
	}
	assertMatrixProduct(t, a, qr.Q, qr.R)
}

func TestQRService_Decompose_Tall_2x1(t *testing.T) {
	svc := services.NewQRService()
	a := [][]float64{{3}, {4}}

	qr, err := svc.Decompose(a)
	if err != nil {
		t.Fatal(err)
	}
	if len(qr.Q) != 2 || len(qr.Q[0]) != 2 {
		t.Fatalf("expected Q 2×2, got %d×%d", len(qr.Q), len(qr.Q[0]))
	}
	if len(qr.R) != 2 || len(qr.R[0]) != 1 {
		t.Fatalf("expected R 2×1, got %d×%d", len(qr.R), len(qr.R[0]))
	}
	assertMatrixProduct(t, a, qr.Q, qr.R)
}

func assertMatrixProduct(t *testing.T, a, q, r [][]float64) {
	t.Helper()

	qMat := newDense(q)
	rMat := newDense(r)
	var product mat.Dense
	product.Mul(qMat, rMat)

	orig := newDense(a)
	rows, cols := orig.Dims()
	const tol = 1e-9

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if math.Abs(product.At(i, j)-orig.At(i, j)) > tol {
				t.Fatalf("Q·R != A at (%d,%d): got %v want %v",
					i, j, product.At(i, j), orig.At(i, j))
			}
		}
	}
}

func newDense(matrix [][]float64) *mat.Dense {
	m := len(matrix)
	n := len(matrix[0])
	d := mat.NewDense(m, n, nil)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			d.Set(i, j, matrix[i][j])
		}
	}
	return d
}
