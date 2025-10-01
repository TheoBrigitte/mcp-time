<p align="center">
    <img src="assets/mcp-time.png" alt="MCP time logo" height="100px">
</p>

<div align="center">

<a href="https://github.com/TheoBrigitte/mcp-time/releases"><img src="https://img.shields.io/github/release/TheoBrigitte/mcp-time.svg" alt="Github release"></a>
<a href="https://github.com/TheoBrigitte/mcp-time/actions/workflows/build.yaml"><img src="https://github.com/TheoBrigitte/mcp-time/actions/workflows/build.yaml/badge.svg?branch=main" alt="Github action"></a>
<a href="https://pkg.go.dev/github.com/TheoBrigitte/mcp-time"><img src="https://pkg.go.dev/badge/github.com/TheoBrigitte/mcp-time.svg" alt="Go reference"></a>
<a href="https://archestra.ai/mcp-catalog/theobrigitte__mcp-time"><img src="https://archestra.ai/mcp-catalog/api/badge/quality/TheoBrigitte/mcp-time" alt="Trust Score"></a>

<strong>Time MCP Server</strong>

*A Model Context Protocol server that enables AI assistants to interact with time*

</div>

## Overview

The Time MCP Server is a [Model Context Protocol (MCP)](https://github.com/modelcontextprotocol) server that provides AI assistants and other MCP clients with standardized tools to perform time and date-related operations. This server acts as a bridge between AI tools and a robust time-handling backend, allowing for complex time manipulations through natural language interactions.

## Features

- **Time Manipulation**: Get current time, convert between timezones, and add or subtract durations.
- **Natural Language Parsing**: Understands relative time expressions like "yesterday" or "next month".
- **Time Comparison**: Compare two different times.
- **Flexible Formatting**: Supports a wide variety of predefined and custom time formats.
- **MCP Compliance**: Fully compatible with the Model Context Protocol standard.
- **Multiple Transports**: Can be run using `stdio` for simple integrations or as an `HTTP stream` server for network access.

## Installation

This MCP server can be integrated with various AI assistants that support the Model Context Protocol.

### On Cursor

Use the link below to install directly in [Cursor](https://cursor.com).

<a href="cursor://anysphere.cursor-deeplink/mcp/install?name=time&config=eyJjb21tYW5kIjoiZG9ja2VyIiwiYXJncyI6WyJydW4iLCItLXJtIiwiLWkiLCJ0aGVvMDEvbWNwLXRpbWU6bGF0ZXN0Il19">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://cursor.com/deeplink/mcp-install-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://cursor.com/deeplink/mcp-install-light.svg">
  <img alt="Add to Cursor" src="https://cursor.com/deeplink/mcp-install-light.svg" width="200">
</picture>
</a>

### Using npx (JavaScript/Node.js)

Using this method will run the MCP server using [`npx`](https://docs.npmjs.com/cli/v8/commands/npx), which requires Node.js to be installed on your system. Then copy the following JSON configuration into your MCP client to run the server:

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

Copy the following JSON configuration into your MCP client to run the server using Docker:

```json
{
  "mcpServers": {
    "time": {
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

Copy the following JSON configuration into your MCP client to run the server using the binary:

```json
{
  "mcpServers": {
    "time": {
      "command": "mcp-time"
    }
  }
}
```

You need to install the `mcp-time` binary on your system. You can do this in several ways:

#### Install from releases

You can download the latest binary from the [releases page](https://github.com/TheoBrigitte/mcp-time/releases).

```bash
# Replace OS-ARCH with your operating system and architecture (e.g., linux-amd64, darwin-arm64)
curl -Lo mcp-time https://github.com/TheoBrigitte/mcp-time/releases/latest/download/mcp-time.OS-ARCH
install -D -m 755 ./mcp-time ~/.local/bin/mcp-time
```

#### Install with Go

```bash
go install github.com/TheoBrigitte/mcp-time/cmd/mcp-time@latest
```

This will install the `mcp-time` binary in your `$GOPATH/bin` directory.

#### Building from Source

```bash
git clone https://github.com/TheoBrigitte/mcp-time.git
cd mcp-time
make install
```

This will build and install the `mcp-time` binary in the `~/.local/bin` directory.

## Usage

### Basic Usage

Start the MCP server with the default `stdio` transport:

```bash
mcp-time
```

Start the MCP server with the `stream` transport:

```bash
mcp-time --transport stream --address "http://localhost:8080/mcp"
```

### Advanced Usage

The server supports several command-line options for more advanced configurations:

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

### `relative_time`

Get a time based on a relative natural language expression.

**Parameters:**
- `text` (required): The natural language expression (e.g., `yesterday`, `5 minutes ago`, `next month`).
- `time` (optional): A reference time for the relative expression. Defaults to current time.
- `timezone` (optional): The target timezone for the output.
- `format` (optional): The output format for the time.

### `convert_timezone`

Convert a given time between timezones.

**Parameters:**
- `time` (required): The input time string. It can be in various formats.
- `input_timezone` (optional): The timezone of the input time.
- `output_timezone` (optional): The target timezone for the output.
- `format` (optional): The output format for the time.

### `current_time`

Get the current time.

**Parameters:**
- `format` (optional): The output format for the time. Can be a predefined format (e.g., `RFC3339`, `Kitchen`) or a custom Go layout.
- `timezone` (optional): The target timezone in IANA format (e.g., `America/New_York`). Defaults to UTC.

### `add_time`

Add or subtract a duration to a given time.

**Parameters:**
- `time` (required): The input time string.
- `duration` (required): The duration to add or subtract (e.g., `2h30m`, `-1h`).
- `timezone` (optional): The target timezone for the output.
- `format` (optional): The output format for the time.

### `compare_time`

Compare two times.

**Parameters:**
- `time_a` (required): The first time to compare.
- `time_b` (required): The second time to compare.

**Returns:**
- `-1` if `time_a` is before `time_b`.
- `0` if `time_a` is equal to `time_b`.
- `1` if `time_a` is after `time_b`.

## Credits

- https://github.com/araddon/dateparse
- https://github.com/tj/go-naturaldate
- https://github.com/mark3labs/mcp-go
