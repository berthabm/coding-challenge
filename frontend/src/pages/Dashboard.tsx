import { useState } from 'react';
import { processMatrix, type QRResponse } from '../api';
import MatrixTable from '../components/MatrixTable';
import { formatNumber } from '../utils/format';

const EXAMPLE = '[[1,2,3],[0,1,4],[5,6,0]]';

const MSG_INVALID_JSON =
  'El formato ingresado no es un JSON válido. ' +
  'Revisa comas, comillas y corchetes. ' +
  'Ejemplo válido: [[1,2],[3,4]]';

const MSG_INVALID_MATRIX =
  'La matriz debe contener solo números y todas las filas deben tener ' +
  'la misma cantidad de columnas. Ejemplo válido: [[1,2,3],[4,5,6]]';

const MSG_RANK_DEFICIENT =
  'La matriz ingresada no es válida para descomposición QR porque sus filas o columnas son ' +
  'dependientes. Prueba con una matriz de rango completo.';

const MSG_GENERIC =
  'Ocurrió un error al procesar la solicitud. Inténtalo de nuevo.';

interface Props {
  token: string;
  onLogout: () => void;
}

function parseMatrix(text: string): number[][] {
  let parsed: unknown;
  try {
    parsed = JSON.parse(text.trim());
  } catch (e) {
    console.error('[parseMatrix] JSON parse error:', e);
    throw new Error(MSG_INVALID_JSON);
  }

  if (
    !Array.isArray(parsed) ||
    parsed.length === 0 ||
    !Array.isArray(parsed[0])
  ) {
    throw new Error(MSG_INVALID_MATRIX);
  }

  const colCount = (parsed[0] as unknown[]).length;
  for (const row of parsed as unknown[]) {
    if (
      !Array.isArray(row) ||
      (row as unknown[]).length !== colCount ||
      !(row as unknown[]).every((v) => typeof v === 'number' && isFinite(v))
    ) {
      throw new Error(MSG_INVALID_MATRIX);
    }
  }

  return parsed as number[][];
}

export default function Dashboard({ token, onLogout }: Props) {
  const [raw, setRaw]         = useState(EXAMPLE);
  const [loading, setLoading] = useState(false);
  const [error, setError]     = useState('');
  const [result, setResult]   = useState<QRResponse | null>(null);

  const handleProcess = async () => {
    setError('');
    setLoading(true);
    try {
      const matrix = parseMatrix(raw);
      const data   = await processMatrix(matrix, token);
      setResult(data);
    } catch (err) {
      console.error('[handleProcess]', err);
      if (err instanceof Error) {
        // Messages from parseMatrix are already friendly; API messages need mapping
        const msg = err.message;
        if (
          msg === MSG_INVALID_JSON ||
          msg === MSG_INVALID_MATRIX
        ) {
          setError(msg);
        } else if (
          msg.toLowerCase().includes('rank') ||
          msg.toLowerCase().includes('rank-deficient') ||
          msg.toLowerCase().includes('decomposition failed')
        ) {
          setError(MSG_RANK_DEFICIENT);
        } else if (
          msg.toLowerCase().includes('invalid') ||
          msg.toLowerCase().includes('payload') ||
          msg.toLowerCase().includes('validation')
        ) {
          setError(MSG_INVALID_MATRIX);
        } else {
          setError(MSG_GENERIC);
        }
      } else {
        setError(MSG_GENERIC);
      }
      setResult(null);
    } finally {
      setLoading(false);
    }
  };

  const stats = result?.statistics;

  return (
    <>
      {/* Topbar */}
      <header className="topbar">
        <div>
          <div className="topbar-brand">Interseguro</div>
          <div className="topbar-sub">QR Matrix Processing Platform</div>
        </div>
        <div className="topbar-right">
          <button className="btn btn-danger btn-sm" onClick={onLogout}>
            Cerrar sesión
          </button>
        </div>
      </header>

      <main className="page-container">
        {/* Input */}
        <div className="card">
          <div className="card-title">Matriz de entrada</div>
          <textarea
            className="matrix-textarea"
            value={raw}
            onChange={e => setRaw(e.target.value)}
            spellCheck={false}
          />
          <p className="form-hint">
            Introduce un array de arrays JSON, p.ej.&nbsp;
            <code>[[1,2],[3,4]]</code>
          </p>
          <div className="form-actions">
            <button
              className="btn btn-secondary"
              onClick={() => setRaw(EXAMPLE)}
              disabled={loading}
            >
              Cargar ejemplo
            </button>
            <button
              className="btn btn-primary"
              onClick={handleProcess}
              disabled={loading || !raw.trim()}
            >
              {loading ? <span className="spinner" /> : null}
              {loading ? 'Procesando…' : 'Procesar QR'}
            </button>
          </div>
        </div>

        {error && (
          <div className="alert alert-error" role="alert">
            <div className="alert-icon" aria-hidden="true">⚠️</div>
            <div>
              <div className="alert-title">Error de validación</div>
              <div className="alert-body">{error}</div>
            </div>
          </div>
        )}

        {result && (
          <>
            {/* Statistics */}
            <div className="card">
              <div className="card-title">Estadísticas sobre Q y R</div>
              <div className="stat-grid">
                <div className="stat-item">
                  <div className="stat-label">Máximo QR</div>
                  <div className="stat-value">{formatNumber(stats!.max)}</div>
                </div>
                <div className="stat-item">
                  <div className="stat-label">Mínimo QR</div>
                  <div className="stat-value">{formatNumber(stats!.min)}</div>
                </div>
                <div className="stat-item">
                  <div className="stat-label">Media QR</div>
                  <div className="stat-value">{formatNumber(stats!.average)}</div>
                </div>
                <div className="stat-item">
                  <div className="stat-label">Suma QR</div>
                  <div className="stat-value">{formatNumber(stats!.sum)}</div>
                </div>
                <div className="stat-item">
                  <div className="stat-label">Q diagonal</div>
                  <div className="stat-value" style={{ paddingTop: '.35rem' }}>
                    {stats!.isQDiagonal ? (
                      <span className="badge badge-green">Sí</span>
                    ) : (
                      <span className="badge badge-red">No</span>
                    )}
                  </div>
                </div>
                <div className="stat-item">
                  <div className="stat-label">R diagonal</div>
                  <div className="stat-value" style={{ paddingTop: '.35rem' }}>
                    {stats!.isRDiagonal ? (
                      <span className="badge badge-green">Sí</span>
                    ) : (
                      <span className="badge badge-red">No</span>
                    )}
                  </div>
                </div>
              </div>
            </div>

            {/* Original matrix */}
            <div className="card">
              <div className="card-title">Matriz original</div>
              <MatrixTable matrix={result.originalMatrix} />
            </div>

            {/* QR decomposition */}
            <div className="card">
              <div className="card-title">Descomposición QR</div>
              <div className="qr-grid">
                <MatrixTable matrix={result.qr.q} label="Matriz Q (ortogonal)" />
                <MatrixTable matrix={result.qr.r} label="Matriz R (triangular superior)" />
              </div>
            </div>
          </>
        )}
      </main>
    </>
  );
}
