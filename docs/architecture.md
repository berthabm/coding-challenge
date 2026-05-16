# Arquitectura complementaria

> Documento corto complementario del `README.md` principal (fuente única para instrucciones y entrega).

## Visión general

```text
Cliente HTTP
    │
    ▼
┌─────────────┐       HTTP POST        ┌─────────────┐
│   api-go    │ ─────────────────────► │  api-node   │
│ Go + Fiber  │    /api/stats (+q,+r)  │Express+ Joi │
│  :8080      │                         │ :3001       │
└─────────────┘                         └─────────────┘

Local:   API2_BASE_URL=http://localhost:3001
Docker: API2_BASE_URL=http://api-node:3001
```

Flujo oficial del challenge: **`POST /api/qr`** (Go).

## Principios por servicio

| Aspecto | api-go | api-node |
|---------|--------|----------|
| Capas típicas | controllers → services → cliente HTTP (`clients`) | routes → middleware → controllers → services |
| Config | `internal/config`, variables `APP_*` y `API2_*` | `src/config` + `dotenv` |
| Errores | `utils.AppError` + Fiber `ErrorHandler` | `AppError` + middleware único |
| Logs | slog JSON (`RequestLogger`) | pino + pino-http |
| Tests (`*_test.go` / jest) | `internal/**`, comandos README | `tests/` |

## Paso de datos

1. **Entrada cliente** `{ "matrix": number[][] }` → validación **`ValidateMatrix`**.
2. **QR** mediante **`gonum.org/v1/gonum/mat`** (requiere `filas ≥ columnas`).
3. **Outbound** `{ "q","r" }` → cliente HTTP configurado contra `/api/stats`.
4. **Salida cliente:** el JSON estadístico de Node se reenvía al cliente en **`POST /api/qr`** sin enriquecimientos extra (contrato estable).

## Observabilidad frente errores típicos

| Situación | Comportamiento resumido |
|-----------|--------------------------|
| JSON inválido / datos no numéricos | 400 (Go) vs 400 Joi (Node) según ubicación fallo |
| Formatos donde **filas `<` columnas** | 422 **`QR_DECOMPOSITION_FAILED`** (servicio QR en Go rechaza ese caso antes de LAPACK por claridad mensaje) |

| Node caído / tiempo excedido | 502 **`DOWNSTREAM_UNAVAILABLE`** Go |
| Endpoint inexistente | 404 formato JSON estable en Node |

## Evoluciones posibles sin reescribir core

Montar mismo stack en Compose con diferentes `API_GO_PORT`; test E2E con `docker compose run` + script curl.
