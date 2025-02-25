# 🏢 Asset Management Service

## 📖 About the Project

**Asset Management Service** is a scalable and efficient microservice designed to **streamline tracking and management of organizational assets**. This service enables users to **add, update, delete, and maintain asset records** while ensuring **efficient lifecycle management**.

---

## ✨ Key Features

- 🆕 **Add New Assets** – Register assets with detailed metadata.
- ✏️ **Update Asset Details** – Modify asset properties efficiently.
- 🗑 **Delete Assets** – Remove obsolete assets from the system.
- 🔧 **Maintenance Scheduling** – Track asset maintenance history and plan upcoming maintenance.
- 📊 **Audit Logging** – Ensure accountability with **logs of asset operations**.
- 🛠 **Containerized Deployment** – Deploy seamlessly with **Docker**.
- 📑 **Database Migrations** – Manage schema updates effortlessly.

---

## 🛠 Technology Stack

- **Backend Framework**: [Go](https://golang.org/) – High-performance microservices architecture.
- **Database**: [PostgreSQL](https://www.postgresql.org/) – Reliable relational database.
- **ORM**: [GORM](https://gorm.io/) – Simplified database interaction.
- **Docker**: [Docker](https://www.docker.com/) – Containerized deployment.

---

## 📦 Installation and Setup

### Prerequisites

- Install **[Go](https://golang.org/doc/install)**.
- Install **[PostgreSQL](https://www.postgresql.org/download/)**.
- (Optional) Install **[Docker](https://www.docker.com/)** for containerized deployment.

### Steps to Run

1. Clone the repository:
   ```bash
   git clone https://github.com/HieronimusMorgan/Asset-Service.git
   cd Asset-Service
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure environment variables:
   - Create a `.env` file to store **database credentials and other configurations**:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=asset_management
   ```

4. Run database migrations:
   ```bash
   go run cmd/migrate.go
   ```

5. Start the application:
   ```bash
   go run cmd/main.go
   ```

---

## 🔗 API Endpoints

### 🔓 Public Routes
- `GET /health` → **Service health check**.

### 🔒 Protected Routes (Require Authentication)
- `POST /v1/assets` → **Register a new asset**.
- `PUT /v1/assets/{id}` → **Update asset details**.
- `DELETE /v1/assets/{id}` → **Remove an asset**.
- `GET /v1/assets/{id}` → **Retrieve asset details**.
- `GET /v1/assets` → **List all assets**.
- `POST /v1/assets/{id}/maintenance` → **Schedule maintenance**.

---

## 📂 Project Structure

```
Asset-Service/
├── cmd/
│   ├── main.go          # Entry point of the application
│   └── migrate.go       # Database migration script
├── config/              # Configuration files
├── internal/
│   ├── models/          # Data models
│   ├── repository/      # Database interactions
│   └── services/        # Business logic
├── migration/           # Database migration files
├── pkg/response/        # API response structures
└── README.md
```

---

## 🤝 Contributing

Contributions are **welcome**! Follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Commit your changes (`git commit -m 'Add YourFeature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a Pull Request.

For major updates, **open an issue** first to discuss your proposal.

---

## 📜 License

This project is licensed under the **MIT License**. See the `LICENSE` file for more details.

---
