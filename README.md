# seats-aero-cli

A command-line interface for searching award flight availability using the [seats.aero](https://seats.aero) API.

## Features

- **Interactive Mode**: Guided prompts for easy searching
- **CLI Mode**: Command-line flags for scripting and automation
- **Multiple Export Formats**: JSON and CSV output
- **Cached Search**: Search for availability between specific airports and dates
- **Bulk Availability**: Retrieve large amounts of availability data for a mileage program
- **Route Listing**: View available routes for a mileage program
- **Trip Details**: Get detailed flight information for specific availability

## Installation

### From Source

```bash
go install github.com/JHill6253/seats-aero-cli/cmd/seats@latest
```

### Build Locally

```bash
git clone https://github.com/JHill6253/seats-aero-cli.git
cd seats-aero-cli
go build -o seats ./cmd/seats
```

## Configuration

### API Key

You need a seats.aero Pro account to use this CLI. Get your API key from the [seats.aero settings page](https://seats.aero/settings).

Set your API key using an environment variable:

```bash
export SEATS_AERO_API_KEY="your-api-key-here"
```

Or create a config file at `~/.config/seats-aero/config.yaml`:

```yaml
api_key: "your-api-key-here"

# Optional defaults
default_sources:
  - aeroplan
  - united
  - alaska

default_cabins:
  - J  # Business
  - F  # First

preferred_airports:
  - SFO
  - LAX
```

## Usage

### Interactive Mode (Default)

Launch the interactive guided CLI:

```bash
seats
```

This will present a menu-driven interface:

```
seats.aero CLI
Search for award flight availability

? What would you like to do?
  > Search for flights
    View bulk availability
    List routes
    Get trip details
    Exit
```

Follow the prompts to enter search criteria, view results, and export data.

### CLI Mode

Use flags for scripting and automation:

#### Search

Search for cached availability between airports:

```bash
# Basic search
seats search --from SFO --to NRT --start-date 2024-06-01

# Multiple airports
seats search --from SFO,LAX --to NRT,HND --start-date 2024-06-01 --end-date 2024-06-15

# Filter by cabin and source
seats search --from SFO --to NRT --cabin J --source united,aeroplan

# Direct flights only
seats search --from SFO --to NRT --direct-only

# Export to JSON
seats search --from SFO --to NRT --output json > results.json

# Export to CSV
seats search --from SFO --to NRT --output csv > results.csv
```

#### Bulk Availability

Get bulk availability for a mileage program:

```bash
# All availability for a program
seats availability --source aeroplan

# Filter by cabin
seats availability --source united --cabin J,F

# Filter by region
seats availability --source delta --origin-region north-america --dest-region europe
```

#### Routes

List available routes:

```bash
# All routes for a program
seats routes --source united

# Routes from a specific origin
seats routes --source aeroplan --origin SFO
```

#### Trip Details

Get detailed flight information:

```bash
seats trips <availability-id>
```

#### Configuration

View current configuration:

```bash
seats config show
```

## Cabin Classes

| Code | Class |
|------|-------|
| Y | Economy |
| W | Premium Economy |
| J | Business |
| F | First |

## Mileage Programs (Sources)

| Source | Program |
|--------|---------|
| `aeroplan` | Air Canada Aeroplan |
| `alaska` | Alaska Mileage Plan |
| `american` | American Airlines |
| `delta` | Delta SkyMiles |
| `united` | United MileagePlus |
| `emirates` | Emirates Skywards |
| `etihad` | Etihad Guest |
| `flyingblue` | Air France/KLM Flying Blue |
| `qantas` | Qantas Frequent Flyer |
| `singapore` | Singapore KrisFlyer |
| `virginatlantic` | Virgin Atlantic Flying Club |
| ... | [See full list](https://developers.seats.aero/reference/concepts-copy) |

## API Limits

- Pro users: 1,000 API calls per day
- Commercial use requires written agreement with seats.aero

## License

MIT License - See [LICENSE](LICENSE) for details.

## Disclaimer

This is an unofficial CLI tool and is not sponsored, endorsed, or approved by the seats.aero team. This is just a personal project for fun. Use of the seats.aero API is subject to their [terms of service](https://seats.aero/terms).
