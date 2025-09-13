## Speech-to-Text Recognition Service

<img width="1524" height="305" alt="image" src="https://github.com/user-attachments/assets/3586c131-5df5-4c08-81d3-01e49b0a9701" />


### Project Overview

**Speech-to-Text** is a modern service for converting speech to text using advanced speech recognition technologies. The project provides accurate audio-to-text conversion capabilities with robust monitoring and management features.

### Technology Stack

* **Programming Language**: Go (Golang)
* **Database**: PostgreSQL
* **Message Broker**: RabbitMQ
* **Monitoring**:
    * Prometheus
    * Grafana
* **External API Integration**
* **Core Components**:
    go-chi

### API Endpoints

* **POST /audio** — submit audio URL for recognition
* **GET /status** — check processing status
* **GET /result** — retrieve recognition result

### Audio Requirements

* **Supported Formats**:
    * WAV (recommended)
    * MP3
    * FLAC
    * OGG
    * AAC

### Installation & Setup

1. **Clone the repository**:
```bash
git clone https://github.com/itsvovovova/speechToText.git
cd speechToText
```

2. **Install dependencies**:
```bash
go mod tidy
```

3. **Configure environment**:
* Set up PostgreSQL connection
* Configure RabbitMQ credentials
* Configure Prometheus scraping

Example .env:
```
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=database
DB_SSL_MODE=disable
LOG_LEVEL=debug
LOG_FORMAT=text

SERVER_HOST=localhost
SERVER_PORT=8080
HOST_PORT=8080

REDIS_HOST=redis:6379

RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USER_PORT=15672
RABBITMQ_USER=GUEST
RABBITMQ_PASSWORD=GUEST
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

DEEPGRAM_API=0uhusec2a8c6e40c3f857621f010e1d6a
```
4. **Run the service**:
```bash
docker-compose up --build
```
### 📜 License

This project is licensed under the MIT License.
