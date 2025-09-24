# Peril

[![Go Version](https://img.shields.io/badge/Go-1.22.1-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-3.13-orange.svg)](https://www.rabbitmq.com/)

A distributed real-time strategy game built with Go and RabbitMQ, featuring pub/sub messaging architecture for scalable multiplayer gameplay.

## Overview

Peril is a multiplayer strategy game where players command armies, move units across territories, and engage in tactical warfare. The game demonstrates advanced distributed systems concepts including:

- **Event-driven architecture** with RabbitMQ pub/sub messaging
- **Real-time multiplayer** synchronization
- **Scalable server architecture** supporting multiple concurrent game instances
- **Microservices design** with separate client and server components

> **Note**: The core game functionality was forked from the starter project in Boot.dev's [Learn Pub/Sub](https://learn.boot.dev/learn-pub-sub) course and has been enhanced with additional features and production-ready architecture.

## Features

- 🎮 **Real-time multiplayer gameplay** with instant move synchronization
- 🏰 **Army management** - spawn, move, and command military units
- ⚔️ **Strategic warfare** with unit rankings and battle mechanics
- 🔄 **Live game state** with pause/resume functionality
- 📊 **Game logging** and activity tracking
- 🐳 **Docker support** for easy deployment
- 🔧 **Horizontal scaling** with multi-server support

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Game Client   │    │   Game Server   │    │   RabbitMQ      │
│                 │    │                 │    │   Message       │
│ • User Input    │◄──►│ • Game Logic    │◄──►│   Broker        │
│ • Game Display  │    │ • State Mgmt    │    │                 │
│ • Move Commands │    │ • Event Handling│    │ • Pub/Sub       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Components

- **Client (`cmd/client/`)**: Player interface handling user input and game state visualization
- **Server (`cmd/server/`)**: Game logic processor managing game state and player interactions
- **Game Logic (`internal/gamelogic/`)**: Core game mechanics including moves, spawning, and warfare
- **Pub/Sub (`internal/pubsub/`)**: RabbitMQ messaging layer for real-time communication
- **Routing (`internal/routing/`)**: Message routing and exchange configuration

## Prerequisites

- **Go 1.22.1** or higher
- **Docker** (for RabbitMQ)
- **RabbitMQ 3.13** (automatically managed via Docker)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/dmitriy-zverev/peril.git
cd peril
```

### 2. Start RabbitMQ

```bash
chmod +x rabbit.sh
./rabbit.sh start
```

This will start a RabbitMQ container with management interface available at http://localhost:15672 (guest/guest).

### 3. Start the Game Server

```bash
go run cmd/server/main.go
```

### 4. Start Game Clients

In separate terminals, start multiple clients:

```bash
go run cmd/client/main.go
```

## Usage

### Server Commands

- `pause` - Pause the game for all players
- `resume` - Resume the game
- `quit` - Shutdown the server

### Client Commands

- `spawn <rank> <location>` - Spawn a new unit
  - Example: `spawn sergeant 42,73`
- `move <from> <to>` - Move units between locations
  - Example: `move 42,73 45,76`
- `status` - Display current game state
- `spam <count>` - Generate test log messages
- `quit` - Disconnect from the game

### Unit Ranks

Units have different capabilities based on rank:
- `private` - Basic infantry unit
- `sergeant` - Enhanced combat unit
- `lieutenant` - Advanced tactical unit
- `captain` - Elite command unit

## Scaling

### Multiple Server Instances

Run multiple game servers for load distribution:

```bash
chmod +x multiserver.sh
./multiserver.sh 3  # Start 3 server instances
```

### Docker Deployment

Build and run with Docker:

```bash
docker build -t peril .
docker run -p 8080:8080 peril
```

## Development

### Project Structure

```
peril/
├── cmd/
│   ├── client/          # Client application
│   └── server/          # Server application
├── internal/
│   ├── gamelogic/       # Core game mechanics
│   ├── pubsub/          # RabbitMQ messaging
│   └── routing/         # Message routing
├── scripts/
│   ├── rabbit.sh        # RabbitMQ management
│   └── multiserver.sh   # Multi-server deployment
└── Dockerfile           # Container configuration
```

### Building

```bash
# Build server
go build -o bin/server cmd/server/main.go

# Build client
go build -o bin/client cmd/client/main.go
```

### Testing

```bash
go test ./...
```

## Configuration

### RabbitMQ Connection

Default connection string: `amqp://guest:guest@localhost:5672/`

To use a different RabbitMQ instance, modify the connection string in:
- `cmd/server/main.go`
- `cmd/client/main.go`

### Message Exchanges

- **Direct Exchange**: `peril_direct` - For targeted messages (pause/resume)
- **Topic Exchange**: `peril_topic` - For broadcast messages (moves, logs)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is free to use without any license.

## Acknowledgments

- **Boot.dev** - Original starter project and pub/sub learning materials
- **RabbitMQ** - Robust messaging broker
- **Go Community** - Excellent ecosystem and libraries

## Support

For questions and support:
- Open an issue on GitHub
- Check the [Boot.dev Learn Pub/Sub course](https://learn.boot.dev/learn-pub-sub) for foundational concepts

---

**Built with ❤️ using Go and RabbitMQ**
