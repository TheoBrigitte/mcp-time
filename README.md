<p align="center">
    <img src="assets/mcp-time.png" alt="MCP time logo" height="100px">
</p>

<div align="center">

<a href="https://github.com/TheoBrigitte/mcp-time/releases"><img src="https://img.shields.io/github/release/TheoBrigitte/mcp-time.svg" alt="Github release"></a>
<a href="https://github.com/TheoBrigitte/mcp-time/actions/workflows/go.yaml"><img src="https://github.com/TheoBrigitte/mcp-time/actions/workflows/go.yaml/badge.svg?branch=main" alt="Github action"></a>
<a href="https://pkg.go.dev/github.com/TheoBrigitte/mcp-time"><img src="https://pkg.go.dev/badge/github.com/TheoBrigitte/mcp-time.svg)](https://pkg.go.dev/github.com/TheoBrigitte/mcp-time"></a>

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

## Prerequisites

- Go 1.24.2 or later

## Installation

### Building from Source

```bash
git clone https://github.com/TheoBrigitte/mcp-time.git
cd mcp-time
make install
```

This will build and install the `mcp-time` binary in the `~/.local/bin` directory, which should be in your `PATH`.

### Using Go Install

```bash
go install github.com/TheoBrigitte/mcp-time/cmd/mcp-time@latest
```

## Integration with AI Assistants

This MCP server can be integrated with various AI assistants that support the Model Context Protocol.

### Example MCP Client Configuration

```json
{
  "servers": {
    "time": {
      "command": "/path/to/mcp-time"
    }
  }
}
```

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

```bash
mcp-time --help
```

```
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

Get the current time.

**Parameters:**
- `format` (optional): The output format for the time. Can be a predefined format (e.g., `RFC3339`, `Kitchen`) or a custom Go layout.
- `timezone` (optional): The target timezone in IANA format (e.g., `America/New_York`). Defaults to UTC.

### `convert_timezone`

Convert a given time between timezones.

**Parameters:**
- `time` (required): The input time string. It can be in various formats.
- `input_timezone` (optional): The timezone of the input time.
- `output_timezone` (optional): The target timezone for the output.
- `format` (optional): The output format for the time.

### `add_time`

Add or subtract a duration to a given time.

**Parameters:**
- `time` (required): The input time string.
- `duration` (required): The duration to add or subtract (e.g., `2h30m`, `-1h`).
- `timezone` (optional): The target timezone for the output.
- `format` (optional): The output format for the time.

### `relative_time`

Get a time based on a relative natural language expression.

**Parameters:**
- `text` (required): The natural language expression (e.g., `yesterday`, `5 minutes ago`, `next month`).
- `time` (optional): A reference time for the relative expression. Defaults to current time.
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
