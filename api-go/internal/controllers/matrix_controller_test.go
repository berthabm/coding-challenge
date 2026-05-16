package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/interseguros/challenge/api-go/internal/controllers"
	"github.com/interseguros/challenge/api-go/internal/middleware"
	"github.com/interseguros/challenge/api-go/internal/models"
	"github.com/interseguros/challenge/api-go/internal/services"
	"github.com/interseguros/challenge/api-go/internal/utils"
)

type mockStatsClient struct{}

func (m *mockStatsClient) SendQRForStats(ctx context.Context, payload *models.API2StatsPayload) (map[string]interface{}, error) {
	return map[string]interface{}{
		"statistics": map[string]interface{}{
			"max":         1.0,
			"min":         0.0,
			"isQDiagonal": false,
			"isRDiagonal": false,
		},
	}, nil
}

func TestProcessQR(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	logger := utils.NewLogger("error")
	svc := services.NewMatrixService(services.NewQRService(), &mockStatsClient{}, logger)
	ctrl := controllers.NewMatrixController(svc, logger)
	app.Post("/api/qr", ctrl.ProcessQR)

	body, _ := json.Marshal(models.MatrixRequest{
		Matrix: [][]float64{{1, 0}, {0, 1}},
	})
	req := httptest.NewRequest("POST", "/api/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status %d body %s", resp.StatusCode, b)
	}
}

// newMatrixTestApp builds a Fiber app with /api/qr wired (no auth middleware).
func newMatrixTestApp() *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	logger := utils.NewLogger("error")
	svc := services.NewMatrixService(services.NewQRService(), &mockStatsClient{}, logger)
	ctrl := controllers.NewMatrixController(svc, logger)
	app.Post("/api/qr", ctrl.ProcessQR)
	return app
}

func TestProcessQR_DiagonalMatrix(t *testing.T) {
	app := newMatrixTestApp()

	matrix := models.MatrixRequest{
		Matrix: [][]float64{{1, 0, 0}, {0, 5, 0}, {0, 0, 9}},
	}
	body, _ := json.Marshal(matrix)
	req := httptest.NewRequest("POST", "/api/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("diagonal matrix: expected 200, got %d: %s", resp.StatusCode, b)
	}
}

func TestProcessQR_RectangularMatrix(t *testing.T) {
	app := newMatrixTestApp()

	matrix := models.MatrixRequest{
		Matrix: [][]float64{{1, 2}, {3, 4}, {5, 6}}, // 3×2 (m > n)
	}
	body, _ := json.Marshal(matrix)
	req := httptest.NewRequest("POST", "/api/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("3x2 matrix: expected 200, got %d: %s", resp.StatusCode, b)
	}
}

func TestProcessQR_InvalidJSON(t *testing.T) {
	app := newMatrixTestApp()

	malformed := strings.NewReader(`{"matrix": [[1,2],[3,4]`) // truncated JSON
	req := httptest.NewRequest("POST", "/api/qr", malformed)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("invalid JSON: expected 400, got %d", resp.StatusCode)
	}
}

func TestProcessQR_JaggedMatrix(t *testing.T) {
	app := newMatrixTestApp()

	// Valid JSON but rows have different lengths — fails ValidateMatrix
	jagged := strings.NewReader(`{"matrix": [[1,2,3],[4,5],[6,7,8]]}`)
	req := httptest.NewRequest("POST", "/api/qr", jagged)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("jagged matrix: expected 400, got %d: %s", resp.StatusCode, b)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	errObj, _ := result["error"].(map[string]interface{})
	if errObj["code"] != "INVALID_MATRIX" {
		t.Fatalf("expected code INVALID_MATRIX, got %v", errObj["code"])
	}
}

func TestProcessQR_EmptyMatrix(t *testing.T) {
	app := newMatrixTestApp()

	body, _ := json.Marshal(models.MatrixRequest{Matrix: [][]float64{}})
	req := httptest.NewRequest("POST", "/api/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("empty matrix: expected 400, got %d: %s", resp.StatusCode, b)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	errObj, _ := result["error"].(map[string]interface{})
	if errObj["code"] != "INVALID_MATRIX" {
		t.Fatalf("expected code INVALID_MATRIX, got %v", errObj["code"])
	}
}

func TestProcessQR_NonNumericValue(t *testing.T) {
	app := newMatrixTestApp()

	// JSON is valid but "hola" cannot unmarshal into float64 — BodyParser fails
	nonNumeric := strings.NewReader(`{"matrix": [[1,2],[3,"hola"]]}`)
	req := httptest.NewRequest("POST", "/api/qr", nonNumeric)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("non-numeric value: expected 400, got %d: %s", resp.StatusCode, b)
	}
}
