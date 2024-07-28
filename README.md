# TextBin 📝

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=for-the-badge&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-🐳-2496ED?style=for-the-badge&logo=docker)

TextBin is a modern, feature-rich pastebin alternative built with Go. It allows users to easily share and manage text snippets with powerful functionality and a clean, intuitive API.

## 🌟 Features

### Implemented ✅

- 🔐 User authentication and authorization
  - User registration and login
  - JWT-based authentication
- 📝 Text snippet management
  - Create, read, update, and delete text snippets
  - Support for public and private snippets
- ⏳ Expiration settings for snippets
- 🎨 Syntax highlighting support
- 👍 Like system for snippets
- 💬 Commenting system
- 🔒 CORS support
- 📊 Basic rate limiting

### Planned Enhancements 🚀

- 🔍 Full-text search capabilities
- 📈 Advanced rate limiting and request throttling
- 📨 Email notifications
- 🔗 Sharing via short URLs
- 📱 Mobile-friendly API endpoints
- 🔄 Version history for snippets
- 🏷️ Tagging system for better organization
- 👥 User groups and collaboration features
- 🔐 Two-factor authentication (2FA)
- 📊 User dashboard with usage statistics
- 🌐 Multi-language support
- 🔌 API for third-party integrations

## 🛠️ Technology Stack

- **Backend:** Go 1.21+
- **Database:** PostgreSQL 15+
- **Authentication:** JWT
- **API Documentation:** Swagger/OpenAPI (planned)
- **Containerization:** Docker

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- PostgreSQL 15 or later
- Docker (optional, for containerized deployment)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/textbin.git
   cd textbin
    ```
2. Set up the PostgreSQL database:
   ```bash
   psql -U postgres -c "CREATE DATABASE textbin"
   ```
3. Copy the example environment file and configure the environment variables:
   ```bash
   DB_DSN=postgres://username:password@localhost/textbin?sslmode=disable
   JWT_SECRET=your_jwt_secret_here
   SMTP_HOST=smtp.example.com
   SMTP_PORT=587
   SMTP_USERNAME=your_username
   SMTP_PASSWORD=your_password
   SMTP_SENDER=TextBin <noreply@textbin.example.com>
   ```
4. Run db migrations:
   ```bash
   go run ./cmd/migrate
   ```
5. Start the server:
   ```bash
    go run ./cmd/api/main.go
    ```
6. Visit `http://localhost:4000/v1/healthcheck` in your browser to see the API status.

## 🤝 Contributing

We welcome contributions! Please see our Contribution Guidelines for more information on how to get started.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgements

Go Programming Language
PostgreSQL
JWT-Go

## 📞 Support

If you encounter any issues or have questions, please open an issue on our GitHub repository.
