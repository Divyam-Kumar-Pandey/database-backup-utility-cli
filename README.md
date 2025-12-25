## 🛡️ Database Backup Utility (CLI)

A **production-ready, cloud-capable database backup and restore CLI** written in **Go**, supporting automated scheduling, compression, secure credential handling, and pluggable storage backends.

---

## 🚀 Key Features

### 🔄 Database Support

* **MySQL** (via native `mysqldump` / `mysql`)
* Architecture designed for future database support:

  * PostgreSQL
  * MongoDB 
  * SQLite

### 📦 Backup Pipeline

```
Database → SQL dump → Gzip compression → Storage backend
```

* Streaming-based (handles large databases)
* No memory-heavy operations
* Native database tooling (reliable & fast)

---

### ☁️ Storage Backends

* **Local filesystem** (`./backups`)
* **AWS S3**
* Cloud-ready adapter pattern (GCS planned)

---

### ⏱️ Automated Scheduling

* Native **cron integration**
* No long-running daemon
* Survives reboots
* Fully manageable from CLI

```bash
backup-tool schedule add my_db "0 2 * * *"
backup-tool schedule remove my_db
```

---

### 🔐 Security Hardening

* No DB credentials exposed in process list
* Environment-variable–based secrets
* Cloud-native credential resolution (IAM / AWS SDK)
* Encryption-ready backup pipeline

---

### 📝 Observability

* Structured logging (INFO / ERROR)
* File-based logs (`db_backup.log`)
* Clean CLI output for operators
* Debug-friendly execution traces

---

## 📦 Installation

### Prerequisites

* Go **1.21+**
* MySQL client tools (`mysqldump`, `mysql`)
* AWS credentials (for S3 storage)

### Build

```bash
git clone https://github.com/Divyam-Kumar-Pandey/database-backup-utility-cli
cd database-backup-utility-cli
go build -o db-backup-cli ./cmd/main.go
sudo mv db-backup-cli /usr/local/bin/

which db-backup-cli
db-backup-cli --help
```

---

## ⚙️ Configuration

📄 `db_backup_config.yaml`

```yaml
storage:
  type: s3        # local | s3
  bucket: my-backup-bucket
  region: ap-south-1
  path: backups

databases:
  my_mysql:
    type: mysql
    host: localhost
    port: 3306
    user: root
    password: ${MYSQL_PASSWORD}

```

Secrets are **never hardcoded** — use environment variables.

---

## 📖 Usage

### Backup a Database

```bash
db-backup-cli backup my_mysql
```

### Restore from Backup

```bash
db-backup-cli restore backups/my_mysql_backup.sql.gz my_mysql
```

### List Stored Backups

```bash
db-backup-cli list
```

### Schedule Automatic Backups

```bash
db-backup-cli schedule add my_mysql "0 2 * * *"
db-backup-cli schedule list
db-backup-cli schedule remove my_mysql
```

---

## 🧠 Architecture Overview

```
cmd/main.go              → CLI entrypoint (Cobra)
pkg/
 ├── core/               → Interfaces (Database, Storage)
 ├── databases/          → MySQL adapter (Postgres/Mongo planned)
 ├── storage/            → Local & AWS S3 adapters
 └── utils/              → Logging, compression, cron, security
```

**Adapter-based design** allows adding new databases or storage providers without modifying CLI logic.

---

## 🔒 Security Design Highlights

* `MYSQL_PWD` environment variable (prevents credential leaks)
* No hardcoded cloud credentials
* Encryption-ready pipeline
* Safe logging (no secret leakage)

---

## 🛠️ Tech Stack

* **Go**
* Cobra & Viper (CLI + config)
* AWS SDK v2
* Native OS tools (cron, mysqldump)
* Gzip compression

---

## 📌 Planned Enhancements (TODOs)

These are **intentional next-phase improvements**, not missing basics:

### 🔄 Database Support

* [ ] PostgreSQL adapter
* [ ] MongoDB adapter
* [ ] SQLite adapter

### ☁️ Storage

* [ ] Google Cloud Storage (GCS)
* [ ] Azure Blob Storage

### 🔐 Security

* [ ] AES encryption for backups
* [ ] AWS KMS integration
* [ ] Backup checksum verification

### ⏱️ Scheduling

* [ ] Schedule listing command
* [ ] Schedule validation
* [ ] Email/Slack alerts on failure

### 🧪 Quality

* [ ] Unit tests with mocks
* [ ] Integration tests
* [ ] CI pipeline

### 📦 Distribution

* [ ] `.deb` / `.rpm` packages
* [ ] Homebrew formula
* [ ] Docker image

---