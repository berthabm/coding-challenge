/**
 * HTTP layer for matrix statistics — delegates to services, no business logic.
 */
import { computeMatrixStats } from '../services/matrixStatsService.js';

/**
 * POST /api/v1/matrices/stats
 * Receives matrix from API 1 (or direct clients) and returns aggregated stats.
 */
export function postMatrixStats(req, res, next) {
  try {
    const { matrix, source } = req.validated;

    const stats = computeMatrixStats(matrix);

    res.status(200).json({
      success: true,
      data: {
        stats,
        meta: { source: source ?? 'unknown' },
      },
    });
  } catch (err) {
    next(err);
  }
}
