import { computeMatrixStats } from '../../src/services/matrixStatsService.js';

describe('computeMatrixStats', () => {
  it('returns scaffold shape with dimensions', () => {
    const matrix = [
      [1, 2],
      [3, 4],
    ];
    const result = computeMatrixStats(matrix);
    expect(result.dimensions).toEqual({ rows: 2, cols: 2 });
    expect(result).toHaveProperty('max');
    expect(result).toHaveProperty('isDiagonal');
  });
});
