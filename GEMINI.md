# Gemini Code Assistant Documentation

This document provides an overview of the `g_static_site_template` project, its structure, and key commands for development.

## Project Overview

`g_static_site_template` is a web application built with Go. It serves a static website and includes a backend for handling API requests. The project uses a SQLite database, with `sqlc` for generating type-safe Go code for database interactions. The frontend is built using the `templ` templating engine and styled with Tailwind CSS.

## Technologies Used

- **Backend:**
  - Go (version 1.24.0)
  - Fiber (web framework)
  - Zerolog (for logging)
  - Viper (for configuration)
- **Database:**
  - SQLite
  - `sqlc` (for code generation)
  - `golang-migrate` (for migrations)
- **Frontend:**
  - `templ` (templating engine)
  - Tailwind CSS
- **Tooling:**
  - Docker (for `sqlc`)
  - Makefile (for command automation)
  - `air` (for live reloading)

## Getting Started

### Prerequisites

- Go (1.24.0 or later)
- Docker
- `migrate` CLI tool
- `templ` CLI tool
- `tailwindcss` CLI tool
- `air` CLI tool

### Installation and Setup

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd g_static_site_template
    ```

2.  **Set up environment variables:**
    Copy the `.env.example` file to `.env` and update the variables as needed.
    ```bash
    cp .env.example .env
    ```

3.  **Run database migrations:**
    The database is a SQLite file, which will be created automatically. You need to run the migrations to create the schema.
    ```bash
    make mgup
    ```

4.  **Generate assets:**
    ```bash
    make sqlc
    make templgen
    make twc
    ```

5.  **Run the application with live reload:**
    ```bash
    air
    ```
    The application will be available at `http://localhost:8088` (or the port specified in your `.env` file).

## Development

### Common Commands

The `Makefile` provides several commands to streamline development:

-   `make mgup`: Apply database migrations.
-   `make mgdown`: Roll back database migrations.
-   `make sqlc`: Generate Go code from SQL queries.
-   `make templfmt`: Format `.templ` files.
-   `make templgen`: Generate Go code from `.templ` files.
-   `make twc`: Watch for changes in `input.css` and rebuild the CSS.

### File Structure

-   `cmd/server/main.go`: The main entry point of the application.
-   `internal/`: Contains the core application logic.
  -   `config/`: Configuration management.
  -   `db/`: Database-related files, including migrations, queries, and generated code.
  -   `modules/`: Application modules, containing `templ` components and pages.
  -   `server/`: The Fiber web server implementation.
-   `static/`: Static assets like CSS, images, and JavaScript.
-   `Makefile`: Automation scripts.
-   `go.mod`, `go.sum`: Go module definitions.
-   `sqlc.yaml`: Configuration for `sqlc`.
-   `.air.toml`: Configuration for `air` (live-reloading).

## Linting and Formatting

-   **Go:**
    -   Formatting: `go fmt ./...`
    -   Linting: Use a tool like `golangci-lint`.
-   **`templ`:**
    -   Formatting: `make templfmt`
-   **Tailwind CSS:**
    -   The `make twc` command will minify the output CSS.