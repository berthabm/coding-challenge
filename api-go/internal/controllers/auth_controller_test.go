package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/interseguros/challenge/api-go/internal/config"
	"github.com/interseguros/challenge/api-go/internal/controllers"
	"github.com/interseguros/challenge/api-go/internal/middleware"
	"github.com/interseguros/challenge/api-go/internal/models"
	"github.com/interseguros/challenge/api-go/internal/services"
	"github.com/interseguros/challenge/api-go/internal/utils"
)

const authTestSecret = "test-jwt-secret-auth-controller-32ch"

// mockQRStatsClient is a separate mock used only in auth/jwt tests.
type mockQRStatsClient struct{}

func (m *mockQRStatsClient) SendQRForStats(_ context.Context, _ *models.API2StatsPayload) (map[string]interface{}, error) {
	return map[string]interface{}{
		"statistics": map[string]interface{}{
			"max":         1.0,
			"isQDiagonal": false,
			"isRDiagonal": false,
		},
	}, nil
}

// newAuthTestApp builds a minimal Fiber app with login + protected /api/qr.
func newAuthTestApp() *fiber.App {
	cfg := &config.Config{
		AuthUsername: "testuser",
		AuthPassword: "testpass",
		JWTSecret:    authTestSecret,
	}
	logger := utils.NewLogger("error")
	svc := services.NewMatrixService(services.NewQRService(), &mockQRStatsClient{}, logger)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	authCtrl := controllers.NewAuthController(cfg, logger)
	matCtrl := controllers.NewMatrixController(svc, logger)

	app.Post("/api/auth/login", authCtrl.Login)

	protected := app.Group("/api", middleware.JWTAuth(authTestSecret))
	protected.Post("/qr", matCtrl.ProcessQR)

	return app
}

// signTestToken generates a valid JWT using the test secret.
func signTestToken() string {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "testuser",
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, _ := tok.SignedString([]byte(authTestSecret))
	return signed
}

// ---- Login tests ----

func TestLogin_Correct(t *testing.T) {
	app := newAuthTestApp()

	body, _ := json.Marshal(map[string]string{"username": "testuser", "password": "testpass"})
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, b)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, ok := result["token"]; !ok {
		t.Fatal("expected 'token' field in response body")
	}
}

func TestLogin_Incorrect(t *testing.T) {
	app := newAuthTestApp()

	body, _ := json.Marshal(map[string]string{"username": "testuser", "password": "wrongpass"})
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	errObj, _ := result["error"].(map[string]interface{})
	if errObj["code"] != "INVALID_CREDENTIALS" {
		t.Fatalf("expected code INVALID_CREDENTIALS, got %v", errObj["code"])
	}
}

// ---- JWT middleware tests on /api/qr ----

func TestProcessQR_NoToken(t *testing.T) {
	app := newAuthTestApp()

	body, _ := json.Marshal(models.MatrixRequest{Matrix: [][]float64{{1, 2}, {3, 4}}})
	req := httptest.NewRequest("POST", "/api/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// deliberately no Authorization header

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	errObj, _ := result["error"].(map[string]interface{})
	if errObj["code"] != "UNAUTHORIZED" {
		t.Fatalf("expected code UNAUTHORIZED, got %v", errObj["code"])
	}
}

func TestProcessQR_ValidToken(t *testing.T) {
	app := newAuthTestApp()
	token := signTestToken()

	body, _ := json.Marshal(models.MatrixRequest{Matrix: [][]float64{{1, 2}, {3, 4}}})
	req := httptest.NewRequest("POST", "/api/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, b)
	}
}
