# Installing Thunderdome

This guide provides detailed instructions for installing and setting up Thunderdome using different methods.

## Prerequisites

- PostgreSQL database
- SMTP server for email functionality
- Docker (optional)
- Kubernetes cluster (optional)
- Helm (optional, for K8s installation)

## Installation Methods

### 1. Docker Installation (Recommended)

1. Pull the latest Docker image:
   \```bash
docker pull stevenweathers/thunderdome-planning-poker
\```

2. Create a `.env` file with your configuration:
   \```properties
THUNDERDOME_DB_HOST=your-db-host
THUNDERDOME_DB_PORT=5432
THUNDERDOME_DB_NAME=thunderdome
THUNDERDOME_DB_USER=your-db-user
THUNDERDOME_DB_PASSWORD=your-db-password
THUNDERDOME_SMTP_HOST=your-smtp-host
THUNDERDOME_SMTP_PORT=587
THUNDERDOME_SMTP_USER=your-smtp-user
THUNDERDOME_SMTP_PASSWORD=your-smtp-password
THUNDERDOME_APP_HOST=your-domain.com
\```

3. Run the container:
   \```bash
docker run -d \
  --name thunderdome \
  --env-file .env \
  -p 8080:8080 \
  stevenweathers/thunderdome-planning-poker
\```

### 2. Binary Installation

1. Download the latest release from [GitHub Releases](https://github.com/StevenWeathers/thunderdome-planning-poker/releases/latest)

2. Extract the binary

3. Set up environment variables (same as Docker configuration above)

4. Run the binary:
   \```bash
./thunderdome
\```

### 3. Kubernetes Installation

1. Navigate to the Helm charts directory:
   \```bash
cd ./build/helm
\```

2. Install PostgreSQL database:
   \```bash
helm install -f thunderdome-db.yaml thunderdome-db ./db
\```

3. Install mail server:
   \```bash
helm install -f thunderdome-mail.yaml thunderdome-mail ./app
\```

4. Install Thunderdome:
   \```bash
helm install -f thunderdome.yaml thunderdome ./app
\```

5. Access the application:
   \```bash
kubectl port-forward svc/thunderdome 8080:8080
\```

## Configuration

For detailed configuration options, refer to the [Configuration Guide](CONFIGURATION.md).

## Upgrading from v2 to v3

If you're upgrading from version 2.x.x to 3.x.x:

1. Review breaking changes in v3 (notably removal of Dicebear Avatars service)
2. If using docker-compose, upgrade PostgreSQL to version 15
3. Run the latest 2.x.x release to complete SQL migrations
4. Run the latest 3.x.x release

## Verification

After installation:

1. Access Thunderdome at `http://localhost:8080` (or your configured domain)
2. Create an admin account on first run
3. Verify email functionality by testing the registration process

## Troubleshooting

- Check logs for any startup errors
- Verify database connection and credentials
- Ensure SMTP server is properly configured
- Confirm proper network connectivity and port availability

For additional help, refer to the [User Guide](GUIDE.md) or open an issue on GitHub.