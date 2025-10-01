<p align="center">
    <img src="assets/mcp-time.png" alt="MCP time logo" height="100px">
</p>

<div align="center">

[![GitHub release](https://img.shields.io/github/release/TheoBrigitte/mcp-time.svg?color)](https://github.com/TheoBrigitte/mcp-time/releases)
[![Build Status](https://github.com/TheoBrigitte/mcp-time/actions/workflows/build.yaml/badge.svg?branch=main)](https://github.com/TheoBrigitte/mcp-time/actions/workflows/build.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/TheoBrigitte/mcp-time.svg)](https://pkg.go.dev/github.com/TheoBrigitte/mcp-time)
[![Trust Score](https://archestra.ai/mcp-catalog/api/badge/quality/TheoBrigitte/mcp-time)](https://archestra.ai/mcp-catalog/theobrigitte__mcp-time)
[![NPM Package](https://img.shields.io/npm/v/@theo.foobar/mcp-time?color)](https://www.npmjs.com/package/@theo.foobar/mcp-time)
[![Docker Image](https://img.shields.io/docker/v/theo01/mcp-time?label=docker)](https://hub.docker.com/r/theo01/mcp-time)

<strong>Time MCP Server</strong>

*A Model Context Protocol server that enables AI assistants to interact with time*

</div>

## Overview

The Time MCP Server is a [Model Context Protocol (MCP)](https://github.com/modelcontextprotocol) server that provides AI assistants and other MCP clients with standardized tools to perform time and date-related operations. This server acts as a bridge between AI tools and a robust time-handling backend, allowing for complex time manipulations through natural language interactions.

## Features

- **⏰ Time Manipulation** - Get current time, convert between timezones, and add or subtract durations
- **🗣️ Natural Language Parsing** - Understands relative time expressions like "yesterday" or "next month"
- **⚖️ Time Comparison** - Compare two different times with ease
- **🎨 Flexible Formatting** - Supports a wide variety of predefined and custom time formats
- **✅ MCP Compliance** - Fully compatible with the Model Context Protocol standard
- **🔄 Multiple Transports** - Supports `stdio` for local integrations and `HTTP stream` for network access

## Installation

This MCP server can be integrated with various AI assistant clients that support the Model Context Protocol, including [Cursor](https://cursor.com/), [Claude Desktop](https://claude.ai/download), [Claude Code](https://www.claude.com/product/claude-code), and [many more](https://modelcontextprotocol.io/clients).

### 🚀 One-Click Install (Cursor)

Click the button below to automatically configure the MCP server using Docker in your [Cursor](https://cursor.com) environment:

<a href="cursor://anysphere.cursor-deeplink/mcp/install?name=time&config=eyJjb21tYW5kIjoiZG9ja2VyIiwiYXJncyI6WyJydW4iLCItLXJtIiwiLWkiLCJ0aGVvMDEvbWNwLXRpbWU6bGF0ZXN0Il19">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://cursor.com/deeplink/mcp-install-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://cursor.com/deeplink/mcp-install-light.svg">
  <img alt="Add to Cursor" src="https://cursor.com/deeplink/mcp-install-light.svg" width="200">
</picture>
</a>

### Using npx (JavaScript/Node.js)

This method runs the MCP server using [`npx`](https://docs.npmjs.com/cli/v8/commands/npx), which requires Node.js to be installed. Copy the following JSON configuration into your MCP client settings:

```json
{
  "mcpServers": {
    "mcp-time": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@theo.foobar/mcp-time"
      ]
    }
  }
}
```

### Using Docker

Run the MCP server in an isolated container. Requires [Docker](https://www.docker.com/get-started/) to be installed. Copy this JSON configuration into your MCP client settings:

```json
{
  "mcpServers": {
    "mcp-time": {
      "type": "stdio",
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "theo01/mcp-time:latest"
      ]
    }
  }
}
```

### Using binary

Install the `mcp-time` binary directly on your system. Choose one of the installation methods below, ensuring the binary is placed in a directory that's in your `PATH`. Then add this JSON configuration to your MCP client settings:

```json
{
  "mcpServers": {
    "mcp-time": {
      "type": "stdio",
      "command": "mcp-time"
    }
  }
}
```

#### Option 1: Download from Releases

Download the latest pre-built binary from the [releases page](https://github.com/TheoBrigitte/mcp-time/releases):

```bash
# Replace OS-ARCH with your platform (e.g., linux-amd64, darwin-arm64, windows-amd64)
curl -Lo mcp-time https://github.com/TheoBrigitte/mcp-time/releases/latest/download/mcp-time.OS-ARCH
install -D -m 755 ./mcp-time ~/.local/bin/mcp-time
```

#### Option 2: Install with Go

For Go developers, install directly using `go install`:

```bash
go install github.com/TheoBrigitte/mcp-time/cmd/mcp-time@latest
```

The binary will be installed in your `$GOPATH/bin` directory.

#### Option 3: Build from Source

Clone and build the project using `make`:

```bash
git clone https://github.com/TheoBrigitte/mcp-time.git
cd mcp-time
make install
```

The binary will be installed in `~/.local/bin/mcp-time`.

## Usage

### Basic Usage

**Start with stdio transport** (default, for MCP clients):
```bash
mcp-time
```

**Start with HTTP stream transport** (for network access):
```bash
mcp-time --transport stream --address "http://localhost:8080/mcp"
```

### Command-Line Options

The server supports the following flags for advanced configurations:

```
$ mcp-time --help
An MCP (Model Context Protocol) server which provides utilities to work with time and dates.

Usage:
  mcp-time [flags]

Flags:
      --address string     Listen address for Stream HTTP Server (only for --transport stream) (default "http://localhost:8080/mcp")
  -h, --help               help for mcp-time
      --log-file string    Path to log file (logs is disabled if not specified)
  -t, --transport string   Transport layer: stdio, stream. (default "stdio")
      --version            Print version information and exit
```

## Available Tools

### `current_time`

Get the current time in any timezone and format.

**Parameters:**
- `format` (optional) - The output format (predefined like `RFC3339`, `Kitchen`, or custom Go layout)
- `timezone` (optional) - Target timezone in IANA format (e.g., `America/New_York`). Defaults to UTC

**Example:** "What time is it in Tokyo?"

### `relative_time`

Get a time based on a relative natural language expression.

**Parameters:**
- `text` (required) - Natural language expression (e.g., `yesterday`, `5 minutes ago`, `next month`)
- `time` (optional) - Reference time for the expression. Defaults to current time
- `timezone` (optional) - Target timezone for the output
- `format` (optional) - Output format for the time

**Example:** "What was the date 3 weeks ago?"

### `convert_timezone`

Convert a given time between timezones.

**Parameters:**
- `time` (required) - Input time string (supports various formats)
- `input_timezone` (optional) - Timezone of the input time
- `output_timezone` (optional) - Target timezone for the output
- `format` (optional) - Output format for the time

**Example:** "Convert 2:30 PM EST to Tokyo time"

### `add_time`

Add or subtract a duration from a given time.

**Parameters:**
- `time` (required) - Input time string
- `duration` (required) - Duration to add/subtract (e.g., `2h30m`, `-1h`, `24h`)
- `timezone` (optional) - Target timezone for the output
- `format` (optional) - Output format for the time

**Example:** "What time will it be in 45 minutes?"

### `compare_time`

Compare two times and determine their relationship. Supports timezone-aware comparisons.

**Parameters:**
- `time_a` (required) - First time to compare
- `time_a_timezone` (optional) - Timezone for `time_a` in IANA format (e.g., `America/New_York`)
- `time_b` (required) - Second time to compare
- `time_b_timezone` (optional) - Timezone for `time_b` in IANA format (e.g., `Europe/London`)

**Returns:**
- `-1` if `time_a` is before `time_b`
- `0` if `time_a` equals `time_b`
- `1` if `time_a` is after `time_b`

**Example:** "Is 3 PM EST before 8 PM GMT?"

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

Built with these excellent libraries:
- [araddon/dateparse](https://github.com/araddon/dateparse) - Parse dates without knowing the format
- [tj/go-naturaldate](https://github.com/tj/go-naturaldate) - Natural language date parsing
- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - Model Context Protocol SDK for Go
