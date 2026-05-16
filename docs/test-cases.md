# Test Cases - QR Matrix Processing

Este documento contiene casos de prueba con entradas y salidas esperadas para validar:

- Login JWT (correcto, incorrecto, token ausente/válido)
- Procesamiento QR (`POST /api/qr`)
- Estadísticas sobre Q y R (`isQDiagonal`, `isRDiagonal`, max, min, avg, sum)
- Matrices diagonales e identidad
- Matrices rectangulares (tall: `m ≥ n`)
- Errores de validación

> **Nota:** Los valores exactos de Q y R pueden variar en la última cifra decimal por precisión de LAPACK. Se indica la estructura esperada y las propiedades cualitativas. Las estadísticas se calculan sobre los valores combinados de Q y R, redondeados a 3 decimales.

---

## Tabla resumen

| # | Caso | Tipo | Input resumido | Resultado esperado |
|---|------|------|----------------|--------------------|
| 1 | Matriz básica 2×2 | Válido | `[[1,2],[3,4]]` | 200 — Q y R no diagonales |
| 2 | Matriz 3×3 | Válido | `[[1,2,3],[4,5,6],[7,8,9]]` | 200 — Q y R no diagonales |
| 3 | Matriz diagonal 3×3 | Válido | `[[1,0,0],[0,5,0],[0,0,9]]` | 200 — `isRDiagonal: true` |
| 4 | Matriz identidad 3×3 | Válido | `[[1,0,0],[0,1,0],[0,0,1]]` | 200 — `isQDiagonal: true`, `isRDiagonal: true` |
| 5 | Matriz rectangular 3×2 | Válido | `[[1,2],[3,4],[5,6]]` | 200 — `m > n` permitido |
| 6 | Matriz con negativos | Válido | `[[1,-2,3],[-4,5,-6],[7,-8,9]]` | 200 — estadísticas con negativos |
| 7 | Matriz con decimales | Válido | `[[1.5,2.2],[3.8,4.1]]` | 200 — valores redondeados a 3 dec. |
| 8 | JSON inválido | Inválido | `[[1,2],[3,4]` (truncado) | 400 — parse error |
| 9 | Filas de longitud distinta | Inválido | `[[1,2,3],[4,5],[6,7,8]]` | 422 — `INVALID_MATRIX` |
| 10 | Valor no numérico | Inválido | `[[1,2],[3,"hola"]]` | 422 — `INVALID_MATRIX` |
| 11 | Matriz vacía | Inválido | `[]` | 422 — `INVALID_MATRIX` |
| 12 | Login correcto | Auth | `{username, password}` válidos | 200 — devuelve `token` |
| 13 | Login incorrecto | Auth | `{username, password}` inválidos | 401 — `INVALID_CREDENTIALS` |
| 14 | Sin token en `/api/qr` | Auth | Sin `Authorization` header | 401 — `UNAUTHORIZED` |
| 15 | Token válido en `/api/qr` | Auth | `Authorization: Bearer <token>` | 200 — respuesta normal |

---

## Casos válidos

### Caso 1 — Matriz básica 2×2

- **Input:** `[[1, 2], [3, 4]]`
- **Endpoint:** `POST /api/qr` (Bearer JWT)
- **HTTP esperado:** `200 OK`
- **Propiedades:**
  - Q: matriz ortogonal 2×2 — **no diagonal**
  - R: matriz triangular superior 2×2 — **no diagonal** (tiene elementos fuera de la diagonal principal)
  - `isQDiagonal: false`
  - `isRDiagonal: false`

```json
{
  "matrix": [[1, 2], [3, 4]]
}
```

**Estructura esperada del response:**

```json
{
  "originalMatrix": [[1, 2], [3, 4]],
  "qr": {
    "q": [
      [-0.316, -0.949],
      [-0.949,  0.316]
    ],
    "r": [
      [-3.162, -4.427],
      [     0,  0.632]
    ]
  },
  "statistics": {
    "max": 0.632,
    "min": -4.427,
    "average": "...",
    "sum": "...",
    "isQDiagonal": false,
    "isRDiagonal": false
  }
}
```

