/**
 * Express application factory — wires middleware, routes, and error handling.
 * Exported for supertest in integration/unit tests without listening on a port.
 */
import express from 'express';
import pinoHttp from 'pino-http';
import { config } from './config/index.js';
import { errorHandler, notFoundHandler } from './middleware/errorHandler.js';
import { matrixRoutes } from './routes/matrixRoutes.js';
import { statsRoutes } from './routes/statsRoutes.js';
import { healthRoutes } from './routes/healthRoutes.js';
import { logger } from './utils/logger.js';

export function createApp() {
  const app = express();

  app.disable('x-powered-by');
  app.use(express.json({ limit: '1mb' }));
  app.use(
    pinoHttp({
      logger,
      autoLogging: config.nodeEnv !== 'test',
    })
  );

  app.use('/health', healthRoutes);
  app.use('/api', statsRoutes);
  app.use('/api/v1/matrices', matrixRoutes);

  app.use(notFoundHandler);
  app.use(errorHandler);

  return app;
}
