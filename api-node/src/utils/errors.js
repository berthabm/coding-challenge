/**
 * Application errors with HTTP status and stable codes for API consumers.
 */
export class AppError extends Error {
  constructor(code, message, statusCode = 500) {
    super(message);
    this.name = 'AppError';
    this.code = code;
    this.statusCode = statusCode;
  }
}

export const Errors = {
  INVALID_MATRIX: new AppError('INVALID_MATRIX', 'matrix payload is invalid', 400),
  STATS_FAILED: new AppError('STATS_FAILED', 'could not compute matrix statistics', 422),
};

export function toErrorResponse(err) {
  return {
    success: false,
    error: {
      code: err.code || 'INTERNAL_ERROR',
      message: err.message || 'internal server error',
    },
  };
}
