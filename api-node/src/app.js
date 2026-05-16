/**
 * Express application factory — wires middleware, routes, and error handling.
 * Exported for supertest in integration/unit tests without listening on a port.
 */
import cors from 'cors';
import express from 'express';
import pinoHttp from 'pino-http';
import { config } from './config/index.js';
import { errorHandler, notFoundHandler } from './middleware/errorHandler.js';
import { matrixRoutes } from './routes/matrixRoutes.js';
import { statsRoutes } from './routes/statsRoutes.js';
import { healthRoutes } from './routes/healthRoutes.js';
import { authRoutes } from './routes/authRoutes.js';
import { qrProxyRoutes } from './routes/qrProxyRoutes.js';
import { logger } from './utils/logger.js';

export function createApp() {
  const app = express();

  const allowedOrigins = process.env.CORS_ALLOWED_ORIGINS
    ? process.env.CORS_ALLOWED_ORIGINS.split(',').map((o) => o.trim())
    : [
        'http://localhost:5173',
        'http://localhost:8080',
        'https://frontend-995881892656.southamerica-east1.run.app',
      ];

  const corsOptions = {
    origin: allowedOrigins,
    methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
    allowedHeaders: ['Content-Type', 'Authorization'],
  };

  app.disable('x-powered-by');
  app.options('*', cors(corsOptions));
  app.use(cors(corsOptions));
  app.use(express.json({ limit: '1mb' }));
  app.use(
    pinoHttp({
      logger,
      autoLogging: config.nodeEnv !== 'test',
    })
  );

  app.use('/health', healthRoutes);
  app.use('/api/auth', authRoutes);
  app.use('/api/qr', qrProxyRoutes);
  app.use('/api', statsRoutes);
  app.use('/api/v1/matrices', matrixRoutes);

  app.use(notFoundHandler);
  app.use(errorHandler);

  return app;
}
