package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/interseguros/challenge/api-go/internal/models"
	"github.com/interseguros/challenge/api-go/internal/services"
	"github.com/interseguros/challenge/api-go/internal/utils"
)

type mockStatsClient struct{}

func (m *mockStatsClient) SendQRForStats(ctx context.Context, payload *models.API2StatsPayload) (map[string]interface{}, error) {
	if len(payload.Q) == 0 || len(payload.R) == 0 {
		return nil, fmt.Errorf("Q and R must not be empty")
	}
	return map[string]interface{}{
		"statistics": map[string]interface{}{
			"max":         1.0,
			"isQDiagonal": false,
			"isRDiagonal": false,
		},
	}, nil
}

func TestProcessMatrixQR(t *testing.T) {
	logger := utils.NewLogger("error")
	svc := services.NewMatrixService(services.NewQRService(), &mockStatsClient{}, logger)

	result, err := svc.ProcessMatrixQR(context.Background(), &models.MatrixRequest{
		Matrix: [][]float64{{1, 2}, {3, 4}, {5, 6}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Statistics == nil {
		t.Fatal("expected statistics in response")
	}
	if result.OriginalMatrix == nil {
		t.Fatal("expected originalMatrix in response")
	}
	if len(result.QR.Q) == 0 || len(result.QR.R) == 0 {
		t.Fatal("expected non-empty Q and R in response")
	}
}

func TestProcessMatrixQR_InvalidMatrix(t *testing.T) {
	logger := utils.NewLogger("error")
	svc := services.NewMatrixService(services.NewQRService(), &mockStatsClient{}, logger)

	_, err := svc.ProcessMatrixQR(context.Background(), &models.MatrixRequest{
		Matrix: [][]float64{{1, 2}, {3}},
	})
	if err != utils.ErrInvalidMatrix {
		t.Fatalf("expected ErrInvalidMatrix, got %v", err)
	}
}
