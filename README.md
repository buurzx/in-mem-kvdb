# In-Memory Key-Value Database

A simple TCP-based in-memory key-value database with CLI client.

## Features

- In-memory key-value storage
- TCP server with configurable connection handling
- Interactive CLI client
- Basic operations: GET, SET, DEL
- Configurable idle timeout and connection limits
- Graceful shutdown handling

## Building

```bash
go build -o kvdb ./cmd/main.go
```

## Running the Server

Start the server:
```bash
./kvdb server
```

Default server configuration:
- Address: localhost:8080
- Max connections: 100
- Idle timeout: 300 seconds
- Buffer size: 4KB

## Using the CLI Client

Start the CLI client:
```bash
./kvdb cli
```

### Available Commands

1. Set a value:
```bash
[in-mem-kvdb] > SET key value
[OK]
```

2. Get a value:
```bash
[in-mem-kvdb] > GET key
value
```

3. Delete a key:
```bash
[in-mem-kvdb] > DEL key
[OK]
```

4. Exit the CLI:
```bash
[in-mem-kvdb] > exit
```

### Error Examples

1. Key not found:
```bash
[in-mem-kvdb] > GET nonexistent
not found
```

2. Invalid command:
```bash
[in-mem-kvdb] > INVALID
Invalid command. Available commands: DEL, GET, SET
```

3. Wrong number of arguments:
```bash
[in-mem-kvdb] > SET key
Invalid SET command. Usage: SET <key> <value>
```

## Configuration

Both server and client can be configured using environment variables:

### Server Configuration
```bash
KVDB_NETWORK_ADDRESS=":3000" \
KVDB_NETWORK_MAX_CONNECTIONS=200 \
KVDB_NETWORK_IDLE_TIMEOUT=600 \
KVDB_NETWORK_MAX_MESSAGE_SIZE=8192 \
./kvdb server
```

### Client Configuration
```bash
KVDB_NETWORK_ADDRESS="localhost:3000" \
KVDB_NETWORK_IDLE_TIMEOUT=600 \
KVDB_NETWORK_MAX_MESSAGE_SIZE=8192 \
./kvdb cli
```

## Error Handling

- Connection timeouts are handled gracefully
- Invalid commands return descriptive error messages
- Connection limits are enforced to prevent overload
- Graceful shutdown on SIGINT/SIGTERM signals
- Automatic connection cleanup for idle clients

## Project Structure
```
.
├── cmd/
│   ├── main.go          # Application entry point
│   ├── server/          # Server implementation
│   └── kvcli/           # CLI client implementation
├── internal/
│   ├── network/         # Network handling
│   │   └── tcp/         # TCP server and client
│   ├── database/        # Database implementation
│   └── initialization/  # Shared initialization code
└── README.md
```

## Development

The project follows a clean architecture approach:
- Separation of concerns between network, database, and client code
- Configuration through environment variables
- Extensive error handling and logging
- Graceful shutdown handling
