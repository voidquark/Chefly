FROM node:20-alpine AS frontend-builder
WORKDIR /build/frontend
COPY frontend/package*.json ./
RUN npm ci --only=production
COPY frontend/ ./
RUN npm run build

FROM golang:1.25.3-alpine3.22 AS backend-builder
RUN apk add --no-cache gcc musl-dev sqlite-dev git wget
WORKDIR /build/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
COPY --from=frontend-builder /build/frontend/dist ./frontend/dist
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o chefly

FROM alpine:3.22
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    sqlite-libs \
    && rm -rf /var/cache/apk/*
ARG USER_UID=1000
ARG USER_GID=1000
RUN addgroup -g ${USER_GID} chefly && \
    adduser -D -u ${USER_UID} -G chefly -h /app -s /bin/sh chefly
WORKDIR /app
COPY --from=backend-builder /build/backend/chefly .
RUN mkdir -p /app/data && \
    chown -R chefly:chefly /app
USER chefly:chefly
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENV PORT=8080 \
    HOST=0.0.0.0 \
    DB_PATH=/app/data/chefly.db \
    AUDIT_LOG_ENABLED=true \
    AUDIT_LOG_LEVEL=info \
    AUDIT_LOG_FORMAT=json \
    REGISTRATION_ENABLED=true \
    ENVIRONMENT=production \
    GIN_MODE=release \
    RECIPE_GENERATION_LIMIT=unlimited

LABEL maintainer="void@voidquark.com" \
      description="Chefly - AI Powered Recipe Generator" \
      version="1.0"

CMD ["/app/chefly"]
