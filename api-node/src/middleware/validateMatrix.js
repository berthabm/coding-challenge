/**
 * Middleware wrapper: runs Joi validation before controller.
 */
import { validateMatrixStatsBody } from '../validations/matrixValidation.js';

export function validateMatrixStats(req, res, next) {
  try {
    req.validated = validateMatrixStatsBody(req.body);
    next();
  } catch (err) {
    next(err);
  }
}
