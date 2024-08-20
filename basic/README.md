# Go Application with Docker and HTMX Integration

## Overview
This service is a lightweight Go web application designed to handle both server-side rendering and dynamic client-side updates using HTMX. It is containerized for easy deployment and development, and is integrated with a PostgreSQL database for data persistence.

This Go application is a simple web service that serves a homepage with an HTMX-enabled form. The application is containerized using Docker, allowing for easy deployment and environment consistency.


## Main Components

### 1. `main.go`
- **HomePage Route**: 
  - Handles GET requests to the root (`/`) by rendering the `home.html` template.
  - Handles POST requests submitted via HTMX to dynamically update the page content without a full reload.
- **Static File Handling**: 
  - Serves static files like JavaScript and CSS from the `static/` directory.

### 2. `home.html`
- **HTMX Form**: 
  - A form that uses HTMX for AJAX-like behavior. It submits data to the server and dynamically updates a part of the page based on the response.
  
### 3. `Dockerfile`
- **Multi-Stage Build**:
  - **Stage 1**: Uses the official Golang image to build the Go application.
  - **Stage 2**: Uses a minimal Alpine Linux image to create a small and efficient final image with the compiled binary and necessary assets.
- **Static Assets**: Copies templates and static files into the final Docker image.

### 4. `docker-compose.yml`
- **App Service**:
  - Builds and runs the Go application inside a Docker container.
  - Exposes the application on port `8080`.
  - Mounts the current project directory into the container to allow live updates.
- **Database Service**:
  - Configures a PostgreSQL database instance using the official Docker image.
  - Defines environment variables for database credentials.

## Key Features

- **HTMX Integration**: Enables dynamic page updates without full reloads by processing POST requests directly on the homepage.
- **Dockerization**: Ensures the application runs consistently across different environments with a multi-stage Docker build process.
- **Live Reloading**: Supports live reloading of the application during development via Docker Compose.

## How to Run

1. **Build and Run**:
   ```bash
   make docker-run-dev

2. **Access the Application**:
    ```bash
   Visit http://localhost:8080 to view the application.
   
3. **Modify and Test**:
   ```bash
   Make changes to the code, and the updates will be reflected within the running container.