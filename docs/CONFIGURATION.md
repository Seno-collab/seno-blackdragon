# Configuration Guide

This document explains how to configure the Seno Blackdragon application using the new Viper-based configuration system.

## Configuration Sources

The application supports multiple configuration sources in the following order of precedence (highest to lowest):

1. **Environment Variables** (highest priority - always override other sources)
2. **Configuration Files** (YAML format - lower priority)
3. **Default Values** (lowest priority - fallback only)

## Configuration File

Create a `config/config.yaml` file in your project root or use the provided example:

```yaml
environment: development

jwt:
  access_secret: "your-super-secret-access-key"
  refresh_secret: "your-super-secret-refresh-key"

redis:
  host: "localhost"
  port: 6379
  db: 0
  password: ""

db:
  host: "localhost"
  port: "5432"
  name: "seno_blackdragon"
  user: "postgres"
  password: "your-db-password"
  sslmode: "disable"

server:
  host: "0.0.0.0"
  port: "8080"
```

## Environment Variables

Environment variables have the **highest priority** and will always override values from configuration files and defaults. This ensures that you can easily override any configuration at runtime without modifying files.

You can use environment variables to override configuration values:

```bash
# JWT Configuration
export JWT_ACCESS_SECRET="your-access-secret"
export JWT_REFRESH_SECRET="your-refresh-secret"

# Redis Configuration
export REDIS_HOST="localhost"
export REDIS_PORT="6379"
export REDIS_PASSWORD=""

# Database Configuration
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_NAME="seno_blackdragon"
export DB_USER="postgres"
export DB_PASSWORD="your-password"

# Server Configuration
export SERVER_HOST="0.0.0.0"
export SERVER_PORT="8080"

# Environment
export ENVIRONMENT="production"
```

## Environment Variable Mapping

The following environment variables map to configuration keys:

| Environment Variable | Configuration Key    | Description                |
| -------------------- | -------------------- | -------------------------- |
| `JWT_ACCESS_SECRET`  | `jwt_access_secret`  | JWT access token secret    |
| `JWT_REFRESH_SECRET` | `jwt_refresh_secret` | JWT refresh token secret   |
| `REDIS_HOST`         | `redis_host`         | Redis server hostname      |
| `REDIS_PORT`         | `redis_port`         | Redis server port          |
| `REDIS_DB`           | `redis_db`           | Redis database number      |
| `REDIS_PASSWORD`     | `redis_password`     | Redis server password      |
| `DB_HOST`            | `db_host`            | PostgreSQL server hostname |
| `DB_PORT`            | `db_port`            | PostgreSQL server port     |
| `DB_NAME`            | `db_name`            | PostgreSQL database name   |
| `DB_USER`            | `db_user`            | PostgreSQL username        |
| `DB_PASSWORD`        | `db_password`        | PostgreSQL password        |
| `DB_SSLMODE`         | `db_sslmode`         | PostgreSQL SSL mode        |
| `SERVER_HOST`        | `server_host`        | Server bind address        |
| `SERVER_PORT`        | `server_port`        | Server listen port         |
| `ENVIRONMENT`        | `environment`        | Application environment    |

## Configuration File Locations

The application searches for configuration files in the following locations:

1. Current working directory (`.`)
2. `./config` subdirectory
3. `/etc/seno-blackdragon` (for system-wide configuration)

## Backward Compatibility

The application still supports `.env` files for backward compatibility, but environment variables take precedence over config file values.

## Validation

The configuration system validates required fields on startup and will fail with descriptive error messages if critical configuration is missing.

## Usage in Code

```go
cfg := config.LoadConfig(logger)

// Access configuration values
port := cfg.Server.Port
dbHost := cfg.DB.Host

// Use helper methods
if cfg.IsDevelopment() {
    // Development-specific logic
}

if cfg.IsProduction() {
    // Production-specific logic
}

// Access raw Viper values
customValue := cfg.GetString("custom.key")

// Debug configuration sources
debug := cfg.DebugConfigSources()
fmt.Printf("Configuration debug info: %+v\n", debug)
```

## Troubleshooting

### Environment Variables Not Working?

1. **Check variable names**: Ensure you're using the correct environment variable names (e.g., `JWT_ACCESS_SECRET`, not `SENO_JWT_ACCESS_SECRET`)
2. **Verify precedence**: Environment variables should override config file values
3. **Check logs**: Look for "Configuration sources loaded" log message to see what sources were loaded
4. **Use debug method**: Call `cfg.DebugConfigSources()` to see which values are set from environment vs defaults

### Common Issues

- **Environment variables ignored**: Make sure you're not setting them after the application starts
- **Config file overriding env vars**: This shouldn't happen - environment variables have highest priority
- **Missing required fields**: Check that all required environment variables are set