---

### Caso 2 — Matriz 3×3 genérica

- **Input:** `[[1, 2, 3], [4, 5, 6], [7, 8, 9]]`
- **Endpoint:** `POST /api/qr`
- **HTTP esperado:** `200 OK`
- **Propiedades:**
  - Q y R no son diagonales
  - `isQDiagonal: false`
  - `isRDiagonal: false`
  - La matriz es singular (det = 0), pero QR sigue siendo válida

```json
{
  "matrix": [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
}
```

**Estructura esperada del response:**

```json
{
  "originalMatrix": [[1, 2, 3], [4, 5, 6], [7, 8, 9]],
  "qr": {
    "q": [["..."], ["..."], ["..."]],
    "r": [["..."], ["..."], ["..."]]
  },
  "statistics": {
    "max": "...",
    "min": "...",
    "average": "...",
    "sum": "...",
    "isQDiagonal": false,
    "isRDiagonal": false
  }
}
```

---

### Caso 3 — Matriz diagonal 3×3

- **Input:** `[[1, 0, 0], [0, 5, 0], [0, 0, 9]]`
- **Endpoint:** `POST /api/qr`
- **HTTP esperado:** `200 OK`
- **Propiedades:**
  - La descomposición QR de una matriz diagonal produce Q = ±I y R = ±A
  - `isQDiagonal: true` — Q resultante es (posiblemente) diagonal
  - `isRDiagonal: true` — R es diagonal (triangular superior sin off-diagonals)

```json
{
  "matrix": [[1, 0, 0], [0, 5, 0], [0, 0, 9]]
}
```

**Estructura esperada del response:**

```json
{
  "originalMatrix": [[1, 0, 0], [0, 5, 0], [0, 0, 9]],
  "qr": {
    "q": [
      [1, 0, 0],
      [0, 1, 0],
      [0, 0, 1]
    ],
    "r": [
      [1, 0, 0],
      [0, 5, 0],
      [0, 0, 9]
    ]
  },
  "statistics": {
    "max": 9,
    "min": 0,
    "average": "...",
    "sum": "...",
    "isQDiagonal": true,
    "isRDiagonal": true
  }
}
```

---

### Caso 4 — Matriz identidad 3×3

- **Input:** `[[1, 0, 0], [0, 1, 0], [0, 0, 1]]`
- **Endpoint:** `POST /api/qr`
- **HTTP esperado:** `200 OK`
- **Propiedades:**
  - QR(I) = I · I → Q = I, R = I
  - `isQDiagonal: true`
  - `isRDiagonal: true`
  - max = 1, min = 0

```json
{
  "matrix": [[1, 0, 0], [0, 1, 0], [0, 0, 1]]
}
```

**Estructura esperada del response:**

```json
{
  "originalMatrix": [[1, 0, 0], [0, 1, 0], [0, 0, 1]],
  "qr": {
    "q": [[1, 0, 0], [0, 1, 0], [0, 0, 1]],
    "r": [[1, 0, 0], [0, 1, 0], [0, 0, 1]]
  },
  "statistics": {
    "max": 1,
    "min": 0,
    "average": "...",
    "sum": "...",
    "isQDiagonal": true,
    "isRDiagonal": true
  }
}
```

---

### Caso 5 — Matriz rectangular 3×2 (tall)

- **Input:** `[[1, 2], [3, 4], [5, 6]]`
- **Endpoint:** `POST /api/qr`
- **HTTP esperado:** `200 OK` — `m (3) ≥ n (2)`, válido
- **Propiedades:**
  - Q resulta 3×2 (o 3×3 según implementación), R es 2×2 triangular superior
  - `isQDiagonal: false`
  - `isRDiagonal: false`

```json
{
  "matrix": [[1, 2], [3, 4], [5, 6]]
}
```

**Estructura esperada del response:**

