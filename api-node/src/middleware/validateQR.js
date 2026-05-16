import { validateQRStatsBody } from '../validations/qrValidation.js';

export function validateQRStats(req, res, next) {
  try {
    req.validated = validateQRStatsBody(req.body);
    next();
  } catch (err) {
    next(err);
  }
}
