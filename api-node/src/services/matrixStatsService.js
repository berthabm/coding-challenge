/**
 * Matrix statistics domain service.
 * Computes max, min, average, sum, and diagonal check — implementation pending.
 */

/**
 * @typedef {Object} MatrixStats
 * @property {number} max
 * @property {number} min
 * @property {number} average
 * @property {number} sum
 * @property {boolean} isDiagonal
 * @property {{ rows: number, cols: number }} dimensions
 */

/**
 * @param {number[][]} matrix
 * @returns {MatrixStats}
 */
export function computeMatrixStats(matrix) {
  const rows = matrix.length;
  const cols = rows > 0 ? matrix[0].length : 0;

  // TODO: implement max, min, average, sum, isDiagonal
  return {
    max: 0,
    min: 0,
    average: 0,
    sum: 0,
    isDiagonal: false,
    dimensions: { rows, cols },
  };
}
