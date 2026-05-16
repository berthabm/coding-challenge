/**
 * Structured logging via Pino (JSON in production, pretty in dev if desired later).
 */
import pino from 'pino';
import { config } from '../config/index.js';

export const logger = pino({
  name: config.appName,
  level: config.logLevel,
});
