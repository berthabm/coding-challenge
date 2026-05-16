# QR Matrix Processing Platform

**Coding Challenge — Interseguro**

Plataforma distribuida de procesamiento matricial con descomposición QR, estadísticas descriptivas, autenticación JWT y dashboard corporativo. Dos APIs independientes (Go + Node.js) coordinadas por HTTP, con frontend React y despliegue cloud-ready en Google Cloud Run.

---

## Índice

1. [Descripción](#1-descripción)
2. [Arquitectura](#2-arquitectura)
3. [Stack tecnológico](#3-stack-tecnológico)
4. [Estructura del proyecto](#4-estructura-del-proyecto)
5. [Prerrequisitos](#5-prerrequisitos)
6. [Ejecución local](#6-ejecución-local)
7. [Docker](#7-docker)
8. [Tests](#8-tests)
9. [Autenticación JWT](#9-autenticación-jwt)
10. [Endpoints](#10-endpoints)
11. [Ejemplos request / response](#11-ejemplos-request--response)
12. [Decisiones técnicas](#12-decisiones-técnicas)
13. [Supuestos](#13-supuestos)
14. [Mejoras futuras](#14-mejoras-futuras)
15. [Cloud Run — Google Cloud](#15-cloud-run--google-cloud)

---

> Para casos de prueba manuales con entradas y salidas esperadas (login JWT, matrices válidas e inválidas, comandos curl listos para usar), revisar [docs/test-cases.md](docs/test-cases.md).

---

## 1. Descripción

El proyecto implementa un sistema de dos APIs que se comunican entre sí para realizar una operación de álgebra lineal completa:

| Servicio | Rol |
|----------|-----|
| **api-go** (Fiber v2, puerto 8080) | Recibe la matriz del cliente, la valida, ejecuta la **descomposición QR** con `gonum/mat` (LAPACK) y envía Q y R a api-node. Ensambla la respuesta final. |
| **api-node** (Express 4, puerto 3001) | Recibe Q y R, los valida con Joi y calcula **estadísticas descriptivas** sobre los valores combinados: max, min, promedio, suma, si Q es diagonal y si R es diagonal. |
| **frontend** (React + Vite, puerto 5173) | Dashboard corporativo con login JWT, entrada de matriz, procesamiento y visualización de resultados. |

El sistema es **stateless**: sin base de datos, sin colas. Cada petición es un ciclo completo e independiente.

---

## 2. Arquitectura

```
  Browser / Cliente HTTP
          │
          │  POST /api/auth/login   (pública)
          │  POST /api/qr           (Bearer JWT)
          ▼
┌─────────────────────────────────────────────────┐
│                  api-go :8080                   │
│               (Go + Fiber v2)                   │
│                                                 │
│  1. CORS middleware                             │
│  2. JWT middleware (rutas protegidas)           │
│  3. ValidateMatrix → rechazo temprano           │
│  4. QRService.Decompose (gonum/mat, LAPACK)     │
│  5. round3() → Q y R a 3 decimales             │
│  6. StatsClient.PostStats ─────────────────►   │
└───────────────────────────────┬─────────────────┘
                                │  POST /api/stats
                                │  { "q": [...], "r": [...] }
                                ▼
┌─────────────────────────────────────────────────┐
│               api-node :3001                    │
│            (Node.js + Express 4)                │
│                                                 │
│  1. Joi validation (q, r required)              │
│  2. flatten(q) + flatten(r) → valores           │
│  3. max, min, avg, sum, isQDiagonal, isRDiag.   │
│  4. round3() → todos a 3 decimales              │
└───────────────────────────────┬─────────────────┘
                                │  { "statistics": { ... } }
                                ▼
  Cliente recibe:
  {
    "originalMatrix": [[...]],
    "qr": { "q": [[...]], "r": [[...]] },
    "statistics": {
      "max": ..., "min": ..., "average": ..., "sum": ...,
      "isQDiagonal": false, "isRDiagonal": false
    }
  }
```

**Flujo de comunicación:**
- En **local**: `api-go → http://localhost:3001`
- En **Docker**: `api-go → http://api-node:3001` (DNS interno de Compose)
- En **Cloud Run**: `api-go → <URL pública de api-node>`

---

## 3. Stack tecnológico

| Capa | Tecnología | Versión | Por qué |
|------|-----------|---------|---------|
| API 1 | **Go + Fiber v2** | Go 1.22 | Cómputo numérico denso, binario estático, alta concurrencia |
| QR | **gonum/mat** (LAPACK) | 0.15.x | Implementación madura, numéricamente estable |
| Auth | **golang-jwt/jwt v5** | v5 | JWT HS256, biblioteca oficial de la comunidad Go |
| Logging Go | **log/slog** (stdlib) | Go 1.21+ | JSON estructurado sin dependencias extra |
| API 2 | **Node.js + Express 4** | Node ≥20 | Validación declarativa, ideal para servicios de métricas |
| Validación | **Joi 17** | 17.x | Esquemas declarativos, mensajes descriptivos, testeable |
| Logging Node | **pino + pino-http** | 9.x | JSON de alta performance, correlación por request |
| Tests Node | **Jest 29 + Supertest** | 29.x | Unitarios e integración HTTP |
| Tests Go | **testing** (stdlib) | — | Sin dependencias adicionales |
| Frontend | **React 18 + Vite 5 + TypeScript** | — | SPA moderna, HMR, tipado estático |
| Servidor estático | **nginx Alpine** | stable | SPA fallback, gzip, headers de seguridad |
| Contenedores | **Docker multi-stage + Compose v2** | Docker 24+ | Builds reproducibles, red aislada `challenge-net` |

---

## 4. Estructura del proyecto

```
.
├── api-go/                          # Servicio 1 — Go + Fiber
│   ├── Dockerfile
│   ├── go.mod
│   ├── cmd/server/main.go
│   └── internal/
│       ├── clients/                 # HTTP client → api-node
│       ├── config/                  # Variables de entorno (JWT, CORS, etc.)
│       ├── controllers/             # Handlers: auth, matrix
│       ├── middleware/              # ErrorHandler, RequestLogger, JWTAuth, CORS
│       ├── models/                  # DTOs request/response
│       ├── routes/                  # DI root + registro de rutas
│       ├── services/                # QRService, MatrixService (con round3)
│       └── utils/                   # Validador, logger, errores
│
├── api-node/                        # Servicio 2 — Node.js + Express
│   ├── Dockerfile
│   ├── package.json
│   ├── src/
│   │   ├── app.js                   # Factory Express (testeable sin puerto)
│   │   ├── index.js
│   │   ├── controllers/
│   │   ├── middleware/              # Joi middleware + ErrorHandler
│   │   ├── routes/
│   │   ├── services/                # qrStatsService, matrixStatsService
│   │   └── validations/            # Esquemas Joi
│   └── tests/
│       ├── integration/
│       ├── services/
│       └── validations/
│
├── frontend/                        # Dashboard React + Vite + TypeScript
│   ├── Dockerfile                   # Multi-stage: build → nginx
│   ├── nginx.conf                   # SPA fallback + seguridad
│   ├── src/
│   │   ├── App.tsx                  # Router token-based
│   │   ├── api.ts                   # login(), processMatrix()
│   │   ├── index.css                # Design system corporativo
│   │   ├── components/MatrixTable.tsx
│   │   ├── pages/Login.tsx
│   │   ├── pages/Dashboard.tsx
│   │   └── utils/format.ts         # formatNumber() — 3 decimales
│
├── docs/
│   ├── architecture.md
│   └── samples/                    # JSON de ejemplo por endpoint
│
├── docker-compose.yml
├── compose.env.example
└── README.md
```

---

## 5. Prerrequisitos

### Ejecución local

| Herramienta | Mínimo |
|-------------|--------|
| Go | 1.22 |
| Node.js | 20.x |
| npm | 10.x |

### Docker

| Herramienta | Mínimo |
|-------------|--------|
| Docker Engine | 24.x |
| Docker Compose | v2 (plugin) |

> Usar siempre `docker compose` (plugin v2). La CLI clásica `docker-compose` v1 no es compatible.

---

## 6. Ejecución local

### 1 — api-node (puerto 3001)

```bash
cd api-node
npm install
npm run dev        # hot-reload con node --watch
```

### 2 — api-go (puerto 8080)

```bash
cd api-go
go mod tidy
go run ./cmd/server
```

### 3 — frontend (puerto 5173)

```bash
cd frontend
npm install
npm run dev
```

### Variables de entorno relevantes

| Variable | Servicio | Descripción | Defecto |
|----------|----------|-------------|---------|
| `APP_PORT` | ambos | Puerto de escucha | `8080` / `3001` |
| `LOG_LEVEL` | ambos | `debug`, `info`, `warn` | `info` |
| `API2_BASE_URL` | api-go | URL de api-node | `http://localhost:3001` |
| `API2_TIMEOUT_SECONDS` | api-go | Timeout HTTP saliente | `10` |
| `AUTH_USERNAME` | api-go | Usuario de login | `admin` |
| `AUTH_PASSWORD` | api-go | Contraseña de login | `admin123` |
| `JWT_SECRET` | api-go | Clave HMAC-SHA256 | `change-me-in-production` |
| `CORS_ALLOWED_ORIGINS` | api-go | Orígenes permitidos (coma-separated) | `http://localhost:5173` |
| `VITE_API_BASE_URL` | frontend | URL de api-go (build-time) | `http://localhost:8080` |

### Verificación rápida

```bash
curl http://localhost:8080/health
curl http://localhost:3001/health
```

---

## 7. Docker

### Levantar todo el stack

```bash
docker compose up --build
```

```bash
# Modo background
docker compose up -d --build
```

Compose levanta los servicios en orden con healthchecks:

```
api-node → (healthy) → api-go → (healthy) → frontend
```

### Personalizar puertos y credenciales

```bash
cp compose.env.example .env
# Editar .env con los valores deseados
docker compose up --build
```

**Variables disponibles en `compose.env.example`:**

```env
API_GO_PORT=8080
API_NODE_PORT=3001
FRONTEND_PORT=5173
NODE_ENV=production
LOG_LEVEL=info
API2_TIMEOUT_SECONDS=10

AUTH_USERNAME=admin
AUTH_PASSWORD=admin123
JWT_SECRET=change-me-in-production   # ⚠️ cambiar en producción

CORS_ALLOWED_ORIGINS=http://localhost:5173
VITE_API_BASE_URL=http://localhost:8080
```

### URLs de acceso

| Servicio | URL local |
|----------|-----------|
| Frontend (dashboard) | http://localhost:5173 |
| API Go | http://localhost:8080 |
| API Node | http://localhost:3001 |

### Detener

```bash
docker compose down
```

---

## 8. Tests

### Node.js — Jest + Supertest

```bash
cd api-node
npm install
npm test                 # ejecución única (CI)
npm run test:watch       # modo watch (dev)
```

| Suite | Cobertura |
|-------|-----------|
| `qrValidation.test.js` | Esquema Joi: campos requeridos, tipos, matrices irregulares, `originalMatrix` opcional |
| `qrStatsService.test.js` | max, min, avg, sum, `isQDiagonal`, `isRDiagonal` sobre Q y R; redondeo a 3 decimales |
| `matrixStatsService.test.js` | Estadísticas sobre matriz directa |
| `statsRoutes.test.js` | Integración: `POST /api/stats` (Supertest) |
| `matrixRoutes.test.js` | Integración: `POST /api/v1/matrices/stats` (Supertest) |

**Resultado esperado:** 21 tests en 5 suites ✓

### Go — testing stdlib

```bash
cd api-go
go test ./...

# Subset específico
go test ./internal/utils/ ./internal/services/ ./internal/controllers/
```

| Suite | Cobertura |
|-------|-----------|
| `validator_test.go` | Dimensiones, valores no numéricos, casos límite |
| `qr_service_test.go` | Propiedad A ≈ Q·R con tolerancia `float64` |
| `matrix_service_test.go` | Orquestación completa con mock del cliente HTTP |
| `matrix_controller_test.go` | Handlers Fiber con peticiones HTTP de prueba |

**Resultado esperado:** 3 paquetes pasan ✓

---

## 9. Autenticación JWT

api-go implementa autenticación **JWT HS256** con tokens de 24 horas de validez.

### Obtener token

```bash
curl -sS -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

```json
{ "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." }
```

### Usar el token en peticiones protegidas

```bash
curl -X POST http://localhost:8080/api/qr \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"matrix": [[1,2],[3,4]]}'
```

### Rutas públicas vs protegidas

| Ruta | Auth requerida |
|------|----------------|
| `GET /health` | No |
| `POST /api/auth/login` | No |
| `POST /api/qr` | **Bearer JWT** |
| `POST /api/v1/matrices/qr` | **Bearer JWT** |

### Uso con Postman

1. `POST /api/auth/login` → copiar campo `token`
2. En cualquier petición protegida → pestaña **Authorization** → tipo **Bearer Token** → pegar token

### Configuración

```env
AUTH_USERNAME=admin
AUTH_PASSWORD=admin123
JWT_SECRET=mi-secreto-seguro-32-chars   # openssl rand -hex 32
```

---

## 10. Endpoints

### api-go — puerto 8080

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| `GET` | `/health` | — | Liveness check |
| `POST` | `/api/auth/login` | — | Devuelve JWT (24h). Body: `{"username","password"}` |
| `POST` | `/api/qr` | JWT | Descomposición QR + estadísticas sobre Q y R. |
| `POST` | `/api/v1/matrices/qr` | JWT | Igual, respuesta envuelta en `{"success":true,"data":{...}}` |

### api-node — puerto 3001

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/health` | Liveness check |
| `POST` | `/api/stats` | Interno (llamado por api-go). Acepta `{q, r}`, calcula estadísticas. |
| `POST` | `/api/v1/matrices/stats` | Estadísticas directas sobre una matriz `{matrix}`. |

### Formato de error uniforme

```json
{
  "success": false,
  "error": {
    "code": "INVALID_MATRIX",
    "message": "matrix must have at least 1 row and 1 column"
  }
}
```

| Código | HTTP | Cuándo |
|--------|------|--------|
| `INVALID_MATRIX` | 422 | Matriz malformada o filas inconsistentes |
| `QR_DECOMPOSITION_FAILED` | 422 | `filas < columnas` o error numérico |
| `DOWNSTREAM_UNAVAILABLE` | 502 | api-node no responde |
| `INVALID_CREDENTIALS` | 401 | Usuario o contraseña incorrectos |
| `UNAUTHORIZED` | 401 | Token ausente, expirado o inválido |
| `VALIDATION_ERROR` | 400 | Payload inválido en api-node |

---

## 11. Ejemplos request / response

### Login

```bash
curl -sS -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

```json
{ "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." }
```

---

### POST /api/qr — Flujo completo

```bash
TOKEN=$(curl -sS -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

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

**Response 200:**

```json
{
  "originalMatrix": [
    [12, -51,   4],
    [ 6, 167, -68],
    [-4,  24, -41]
  ],
  "qr": {
    "q": [
      [ 0.857, -0.394,  0.331],
      [ 0.429,  0.902, -0.034],
      [-0.286,  0.171,  0.943]
    ],
    "r": [
      [ 14,  21,   -2],
      [  0, 175,  -35],
      [  0,   0,  -98]
    ]
  },
  "statistics": {
    "max": 175,
    "min": -98,
    "average": 8.433,
    "sum": 253,
    "isQDiagonal": false,
    "isRDiagonal": false
  }
}
```

> Q y R se presentan redondeados a **3 decimales**. Las estadísticas se calculan sobre los **valores combinados de Q y R**.

---

### Error — token ausente

```json
{
  "success": false,
  "error": { "code": "UNAUTHORIZED", "message": "missing or invalid authorization token" }
}
```

---

### Error — matriz ancha (filas < columnas)

```bash
curl -sS -X POST http://localhost:8080/api/qr \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"matrix":[[1,2,3]]}'
```

```json
{
  "success": false,
  "error": { "code": "QR_DECOMPOSITION_FAILED", "message": "QR decomposition requires rows >= columns" }
}
```

---

### POST /api/stats — Debug directo a api-node

```bash
curl -sS -X POST http://localhost:3001/api/stats \
  -H "Content-Type: application/json" \
  -d '{ "q": [[1,0],[0,1],[0,0]], "r": [[2,3],[0,1]] }'
```

```json
{
  "statistics": {
    "max": 3, "min": 0, "average": 0.7, "sum": 7,
    "isQDiagonal": false, "isRDiagonal": false
  }
}
```

---

## 12. Decisiones técnicas

### Go para QR, Node para estadísticas

La descomposición QR requiere aritmética de punto flotante densa y numéricamente estable. Go + `gonum/mat` (LAPACK) proporciona resultados reproducibles con tipado estático. Node.js es más natural para el rol de "servicio de reporting": validación declarativa con Joi y stack de test rápido con Jest. La separación define una **frontera de contrato clara** equivalente al modelo de microservicios con equipos paralelos.

### Estadísticas sobre Q y R

Las estadísticas se calculan sobre los **valores combinados de Q y R** — el resultado directo del algoritmo de Go. Esto cumple el enunciado del reto: api-node recibe el resultado de api-go y realiza una operación adicional sobre esos datos. La matriz original se incluye en la respuesta como referencia de entrada.

### Redondeo a 3 decimales

- **Go** (`matrix_service.go`): `roundMatrix()` redondea Q y R antes de la respuesta final al cliente.
- **Node** (`qrStatsService.js`): `round3()` redondea max, min, average y sum de estadísticas.

Garantiza respuestas limpias independientemente de la precisión interna de LAPACK.

### JWT HS256

HMAC-SHA256 con secreto configurable por variable de entorno. Tokens de 24 horas. Sin base de datos. Apropiado para un único usuario administrador.

### CORS configurable

`CORS_ALLOWED_ORIGINS` es una variable de entorno. En local apunta a `http://localhost:5173`. En Cloud Run se actualiza con la URL del frontend desplegado.

### Validación en doble capa

1. **Go** valida la matriz antes de invocar gonum (dimensiones, valores numéricos, filas regulares).
2. **Node** valida Q y R con Joi. Sin confianza implícita en el upstream.

### Orden de arranque con healthchecks

`depends_on: condition: service_healthy` garantiza que `api-go` solo arranque cuando `api-node` responde en `/health`, y `frontend` espera a `api-go`.

### Multi-stage Dockerfiles + usuario no-root

Imágenes finales sin compiladores ni código fuente. El usuario `app` (uid 10001) se crea **antes** de `WORKDIR` para evitar el error `unknown user/group` de Alpine.

### Tabla de decisiones de librerías

| Decisión | Alternativa | Razón |
|----------|-------------|-------|
| `gonum/mat` QR | Gram-Schmidt manual | Auditado, LAPACK-backed, menos bugs numéricos |
| `golang-jwt/jwt` v5 | `dgrijalva/jwt-go` (deprecado) | Mantenida, API idiomática |
| `log/slog` stdlib | `zap`, `logrus` | JSON estructurado sin dependencias extra |
| `pino-http` | `morgan` | Correlación request/log, formato JSON nativo |
| Joi | Zod, Yup | Ecosistema maduro, mensajes descriptivos |
| CSS puro | Tailwind, MUI | Sin dependencias UI, diseño corporativo propio |

---

## 13. Supuestos

1. Matriz en cuerpo JSON (`Content-Type: application/json`). Sin multipart ni query params.
2. QR requiere `filas ≥ columnas`. Matrices anchas devuelven `422`.
3. Credenciales en variables de entorno. Sin base de datos. Un único usuario administrador.
4. `JWT_SECRET` por defecto solo para desarrollo. Producción requiere secreto fuerte.
5. Sin persistencia. Cada petición es stateless.
6. `api-node` es un servicio interno. Sin protección adicional más allá de la red Docker.
7. Logs como único mecanismo de observabilidad. Sin métricas ni trazas distribuidas.

---

## 14. Mejoras futuras

1. **OpenAPI/Swagger** — Contrato declarativo compartido entre servicios.
2. **Tests E2E en CI** — `docker compose run` efímero + curl en cada PR.
3. **Soporte matrices anchas** — QR con formulación alternativa o pseudoinversa.
4. **Autenticación multi-usuario** — Refresh tokens, revocación, base de datos.
5. **Observabilidad** — Prometheus/Grafana + OpenTelemetry en el hop Go → Node.
6. **Rate limiting** — Protección ante abuso en exposición pública.
7. **Escaneo de imágenes** — Trivy en CI.
8. **Tests del frontend** — Vitest + React Testing Library.
9. **Secret Manager (GCP)** — Variables sensibles fuera de las env vars planas.

---

## 15. Cloud Run — Google Cloud

El reto solicita utilizar servicios en la nube. La estrategia elegida es **Google Cloud Run**: contenedores serverless completamente gestionados, sin administrar infraestructura.

### Arquitectura cloud

```
  Browser
    │
    ▼
┌──────────────────────────────────┐
│  Cloud Run — frontend            │
│  https://frontend-xxx.run.app    │
│  (nginx, puerto 80)              │
└──────────────┬───────────────────┘
               │  fetch(VITE_API_BASE_URL)
               ▼
┌──────────────────────────────────┐
│  Cloud Run — api-go              │
│  https://api-go-xxx.run.app      │
│  JWT + CORS + QR                 │
└──────────────┬───────────────────┘
               │  API2_BASE_URL  (NO localhost)
               ▼
┌──────────────────────────────────┐
│  Cloud Run — api-node            │
│  https://api-node-xxx.run.app    │
│  Estadísticas Q y R              │
└──────────────────────────────────┘
```

> `api-go` **no usa** `localhost` para llamar a `api-node` en cloud. Usa la URL pública configurada en `API2_BASE_URL`.

### Prerrequisitos

```bash
gcloud auth login
gcloud config set project <PROJECT_ID>
gcloud services enable artifactregistry.googleapis.com run.googleapis.com
gcloud auth configure-docker <REGION>-docker.pkg.dev
```

### Placeholders de sustitución

| Variable | Descripción |
|----------|-------------|
| `<PROJECT_ID>` | ID del proyecto GCP |
| `<REGION>` | Región, ej. `us-central1` |
| `<REPO>` | Nombre del repositorio Artifact Registry |
| `<API_NODE_URL>` | URL obtenida al desplegar api-node |
| `<API_GO_URL>` | URL obtenida al desplegar api-go |
| `<FRONTEND_URL>` | URL obtenida al desplegar frontend |
| `<JWT_SECRET>` | `openssl rand -hex 32` |

### Paso 1 — Artifact Registry

```bash
gcloud artifacts repositories create <REPO> \
  --repository-format=docker --location=<REGION>
```

### Paso 2 — api-node

```bash
docker build -t <REGION>-docker.pkg.dev/<PROJECT_ID>/<REPO>/api-node:latest ./api-node
docker push  <REGION>-docker.pkg.dev/<PROJECT_ID>/<REPO>/api-node:latest

gcloud run deploy api-node \
  --image=<REGION>-docker.pkg.dev/<PROJECT_ID>/<REPO>/api-node:latest \
  --region=<REGION> --platform=managed --allow-unauthenticated --port=3001 \
  --set-env-vars="NODE_ENV=production,APP_PORT=3001"

# Guardar URL
gcloud run services describe api-node --region=<REGION> --format="value(status.url)"
```

### Paso 3 — api-go

```bash
docker build -t <REGION>-docker.pkg.dev/<PROJECT_ID>/<REPO>/api-go:latest ./api-go
docker push  <REGION>-docker.pkg.dev/<PROJECT_ID>/<REPO>/api-go:latest

gcloud run deploy api-go \
  --image=<REGION>-docker.pkg.dev/<PROJECT_ID>/<REPO>/api-go:latest \
  --region=<REGION> --platform=managed --allow-unauthenticated --port=8080 \
  --set-env-vars="\
APP_ENV=production,APP_PORT=8080,\
API2_BASE_URL=<API_NODE_URL>,API2_MATRIX_PATH=/api/stats,\
AUTH_USERNAME=admin,AUTH_PASSWORD=admin123,\
JWT_SECRET=<JWT_SECRET>,CORS_ALLOWED_ORIGINS=<FRONTEND_URL>"

gcloud run services describe api-go --region=<REGION> --format="value(status.url)"
```

### Paso 4 — frontend

```bash
docker build \
  --build-arg VITE_API_BASE_URL=<API_GO_URL> \
  -t <REGION>-docker.pkg.dev/<PROJECT_ID>/<REPO>/frontend:latest ./frontend
docker push <REGION>-docker.pkg.dev/<PROJECT_ID>/<REPO>/frontend:latest

gcloud run deploy frontend \
  --image=<REGION>-docker.pkg.dev/<PROJECT_ID>/<REPO>/frontend:latest \
  --region=<REGION> --platform=managed --allow-unauthenticated --port=80

gcloud run services describe frontend --region=<REGION> --format="value(status.url)"
```

### Paso 5 — Actualizar CORS

```bash
gcloud run services update api-go \
  --region=<REGION> \
  --update-env-vars="CORS_ALLOWED_ORIGINS=<FRONTEND_URL>"
```

### Paso 6 — Verificar

```bash
curl https://api-node-xxx.run.app/health
curl https://api-go-xxx.run.app/health

TOKEN=$(curl -sS -X POST https://api-go-xxx.run.app/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

curl -sS https://api-go-xxx.run.app/api/qr \
  -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"matrix":[[1,2],[3,4]]}'
```

### Notas de producción

| Aspecto | Recomendación |
|---------|--------------|
| Secretos | **Secret Manager**: `gcloud secrets create jwt-secret` |
| Cold start | `--min-instances=1` si se necesita latencia constante |
| Red interna | VPC Connector + URLs `.internal` para tráfico privado |
| api-node privado | `--no-allow-unauthenticated` + autenticación service-to-service |
| CI/CD | Cloud Build o GitHub Actions con `gcloud run deploy` |

---

## Licencia

Este repositorio se presenta únicamente con fines de evaluación técnica (*coding challenge*).
