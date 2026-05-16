import jwt from 'jsonwebtoken';

export async function postLogin(req, res, next) {
  try {
    const { username, password } = req.body ?? {};

    const validUser = process.env.AUTH_USERNAME ?? 'admin';
    const validPass = process.env.AUTH_PASSWORD ?? 'admin123';
    const secret    = process.env.JWT_SECRET     ?? 'change-me-in-production';

    if (!username || !password || username !== validUser || password !== validPass) {
      return res.status(401).json({
        success: false,
        error: { code: 'INVALID_CREDENTIALS', message: 'Invalid username or password' },
      });
    }

    const token = jwt.sign({ sub: username }, secret, { expiresIn: '24h' });
    return res.json({ token });
  } catch (err) {
    next(err);
  }
}
