# Koito Proxy

![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge\&logo=go)
![Gin](https://img.shields.io/badge/Gin-Framework-00ACD7?style=for-the-badge)
![SQLite](https://img.shields.io/badge/SQLite-Embedded-003B57?style=for-the-badge\&logo=sqlite)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge\&logo=docker)

A **Metadata Correction Transparent Proxy** for Koito that intercepts music playback and scrobble requests, applies user-defined metadata correction rules, and forwards the modified requests to the upstream Koito service.

Built with **Go** and the **Gin Web Framework**, Koito Proxy acts as a lightweight middleware layer between your music client and Koito.

> This project is actively maintained. Internal rule handling and storage mechanisms are still experimental and may change between releases.

---

## The Idea

Instead of relying on manual fixes or post-scrobble corrections in Koito:

> Metadata correction is applied before requests reach the upstream API.

This allows track metadata to be normalized, merged, or rewritten automatically before Koito processes it.

---

## Overview

Koito Proxy sits between your music client and Koito:

* Intercepts scrobble API requests
* Applies metadata correction rules in real time
* Normalizes or merges track metadata automatically
* Transparently forwards the modified request upstream
* Returns the upstream response unchanged

---

## Features

* ⚡ Transparent reverse proxy for Koito API traffic
* 🧠 Rule-based metadata transformation engine
* 🔁 Create rules from track merges performed through Koito's Web UI
* 🗄️ SQLite-backed rule storage
* 🐳 Docker-first deployment
* 🚀 Lightweight and resource-efficient

---

## Architecture

```text
Music Client (e.g. Navidrome)
          │
          ▼
     Koito Proxy
          │
          ▼
Upstream Koito API Service
```

---

## Tech Stack

* Go (Golang)
* Gin Web Framework
* SQLite
* GORM
* Docker
* Docker Compose

---

## Installation

### Prerequisites

* Docker
* Docker Compose

---

### Run with Docker Compose (Recommended)

```bash
mkdir -p koito-proxy
cd koito-proxy

wget https://raw.githubusercontent.com/DarknessKiller/koito_proxy/refs/heads/master/compose.yml

docker compose up -d
```

---

## Docker Compose Configuration

```yaml
services:
  koito-proxy:
    container_name: koito-proxy
    image: ghcr.io/darknesskiller/koito_proxy:latest

    environment:
      # Required
      - KOITO_URL=http://localhost:4110

      # Optional
      - PROXY_PORT=4112
      - PROXY_DB=/app/data/koito_proxy.db

    volumes:
      - ./data:/app/data

    ports:
      - "4112:4112"

    restart: unless-stopped
```

---

## Environment Variables

| Variable     | Required | Default              | Description                                                            |
| ------------ | -------- | -------------------- | ---------------------------------------------------------------------- |
| `KOITO_URL`  | Yes      | -                    | URL of the upstream Koito instance that requests will be forwarded to. |
| `PROXY_PORT` | No       | `4112`               | Port the proxy listens on inside the container.                        |
| `PROXY_DB`   | No       | `/app/data/koito_proxy.db` | Path to the SQLite database file used for rule storage.                |

### Example: Custom Database File

```yaml
environment:
  - KOITO_URL=http://koito:4110
  - PROXY_DB=/app/data/koito_proxy.db
```

### Example: Custom Listening Port

```yaml
environment:
  - KOITO_URL=http://koito:4110
  - PROXY_PORT=5000

ports:
  - "5000:5000"
```

### Example: Custom Port and Database

```yaml
environment:
  - KOITO_URL=http://koito:4110
  - PROXY_PORT=5000
  - PROXY_DB=/app/data/custom.db

ports:
  - "5000:5000"
```

---

## How It Works

1. A client sends a scrobble request intended for Koito.
2. Koito Proxy intercepts the request.
3. Relevant metadata is extracted from the payload.
4. The rule engine evaluates and applies matching transformations.
5. The modified request is forwarded to the upstream Koito instance.
6. The upstream response is returned to the client transparently.

---

## Rule Processing

Rules allow metadata to be rewritten before Koito receives it.

Typical use cases include:

* Correcting artist names
* Merging duplicate artists
* Fixing album titles
* Normalizing track metadata
* Handling inconsistent tags from different music sources

All transformations occur before the request reaches Koito.

---

## Data Storage

Koito Proxy stores rule definitions in a SQLite database.

By default:

```text
/app/data/koito_proxy.db
```

The database location can be customized using the `PROXY_DB` environment variable.

---

## Pending Tasks

* [x] Rewrite rule engine
* [x] Replace database layer with GORM
* [ ] Implement Rules CRUD API

---

## Design Principles

* Fully transparent proxy behavior
* Deterministic metadata transformations
* Lightweight runtime footprint
* Simple deployment and maintenance
* Minimal configuration requirements