```json
{
  "originalMatrix": [[1, 2], [3, 4], [5, 6]],
  "qr": {
    "q": [["..."], ["..."], ["..."]],
    "r": [["..."], ["..."]]
  },
  "statistics": {
    "max": "...",
    "min": "...",
    "average": "...",
    "sum": "...",
    "isQDiagonal": false,
    "isRDiagonal": false
  }
}
```

---

### Caso 6 — Matriz con negativos

- **Input:** `[[1, -2, 3], [-4, 5, -6], [7, -8, 9]]`
- **Endpoint:** `POST /api/qr`
- **HTTP esperado:** `200 OK`
- **Propiedades:**
  - Estadísticas con valores negativos — `min` negativo
  - `isQDiagonal: false`, `isRDiagonal: false`

```json
{
  "matrix": [[1, -2, 3], [-4, 5, -6], [7, -8, 9]]
}
```

**Estructura esperada del response:**

```json
{
  "originalMatrix": [[1, -2, 3], [-4, 5, -6], [7, -8, 9]],
  "qr": {
    "q": [["..."], ["..."], ["..."]],
    "r": [["..."], ["..."], ["..."]]
  },
  "statistics": {
    "max": "... (positivo)",
    "min": "... (negativo)",
    "average": "...",
    "sum": "...",
    "isQDiagonal": false,
    "isRDiagonal": false
  }
}
```

---

### Caso 7 — Matriz con decimales

- **Input:** `[[1.5, 2.2], [3.8, 4.1]]`
- **Endpoint:** `POST /api/qr`
- **HTTP esperado:** `200 OK`
- **Propiedades:**
  - Todos los valores de Q, R y estadísticas redondeados a **máximo 3 decimales**
  - `isQDiagonal: false`, `isRDiagonal: false`

```json
{
  "matrix": [[1.5, 2.2], [3.8, 4.1]]
}
```

**Estructura esperada del response:**

```json
{
  "originalMatrix": [[1.5, 2.2], [3.8, 4.1]],
  "qr": {
    "q": [
      [-0.367, -0.930],
      [-0.930,  0.367]
    ],
    "r": [
      [-4.085, -5.054],
      [     0,  0.049]
    ]
  },
  "statistics": {
    "max": 0.367,
    "min": -5.054,
    "average": "...",
    "sum": "...",
    "isQDiagonal": false,
    "isRDiagonal": false
  }
}
```

---

## Casos inválidos

### Caso 8 — JSON malformado

- **Input:** `[[1,2],[3,4]` (cuerpo JSON truncado, falta el cierre)
- **HTTP esperado:** `400 Bad Request`
- **Error esperado:** parse error del framework (Fiber) antes de llegar al handler

```json
{
  "success": false,
  "error": {
    "code": "INVALID_BODY",
    "message": "invalid JSON body"
  }
}
```

---

### Caso 9 — Filas de longitud distinta

- **Input:** `[[1, 2, 3], [4, 5], [6, 7, 8]]`
- **HTTP esperado:** `422 Unprocessable Entity`
- **Validación:** api-go detecta filas irregulares antes de invocar gonum

```json
{
  "matrix": [[1, 2, 3], [4, 5], [6, 7, 8]]
}
```

**Response esperado:**

```json
{
  "success": false,
  "error": {
    "code": "INVALID_MATRIX",
    "message": "all rows must have the same number of columns"
  }
}
```

---

### Caso 10 — Valor no numérico

- **Input:** `[[1, 2], [3, "hola"]]`
- **HTTP esperado:** `422 Unprocessable Entity`

```json
{
  "matrix": [[1, 2], [3, "hola"]]
}
```

**Response esperado:**

```json
{
  "success": false,
  "error": {
    "code": "INVALID_MATRIX",
    "message": "matrix values must be numeric"
  }
}
```

---

### Caso 11 — Matriz vacía

- **Input:** `[]`
- **HTTP esperado:** `422 Unprocessable Entity`

```json
{
  "matrix": []
}
```

**Response esperado:**

```json
{
  "success": false,
  "error": {
    "code": "INVALID_MATRIX",
    "message": "matrix must have at least 1 row and 1 column"
  }
}
```

---

### Caso adicional — Matriz ancha (columnas > filas)

