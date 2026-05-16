/**
 * Validation for POST /api/stats — Q and R matrices from API Go.
 * originalMatrix is accepted but optional (may still be sent by api-go).
 */
import Joi from 'joi';
import { AppError } from '../utils/errors.js';

const rowSchema = Joi.array().items(Joi.number()).min(1);

const qrSchema = Joi.object({
  q: Joi.array().items(rowSchema).min(1).required(),
  r: Joi.array().items(rowSchema).min(1).required(),
  originalMatrix: Joi.array().items(rowSchema).min(1).optional(),
});

function assertRectangular(matrix, label) {
  const cols = matrix[0].length;
  for (const row of matrix) {
    if (row.length !== cols) {
      throw new AppError(
        'VALIDATION_ERROR',
        `${label} must be a rectangular matrix`,
        400
      );
    }
  }
}

/**
 * @param {unknown} body
 */
export function validateQRStatsBody(body) {
  const { error, value } = qrSchema.validate(body, {
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

  assertRectangular(value.q, 'q');
  assertRectangular(value.r, 'r');
  if (value.originalMatrix) {
    assertRectangular(value.originalMatrix, 'originalMatrix');
  }

  return value;
}
