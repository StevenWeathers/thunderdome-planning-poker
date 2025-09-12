# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Using Task (Recommended)
```bash
# Install dependencies
task install

# Build and run development server
task dev                # Runs on http://localhost:8080
task dev-secure         # Runs on https://thunderdome.localhost with Caddy

# Build only
task build              # Full build with frontend and backend
task build-ui           # Build frontend only
task build-go           # Build backend only

# Testing
task test-go            # Run Go tests

# Code formatting (required before commit)
task format             # Format all code (Go, UI, E2E)
task fmt-go             # Format Go code only
task fmt-ui             # Format UI code only

# Generate documentation
task gen-swag           # Generate Swagger API docs
task gen-i8n            # Generate i18n localizations

# Clean build artifacts
task clean
```

### Using Make (Alternative)
```bash
make install            # Install dependencies
make dev                # Build and run development server
make build              # Build application
make format             # Format all code
make testgo             # Run Go tests
make generate           # Generate docs and localizations
```

### Frontend Commands (in ui/ directory)
```bash
npm run dev             # Start Vite dev server
npm run build           # Build frontend
npm run check           # Run Svelte type checking
npm test                # Run Jest tests
npm run format          # Format frontend code
npm run locales         # Generate i18n types
```

### Database Migrations
```bash
# Create new migration (replace NAME with descriptive name)
go tool goose -dir internal/db/migrations create NAME sql
```

## Architecture Overview

### Project Structure
- **Go Backend** (`/internal`): RESTful API server with WebSocket support
  - `/internal/http`: HTTP handlers and routing
  - `/internal/db`: Database layers for different domains (poker, retro, storyboard, team, user)
  - `/internal/wshub`: WebSocket hub for real-time communication
  - `/internal/config`: Application configuration
  - `/internal/email`: Email notification service
  - `/internal/oauth`: OAuth authentication providers
  - `/internal/atlassian`: Jira integration
  - `/internal/webhook`: Webhook support for subscriptions

- **Svelte Frontend** (`/ui`): Single-page application
  - `/ui/src/pages`: Route components for different features
  - `/ui/src/components`: Reusable UI components
  - `/ui/src/i18n`: Internationalization support
  - `/ui/src/stores.ts`: Svelte stores for state management
  - `/ui/src/apiclient.ts`: API client for backend communication

- **Database**: PostgreSQL with Goose migrations
  - Migrations in `/internal/db/migrations`
  - Domain-specific repositories in `/internal/db/*`

### Key Features & Domains
1. **Planning Poker** (`/internal/db/poker`): Agile estimation sessions
2. **Retrospectives** (`/internal/db/retro`): Sprint retrospective meetings
3. **Story Mapping** (`/internal/db/storyboard`): User story organization
4. **Team Check-ins**: Async stand-up functionality
5. **Team Management** (`/internal/db/team`): Organization and team structures
6. **User Management** (`/internal/db/user`): User accounts and authentication

### Real-time Communication
- WebSocket hub pattern in `/internal/wshub`
- Handles real-time updates for poker sessions, retros, and storyboards
- Uses Sockette on frontend for WebSocket connection management

### Authentication & Security
- Cookie-based session management
- OAuth support (multiple providers)
- API key authentication for external access
- LDAP integration support

### Frontend Stack
- **Svelte 5** with TypeScript
- **Vite** for build tooling
- **TailwindCSS** for styling
- **Navaid** for routing
- **TypeSafe-i18n** for internationalization
- **Lucide** for icons

### API Documentation
- Swagger/OpenAPI documentation auto-generated from code comments
- Located at `/docs/swagger` after generation
- Accessible at `/api/docs` when running

## Development Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL
- Task (recommended) or Make
- Goose (for migrations)
- Caddy (optional, for HTTPS development)

## Database Configuration
Default development database settings (from docker-compose.yml):
- Host: localhost
- Database: thunderdome
- User: thor
- Password: odinson

## Important Development Notes
- Always run `task format` before committing to avoid CI failures
- API changes require regenerating Swagger docs with `task gen-swag`
- Frontend changes to i18n require running `task gen-i8n`
- The application embeds the UI build into the Go binary
- WebSocket connections are managed through a hub pattern for scalability