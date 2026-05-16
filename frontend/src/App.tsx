import { useState } from 'react';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';

const TOKEN_KEY = 'interseguros_token';

function App() {
  const [token, setToken] = useState<string | null>(
    () => localStorage.getItem(TOKEN_KEY),
  );

  const handleLogin = (jwt: string) => {
    localStorage.setItem(TOKEN_KEY, jwt);
    setToken(jwt);
  };

  const handleLogout = () => {
    localStorage.removeItem(TOKEN_KEY);
    setToken(null);
  };

  if (!token) {
    return <Login onLogin={handleLogin} />;
  }

  return <Dashboard token={token} onLogout={handleLogout} />;
}

export default App;
