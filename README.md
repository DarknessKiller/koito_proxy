# Koito Proxy

![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge\&logo=go)
![Gin](https://img.shields.io/badge/Gin-Framework-00ACD7?style=for-the-badge)
![SQLite](https://img.shields.io/badge/SQLite-Embedded-003B57?style=for-the-badge\&logo=sqlite)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge\&logo=docker)

A **Metadata Correction Transparent Proxy** for Koito that intercepts music playback scrobble requests and applies user-defined metadata rules before they reach the upstream service.

Built in **Go (Gin framework)**, Koito Proxy acts as a lightweight middleware layer between your music client and the Koito backend.

> This project is actively maintaining. Internal rule handling and storage mechanisms are still experimental and may change frequently.

---

## The Idea

Instead of relying on manual fixes or post-scrobble corrections in Koito:

> Metadata correction is applied before requests reach the upstream API.

---

## Overview

Koito Proxy sits between your music client and Koito:

* Intercepts scrobble / playback API requests
* Applies metadata correction rules in real time
* Normalizes or merges track metadata automatically
* Forwards the modified request upstream transparently

---

## Features

* ⚡ Transparent reverse proxy for Koito API traffic
* 🧠 Rule-based metadata transformation engine
* 🔁 Creates Rules By Track merging via Koito's WebUI or manual rules insertion to database
* 🐳 Docker-first deployment

---

## Architecture

```text
Music Client (e.g. Navidrome)
          ↓
     Koito Proxy
          ↓
Upstream Koito API Service
```

---

## Tech Stack

* Go (Golang)
* Gin Web Framework
* SQLite
* Docker / Docker Compose

---

## Installation

### Prerequisites

* Docker + Docker Compose

---

### Run with Docker Compose (Recommended)

```bash
git clone https://github.com/your-org/koito-proxy.git
cd koito-proxy

docker compose up --build -d
```

---

## Docker Compose Configuration

```yaml
services:
  koito-proxy:
    container_name: koito-proxy
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      - KOITO_URL=http://localhost:4110
    volumes:
      - ./data:/app/data
    ports:
      - "8080:4112"
    restart: unless-stopped
```

---

## How It Works

1. Client sends a request intended for Koito
2. Koito Proxy intercepts the request
3. Metadata is extracted from payload
4. Rule engine applies transformations
5. Modified request is forwarded upstream
6. Response is returned transparently

---

## Pending Tasks

* [x] Rewrite rule engine
* [x] Replace DB layer with GORM
* [ ] Implement Rules CRUD API (management layer)

---

## Design Principles

* Fully transparent proxy behavior
* Deterministic metadata transformations
* Lightweight and predictable runtime behavior

