import { validateQRStatsBody } from '../../src/validations/qrValidation.js';
import { AppError } from '../../src/utils/errors.js';

const basePayload = {
  q: [[1, 0], [0, 1]],
  r: [[2, 3], [0, 4]],
  originalMatrix: [[1, 2], [3, 4]],
};

describe('validateQRStatsBody', () => {
  it('accepts valid q, r and originalMatrix', () => {
    const value = validateQRStatsBody(basePayload);
    expect(value.q).toHaveLength(2);
    expect(value.r[0][0]).toBe(2);
    expect(value.originalMatrix[1][0]).toBe(3);
  });

  it('rejects payload without q', () => {
    expect(() =>
      validateQRStatsBody({ r: [[1]], originalMatrix: [[1]] })
    ).toThrow(AppError);
    try {
      validateQRStatsBody({ r: [[1]], originalMatrix: [[1]] });
    } catch (e) {
      expect(e.code).toBe('VALIDATION_ERROR');
    }
  });

  it('rejects payload without r', () => {
    expect(() =>
      validateQRStatsBody({ q: [[1]], originalMatrix: [[1]] })
    ).toThrow(AppError);
  });

  it('accepts payload without originalMatrix (it is now optional)', () => {
    const value = validateQRStatsBody({ q: [[1]], r: [[1]] });
    expect(value.q).toHaveLength(1);
    expect(value.originalMatrix).toBeUndefined();
  });

  it('rejects jagged q matrix', () => {
    expect(() =>
      validateQRStatsBody({
        q: [[1, 2], [3]],
        r: [[1, 0], [0, 1]],
        originalMatrix: [[1, 2], [3, 4]],
      })
    ).toThrow(AppError);
    try {
      validateQRStatsBody({
        q: [[1, 2], [3]],
        r: [[1, 0], [0, 1]],
        originalMatrix: [[1, 2], [3, 4]],
      });
    } catch (e) {
      expect(e.message).toMatch(/q must be a rectangular matrix/);
    }
  });

  it('rejects jagged r matrix', () => {
    expect(() =>
      validateQRStatsBody({
        q: [[1, 0], [0, 1]],
        r: [[1, 2, 3], [4]],
        originalMatrix: [[1, 2], [3, 4]],
      })
    ).toThrow(AppError);
  });

  it('rejects jagged originalMatrix', () => {
    expect(() =>
      validateQRStatsBody({
        q: [[1, 0], [0, 1]],
        r: [[2, 3], [0, 4]],
        originalMatrix: [[1, 2], [3]],
      })
    ).toThrow(AppError);
    try {
      validateQRStatsBody({
        q: [[1, 0], [0, 1]],
        r: [[2, 3], [0, 4]],
        originalMatrix: [[1, 2], [3]],
      });
    } catch (e) {
      expect(e.message).toMatch(/originalMatrix must be a rectangular matrix/);
    }
  });
});
