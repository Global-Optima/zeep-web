# zeep-backend

This template provides a foundation for developing the backend of **ZeepApp**, a real-time drink ordering application, using **Golang**. The backend service manages order flows, tracks inventory, and enables real-time communication with the frontend.

## Recommended IDE Setup

- **GoLand** or **Visual Studio Code** with the **Go** plugin for optimal Golang development.

## Customize Configuration

Configurations are specified in the `config.yaml` file, which includes database settings, server ports, and third-party integrations as required.

## Project Setup

1. **Install Dependencies**
   ```bash
   go mod tidy
   ```

2. **Database Setup and Migrations**  
   Ensure your database is running (PostgreSQL recommended). Database migrations are handled using **GORM**. Run the initial migrations by executing:
   ```bash
   go run ./cmd/main.go
   ```

3. **Run Server in Development Mode**
   The main entry file for the server is located at `./cmd/main.go`. Start the server in development mode by running:
   ```bash
   go run ./cmd/main.go
   ```

4. **Environment Setup**  
   Configure your `config.yaml` file to set environment-specific variables, such as database credentials and API keys.

## Compile and Run for Production

To compile and run the backend in production mode, execute:

```bash
go build -o zeep-backend ./cmd/main.go
./zeep-backend
```

## Run Unit Tests

Run all unit tests to ensure your code is functioning as expected:

```bash
go test ./...
```

## Run Integration Tests

Run integration tests to verify backend functionality with other systems:

```bash
go test -tags=integration ./test/integration
```

---

This structured setup should help you get started with ZeepApp’s backend development in Golang.