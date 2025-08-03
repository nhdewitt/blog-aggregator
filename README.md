# Gator - RSS Feed Aggregator

A command-line RSS feed aggregator built in Go that allows you to manage feeds, follow your favorite blogs, and browse posts in your terminal.

## Features

- **User Management**: Create accounts and switch between users
- **Feed Management**: Add RSS feeds and browse all available feeds
- **Feed Following**: Follow/unfollow specific feeds
- **Post Aggregation**: Automatically fetch new posts from followed feeds
- **Post Browsing**: View latest posts with titles, descriptions, and publication dates
- **Multi-format Date Support**: Handles various RSS date formats automatically
- **HTML Entity Decoding**: Properly displays special characters in titles and descriptions

## Installation

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- [Goose](https://github.com/pressly/goose) for database migrations

### Install Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Build from Source

```bash
git clone https://github.com/nhdewitt/blog-aggregator.git
cd blog-aggregator

# Install dependencies
go mod download

# Build the application
go build -o gator ./cmd/gator
```

### Install with Go

```bash
go install github.com/nhdewitt/blog-aggregator/cmd/gator@latest
```

## Configuration

### Database Setup

1. **Create a PostgreSQL database:**
```sql
CREATE DATABASE gator;
```

2. **Run database migrations:**
```bash
# Navigate to your project directory
cd blog-aggregator

# Run migrations with Goose
goose -dir sql/schema postgres "postgres://username:password@localhost/gator?sslmode=disable" up
```

3. **Create configuration file** at `~/.gatorconfig.json`:
```json
{
  "db_url": "postgres://username:password@localhost/gator?sslmode=disable",
  "current_user_name": ""
}
```

### Environment Variables (Alternative)

You can also use environment variables instead of a config file:

```bash
export GATOR_DB_URL="postgres://username:password@localhost/gator?sslmode=disable"
```

## Usage

### User Management

```bash
# Create a new user
gator register alice

# Login as an existing user
gator login alice

# List all users
gator users

# Reset all user data (destructive!)
gator reset
```

### Feed Management

```bash
# Add a new RSS feed
gator addfeed "Tech Blog" https://example.com/feed.xml

# List all feeds in the system
gator feeds

# Follow an existing feed
gator follow https://example.com/feed.xml

# See feeds you're following
gator following

# Unfollow a feed
gator unfollow https://example.com/feed.xml
```

### Post Aggregation and Browsing

```bash
# Start the aggregator (runs continuously)
gator agg 30s  # Fetch feeds every 30 seconds
gator agg 5m   # Fetch feeds every 5 minutes

# Browse latest posts (default: 2 posts)
gator browse

# Browse specific number of posts
gator browse 10
```

## Commands Reference

| Command    | Usage          | Description                       |
|------------|----------------|-----------------------------------|
| `register` | `<username>`   | Create a new user account         |
| `login`    | `<username>`   | Switch to an existing user        |
| `users`    |                | List all registered users         |
| `reset`    |                | Delete all users and data         |
| `addfeed`  | `<name> <url>` | Add a new RSS feed                |
| `feeds`    |                | Show all feeds in the system      |
| `follow`   | `<url>`        | Follow an existing feed           |
| `following`|                | Show feeds you're following       |
| `unfollow` | `<url>`        | Stop following a feed             |
| `agg`      | `<duration>`   | Start continuous feed aggregation |
| `browse`   | `[limit]`      | Browse your latest posts          |
└-----------------------------------------------------------------┘

## Project Structure

```
blog-aggregator/
├── cmd/gator/           # Application entry point
│   └── main.go
├── internal/
│   ├── app/             # Application state and services
│   │   ├── fetch_feed.go   # RSS feed fetching logic
│   │   ├── state.go        # Shared state definition
│   │   ├── scrape_feeds.go # RSS feed scraping logic
|   |   └── rss_feed.go     # RSS data structures
│   ├── commands/        # CLI command system
│   │   ├── commands.go      # Command registration
│   │   ├── middleware.go    # Authentication middleware
│   │   ├── user_handlers.go # User management commands
│   │   ├── feed_handlers.go # Feed management commands
│   │   └── aggregator_handlers.go # Aggregation commands
│   ├── config/          # Configuration management
│   └── database/        # Database layer
├── sql/
│   └── schema/          # Goose database migrations
│       ├── 0100_users.sql
│       ├── 0101_feeds.sql
│       ├── 0102_feed_follows.sql
│       └── 0103_posts.sql
├── go.mod
├── go.sum
└── README.md
```

## Development

### Database Setup

1. **Create a PostgreSQL database:**
```sql
CREATE DATABASE gator;
```

2. **Run migrations:**
```bash
goose -dir sql/schema postgres "postgres://user:pass@localhost/gator?sslmode=disable" up
```

3. **Create new migrations (if needed):**
```bash
# Create a new migration file
goose -dir sql/schema create add_new_table sql

# This creates: sql/schema/YYYYMMDDHHMMSS_add_new_table.sql
```

4. **Migration commands:**
```bash
# Check migration status
goose -dir sql/schema postgres "your-db-url" status

# Roll back last migration
goose -dir sql/schema postgres "your-db-url" down

# Reset database (careful!)
goose -dir sql/schema postgres "your-db-url" reset
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Build for current platform
go build -o bin/gator ./cmd/gator

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o bin/gator-linux ./cmd/gator
GOOS=windows GOARCH=amd64 go build -o bin/gator-windows.exe ./cmd/gator
GOOS=darwin GOARCH=amd64 go build -o bin/gator-macos ./cmd/gator
```

## Examples

### Basic Workflow

```bash
# 1. Set up the database
createdb gator
goose -dir sql/schema postgres "postgres://user:pass@localhost/gator?sslmode=disable" up

# 2. Create a user
gator register john

# 3. Add some feeds
gator addfeed "Go Blog" https://blog.golang.org/feed.atom
gator addfeed "Hacker News" https://hnrss.org/frontpage

# 4. Start aggregating posts in background
gator agg 1m &

# 5. Browse latest posts
gator browse 5
```

### Following Feeds Added by Others

```bash
# See all available feeds
gator feeds

# Follow a feed by URL
gator follow https://blog.golang.org/feed.atom

# Check what you're following
gator following
```

## Troubleshooting

### Common Issues

**Database Connection Errors**
- Verify your database URL in the config file
- Ensure PostgreSQL is running
- Check that the database exists

**Feed Parsing Errors**
- The aggregator logs parsing errors but continues processing
- Some feeds may use non-standard date formats
- Check feed URLs are accessible and return valid RSS/Atom

**No Posts Appearing**
- Make sure you're following feeds (`gator following`)
- Run the aggregator to fetch posts (`gator agg 1m`)
- Some feeds may not have recent posts

### Logs

The aggregator outputs logs during feed processing:
```
Collecting feeds every 1m
Could not parse pubDate "invalid date" from feed "Example Blog": parsing time "invalid date": ...
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with Go's standard library
- Uses PostgreSQL for data storage
- RSS parsing handles multiple date formats commonly found in feeds