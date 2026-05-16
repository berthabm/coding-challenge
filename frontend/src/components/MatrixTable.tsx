import { formatNumber } from '../utils/format';

interface Props {
  matrix: number[][];
  label?: string;
}

export default function MatrixTable({ matrix, label }: Props) {
  return (
    <div>
      {label && <p className="matrix-label">{label}</p>}
      <div className="matrix-scroll">
        <table className="matrix-table">
          <thead>
            <tr>
              <th>#</th>
              {matrix[0].map((_, i) => <th key={i}>c{i}</th>)}
            </tr>
          </thead>
          <tbody>
            {matrix.map((row, r) => (
              <tr key={r}>
                <th>r{r}</th>
                {row.map((v, c) => <td key={c}>{formatNumber(v)}</td>)}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
