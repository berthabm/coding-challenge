const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export interface QRResponse {
  originalMatrix: number[][];
  qr: { q: number[][]; r: number[][] };
  statistics: {
    max: number;
    min: number;
    average: number;
    sum: number;
    isQDiagonal: boolean;
    isRDiagonal: boolean;
  };
}

export async function login(username: string, password: string): Promise<string> {
  const res = await fetch(`${API_BASE}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });

  if (!res.ok) {
    const data = (await res.json().catch(() => ({}))) as { error?: { message?: string } };
    throw new Error(data?.error?.message ?? 'Credenciales inválidas');
  }

  const data = (await res.json()) as { token: string };
  return data.token;
}

export async function processMatrix(
  matrix: number[][],
  token: string,
): Promise<QRResponse> {
  const res = await fetch(`${API_BASE}/api/qr`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ matrix }),
  });

  if (!res.ok) {
    const data = (await res.json().catch(() => ({}))) as { error?: { message?: string } };
    throw new Error(data?.error?.message ?? `Error ${res.status}`);
  }

  return res.json() as Promise<QRResponse>;
}
