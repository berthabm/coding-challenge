/** @type {import('jest').Config} */
export default {
  testEnvironment: 'node',
  roots: ['<rootDir>/tests'],
  transform: {},
  moduleNameMapper: {},
  collectCoverageFrom: ['src/**/*.js', '!src/index.js'],
};
