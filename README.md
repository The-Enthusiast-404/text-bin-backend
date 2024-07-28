# Text Bin API 🚀

![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-336791?style=for-the-badge&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-🐳-2496ED?style=for-the-badge&logo=docker)

Text Bin API is a robust and scalable backend service for the Text Bin application, providing secure text storage, retrieval, and management functionalities.

## 🌟 Features

- 📊 RESTful API endpoints
- 🔐 User authentication and authorization
- 📝 CRUD operations for text snippets
- ⏳ Text expiration management
- 🔍 Advanced text search capabilities
- 📊 Rate limiting and request throttling
- 🔒 CORS support

## 🛠️ Technologies

- **Language:** Go 1.20+
- **Web Framework:** Custom (using net/http)
- **Database:** PostgreSQL 13+
- **ORM:** Custom SQL with database/sql
- **Authentication:** JWT
- **Containerization:** Docker
- **API Documentation:** Swagger/OpenAPI

## 🚀 Getting Started

### Prerequisites

- Go 1.20 or later
- PostgreSQL 13 or later
- Docker (optional)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/The-Enthusiast-404/text-bin-backend.git
   cd text-bin-backend
    ```
2. Set up the database:
```bash
    psql -U postgres
    CREATE DATABASE textbin;
```
3. Set up the environment variables:
```bash
DB_DSN=postgres://username:password@localhost/textbin?sslmode=disable
```

4. Run DB migrations:
```bash
go run ./cmd/migrate
```

5. Build and run the application
```bash
go build ./cmd/api
./api
```
Certainly! Here's the backend README in markdown format:
markdownCopy# Text Bin API 🚀

![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-336791?style=for-the-badge&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-🐳-2496ED?style=for-the-badge&logo=docker)

Text Bin API is a robust and scalable backend service for the Text Bin application, providing secure text storage, retrieval, and management functionalities.

## 🌟 Features

- 📊 RESTful API endpoints
- 🔐 User authentication and authorization
- 📝 CRUD operations for text snippets
- ⏳ Text expiration management
- 🔍 Advanced text search capabilities
- 📊 Rate limiting and request throttling
- 🔒 CORS support

## 🛠️ Technologies

- **Language:** Go 1.20+
- **Web Framework:** Custom (using net/http)
- **Database:** PostgreSQL 13+
- **ORM:** Custom SQL with database/sql
- **Authentication:** JWT
- **Containerization:** Docker
- **API Documentation:** Swagger/OpenAPI

## 🚀 Getting Started

### Prerequisites

- Go 1.20 or later
- PostgreSQL 13 or later
- Docker (optional)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/The-Enthusiast-404/text-bin-backend.git
   cd text-bin-backend

Set up the database:
bashCopypsql -U postgres
CREATE DATABASE textbin;

Configure environment variables:
Create a .env file in the root directory and add:
CopyDB_DSN=postgres://username:password@localhost/textbin?sslmode=disable
JWT_SECRET=your_jwt_secret_here

Run database migrations:
bashCopygo run ./cmd/migrate

Build and run the application:
bashCopygo build ./cmd/api
./api




## 🤝 Contributing
We welcome contributions! Please see our Contribution Guidelines for more information.

## 📄 License
This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgements

Go Programming Language
PostgreSQL
JWT
