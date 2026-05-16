/**
 * Request validation with Joi — keeps controllers thin and rules declarative.
 */
import Joi from 'joi';
import { AppError } from '../utils/errors.js';

const matrixSchema = Joi.object({
  matrix: Joi.array()
    .items(Joi.array().items(Joi.number()).min(1))
    .min(1)
    .required(),
  source: Joi.string().optional(),
});

/**
 * Validates POST body for /api/v1/matrices/stats.
 * Throws AppError on failure (handled by error middleware).
 */
export function validateMatrixStatsBody(body) {
  const { error, value } = matrixSchema.validate(body, {
    abortEarly: false,
    stripUnknown: true,
  });

  if (error) {
    throw new AppError(
      'VALIDATION_ERROR',
      error.details.map((d) => d.message).join('; '),
      400
    );
  }

  return value;
}
