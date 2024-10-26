# API Gateway

## Description

This project is an API Gateway implemented in Go using the Gin web framework. It serves as a central entry point for managing and routing requests to various microservices in a distributed system architecture.

## Features

- User authentication and authorization
- Token refresh mechanism
- User registration and update
- Routing to different microservices (Users, Auth, etc.)
- Environment configuration management
- Error handling

## Architecture and Design Patterns

### Layered Architecture

The project follows a layered architecture pattern:

1. **Controllers**: Handle incoming HTTP requests and responses.
2. **Services**: Contain business logic and communicate with external APIs.
3. **DTOs (Data Transfer Objects)**: Define structures for data exchange between layers.
4. **Routes**: Define API endpoints and map them to controller methods.

### Dependency Injection

The project uses dependency injection to manage dependencies between components, improving modularity and testability.

### Interface-based Design

Interfaces are used extensively to define contracts between different parts of the application, allowing for easier mocking and testing.


### Error Handling

A custom error type is implemented to standardize error responses across the application.

## Project Structure

- `src/`
  - `config/`: Configuration-related code
  - `controllers/`: Request handlers
  - `dto/`: Data Transfer Objects
  - `errors/`: Custom error types
  - `routes/`: API route definitions
  - `services/`: Business logic and external API communication
  - `utils/`: Utility functions (e.g., JWT handling)
- `main.go`: Application entry point

## How It Works

1. The application starts in `main.go`, loading environment variables and building the application.
2. Gin router is set up with defined routes.
3. Incoming requests are handled by appropriate controllers.
4. Controllers call corresponding services to process the request.
5. Services communicate with external APIs or perform business logic.
6. Responses are sent back through the controller to the client.

## Setup and Running

1. Ensure Go is installed on your system.
2. Clone the repository.
3. Copy `.env.example` to `.env` and fill in the required values.
4. Run `go mod tidy` to install dependencies.
5. Start the application with `go run main.go`.

The API Gateway will be available at `http://localhost:4000` (or the port specified in your `.env` file).

