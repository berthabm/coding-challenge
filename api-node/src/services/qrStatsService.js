/**
 * Statistics computed over the Q and R matrices returned by the QR decomposition.
 */

function flatten(matrix) {
  return matrix.flat();
}

/**
 * A square matrix is diagonal when all off-diagonal entries are zero.
 */
export function isDiagonal(matrix) {
  const rows = matrix.length;
  if (rows === 0) return false;
  const cols = matrix[0].length;
  if (rows !== cols) return false;

  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      if (i !== j && matrix[i][j] !== 0) {
        return false;
      }
    }
  }
  return true;
}

/**
 * Rounds a number to at most 3 decimal places, stripping trailing zeros.
 * @param {number} v
 * @returns {number}
 */
function round3(v) {
  return parseFloat(v.toFixed(3));
}

/**
 * Computes descriptive statistics on the combined values of matrices Q and R.
 * @param {number[][]} q  - Orthogonal factor from QR decomposition
 * @param {number[][]} r  - Upper-triangular factor from QR decomposition
 */
export function computeQRStatistics(q, r) {
  const values = [...flatten(q), ...flatten(r)];

  if (values.length === 0) {
    return {
      max: 0,
      min: 0,
      average: 0,
      sum: 0,
      isQDiagonal: isDiagonal(q),
      isRDiagonal: isDiagonal(r),
    };
  }

  const sum     = values.reduce((acc, n) => acc + n, 0);
  const max     = Math.max(...values);
  const min     = Math.min(...values);
  const average = sum / values.length;

  return {
    max:         round3(max),
    min:         round3(min),
    average:     round3(average),
    sum:         round3(sum),
    isQDiagonal: isDiagonal(q),
    isRDiagonal: isDiagonal(r),
  };
}
