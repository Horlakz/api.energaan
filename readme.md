# Energaan API

Energaan API is a robust backend service designed to power modern web and mobile applications. Built with Go and leveraging the high-performance Fiber framework, this API provides a scalable and efficient solution for handling web requests. The project structure is meticulously organized to support a wide range of features, including authentication, product management, and more.

## Features

- **Environment Configuration**: Utilizes `godotenv` for loading environment variables from a `.env` file, ensuring secure and flexible configuration.
- **Fiber Framework**: Built on top of Fiber, a high-performance web framework for Go, offering express-like syntax and middleware support.
- **Middleware Integration**: Incorporates essential middleware for logging, CORS, rate limiting, and error handling to ensure a robust and secure API.
- **Database Connectivity**: Features a modular database client setup in `database.StartDatabaseClient`, facilitating easy integration with various databases.
- **Dynamic Routing**: Implements a dynamic router setup in `router.InitializeRouter`, allowing for scalable and maintainable route management.
- **DTO Pattern**: Utilizes the Data Transfer Object (DTO) pattern for efficient data handling and validation, as seen in `database/dto/app`.

## Technologies Used

- **Go**: A statically typed, compiled programming language designed for simplicity and efficiency.
- **Fiber**: An express-inspired web framework for Go, built on top of Fasthttp.
- **godotenv**: A Go port of Ruby's dotenv library, which loads env vars from a `.env` file.
- **UUID**: Utilizes Google's UUID library for generating unique identifiers for various entities.
- **Docker**: Supports containerization with Docker for easy deployment and scaling.

## Getting Started

To get started with the Energaan API, clone the repository and install the necessary dependencies:

```sh
git clone https://github.com/horlakz/energaan-api.git
cd energaan-api
go mod tidy
```
