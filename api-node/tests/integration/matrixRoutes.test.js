import request from 'supertest';
import { createApp } from '../../src/app.js';

describe('POST /api/v1/matrices/stats', () => {
  const app = createApp();

  it('accepts valid matrix payload (scaffold)', async () => {
    const res = await request(app)
      .post('/api/v1/matrices/stats')
      .send({ matrix: [[1, 0], [0, 1]], source: 'api-go' })
      .expect(200);

    expect(res.body.success).toBe(true);
    expect(res.body.data.stats.dimensions).toEqual({ rows: 2, cols: 2 });
  });

  it('rejects invalid payload', async () => {
    await request(app).post('/api/v1/matrices/stats').send({}).expect(400);
  });
});
