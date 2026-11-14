# LiteBank
Full-stack банковское приложение, которое реализует безопасные финансовые операции.
Backend приложения написан на Go с использованием чистой архитектуры
Frontend приложения написан на React

# Technology
* **Backend**: Go, PostgreSQL, `pgx`, JWT, Goose
* **Frontend**: React, TypeScript, Vite, Axios, React Router
* **Infrastructure**: Docker, Docker compose, CI (GitHub Actions)
* **Documentation**: Swagger

# Features
* **Защита от Race Condition**
* **Idempotency (Идемпотентность)**
* **JWT Авторизация**
* **Double-Entry Ledger**

# How to start
1. Клонировать репозиторий:
    ```bash
    git clone https://github.com/oopchoop/litebank.git
    cd litebank
2. Клонируйте .env.example и назовите .env
3. Поднимите контейнеры,
    ```bash
    docker compose up -d --build
4. Запустите frontend:
    ```bash
    cd frontend
    npm install
    npm run dev