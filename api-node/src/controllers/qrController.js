/**
 * HTTP handlers for Q/R statistics (POST /api/stats).
 */
import { computeQRStatistics } from '../services/qrStatsService.js';
import { logger } from '../utils/logger.js';

export function postQRStats(req, res, next) {
  try {
    const { q, r } = req.validated;

    logger.info(
      { qRows: q.length, rRows: r.length },
      'POST /api/stats — computing statistics over Q and R'
    );

    const statistics = computeQRStatistics(q, r);

    logger.info({ statistics }, 'statistics computed');

    res.status(200).json({ statistics });
  } catch (err) {
    next(err);
  }
}
