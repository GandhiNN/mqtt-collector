# mqtt-collector

## Key Features

- **Multi-broker Support**: Concurrent collection from 26+ brokers using goroutines
- **Payload Classification**: Automatic detection of JSON (objects/arrays only), XML, text, and binary payloads
- **Database Flexibility**: Supports both SQLite and PostgreSQL with automatic driver selection
- **Stateless Design**: Collectors can be restarted without state loss
- **Graceful Shutdown**: Proper cleanup on SIGTERM/SIGINT signals
- **API Endpoints**:
-- `POST /api/samples` - Store topic samples
-- `GET /api/topics` - List all topics with pagination
-- `GET /api/topics/search` - Find specific topic by broker+topic
-- `GET /health` - Health check

## Data Flow

MQTT Brokers → BrokerCollector → PayloadDetector → DBClient → API Server → Repository → Database
                     ↓
              [Topic Sampling & Classification]
                     ↓
              [HTTP POST to /api/samples]
                     ↓
              [Upsert to topics table]

### Collector Backend

Stateless MQTT topic discovery service that connects to multiple brokers simultaneously

1. Loads configuration from environment/files
2. Creates MultiCollector that manages multiple broker connections
3. Spawns individual BrokerCollector goroutines for each configured broker
4. Each collector subscribes to all topics (#) and samples unique topics once
5. Detected payloads are classified (JSON/XML/Text/Binary) using internal/payload
6. Samples are sent to the API server via HTTP client
7. Runs for configured duration with graceful shutdown on signals

### API Server

RESTful HTTP service for storing and retrieving topic data

1. Initializes database connection (SQLite/PostgreSQL) with migrations
2. Sets up repository layer for data access
3. Creates HTTP router with CORS and logging middleware
4. Exposes endpoints for topic CRUD operations
5. Handles graceful shutdown with connection draining

### Database Schema

The system uses a single `topics` table with upsert logic to maintain the latest sample for each broker-topic combination, tracking payload type, sample data, and timestamps.

This architecture enables scalable MQTT topic discovery across multiple brokers while providing a clean API for topic exploration and analysis.