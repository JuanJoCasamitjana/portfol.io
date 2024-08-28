# Etapa 1: Construcción
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN apk add --no-cache gcc g++ musl-dev
RUN go mod tidy
RUN CGO_ENABLED=1 go build -o myapp

# Etapa 2: Imagen mínima para ejecutar
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/myapp .

ENV PORT=8080

EXPOSE ${PORT}

# Usar la variable de entorno para ejecutar la aplicación en el puerto correcto
CMD ["sh", "-c", "./myapp"]
