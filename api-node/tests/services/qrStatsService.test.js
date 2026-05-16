import { computeQRStatistics, isDiagonal } from '../../src/services/qrStatsService.js';

describe('computeQRStatistics', () => {
  it('calcula max, min, suma y promedio sobre los valores combinados de Q y R', () => {
    // Q = identidad 2x2, R = [[2,3],[0,1]]
    const q = [[1, 0], [0, 1]];
    const r = [[2, 3], [0, 1]];
    // valores combinados: 1, 0, 0, 1, 2, 3, 0, 1 = [0,0,0,1,1,1,2,3]
    const stats = computeQRStatistics(q, r);
    expect(stats.max).toBe(3);
    expect(stats.min).toBe(0);
    expect(stats.sum).toBeCloseTo(8);
    expect(stats.average).toBeCloseTo(1);
    expect(stats.isQDiagonal).toBe(true);
    expect(stats.isRDiagonal).toBe(false);
  });

  it('detecta ambas matrices diagonales', () => {
    const q = [[1, 0], [0, 1]];
    const r = [[2, 0], [0, 3]];
    const stats = computeQRStatistics(q, r);
    expect(stats.isQDiagonal).toBe(true);
    expect(stats.isRDiagonal).toBe(true);
  });

  it('maneja valores negativos', () => {
    const q = [[-1, 0], [0, 1]];
    const r = [[2, -3], [0, 1]];
    const stats = computeQRStatistics(q, r);
    expect(stats.max).toBe(2);
    expect(stats.min).toBe(-3);
    expect(stats.isQDiagonal).toBe(true);
    expect(stats.isRDiagonal).toBe(false);
  });

  it('redondea a máximo 3 decimales', () => {
    // Use values that produce non-zero floats in all positions
    const q = [[1 / 3, 1 / 7], [1 / 7, 1 / 3]];
    const r = [[1 / 6, 1 / 9], [1 / 9, 1 / 6]];
    const stats = computeQRStatistics(q, r);
    // All values > 0 so min > 0; max should be 1/3 ≈ 0.333
    expect(stats.max).toBe(0.333);   // round3(1/3)
    expect(String(stats.min).split('.')[1]?.length ?? 0).toBeLessThanOrEqual(3);
    expect(String(stats.average).split('.')[1]?.length ?? 0).toBeLessThanOrEqual(3);
  });

  it('maneja matrices rectangulares en Q', () => {
    const q = [[1, 0], [0, 1], [0, 0]];
    const r = [[2, 3], [0, 1]];
    const stats = computeQRStatistics(q, r);
    expect(stats.isQDiagonal).toBe(false); // Q no es cuadrada
    expect(typeof stats.max).toBe('number');
  });
});

describe('isDiagonal', () => {
  it('returns false para matriz vacía o no cuadrada', () => {
    expect(isDiagonal([])).toBe(false);
    expect(isDiagonal([[1, 2, 3]])).toBe(false);
    expect(isDiagonal([[1, 0], [0, 1], [0, 0]])).toBe(false);
  });

  it('returns true para identidad y diagonal con ceros', () => {
    expect(isDiagonal([[1, 0], [0, 1]])).toBe(true);
    expect(isDiagonal([[2, 0], [0, 3]])).toBe(true);
    expect(isDiagonal([[0, 0], [0, 0]])).toBe(true);
  });

  it('returns false cuando un elemento fuera de la diagonal no es exactamente cero', () => {
    expect(isDiagonal([[1, 2], [0, 3]])).toBe(false);
    expect(isDiagonal([[1, 0.0000001], [0, 1]])).toBe(false);
  });
});
