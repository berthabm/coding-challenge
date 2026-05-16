/**
 * Centralized configuration from environment variables.
 * Avoids scattering process.env across controllers/services (easier testing).
 */
import dotenv from 'dotenv';

dotenv.config();

const num = (value, fallback) => {
  const n = Number(value);
  return Number.isFinite(n) ? n : fallback;
};

export const config = {
  appName: process.env.APP_NAME || 'api-node',
  nodeEnv: process.env.NODE_ENV || 'development',
  host: process.env.APP_HOST || '0.0.0.0',
  port: num(process.env.APP_PORT, 3000),
  logLevel: process.env.LOG_LEVEL || 'info',
};
