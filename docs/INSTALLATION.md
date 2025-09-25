# Installing Thunderdome

This guide provides detailed instructions for installing and setting up Thunderdome using different methods.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Methods](#installation-methods)
- [Configuration](#configuration)
- [Verification](#verification)
- [How to run DB Migrations and when necessary](#how-to-run-db-migrations-and-when-necessary)
- [Troubleshooting](#troubleshooting)
- [Upgrading from v2 to v3](#upgrading-from-v2-to-v3)

## Prerequisites

- PostgreSQL database
- SMTP server for email functionality
- Docker (optional)
- Kubernetes cluster (optional)
- Helm (optional, for K8s installation)

## Installation Methods

### 1. Docker Installation

1. Pull the latest Docker image:
```bash
docker pull stevenweathers/thunderdome-planning-poker
```

2. Create an `.env` file with your configuration:
```properties
DB_HOST=your-db-host
DB_PORT=5432
DB_NAME=thunderdome
DB_USER=your-db-user
DB_PASSWORD=your-db-password
SMTP_HOST=your-smtp-host
SMTP_PORT=587
SMTP_USER=your-smtp-user
SMTP_PASSWORD=your-smtp-password
APP_HOST=your-domain.com
```

3. Run the container:
```bash
docker run -d \
  --name thunderdome \
  --env-file .env \
  -p 8080:8080 \
  stevenweathers/thunderdome-planning-poker
```

### 2. Binary Installation

1. Download the latest release from [GitHub Releases](https://github.com/StevenWeathers/thunderdome-planning-poker/releases/latest)

2. Extract the binary

3. Set up environment variables (same as Docker configuration above)

4. Run the binary:
```bash
./thunderdome serve
```

### 3. Kubernetes Installation

1. Navigate to the Helm charts directory:
```bash
cd ./build/helm
```

2. Install PostgreSQL database:
```bash
helm install -f thunderdome-db.yaml thunderdome-db ./db
```

3. Install mail server:
```bash
helm install -f thunderdome-mail.yaml thunderdome-mail ./app
```

4. Install Thunderdome:
```bash
helm install -f thunderdome.yaml thunderdome ./app
```

5. Access the application:
```bash
kubectl port-forward svc/thunderdome 8080:8080
```

## Configuration

For detailed configuration options, refer to the [Configuration Guide](CONFIGURATION.md).

## Verification

After installation:

1. Access Thunderdome at `http://localhost:8080` (or your configured domain)
2. Create an admin account on first run
3. Verify email functionality by testing the registration process

## How to run DB Migrations and when necessary

By default Thunderdome will automatically attempt to run the db migrations, however there may be usecases where you want to manually run the migrations.

### Migration subcommands:

`./thunderdome-planning-poker migrate up` to Run all pending migrations
`./thunderdome-planning-poker migrate down` to Rollback the last migration
`./thunderdome-planning-poker migrate status` to Show migration status

## Troubleshooting

- Check logs for any startup errors
- Verify database connection and credentials
- Ensure SMTP server is properly configured
- Confirm proper network connectivity and port availability

For additional help, refer to the [User Guide](GUIDE.md) or open an issue on GitHub.

## Upgrading from v2 to v3

If you're upgrading from version 2.x.x to 3.x.x:

1. Review breaking changes in v3 (notably removal of Dicebear Avatars service)
2. If using docker-compose, upgrade PostgreSQL to version 15
3. Run the latest 2.x.x release to complete SQL migrations
4. Run the latest 3.x.x release