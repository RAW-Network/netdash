# 🚀 NetDash

NetDash is a lightweight, self-hosted **internet speedtest tracker** built with **Go**, **HTMX**, and **SQLite**. It allows you to monitor your network performance over time with automated scheduling and a sleek dashboard

---

🔍 **Preview**

<img src="https://files.catbox.moe/9xi6hc.png" alt="NetDash Dashboard Preview" width="800">

---

## ✨ Features

* 📊 **Smooth Analytics** – Visualize **Download**, **Upload**, **Ping**, and **Packet Loss**
* 📅 **Automated Scheduling** – Run automatically using Cron expressions directly from the settings interface
* 💾 **Persistent Storage** – All results are saved in  database, ensuring your history is preserved
* 🐳 **Docker Ready** – Multi-stage build with support for **amd64** and **arm** architectures
* 🌚 **Modern UI** – Clean, dark-themed interface designed for continuous monitoring

---

## 🚀 Installation and Usage Guide

### ▶️ Run with Docker (Recommended)

```yaml
services:
  netdash:
    image: ghcr.io/raw-network/netdash:latest
    container_name: netdash
    ports:
      - "8080:80"
    volumes:
      - ./data:/netdash/data
    environment:
      - TZ=UTC
    restart: unless-stopped
```

Start the application:

```bash
docker compose up -d
```

Access the interface at:

```
http://localhost:8080
```

Stop the application:

```bash
docker compose down
```

---

## ⚙️ Configuration

Customize behavior through the **web interface** and environment variables if needed:

| Option / Variable | Description                                    | Default     |
| ----------------- | ---------------------------------------------- | ----------- |
| `TZ`              | Container timezone                             | `UTC`       |
| Schedule          | Cron expression for automated speedtests       | `0 * * * *` |
| History Limit     | Number of data points shown on chart and table | `10`        |
| Server ID         | Force a specific Ookla Speedtest server        | `Auto`      |

---

## 💻 Option 2: Run Locally (Development)

Clone the repository:

```bash
git clone https://github.com/raw-network/netdash.git
cd netdash
```

Install dependencies:

```bash
go mod download
```

Run the application:

```bash
go run main.go
```

Access the app at:

```
http://localhost:80
```

---

## 🛠️ Tech Stack

* **Backend**: Go (Echo Framework)
* **Frontend**: HTMX, Tailwind CSS, Chart.js
* **Database**: SQLite (GORM)
* **Scheduler**: robfig/cron
* **Speedtest Engine**: Ookla Speedtest CLI

---

## 📂 Project Structure

```plaintext
netdash/
├── data/                         # Persistent storage
├── internal/
│   ├── handler/
│   │   └── http.go               # HTTP routes and UI controllers
│   ├── logger/
│   │   └── logger.go             # Application logging utility
│   ├── model/
│   │   └── types.go              # GORM models and data structures
│   ├── repository/
│   │   ├── repository.go         # Data access layer interface
│   │   └── sqlite.go             # SQLite implementation and migrations
│   ├── server/
│   │   └── server.go             # Echo server configuration
│   ├── service/
│   │   ├── scheduler.go          # Background cron job management
│   │   └── speedtest.go          # Ookla CLI execution and JSON parsing
│   └── utils/
│       └── network.go            # Network helper functions
├── web/
│   ├── static/
│   │   ├── dashboard.js          # Chart logic and silent status polling
│   │   ├── settings.js           # Settings page form interactions
│   │   └── style.css             # Tailwind and custom layout styles
│   ├── template/
│   │   ├── layout/
│   │   │   └── base.html         # Base HTML template wrapper
│   │   ├── partials/
│   │   │   ├── result_row.html   # Table row partial
│   │   │   └── status_button.html# Dynamic button state partial
│   │   ├── index.html            # Main Dashboard view
│   │   └── settings.html         # Configuration view
│   └── efs.go                    # Go embed for serving static assets
├── .dockerignore                 # Docker build context exclusions
├── .gitignore                    # Git ignore rules
├── Dockerfile                    # Multi-stage Docker build (Port 80)
├── go.mod                        # Go module manifest
├── go.sum                        # Go module checksums
├── main.go                       # Application entry point
├── LICENSE                       # MIT License
├── VERSION                       # Application version
└── README.md                     # Main documentation
```

---

## 📄 License

This project is licensed under the **MIT License**
See the [LICENSE](./LICENSE) file for details