- **Input:** `[[1, 2, 3]]` — 1 fila, 3 columnas (`m < n`)
- **HTTP esperado:** `422 Unprocessable Entity`

```json
{
  "matrix": [[1, 2, 3]]
}
```

**Response esperado:**

```json
{
  "success": false,
  "error": {
    "code": "QR_DECOMPOSITION_FAILED",
    "message": "QR decomposition requires rows >= columns"
  }
}
```

---

## Casos de autenticación JWT

### Caso 12 — Login correcto

- **Endpoint:** `POST /api/auth/login`
- **Auth:** No requerida
- **HTTP esperado:** `200 OK`

```bash
curl -sS -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

**Response esperado:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

> El token expira a las **24 horas** de su emisión.

---

### Caso 13 — Login incorrecto

- **Endpoint:** `POST /api/auth/login`
- **HTTP esperado:** `401 Unauthorized`

```bash
curl -sS -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "wrong"}'
```

**Response esperado:**

```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "invalid username or password"
  }
}
```

---

### Caso 14 — Request a /api/qr sin token

- **Endpoint:** `POST /api/qr`
- **HTTP esperado:** `401 Unauthorized`

```bash
curl -sS -X POST http://localhost:8080/api/qr \
  -H "Content-Type: application/json" \
  -d '{"matrix": [[1, 2], [3, 4]]}'
```

**Response esperado:**

```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "missing or invalid authorization token"
  }
}
```

---

### Caso 15 — Request a /api/qr con token válido

- **Endpoint:** `POST /api/qr`
- **HTTP esperado:** `200 OK`

```bash
# 1. Obtener token
TOKEN=$(curl -sS -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# 2. Usar token
curl -sS -X POST http://localhost:8080/api/qr \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"matrix": [[1, 2], [3, 4]]}'
```

**Response esperado:** `200 OK` con `originalMatrix`, `qr` y `statistics`.

---

## Comandos curl de referencia rápida

### Login y obtener token

```bash
curl -sS -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

### Guardar token en variable

```bash
TOKEN=$(curl -sS -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

echo "Token: $TOKEN"
```

### Procesar una matriz válida

```bash
curl -sS -X POST http://localhost:8080/api/qr \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "matrix": [
      [12, -51,   4],
      [ 6, 167, -68],
      [-4,  24, -41]
    ]
  }'
```

### Probar error — sin token

```bash
curl -sS -X POST http://localhost:8080/api/qr \
  -H "Content-Type: application/json" \
  -d '{"matrix": [[1, 2], [3, 4]]}'
```

### Probar error — matriz inválida (filas distintas)

```bash
curl -sS -X POST http://localhost:8080/api/qr \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"matrix": [[1, 2, 3], [4, 5], [6, 7, 8]]}'
```

### Probar error — matriz ancha

```bash
curl -sS -X POST http://localhost:8080/api/qr \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"matrix": [[1, 2, 3]]}'
```

### Probar estadísticas directamente en api-node

```bash
curl -sS -X POST http://localhost:3001/api/stats \
  -H "Content-Type: application/json" \
  -d '{
    "q": [[1, 0], [0, 1]],
    "r": [[2, 3], [0, 4]]
  }'
```

### Healthcheck de ambas APIs

```bash
curl http://localhost:8080/health
curl http://localhost:3001/health
```

---

## Propiedades cualitativas de referencia

| Matriz | `isQDiagonal` | `isRDiagonal` | Notas |
|--------|:-------------:|:-------------:|-------|
| Identidad n×n | `true` | `true` | QR(I) = I · I |
| Diagonal n×n | `true` (posible) | `true` | R absorbe la diagonal |
| Genérica densa | `false` | `false` | Q ortogonal, R triangular superior |
| Rectangular m×n (m>n) | `false` | `false` | Q tall, R cuadrada |
| Con negativos | `false` | `false` | Signos pueden variar |

> El convenio de signo de la descomposición QR no es único. Los valores exactos de Q y R pueden diferir entre implementaciones manteniendo la propiedad A = Q·R.
