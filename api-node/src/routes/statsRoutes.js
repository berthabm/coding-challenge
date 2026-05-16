import { Router } from 'express';
import { postQRStats } from '../controllers/qrController.js';
import { validateQRStats } from '../middleware/validateQR.js';

export const statsRoutes = Router();

statsRoutes.post('/stats', validateQRStats, postQRStats);
