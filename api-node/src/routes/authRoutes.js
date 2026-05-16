import { Router } from 'express';
import { postLogin } from '../controllers/authController.js';

export const authRoutes = Router();

authRoutes.post('/login', postLogin);
