import { Router } from 'express';
import jwt from 'jsonwebtoken';
import { qrDecompose } from '../services/qrDecompService.js';
import { computeQRStatistics } from '../services/qrStatsService.js';

export const qrProxyRoutes = Router();

/**
 * POST /api/qr
 *
 * Accepts { matrix: number[][] }, performs QR decomposition (Gram-Schmidt)
 * and returns statistics over the combined Q and R values.
 *
 * JWT validation uses the same JWT_SECRET env var as /api/auth/login.
 */
qrProxyRoutes.post('/', async (req, res, next) => {
  try {
    // ── 1. Validate JWT ──────────────────────────────────────────────────────
    const auth  = req.headers['authorization'] ?? '';
    const token = auth.startsWith('Bearer ') ? auth.slice(7) : null;

    if (!token) {
      return res.status(401).json({
        success: false,
        error: { code: 'UNAUTHORIZED', message: 'missing or invalid authorization token' },
      });
    }

    const secret = process.env.JWT_SECRET ?? 'change-me-in-production';
    try {
      jwt.verify(token, secret);
    } catch {
      return res.status(401).json({
        success: false,
        error: { code: 'UNAUTHORIZED', message: 'missing or invalid authorization token' },
      });
    }

    // ── 2. Validate matrix input ─────────────────────────────────────────────
    const { matrix } = req.body ?? {};

    if (!Array.isArray(matrix) || matrix.length === 0) {
      return res.status(422).json({
        success: false,
        error: { code: 'INVALID_MATRIX', message: 'matrix must have at least 1 row and 1 column' },
      });
    }

    const colLen = matrix[0].length;
    if (!Array.isArray(matrix[0]) || colLen === 0) {
      return res.status(422).json({
        success: false,
        error: { code: 'INVALID_MATRIX', message: 'matrix must have at least 1 row and 1 column' },
      });
    }

    for (const row of matrix) {
      if (!Array.isArray(row) || row.length !== colLen) {
        return res.status(422).json({
          success: false,
          error: { code: 'INVALID_MATRIX', message: 'all rows must have the same number of columns' },
        });
      }
      if (row.some((v) => typeof v !== 'number' || !isFinite(v))) {
        return res.status(422).json({
          success: false,
          error: { code: 'INVALID_MATRIX', message: 'matrix values must be finite numbers' },
        });
      }
    }

    // ── 3. QR decomposition ──────────────────────────────────────────────────
    let q, r;
    try {
      ({ q, r } = qrDecompose(matrix));
    } catch (decompErr) {
      return res.status(422).json({
        success: false,
        error: { code: 'QR_DECOMPOSITION_FAILED', message: decompErr.message },
      });
    }

    // ── 4. Statistics over Q and R ───────────────────────────────────────────
    const statistics = computeQRStatistics(q, r);

    return res.json({ originalMatrix: matrix, qr: { q, r }, statistics });
  } catch (err) {
    next(err);
  }
});

