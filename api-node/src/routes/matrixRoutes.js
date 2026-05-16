import { Router } from 'express';
import { postMatrixStats } from '../controllers/matrixController.js';
import { validateMatrixStats } from '../middleware/validateMatrix.js';

export const matrixRoutes = Router();

matrixRoutes.post('/stats', validateMatrixStats, postMatrixStats);
