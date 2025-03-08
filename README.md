```ascii
 _    _ _____ _                 _ 
| |  | /  __ \ |               | |
| |  | | /  \/ | ___  _   _  __| |
| |  | | |   | |/ _ \| | | |/ _` |
| |__| | \__/\ | (_) | |_| | (_| |
 \____/ \____/_|\___/ \__,_|\__,_|
                                  
 __  __  _____ _____   _____                          
|  \/  |/  __ \  __ \ /  ___|                         
| .  . || /  \/| |_/ / \ `--.  ___ _ ____   _____ _ __
| |\/| || |    |  __/   `--. \/ _ \ '__\ \ / / _ \ '__|
| |  | || \__/\| |     /\__/ /  __/ |   \ V /  __/ |   
\_|  |_/ \____/\_|     \____/ \___|_|    \_/ \___|_|   
```

# UCloud MCP Server

A cloud instance management server based on MCP-Go and UCloud SDK, supporting UCloud instance management through the MCP protocol.

## Features

- Query instance information
- List all instances
- Get instance status
- Monitor instance performance metrics
- MCP protocol support
- Configuration file support

## Requirements

- Go 1.23 or higher
- UCloud account and API credentials

## Configuration

The service supports two configuration methods:

### 1. Configuration File (Recommended)

Create a `config.json` file:

```json
{
    "region": "cn-bj2",
    "project_id": "your-project-id",
    "public_key": "your-public-key",
    "private_key": "your-private-key"
}
```

### 2. Environment Variables

If not specified in the configuration file, the service will try to read from environment variables:

```bash
export UCLOUD_REGION="cn-bj2"        # UCloud region
export UCLOUD_PROJECT_ID="your-project-id"  # Project ID
export UCLOUD_PUBLIC_KEY="your-public-key"  # API public key
export UCLOUD_PRIVATE_KEY="your-private-key"  # API private key
```

Configuration priority: Configuration file > Environment variables

## Installation and Running

1. Clone the repository:
```bash
git clone https://github.com/renzheng.wang/ucloud-mcp-server.git
cd ucloud-mcp-server
```

2. Install dependencies:
```bash
go mod download
```

3. Build the service:
```bash
go build -o ucloud-mcp-server
```

4. Run the service:

Basic usage:
```bash
./ucloud-mcp-server
```

With custom configuration:
```bash
./ucloud-mcp-server --config /path/to/config.json --port 8080
```

Available startup options:
- `--config`: Specify the path to your configuration file (default: ./config.json)
- `--port`: Specify the port to listen on (default: 8080)

Examples:
```bash
# Use custom config file
./ucloud-mcp-server --config /etc/ucloud/config.json

# Use custom port
./ucloud-mcp-server --port 9000

# Use both custom config and port
./ucloud-mcp-server --config /etc/ucloud/config.json --port 9000
```

The service will provide MCP protocol service through standard input/output.

## Available Operations

### Instance Information
Get detailed information about a specific instance, including:
- Basic instance details
- Configuration information
- Current status
- Resource allocation

### Instance Status
Monitor the current operational status of any instance in real-time.

### Instance Metrics
Access comprehensive monitoring metrics for instances, including:
- CPU utilization
- Disk I/O operations
- Network traffic statistics
- System performance data

### Instance List
View a complete list of all available instances in your account, including their basic information and current status.

## Monitoring Metrics

The system provides the following monitoring metrics:

- **CPU Metrics**
  - CPUUtilization: CPU usage percentage (%)

- **Disk Metrics**
  - IORead: Disk read rate
  - IOWrite: Disk write rate
  - DiskReadOps: Number of disk read operations
  - DiskWriteOps: Number of disk write operations

- **Network Metrics**
  - NICIn: Network inbound traffic (bytes/s)
  - NICOut: Network outbound traffic (bytes/s)
  - NetPacketIn: Number of inbound network packets
  - NetPacketOut: Number of outbound network packets

## Important Notes

- Keep your UCloud API credentials secure
- Use configuration files or key management services for sensitive information in production environments
- All operations are performed through the MCP protocol with standard I/O support
- Monitoring data may have a few minutes delay
- Regularly check monitoring metrics to identify potential issues early
