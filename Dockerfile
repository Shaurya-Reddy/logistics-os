# syntax=docker/dockerfile:1.7

FROM node:24-alpine AS web-build
WORKDIR /src/web
COPY web/package.json ./
RUN npm install --ignore-scripts
COPY web/ ./
RUN npm run build

FROM golang:1.24-alpine AS go-build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/logistics-os ./cmd/logistics-os

FROM alpine:3.22
RUN addgroup -S logistics && adduser -S -G logistics logistics
WORKDIR /app
COPY --from=go-build /out/logistics-os /app/logistics-os
COPY --from=web-build /src/web/dist /app/web/dist
COPY db/migrations /app/db/migrations
COPY scripts/container-entrypoint.sh /app/container-entrypoint.sh
RUN chmod 0555 /app/logistics-os /app/container-entrypoint.sh && chown -R logistics:logistics /app
USER logistics
ENV HTTP_ADDR=:8080 STATIC_DIR=/app/web/dist MIGRATIONS_DIR=/app/db/migrations
EXPOSE 8080
ENTRYPOINT ["/app/container-entrypoint.sh"]
