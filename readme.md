# gator

A RSS feed aggregator built with Go and PostgreSQL. Following bootdev course.

## Prerequisites

- [Go](https://golang.org/dl/) 1.21+
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

```bash
go install github.com/tmuelli/gator@latest
```

## Configuration

Create a config file at `~/.gator.json`:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

| Field | Description |
|---|---|
| `db_url` | PostgreSQL connection string |
| `current_user_name` | The currently logged in user |

### PostgreSQL setup

```bash
# Create the database
createdb gator

# Or via psql
psql -c "CREATE DATABASE gator;"
```

## Usage

```bash
# Register a new user
gator register <username>

# Login as a user
gator login <username>

# Add a feed
gator addfeed <name> <url>

# List all feeds
gator feeds

# Follow a feed
gator follow <url>

# List followed feeds
gator following

# Aggregate posts
gator agg

# Browse posts
gator browse [limit]
```

## License

MIT