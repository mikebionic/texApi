# texApi

**texApi** is a logistics-oriented API ecosystem built with **Go (Golang)** using the **Gin web framework** and **PostgreSQL** as the database.
It provides services for logistics operations, including user management, offers, cargo, vehicles, analytics, chat, notifications, and more.

---

## Features

* Written in **Go** with a modular architecture (`controllers`, `services`, `dto`, `repo`, `queries`).
* REST API endpoints for logistics-related functionality (auth, users, companies, vehicles, cargo, offers, price quotes, news, etc.).
* **PostgreSQL** as the main data store, with versioned SQL schema migrations under `/schemas`.
* **WebSocket-based chat** with support for conversations, messaging, and notifications.
* **Firebase push notifications** integration.
* **Analytics scheduler** for reporting and statistics.
* Utility modules for file uploads, media handling, email (SMTP), JWT, and more.
* Deployment ready with **Docker Compose** and `systemd`.

---

## Project Structure (high-level)

* `cmd/tex/` — Application entry point (`main.go`).
* `config/` — Application configuration and environment loading.
* `database/` — Database initialization logic.
* `internal/`

    * `controllers/` — HTTP controllers mapping routes to services.
    * `services/` — Core business logic for each domain (auth, user, company, cargo, vehicles, chat, etc.).
    * `dto/` — Data Transfer Objects for request/response mapping.
    * `repo/` — Database repository layer (queries, CRUD).
    * `queries/` — SQL query definitions and helpers.
    * `chat/` — WebSocket chat system.
    * `firebasePush/` — Firebase notification integration.
    * `scheduler/` — Background jobs (e.g., analytics).
* `pkg/` — Shared utilities (middlewares, file handling, utils, smtp, sql safety).
* `schemas/` — SQL migration files for PostgreSQL.
* `scripts/` — Helper scripts for DB initialization, updates, and service setup.
* `docs/` — API workflow and feature documentation.
* `assets/` — Static files (logos, email templates).

---

## Setup

### Docker Compose

Using Docker Compose, you can run the app and database:

```bash
# To rebuild the app container
docker compose up --build app

# To rebuild all containers
docker compose up --build

# With DB initialization
INIT_DB=true docker compose up --build -d

# App-only update
INIT_DB=false docker compose up --build app -d
```

---

### Manual Setup

Clone and place the project under:

```bash
~/tex_backend/texApi
```

Run initial configuration:

```bash
make init-sys
```

Then configure the systemd service file:

```
systemd/system/texApi.service
```

To update the app (requires GitHub access keys):

```bash
bash ~/tex_backend/texApi/scripts/update_tex.sh
```

Set the absolute path for the uploads directory in `.env`, then run:

```bash
make upload-dir
```

Database and development setup:

```bash
make db
make dev
```

To build the application:

```bash
make build
```

---

## Requirements

* Go 1.20+
* PostgreSQL 14+
* Docker & Docker Compose (if using container setup)
* make & systemd (for manual setup)

---

## Documentation

Additional technical docs can be found under `/docs/`, e.g.:

* `docs/auth_workflow.md` — Authentication and OAuth flow
* `docs/analytics.md` — Analytics system
* `docs/call_room_docs.md` — Call/room handling
* `docs/firebaseNotification.md` — Push notifications
* `docs/Price Quote Docs tbl_price_quote.md` — Price quote module

