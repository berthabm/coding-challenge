/**
 * QR decomposition via modified Gram-Schmidt orthogonalization.
 *
 * Given an m×n matrix A (m ≥ n), returns:
 *   Q — m×n orthonormal matrix (columns are orthonormal)
 *   R — n×n upper-triangular matrix
 * such that A ≈ Q·R (up to floating-point rounding).
 *
 * Results are rounded to 3 decimal places to match api-go behaviour.
 */

function getColumn(matrix, j) {
  return matrix.map((row) => row[j]);
}

function dot(a, b) {
  return a.reduce((sum, val, i) => sum + val * b[i], 0);
}

function vectorNorm(v) {
  return Math.sqrt(dot(v, v));
}

function round3(v) {
  return parseFloat(v.toFixed(3));
}

/**
 * @param {number[][]} matrix  — input matrix, rows ≥ columns
 * @returns {{ q: number[][], r: number[][] }}
 * @throws {Error} if the matrix is degenerate or dimensions are invalid
 */
export function qrDecompose(matrix) {
  const m = matrix.length;
  if (m === 0) throw new Error('matrix must have at least 1 row and 1 column');

  const n = matrix[0].length;
  if (n === 0) throw new Error('matrix must have at least 1 row and 1 column');

  if (m < n) {
    throw new Error('QR decomposition requires rows >= columns');
  }

  // ── Modified Gram-Schmidt ────────────────────────────────────────────────
  // Build Q as list of n orthonormal column vectors (each of length m)
  const qCols = [];

  for (let j = 0; j < n; j++) {
    // Start with column j of A
    let v = getColumn(matrix, j);

    // Subtract projections onto already-computed basis vectors
    for (let i = 0; i < j; i++) {
      const proj = dot(v, qCols[i]); // qCols[i] is already normalised → denom = 1
      v = v.map((val, idx) => val - proj * qCols[i][idx]);
    }

    const nv = vectorNorm(v);
    if (nv < 1e-10) {
      throw new Error('QR decomposition failed: matrix is rank-deficient');
    }
    qCols.push(v.map((x) => x / nv));
  }

  // Convert Q columns → row-major m×n matrix, rounded to 3 dp
  const Q = Array.from({ length: m }, (_, i) =>
    qCols.map((col) => round3(col[i]))
  );

  // R[i][j] = qCols[i] · column_j(A)  for i ≤ j, else 0  (n×n)
  const R = Array.from({ length: n }, (_, i) =>
    Array.from({ length: n }, (_, j) =>
      i <= j ? round3(dot(qCols[i], getColumn(matrix, j))) : 0
    )
  );

  return { q: Q, r: R };
}
