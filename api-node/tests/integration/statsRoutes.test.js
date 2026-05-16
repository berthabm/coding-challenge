import request from 'supertest';
import { createApp } from '../../src/app.js';

describe('POST /api/stats', () => {
  const app = createApp();

  const payload = {
    q: [
      [1, 0],
      [0, 1],
      [0, 0],
    ],
    r: [
      [1, 2],
      [0, 3],
    ],
  };

  it('returns statistics based on Q and R matrices', async () => {
    const res = await request(app).post('/api/stats').send(payload).expect(200);

    // combined values: 1,0,0,1,0,0 (Q) + 1,2,0,3 (R) = 10 values
    expect(res.body.statistics).toMatchObject({
      isQDiagonal: false,
      isRDiagonal: false,
    });
    expect(typeof res.body.statistics.max).toBe('number');
    expect(typeof res.body.statistics.min).toBe('number');
    expect(typeof res.body.statistics.average).toBe('number');
    expect(typeof res.body.statistics.sum).toBe('number');
  });

  it('accepts payload without originalMatrix', async () => {
    await request(app).post('/api/stats').send(payload).expect(200);
  });

  it('rejects missing q or r', async () => {
    await request(app).post('/api/stats').send({ q: [[1]] }).expect(400);
    await request(app).post('/api/stats').send({ r: [[1]] }).expect(400);
  });
});
