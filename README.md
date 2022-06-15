# gotodo

It is a personal project, a simple todo app in golang. I have used standard go library for most of the work to understand key concepts of backend development like request-response cycle, authentication, authorization, restful architecture, and so on. I have used following extra packages to help me in the development process: 
- [github.com/google/uuid](https://github.com/google/uuid) for generating primary keys
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migration
- [github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx) for database integration
- [github.com/lib/pq](https://github.com/lib/pq) for database integration
- [github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt) for authentication

## Setup

1. Install the required packages using
    ```bash
    go mod download
    ```
2. Create a ```.env``` file in the project root directory and copy the variable names from ```.env.example``` file. Replace ```{{insert}}``` with environment variable values in the ```.env``` file.
  - Example: ```BACKEND_PORT=8000```

3. Run the server using ```make server``` or using docker ```docker-compose up -d```.
