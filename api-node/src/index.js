/**
 * Entry point: loads config, creates Express app, starts HTTP server.
 * Business logic lives under services/; this file only bootstraps.
 */
import { createApp } from './app.js';
import { config } from './config/index.js';
import { logger } from './utils/logger.js';

const app = createApp();

const server = app.listen(config.port, config.host, () => {
  logger.info(
    { port: config.port, host: config.host, env: config.nodeEnv },
    'server started'
  );
});

const shutdown = (signal) => {
  logger.info({ signal }, 'shutting down');
  server.close(() => process.exit(0));
};

process.on('SIGTERM', () => shutdown('SIGTERM'));
process.on('SIGINT', () => shutdown('SIGINT'));
