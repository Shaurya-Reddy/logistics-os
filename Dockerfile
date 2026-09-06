FROM node:24.20.0-alpine3.23 AS frontend
WORKDIR /src/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run check && npm run build

FROM golang:1.27.1-alpine3.23 AS backend
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /src/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/logistics-os ./cmd/logistics-os

FROM alpine:3.23.5
RUN apk add --no-cache --upgrade libcrypto3=3.5.8-r0 libssl3=3.5.8-r0 ca-certificates tzdata && addgroup -g 10001 app && adduser -D -H -u 10001 -G app app
COPY --from=backend /out/logistics-os /app/logistics-os
USER 10001:10001
EXPOSE 8080
ENTRYPOINT ["/app/logistics-os"]
CMD ["serve"]
