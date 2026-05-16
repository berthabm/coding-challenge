import { useState } from 'react';
import { login } from '../api';
import logo from '../assets/interseguro.jpg';

interface Props {
  onLogin: (token: string) => void;
}

export default function Login({ onLogin }: Props) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError]     = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const token = await login(username, password);
      onLogin(token);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error desconocido');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-wrapper">
      <div className="login-card">
        <div className="login-logo-container">
          <img src={logo} alt="Interseguro" className="login-logo" />
        </div>
        <div className="login-brand">
          <h1 className="login-title">Interseguro</h1>
          <div className="login-divider" />
          <p>QR Matrix Processing Platform</p>
        </div>

        <form onSubmit={handleSubmit}>
          <label htmlFor="username">Usuario</label>
          <input
            id="username"
            type="text"
            value={username}
            onChange={e => setUsername(e.target.value)}
            autoComplete="username"
            required
            disabled={loading}
            placeholder="admin"
          />

          <label htmlFor="password">Contraseña</label>
          <input
            id="password"
            type="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            autoComplete="current-password"
            required
            disabled={loading}
            placeholder="••••••••"
          />

          {error && <p className="login-error">{error}</p>}

          <div className="login-btn-wrap">
            <button className="btn btn-primary" type="submit" disabled={loading}>
              {loading ? <span className="spinner" /> : null}
              {loading ? 'Autenticando…' : 'Iniciar sesión'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
