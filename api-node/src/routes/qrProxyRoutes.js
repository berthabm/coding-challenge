import { Router } from 'express';
import jwt from 'jsonwebtoken';

export const qrProxyRoutes = Router();

/**
 * POST /api/qr
 *
 * Validates the Bearer JWT, then proxies the request to api-go (API_GO_URL).
 * This allows the frontend to use a single backend URL (api-node) for both
 * auth and matrix processing.
 *
 * Required env vars:
 *   API_GO_URL   — URL of the api-go service (default: http://localhost:8080)
 *   JWT_SECRET   — same secret used to sign tokens in authController
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

    // ── 2. Proxy to api-go ───────────────────────────────────────────────────
    const apiGoUrl = (process.env.API_GO_URL ?? 'http://localhost:8080').replace(/\/$/, '');

    const upstream = await fetch(`${apiGoUrl}/api/qr`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: auth,
      },
      body: JSON.stringify(req.body),
    });

    const data = await upstream.json();
    return res.status(upstream.status).json(data);
  } catch (err) {
    next(err);
  }
});
