# QR Matrix Processing Platform

**Coding Challenge — Interseguro**

Plataforma distribuida de procesamiento matricial con descomposición QR, estadísticas descriptivas, autenticación JWT y dashboard corporativo. Compuesta por un frontend React/TypeScript, una API Node.js de orquestación y una API Go de cómputo matemático, todo desplegado en Google Cloud Run.

---

## Índice

1. [Descripción](#1-descripción)
2. [Arquitectura](#2-arquitectura)
3. [Stack tecnológico](#3-stack-tecnológico)
4. [Estructura del proyecto](#4-estructura-del-proyecto)
5. [Prerrequisitos](#5-prerrequisitos)
6. [Variables de entorno](#6-variables-de-entorno)
7. [Ejecución local](#7-ejecución-local)
8. [Docker](#8-docker)
9. [Tests](#9-tests)
10. [Autenticación JWT](#10-autenticación-jwt)
11. [Endpoints](#11-endpoints)
12. [Casos de prueba](#12-casos-de-prueba)
13. [Manejo de errores](#13-manejo-de-errores)
14. [Decisiones técnicas](#14-decisiones-técnicas)
15. [CI/CD — Cloud Build](#15-cicd--cloud-build)
16. [Despliegue — Google Cloud Run](#16-despliegue--google-cloud-run)
17. [Mejoras futuras](#17-mejoras-futuras)

---

> Para casos de prueba manuales con curl y salidas esperadas, revisar [docs/test-cases.md](docs/test-cases.md).

---

## 1. Descripción

El proyecto implementa un sistema distribuido de tres capas que realiza descomposición matricial QR completa:

| Capa | Servicio | Rol |
|------|----------|-----|
| **Presentación** | **frontend** (React + Vite, puerto 5173) | Dashboard corporativo con login JWT, entrada de matriz, procesamiento y visualización de resultados con estadísticas. |
| **Orquestación** | **api-node** (Express 4, puerto 3001) | Autenticación JWT, validación de entrada, descomposición QR (Gram-Schmidt) y cálculo de estadísticas descriptivas sobre Q y R. |
| **Cómputo matemático** | **api-go** (Fiber v2, puerto 8080) | Descomposición QR con `gonum/mat` (LAPACK), autenticación JWT y llamada a api-node para estadísticas. Servicio standalone disponible para integración directa. |

El sistema es **stateless**: sin base de datos, sin colas. Cada petición es un ciclo completo e independiente.

---

## 2. Arquitectura

```
  ┌─────────────────────────────────┐
  │           Browser               │
  └────────────────┬────────────────┘
                   │  interacción del usuario
                   ▼
  ┌─────────────────────────────────┐
  │     Frontend — React + Vite     │
  │         (nginx, puerto 8080)    │
  └────────────────┬────────────────┘
                   │  POST /api/auth/login
                   │  POST /api/qr  (Bearer JWT)
                   ▼
  ┌─────────────────────────────────┐
  │    api-node — Express.js        │
  │         (puerto 3001)           │
  │                                 │
  │  • Authentication / JWT         │
  │  • Request validation           │
  │  • CORS & error handling        │
  │  • QR orchestration             │
  │  • Statistics calculation       │
  └────────────────┬────────────────┘
                   │  Internal HTTP — POST /api/qr
                   ▼
  ┌─────────────────────────────────┐
  │     api-go — Fiber v2           │
  │         (puerto 8080)           │
  │                                 │
  │  • Matrix validation            │
  │  • QR decomposition (LAPACK)    │
  │  • Numerical processing         │
  └─────────────────────────────────┘
```

- **Frontend** es la interfaz pública del usuario. Gestiona autenticación, entrada de matrices y visualización de resultados.
- **api-node** actúa como gateway y orquestador: valida las peticiones, gestiona el ciclo de vida del JWT y calcula estadísticas descriptivas sobre los resultados.
- **api-go** encapsula el procesamiento matemático especializado: descomposición QR con `gonum/mat` (LAPACK), numéricamente estable y de alta performance.
- La comunicación entre servicios se realiza mediante **HTTP REST**.
- En **Cloud Run** cada componente se despliega como un servicio independiente con su propia URL pública.

---

## 3. Stack tecnológico

| Capa | Tecnología | Versión | Por qué |
|------|-----------|---------|---------|
| Frontend | **React 18 + TypeScript** | 18.x | SPA moderna, tipado estático, componentes reutilizables |
| Build tool | **Vite 5** | 5.x | HMR, builds rápidos, soporte nativo ESM |
| Estilos | **CSS puro** | — | Design system corporativo propio, sin dependencias UI |
| Servidor estático | **nginx Alpine** | stable | SPA fallback, gzip, headers de seguridad |
| API orquestación | **Node.js + Express 4** | Node ≥20 | Validación declarativa, JWT, QR Gram-Schmidt |
| Validación | **Joi 17** | 17.x | Esquemas declarativos, mensajes descriptivos |
| Auth Node | **jsonwebtoken 9** | 9.x | JWT HS256 en Node.js |
| Logging Node | **pino + pino-http** | 9.x | JSON estructurado, correlación por request |
| Tests Node | **Jest 29 + Supertest** | 29.x | Unitarios e integración HTTP |
| API cómputo | **Go + Fiber v2** | Go 1.22 | Cómputo numérico denso, binario estático |
| QR Go | **gonum/mat** (LAPACK) | 0.15.x | Implementación madura, numéricamente estable |
| Auth Go | **golang-jwt/jwt v5** | v5 | JWT HS256, biblioteca oficial de la comunidad Go |
| Logging Go | **log/slog** (stdlib) | Go 1.21+ | JSON estructurado sin dependencias extra |
| Tests Go | **testing** (stdlib) | — | Sin dependencias adicionales |
| Contenedores | **Docker multi-stage + Compose v2** | Docker 24+ | Builds reproducibles, red aislada `challenge-net` |
| Cloud | **Google Cloud Run** | — | Serverless, sin infraestructura que administrar |
| CI/CD | **Cloud Build** | — | Deploy automático en cada push a `master` |

---

## 4. Estructura del proyecto

```
coding-challenge/
├── api-go/                          # API Go — cómputo QR con LAPACK
│   ├── Dockerfile
│   ├── go.mod
│   ├── cmd/server/main.go           # Punto de entrada
│   └── internal/
│       ├── clients/                 # HTTP client → api-node/stats
│       ├── config/                  # Variables de entorno
│       ├── controllers/             # Handlers: auth, matrix
│       ├── middleware/              # CORS, JWT, ErrorHandler, Logger
│       ├── models/                  # DTOs request/response
│       ├── routes/                  # Registro de rutas con DI
│       ├── services/                # QRService (LAPACK), MatrixService
│       └── utils/                   # Validador, logger, errores
│
├── api-node/                        # API Node.js — orquestación y auth
│   ├── Dockerfile
│   ├── package.json
│   ├── src/
│   │   ├── app.js                   # Factory Express (testeable sin puerto)
│   │   ├── index.js                 # Punto de entrada
│   │   ├── controllers/             # auth, health, matrix, qr
│   │   ├── middleware/              # CORS, Joi validators, ErrorHandler
│   │   ├── routes/                  # auth, health, stats, matrix, qr
│   │   ├── services/                # qrDecompService, qrStatsService, matrixStatsService
│   │   └── validations/            # Esquemas Joi
│   └── tests/
│       ├── integration/             # statsRoutes, matrixRoutes (Supertest)
│       ├── services/                # qrStatsService, matrixStatsService
│       └── validations/            # qrValidation
│
├── frontend/                        # Dashboard React + Vite + TypeScript
│   ├── Dockerfile                   # Multi-stage: node build → nginx
│   ├── nginx.conf                   # SPA fallback, puerto 8080
│   ├── index.html
│   ├── public/favicon.jpg
│   └── src/
│       ├── App.tsx                  # Router token-based
│       ├── api.ts                   # login(), processMatrix()
│       ├── index.css                # Design system corporativo (glassmorphism)
│       ├── assets/interseguro.jpg   # Logo
│       ├── components/MatrixTable.tsx
│       ├── pages/Login.tsx
│       ├── pages/Dashboard.tsx
│       └── utils/format.ts         # formatNumber() — 3 decimales
│
├── docs/
│   ├── architecture.md
│   ├── test-cases.md               # 15 casos de prueba con curl
│   └── samples/                    # JSON de ejemplo por endpoint
│
├── docker-compose.yml
├── compose.env.example             # Plantilla de variables de entorno
└── README.md
```

---

## 5. Prerrequisitos

### Ejecución local sin Docker

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

## 6. Variables de entorno

### Archivo de plantilla

```bash
cp compose.env.example .env
```

El archivo `.env` es leído automáticamente por Docker Compose para interpolar valores en `docker-compose.yml`.

### Variables disponibles

| Variable | Servicio | Descripción | Defecto |
|----------|----------|-------------|---------|
| `API_GO_PORT` | compose | Puerto local de api-go | `8080` |
| `API_NODE_PORT` | compose | Puerto local de api-node | `3001` |
| `FRONTEND_PORT` | compose | Puerto local del frontend | `5173` |
| `NODE_ENV` | api-node | Entorno de ejecución | `production` |
| `LOG_LEVEL` | ambas APIs | `debug`, `info`, `warn` | `info` |
| `AUTH_USERNAME` | api-node, api-go | Usuario de login | `admin` |
| `AUTH_PASSWORD` | api-node, api-go | Contraseña de login | `admin123` |
| `JWT_SECRET` | api-node, api-go | Clave HMAC-SHA256 — **cambiar en producción** | `change-me-in-production` |
| `CORS_ALLOWED_ORIGINS` | api-node, api-go | Orígenes CORS (coma-separated) | `http://localhost:5173` |
| `API2_BASE_URL` | api-go | URL de api-node para estadísticas | `http://localhost:3001` |
| `API2_TIMEOUT_SECONDS` | api-go | Timeout HTTP saliente | `10` |
| `VITE_API_BASE_URL` | frontend | URL de api-node — **build-time** (Vite) | `http://localhost:8080` |

> `VITE_API_BASE_URL` se embebe en el bundle de Vite en tiempo de compilación. No es una variable de runtime. Debe pasarse como `--build-arg` al construir la imagen Docker del frontend.

---

## 7. Ejecución local

### Con Docker Compose (recomendado)

```bash
docker compose up --build
```

Compose levanta los servicios en orden usando healthchecks:

```
api-node → (healthy) → api-go → (healthy) → frontend
```

| Servicio | URL |
|----------|-----|
| Frontend | http://localhost:5173 |
| api-node | http://localhost:3001 |
| api-go | http://localhost:8080 |

### Sin Docker — cada servicio por separado

**1. api-node (puerto 3001)**

```bash
cd api-node
npm install
npm run dev        # hot-reload con node --watch
```

**2. api-go (puerto 8080)**

```bash
cd api-go
go mod tidy
go run ./cmd/server
```

**3. Frontend (puerto 5173)**

```bash
cd frontend
npm install
npm run dev
```

### Verificación rápida

```bash
curl http://localhost:3001/health
curl http://localhost:8080/health
```

---

## 8. Docker

### Levantar el stack

```bash
# Modo interactivo
docker compose up --build

# Modo background
docker compose up -d --build

# Detener
docker compose down
```

### Estructura de imágenes

| Servicio | Base build | Base runtime | Resultado |
|----------|-----------|--------------|-----------|
| api-go | `golang:1.22-alpine` | `alpine:3.19` | Binario estático ~10 MB |
| api-node | `node:20-alpine` | `node:20-alpine` | Deps de producción |
| frontend | `node:20-alpine` | `nginx:stable-alpine` | Archivos estáticos en nginx |

Todas las imágenes usan **usuario no-root** (`app`, uid 10001) y builds multi-stage para excluir compiladores y código fuente del artefacto final.

---

## 9. Tests

### Node.js — Jest + Supertest

```bash
cd api-node
npm install
npm test                 # ejecución única (CI)
npm run test:watch       # modo watch (dev)
```

| Suite | Archivo | Cobertura |
|-------|---------|-----------|
| Validación Joi | `tests/validations/qrValidation.test.js` | Campos requeridos, tipos, matrices irregulares |
| Servicio QR stats | `tests/services/qrStatsService.test.js` | max, min, avg, sum, `isQDiagonal`, `isRDiagonal`, redondeo |
| Servicio matrix stats | `tests/services/matrixStatsService.test.js` | Estadísticas sobre matriz directa |
| Integración stats | `tests/integration/statsRoutes.test.js` | `POST /api/stats` end-to-end |
| Integración matrix | `tests/integration/matrixRoutes.test.js` | `POST /api/v1/matrices/stats` end-to-end |

**Resultado esperado:** `21 passed, 5 suites` ✓

### Go — testing stdlib

```bash
cd api-go
go test ./...

# Paquetes específicos
go test ./internal/utils/
go test ./internal/services/
go test ./internal/controllers/
```

| Suite | Cobertura |
|-------|-----------|
| `validator_test.go` | Dimensiones, valores no numéricos, casos límite |
| `qr_service_test.go` | Propiedad A ≈ Q·R con tolerancia `float64` |
| `matrix_service_test.go` | Orquestación completa con mock del cliente HTTP |
| `matrix_controller_test.go` | Handlers Fiber con peticiones HTTP de prueba |
| `auth_controller_test.go` | Login correcto/incorrecto, JWT en rutas protegidas |

**Resultado esperado:** `3 packages — ok` ✓

---

## 10. Autenticación JWT

Ambas APIs implementan autenticación **JWT HS256** con tokens de 24 horas.

### Credenciales demo

| Campo | Valor |
|-------|-------|
| Usuario | `admin` |
| Contraseña | `admin123` |

> En producción, cambiar `AUTH_USERNAME`, `AUTH_PASSWORD` y `JWT_SECRET` por valores seguros.

### Obtener token

```bash
curl -sS -X POST https://api-node-995881892656.southamerica-east1.run.app/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

```json
{ "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." }
```

### Usar el token

```bash
curl -sS -X POST https://api-node-995881892656.southamerica-east1.run.app/api/qr \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"matrix": [[1,2],[3,4]]}'
```

### Rutas públicas vs protegidas

| Servicio | Ruta | Auth requerida |
|----------|------|----------------|
| api-node | `GET /health` | No |
| api-node | `POST /api/auth/login` | No |
| api-node | `POST /api/qr` | **Bearer JWT** |
| api-node | `POST /api/stats` | No (interno) |
| api-go | `GET /health` | No |
| api-go | `POST /api/auth/login` | No |
| api-go | `POST /api/qr` | **Bearer JWT** |

---

## 11. Endpoints

### api-node — puerto 3001 / Cloud Run

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| `GET` | `/health` | — | Liveness check |
| `POST` | `/api/auth/login` | — | Devuelve JWT 24h |
| `POST` | `/api/qr` | JWT | Descomposición QR + estadísticas |
| `POST` | `/api/stats` | — | Estadísticas sobre Q y R (interno) |
| `POST` | `/api/v1/matrices/stats` | — | Estadísticas directas sobre una matriz |

### api-go — puerto 8080 / Cloud Run

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| `GET` | `/health` | — | Liveness check |
| `POST` | `/api/auth/login` | — | Devuelve JWT 24h |
| `POST` | `/api/qr` | JWT | Descomposición QR con LAPACK + estadísticas (vía api-node) |
| `POST` | `/api/v1/matrices/qr` | JWT | Igual, respuesta envuelta en `{"success":true,"data":{...}}` |

### Formato de respuesta de error

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
| `INVALID_MATRIX` | 422 | Matriz malformada, filas inconsistentes o valores no numéricos |
| `QR_DECOMPOSITION_FAILED` | 422 | Matriz rank-deficient o `filas < columnas` |
| `DOWNSTREAM_UNAVAILABLE` | 502 | api-node no responde (solo en api-go) |
| `INVALID_CREDENTIALS` | 401 | Usuario o contraseña incorrectos |
| `UNAUTHORIZED` | 401 | Token ausente, expirado o inválido |
| `VALIDATION_ERROR` | 400 | Payload inválido según esquema Joi |

---

## 12. Casos de prueba

### Matrices válidas ✓

```json
[[1,0],[0,1]]
```
Matriz identidad 2×2. QR trivial: Q = I, R = I.

```json
[[1,2],[3,4]]
```
Matriz 2×2 de rango completo. Descomposición QR estándar.

```json
[[1,2,3],[0,1,4],[5,6,0]]
```
Matriz 3×3 — **ejemplo precargado en el dashboard**. Rango completo, QR válido.

```json
[[12,-51,4],[6,167,-68],[-4,24,-41]]
```
Ejemplo clásico de Golub & Van Loan. Referencia estándar para verificar correctitud numérica.

### Matrices inválidas ✗

| Entrada | Error | Código |
|---------|-------|--------|
| `[[1,2,3],[4,5,6],[7,8,9]]` | Matriz rank-deficient (fila 3 = fila1 + fila2) | `QR_DECOMPOSITION_FAILED` |
| `[[1,2,3],[4,5]]` | Filas de distinto largo | `INVALID_MATRIX` |
| `[[1,0],[0,"w"]]` | Valor no numérico | `INVALID_MATRIX` |
| `[[1,2,3]]` | 1 fila, 3 columnas — `filas < columnas` | `QR_DECOMPOSITION_FAILED` |
| `[]` | Matriz vacía | `INVALID_MATRIX` |

### Flujo completo con curl

```bash
# 1. Login
TOKEN=$(curl -sS -X POST https://api-node-995881892656.southamerica-east1.run.app/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# 2. Procesar matriz
curl -sS -X POST https://api-node-995881892656.southamerica-east1.run.app/api/qr \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"matrix":[[1,2,3],[0,1,4],[5,6,0]]}'
```

**Response 200:**

```json
{
  "originalMatrix": [[1,2,3],[0,1,4],[5,6,0]],
  "qr": {
    "q": [
      [-0.196, 0.169, 0.966],
      [-0.0,   0.985,-0.169],
      [-0.981,-0.034,-0.193]
    ],
    "r": [
      [-5.099,-7.06, -0.196],
      [0,      1.015, 3.896],
      [0,      0,     3.679]
    ]
  },
  "statistics": {
    "max": 3.896,
    "min": -7.06,
    "average": -0.282,
    "sum": -5.071,
    "isQDiagonal": false,
    "isRDiagonal": false
  }
}
```

---

## 13. Manejo de errores

### En el frontend

El dashboard muestra mensajes amigables según el tipo de error:

| Situación | Mensaje al usuario |
|-----------|-------------------|
| JSON inválido (`[[1,0],[0,w]]`) | *"El formato ingresado no es un JSON válido. Revisa comas, comillas y corchetes."* |
| Matriz mal formada (filas desiguales, valores no numéricos) | *"La matriz debe contener solo números y todas las filas deben tener la misma cantidad de columnas."* |
| Matriz rank-deficient | *"La matriz ingresada no es válida para descomposición QR porque sus filas o columnas son dependientes. Prueba con una matriz de rango completo."* |
| Error inesperado de red/servidor | *"Ocurrió un error al procesar la solicitud. Inténtalo de nuevo."* |

Los errores de validación de input del usuario no generan stack traces en consola — solo `console.warn` con el mensaje.

### En las APIs

Todas las APIs devuelven errores en formato uniforme `{ success: false, error: { code, message } }`. Nunca se exponen stack traces ni detalles internos en producción.

---

## 14. Decisiones técnicas

### Node.js como capa de orquestación

api-node es el punto de entrada único para el frontend. Concentra: CORS, autenticación JWT, validación de matrices, descomposición QR (Gram-Schmidt modificado) y cálculo de estadísticas. Esto permite desplegar el frontend + api-node como un sistema funcional completo en Cloud Run sin depender de api-go.

### api-go como servicio de cómputo de alta precisión

api-go usa `gonum/mat` (LAPACK) para la descomposición QR — más rápido y numéricamente estable que Gram-Schmidt en JavaScript para matrices grandes. Está disponible como servicio standalone para integración directa o para reemplazar el cómputo de api-node configurando `API_GO_URL`.

### Gram-Schmidt vs LAPACK

| | Gram-Schmidt (api-node) | LAPACK via gonum (api-go) |
|--|------------------------|--------------------------|
| Lenguaje | JavaScript | Go |
| Precisión | Alta (double) | Alta (double, pivoting) |
| Matrices grandes | Adecuado | Superior |
| Dependencias | Ninguna | `gonum/mat` |
| Resultado 3×3 | Idéntico | Idéntico |

Para matrices de tamaño moderado (≤ 100×100), los resultados son prácticamente equivalentes a 3 decimales.

### JWT HS256

HMAC-SHA256 con secreto configurable. Tokens de 24 horas. Sin base de datos. Un único usuario administrador. Mismo `JWT_SECRET` en ambas APIs para compatibilidad cruzada.

### CORS configurable

`CORS_ALLOWED_ORIGINS` es una variable de entorno en ambas APIs. Lista de orígenes separada por comas. Soporta múltiples orígenes para desarrollo local y producción simultáneamente.

### CSS puro — Design System corporativo

Sin frameworks UI. El diseño corporativo se implementa con variables CSS, glassmorphism, fuente Sofia Sans y paleta Interseguro. Reduce dependencias y da control total sobre cada pixel.

### Redondeo a 3 decimales

Aplicado en ambos servicios antes de devolver la respuesta. Garantiza resultados limpios independientemente de la precisión interna de los algoritmos.

---

## 15. CI/CD — Cloud Build

Cloud Build está conectado al repositorio GitHub `berthabm/coding-challenge`. Cada `git push` a `master` activa automáticamente el pipeline de build y deploy en Cloud Run.

```bash
# El pipeline se dispara automáticamente con:
git push origin master
```

### Flujo del pipeline

```
GitHub (master) → Cloud Build trigger
    │
    ├─ Build imagen api-node  → push a Artifact Registry
    ├─ Deploy api-node        → Cloud Run
    │
    ├─ Build imagen api-go    → push a Artifact Registry
    ├─ Deploy api-go          → Cloud Run
    │
    ├─ Build imagen frontend  → push a Artifact Registry
    └─ Deploy frontend        → Cloud Run
```

### Deploy manual (sin CI/CD)

```bash
REGION=southamerica-east1
PROJECT_ID=<tu-proyecto>
REPO=<tu-repo-artifact>

# api-node
docker build -t ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/api-node:latest ./api-node
docker push  ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/api-node:latest
gcloud run deploy api-node \
  --image=${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/api-node:latest \
  --region=${REGION} --platform=managed --allow-unauthenticated --port=3001 \
  --set-env-vars="NODE_ENV=production,AUTH_USERNAME=admin,AUTH_PASSWORD=admin123,JWT_SECRET=<secreto>,CORS_ALLOWED_ORIGINS=https://frontend-995881892656.southamerica-east1.run.app"

# api-go
docker build -t ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/api-go:latest ./api-go
docker push  ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/api-go:latest
gcloud run deploy api-go \
  --image=${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/api-go:latest \
  --region=${REGION} --platform=managed --allow-unauthenticated --port=8080 \
  --set-env-vars="AUTH_USERNAME=admin,AUTH_PASSWORD=admin123,JWT_SECRET=<secreto>,CORS_ALLOWED_ORIGINS=https://frontend-995881892656.southamerica-east1.run.app,API2_BASE_URL=https://api-node-995881892656.southamerica-east1.run.app"

# frontend
docker build \
  --build-arg VITE_API_BASE_URL=https://api-node-995881892656.southamerica-east1.run.app \
  -t ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/frontend:latest ./frontend
docker push ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/frontend:latest
gcloud run deploy frontend \
  --image=${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/frontend:latest \
  --region=${REGION} --platform=managed --allow-unauthenticated --port=8080
```

---

## 16. Despliegue — Google Cloud Run

### URLs en producción

| Servicio | URL |
|----------|-----|
| **Frontend** | https://frontend-995881892656.southamerica-east1.run.app |
| **API Node** | https://api-node-995881892656.southamerica-east1.run.app |
| **API Go** | https://api-go-995881892656.southamerica-east1.run.app |

### Credenciales demo

| Campo | Valor |
|-------|-------|
| Usuario | `admin` |
| Contraseña | `admin123` |

### Verificación del estado

```bash
curl https://api-node-995881892656.southamerica-east1.run.app/health
curl https://api-go-995881892656.southamerica-east1.run.app/health
```

### Arquitectura cloud

```
  Browser
    │
    ▼
┌──────────────────────────────────┐
│  Cloud Run — frontend            │
│  nginx, puerto 8080              │
│  https://frontend-xxx.run.app    │
└──────────────┬───────────────────┘
               │  VITE_API_BASE_URL
               ▼
┌──────────────────────────────────┐
│  Cloud Run — api-node            │
│  Express 4, puerto 3001          │
│  https://api-node-xxx.run.app    │
│  Auth · QR · Stats · CORS        │
└──────────────────────────────────┘

┌──────────────────────────────────┐
│  Cloud Run — api-go              │
│  Fiber v2, puerto 8080           │
│  https://api-go-xxx.run.app      │
│  QR LAPACK · Auth (standalone)   │
└──────────────────────────────────┘
```

### Notas de producción

| Aspecto | Recomendación |
|---------|--------------|
| Secretos | **Secret Manager**: `gcloud secrets create jwt-secret --data-file=-` |
| Cold start | `--min-instances=1` si se necesita latencia constante |
| Red privada | VPC Connector + URLs `.internal` para tráfico entre servicios |
| api-go privado | `--no-allow-unauthenticated` + Cloud Run service-to-service auth |
| Monitoreo | Cloud Logging + Cloud Monitoring (métricas automáticas de Cloud Run) |

---

## 17. Mejoras futuras

1. **JWT con refresh tokens** — Renovación automática de sesión sin re-login.
2. **Base de datos** — PostgreSQL/Firestore para usuarios, historial de matrices y audit trail.
3. **Roles y permisos** — Admin, readonly, servicio. Control de acceso granular.
4. **Observabilidad completa** — OpenTelemetry + Prometheus/Grafana. Trazas distribuidas en el hop Node → Go.
5. **Soporte matrices anchas** — QR con pseudo-inversa o descomposición alternativa para `filas < columnas`.
6. **Kubernetes** — Migrar a GKE para mayor control de recursos, autoescalado avanzado y red privada entre servicios.
7. **Pipeline cloudbuild.yaml avanzado** — Tests automatizados en CI, revisión de seguridad con Trivy, deploy canary.
8. **Tests del frontend** — Vitest + React Testing Library para componentes y flujos de usuario.
9. **OpenAPI / Swagger** — Contrato declarativo compartido entre servicios, generación de clientes.
10. **Rate limiting** — Protección ante abuso en endpoints públicos.
11. **Secret Manager integrado** — Variables sensibles fuera de las env vars planas en Cloud Run.

