package services

import (
	"context"
	"math"

	"github.com/interseguros/challenge/api-go/internal/clients"
	"github.com/interseguros/challenge/api-go/internal/models"
	"github.com/interseguros/challenge/api-go/internal/utils"
)

// round3 rounds a float64 to 3 decimal places.
func round3(v float64) float64 {
	return math.Round(v*1000) / 1000
}

// roundMatrix returns a new matrix with every element rounded to 3 decimal places.
func roundMatrix(m [][]float64) [][]float64 {
	out := make([][]float64, len(m))
	for i, row := range m {
		out[i] = make([]float64, len(row))
		for j, v := range row {
			out[i][j] = round3(v)
		}
	}
	return out
}

// MatrixService orchestrates QR decomposition and downstream stats (API 2).
type MatrixService interface {
	ProcessMatrixQR(ctx context.Context, req *models.MatrixRequest) (*models.QRResponse, error)
}

type matrixService struct {
	qrService   QRService
	statsClient clients.StatsClient
	logger      utils.Logger
}

// NewMatrixService wires dependencies (constructor injection for tests).
func NewMatrixService(qr QRService, stats clients.StatsClient, logger utils.Logger) MatrixService {
	return &matrixService{
		qrService:   qr,
		statsClient: stats,
		logger:      logger,
	}
}

// ProcessMatrixQR validates input, decomposes QR, forwards q/r/originalMatrix to API 2,
// and assembles the final structured response.
func (s *matrixService) ProcessMatrixQR(ctx context.Context, req *models.MatrixRequest) (*models.QRResponse, error) {
	if err := utils.ValidateMatrix(req.Matrix); err != nil {
		s.logger.Warn("invalid matrix payload")
		return nil, err
	}

	rows, cols := utils.MatrixDimensions(req.Matrix)
	s.logger.Info("processing matrix", "rows", rows, "cols", cols)

	qr, err := s.qrService.Decompose(req.Matrix)
	if err != nil {
		s.logger.Error("qr decomposition failed", "error", err)
		return nil, utils.ErrQRDecomposition
	}

	s.logger.Info("QR decomposition completed", "q_rows", len(qr.Q), "r_rows", len(qr.R))

	payload := &models.API2StatsPayload{
		Q:              qr.Q,
		R:              qr.R,
		OriginalMatrix: req.Matrix,
	}
	nodeResp, err := s.statsClient.SendQRForStats(ctx, payload)
	if err != nil {
		s.logger.Error("api2 call failed", "error", err)
		return nil, utils.ErrDownstreamAPI
	}

	// Extract the statistics map from Node's response.
	statsRaw, _ := nodeResp["statistics"]
	stats, _ := statsRaw.(map[string]interface{})

	qRounded := roundMatrix(qr.Q)
	rRounded := roundMatrix(qr.R)

	s.logger.Info("statistics received from api2")
	return &models.QRResponse{
		OriginalMatrix: req.Matrix,
		QR:             models.QRResult{Q: qRounded, R: rRounded},
		Statistics:     stats,
	}, nil
}
