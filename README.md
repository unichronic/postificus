# Postificus: Centralized Content Distribution Engine

![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go) ![React](https://img.shields.io/badge/Frontend-React-61DAFB?style=flat&logo=react) ![Docker](https://img.shields.io/badge/Deploy-Docker-2496ED?style=flat&logo=docker) ![Redis](https://img.shields.io/badge/Queue-Asynq%20%2F%20Redis-DC382D?style=flat&logo=redis) ![Status](https://img.shields.io/badge/Status-Active_Development-success)

**Postificus** is a high-performance, event-driven content engine designed to solve the fragmentation of technical blogging. It allows authors to write content once and orchestrate its distribution across a disparate ecosystem of platforms (Medium, Dev.to, LinkedIn) while strictly enforcing SEO authority via Canonical URLs.

## 🚀 Architecture

The system moves beyond simple API scripting by implementing a **Multi-Stage Resilience Pipeline**.

### 1. The Core (Backend)
Built with **Go** and **Echo**, the backend is split into two services for scalability:
*   **API Service (`cmd/api`)**: Handles REST endpoints, authentication, and job enqueuing. High throughput, low latency.
*   **Worker Service (`cmd/worker`)**: Consumes jobs from **Redis** (via `Asynq`) and executes heavy automation tasks. This isolation prevents browser automation from blocking HTTP requests.

### 2. The Interface (Frontend)
A modern, distraction-free writing experience built with **React**, **Vite**, and **TailwindCSS**.
*   **Editor**: Powered by **Tiptap**, offering a Notion-like rich text experience.
*   **Real-time Status**: Polls the backend for publication status across all platforms.

## 🛡️ The "Unstoppable" Delivery System

Unlike standard cross-posters that fail when an API goes down, Postificus implements a **Hierarchical Fallback Strategy**:

| Priority | Method | Description |
| :--- | :--- | :--- |
| **1. Primary** | **Direct API** | Uses official REST APIs (e.g., Dev.to). Fastest and most reliable. |
| **2. Last Resort** | **Stealth Automation** | The "Nuclear Option." Uses **Go-Rod** to launch a headless browser, bypass bot detection, log in, and physically type the content into the editor. Used for platforms without write APIs (like LinkedIn personal profiles). |

## ✨ Key Features

### 🧠 Core Backend
*   **Event-Driven Architecture:** Decoupled ingestion and processing layers.
*   **SEO Guardrails:** Automatically manages `rel=canonical` tags to protect your domain authority.
*   **Concurrency Control:** Worker pools are rate-limited per domain to prevent IP bans.

### 🕵️ Browser Automation
*   **Stealth Mode:** Uses `rod-stealth` to strip `navigator.webdriver` flags, allowing the bot to pass as a human user on Single Page Applications (SPAs).
*   **Headless Production:** Automatically detects production environments (Docker/Render) to run headlessly, while keeping the UI visible for local debugging.

## 🚀 Getting Started

Follow these instructions to set up the project locally.

### Prerequisites

*   **Go**: v1.24 or higher
*   **Node.js**: v20 or higher (for frontend)
*   **Docker**: (Optional) For containerized execution
*   **Redis**: Required for the task queue (if running locally without Docker)
*   **PostgreSQL**: Required for the database (if running locally without Docker)

### Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/unichronic/postificus.git
    cd postificus
    ```

2.  **Configure Environment**
    Copy the example environment file and fill in your credentials:
    ```bash
    cp .env.example .env
    ```
    > **Note:** You will need to obtain session cookies/tokens for Dev.to, Medium, and LinkedIn manually from your browser dev tools if you plan to use the automation features.

### Running with Docker (Recommended)

The easiest way to run the entire stack (DB, Redis, API, Worker, Frontend) is using Docker Compose.

```bash
docker-compose up --build
```

The application will be available at:
*   **Frontend**: http://localhost:5173 (or the port exposed by Vite)
*   **API**: http://localhost:8080

### Local Development

If you prefer to run services individually for development:

1.  **Start Dependencies (DB & Redis)**
    ```bash
    docker-compose up -d db redis minio
    ```

2.  **Run the Backend (API)**
    ```bash
    go run cmd/api/main.go
    ```

3.  **Run the Worker (in a separate terminal)**
    ```bash
    go run cmd/worker/main.go
    ```

4.  **Run the Frontend**
    ```bash
    cd frontend
    npm install
    npm run dev
    ```

## ☁️ Deployment

The system is architected for a modern cloud stack:

*   **Backend**: Deployed as a Docker container on **Render**.
    *   *Note*: Requires a custom Dockerfile to include Chromium dependencies.
*   **Frontend**: Deployed as a static SPA on **Vercel**.


## 🛠️ Tech Stack

*   **Language:** Golang 1.24
*   **Frontend:** React, Vite, TailwindCSS
*   **Web Framework:** Echo v4
*   **Task Queue:** Asynq (Redis)
*   **Browser Automation:** Go-Rod (Stealth + CDP)
*   **Database:** PostgreSQL (GORM) & Redis
*   **Infrastructure:** Docker, Render, Vercel


