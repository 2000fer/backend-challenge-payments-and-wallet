# Etapa de construcción
FROM golang:1.25.4 AS builder

WORKDIR /app

# Copia los archivos necesarios
COPY go.mod go.sum ./
RUN go mod download

# Copia el código fuente
COPY . .

# Construye la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o wallet-api cmd/api/main.go

# Etapa final
FROM alpine:latest

# Instala dependencias necesarias
RUN apk --no-cache add postgresql-client

WORKDIR /app

# Copia el binario desde el builder
COPY --from=builder /app/wallet-api .
# Copia las migraciones
COPY migrations ./migrations
# Copia el script de entrada
COPY docker-entrypoint.sh /docker-entrypoint.sh

# Hace ejecutable el script
RUN chmod +x /docker-entrypoint.sh

# Puerto expuesto
EXPOSE 8080

# Punto de entrada personalizado
ENTRYPOINT ["/docker-entrypoint.sh"]

# Comando por defecto
CMD ["./wallet-api"]