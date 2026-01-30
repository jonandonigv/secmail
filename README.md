# Secmail

Secmail is an early-stage secure email system built in Go, designed to provide end-to-end encryption for email messages. It aims to demonstrate secure communication concepts, including asymmetric key exchange and symmetric body encryption, without relying on external email protocols like SMTP/IMAP.

## Features

- **User Management**: Registration and authentication with password hashing and JWT tokens.
- **End-to-End Encryption**: Messages are encrypted using a combination of symmetric encryption (via age) for the body and asymmetric encryption (RSA) for session keys, ensuring only recipients can decrypt.
- **Multi-Recipient Support**: Send encrypted emails to multiple users (up to 10 recipients).
- **Inbox Retrieval**: Authenticated users can view and decrypt their received messages.
- **Sent Folder**: View all emails you've sent.
- **Delete/Archive**: Manage your emails with soft delete support.
- **Database Storage**: PostgreSQL for persistent storage of users and encrypted messages.

## Current Phase

**Phase 2 (Core Email)**: Completed. The system supports basic send/receive functionality with full end-to-end encryption. Recent updates include security hardening (environment variables for secrets, input validation/sanitization) and basic unit tests for crypto and auth functions. Users can register, login, send encrypted emails, and retrieve them from their inbox via REST API.

## Upcoming Phases

- **Phase 3 (Advanced Features)**: Add support for attachments, email conversations/threading, and full-text search across messages.
- **Phase 4 (Web UI & Polish)**: Implement a simple web interface for email composition and inbox viewing, along with error handling and logging improvements.
- **Phase 5 (Testing & Demo)**: Expand unit tests to cover all components, add integration tests, perform security audit, and prepare for demo deployment.

**Project Status**: Active development resumed. Building core features and improving documentation.

## Prerequisites

- Go 1.25+
- PostgreSQL 12+
- Environment variables configured

## Setup

### Option 1: Docker (Recommended)

The easiest way to run Secmail is using Docker Compose:

```bash
# Clone the repository
git clone <repo-url>
cd secmail

# Set JWT secret (optional, defaults to a placeholder)
export JWT_SECRET="your-super-secret-jwt-key-min-32-characters-long"

# Start the services
docker-compose up -d

# View logs
docker-compose logs -f app
```

The API will be available at `http://localhost:8080`

To stop:
```bash
docker-compose down
```

To stop and remove data:
```bash
docker-compose down -v
```

### Option 2: Manual Setup

#### 1. Clone the repository:
```bash
git clone <repo-url>
cd secmail
```

#### 2. Install dependencies:
```bash
go mod tidy
```

#### 3. Set up PostgreSQL:
Create a database named `secmail`:
```bash
createdb secmail
# Or using psql:
psql -U postgres -c "CREATE DATABASE secmail;"
```

#### 4. Set required environment variables:
```bash
export DATABASE_URL="host=localhost user=postgres password=postgres dbname=secmail port=5432 sslmode=disable"
export JWT_SECRET="your-super-secret-jwt-key-min-32-characters-long"
```

#### 5. Run the server:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## API Documentation

### Authentication

All protected endpoints require an `Authorization` header with a Bearer token:
```
Authorization: Bearer <your-jwt-token>
```

#### Register
```http
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "message": "User registered successfully"
}
```

#### Login
```http
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Email Operations

#### Send Email
```http
POST /emails/send
Content-Type: application/json
Authorization: Bearer <token>

{
  "recipients": [2, 3],
  "subject": "Hello",
  "body": "This is an encrypted message!"
}
```

**Response:**
```json
{
  "message": "Email sent successfully"
}
```

#### Get Inbox
```http
GET /emails/inbox
Authorization: Bearer <token>
```

**Response:**
```json
{
  "messages": [
    {
      "id": 1,
      "conversation_id": 0,
      "sender_id": 2,
      "subject": "Hello",
      "body": "This is an encrypted message!",
      "status": "sent",
      "sent_at": "2026-01-30T10:00:00Z"
    }
  ]
}
```

#### Get Sent Emails
```http
GET /emails/sent
Authorization: Bearer <token>
```

**Response:**
```json
{
  "messages": [
    {
      "id": 1,
      "recipients": [2, 3],
      "subject": "Hello",
      "body": "This is an encrypted message!",
      "status": "sent",
      "sent_at": "2026-01-30T10:00:00Z"
    }
  ]
}
```

#### Reply to Email
```http
POST /emails/reply
Content-Type: application/json
Authorization: Bearer <token>

{
  "message_id": 1,
  "body": "Thanks for your message!"
}
```

**Response:**
```json
{
  "message": "Reply sent successfully"
}
```

#### Get Conversation
```http
GET /emails/conversation/:id
Authorization: Bearer <token>
```

**Response:**
```json
{
  "messages": [
    {
      "id": 1,
      "conversation_id": 0,
      "sender_id": 2,
      "subject": "Hello",
      "body": "Original message",
      "status": "sent",
      "sent_at": "2026-01-30T10:00:00Z"
    },
    {
      "id": 2,
      "conversation_id": 1,
      "sender_id": 3,
      "subject": "Re: Hello",
      "body": "Thanks for your message!",
      "status": "sent",
      "sent_at": "2026-01-30T10:05:00Z"
    }
  ]
}
```

#### Delete Email
```http
DELETE /emails/:id
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "Email deleted successfully"
}
```

## Development

### Running Tests
```bash
# Set required environment variable
export JWT_SECRET=test-secret

# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/crypto
go test ./internal/auth
```

### Code Quality
```bash
# Format code
gofmt -w .

# Run linter
go vet ./...

# Build
go build
```

### Project Structure
```
secmail/
├── main.go                 # Application entry point
├── Dockerfile              # Docker image definition
├── docker-compose.yml      # Docker Compose configuration
├── .dockerignore           # Docker ignore rules
├── internal/
│   ├── auth/              # Authentication (JWT, bcrypt)
│   ├── crypto/            # Encryption utilities (RSA, age)
│   ├── database/          # Database connection
│   ├── email/             # Email send/receive logic
│   ├── handlers/          # HTTP request handlers
│   └── models/            # Database models
├── go.mod
├── go.sum
└── README.md
```

## Security Notes

- Private keys are stored encrypted in the database (demo purposes only; in production, use secure key management like HSM or client-side storage).
- Always use HTTPS in production.
- Never commit secrets to version control.
- This is a prototype for educational purposes and not suitable for real-world use without additional security audits and features like key rotation, TLS, and compliance.

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License
