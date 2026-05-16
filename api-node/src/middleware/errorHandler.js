/**
 * Centralized error handling — maps AppError and unknown errors to JSON.
 */
import { AppError, toErrorResponse } from '../utils/errors.js';
import { logger } from '../utils/logger.js';

export function notFoundHandler(req, res) {
  res.status(404).json({
    success: false,
    error: { code: 'NOT_FOUND', message: `route ${req.method} ${req.path} not found` },
  });
}

export function errorHandler(err, req, res, next) {
  if (res.headersSent) {
    return next(err);
  }

  const status = err instanceof AppError ? err.statusCode : 500;
  if (status >= 500) {
    logger.error({ err }, 'unhandled error');
  }

  res.status(status).json(toErrorResponse(err));
}
