# Backend Challenge - Payments and Wallet API

## Descripción del Proyecto
API RESTful para un sistema de billetera digital que permite a los usuarios realizar pagos, consultar saldos y ver el historial de transacciones. El sistema está construido con Go y utiliza PostgreSQL como base de datos.

## Decisiones Arquitectónicas

### Arquitectura en Capas
- **API Layer**: Maneja las solicitudes HTTP y el enrutamiento usando Gin
- **Service Layer**: Contiene la lógica de negocio
- **Repository Layer**: Maneja la persistencia de datos
- **Domain Layer**: Define los modelos de negocio

### Patrones de Diseño
- **Dependency Injection**: Para facilitar el testing y la mantenibilidad
- **Repository Pattern**: Para abstraer el acceso a datos
- **Factory Pattern**: Para la creación de servicios

### Manejo de Transacciones
- Uso de transacciones de base de datos para operaciones atómicas
- Rollback automático en caso de errores

## Stack Tecnológico

### Backend
- **Lenguaje**: Go 1.21+
- **Framework Web**: Gin
- **Base de Datos**: PostgreSQL 14+
- **ORM**: pgx (driver nativo de PostgreSQL para Go)
- **Manejo de Configuración**: Variables de entorno

### Herramientas de Desarrollo
- **Go Modules**: Para gestión de dependencias
- **golangci-lint**: Para análisis estático de código
- **Docker**: Para el entorno de desarrollo

## Instalación y Ejecución

### Requisitos Previos
- Go 1.21 o superior
- PostgreSQL 14 o superior
- Git

### Configuración Inicial

1. Clonar el repositorio:
   ```bash
   git clone https://github.com/2000fer/backend-challenge-payments-and-wallet.git
   cd backend-challenge-payments-and-wallet
   ```

2. Configurar variables de entorno:
   ```bash
   cp .env.example .env
   # Editar el archivo .env con tus credenciales
   ```

3. Instalar dependencias:
   ```bash
   go mod download
   ```

4. Configurar la base de datos:
   ```bash
   psql -U tu_usuario -d tu_base_de_datos -f migrations/local/init_db.sql
   ```

### Ejecución

```bash
# Modo desarrollo
go run cmd/api/main.go

# Compilar y ejecutar
go build -o wallet-api cmd/api/main.go
./wallet-api
```

## Ejecución con Docker

```bash
# Construir la imagen
docker-compose build

# Iniciar los servicios
docker-compose up -d

# Ver logs
docker-compose logs -f
```

## Ejecución de Tests

```bash
# Ejecutar todos los tests
go test ./...

# Ejecutar tests con cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Ejecutar linters
golangci-lint run
```

## Endpoints Disponibles

### 1. Health Check (Ping)
- `GET /ping`
  - Verifica que el servicio esté en funcionamiento
  - **Ejemplo de respuesta**:
    ```json
    {
      "status": "OK",
      "message": "Service is running"
    }
    ```

### 2. Obtener Saldo
- `GET /api/v1/wallets/:user_id/balance`
  - Obtiene el saldo actual de un usuario
  - **Ejemplo de respuesta exitosa**:
    ```json
    {
      "user_id": 123,
      "balance": 1500.50
    }
    ```

### 3. Realizar Pago
- `POST /api/v1/wallets/:user_id/payments`
  - Crea un nuevo pago
  - **Cuerpo de la solicitud**:
    ```json
    {
      "amount": 100.50,
      "description": "Pago de servicio"
    }
    ```
  - **Ejemplo de respuesta exitosa**:
    ```json
    {
      "transaction_id": "550e8400-e29b-41d4-a716-446655440000",
      "status": "ok"
    }
    ```

### 4. Obtener Historial de Transacciones
- `GET /api/v1/wallets/:user_id/transactions`
  - Obtiene el historial de transacciones
  - **Parámetros de consulta opcionales**:
    - `limit`: Número de transacciones a devolver (por defecto: 10)
  
  - **Ejemplo de respuesta**:
    ```json
    {
      "transactions": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440000",
          "user_id": 123,
          "amount": -100.50,
          "type": "payment",
          "status": "completed",
          "created_at": "2025-11-30T14:30:00Z"
        }
      ]
    }
    ```

## Mejoras Futuras

1. **Autenticación y Autorización**
   - Implementar JWT para autenticación
   - Control de acceso basado en roles (RBAC)

2. **Documentación de la API**
   - Integrar Swagger/OpenAPI
   - Documentación detallada de endpoints

3. **Manejo de Monedas**
   - Soporte para múltiples monedas
   - Conversión de divisas

4. **Notificaciones**
   - Integración con servicio de notificaciones
   - Webhooks para eventos importantes

5. **Monitoreo y Métricas**
   - Integración con Prometheus
   - Logging estructurado
   - Trazabilidad distribuida

6. **Escalabilidad**
   - Implementar colas para procesamiento asíncrono
   - Caché con Redis

7. **Seguridad Adicional**
   - Rate limiting
   - Validación de entrada más robusta
   - Protección contra ataques comunes

8. **Testing**
   - Pruebas de integración
   - Pruebas de carga
   - Pruebas de seguridad