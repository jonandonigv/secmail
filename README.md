# Secmail

Secmail is an early-stage secure email system built in Go, designed to provide end-to-end encryption for email messages. It aims to demonstrate secure communication concepts, including asymmetric key exchange and symmetric body encryption, without relying on external email protocols like SMTP/IMAP.

## Features

- **User Management**: Registration and authentication with password hashing and JWT tokens.
- **End-to-End Encryption**: Messages are encrypted using a combination of symmetric encryption (via age) for the body and asymmetric encryption (RSA) for session keys, ensuring only recipients can decrypt.
- **Multi-Recipient Support**: Send encrypted emails to multiple users.
- **Inbox Retrieval**: Authenticated users can view and decrypt their received messages.
- **Database Storage**: PostgreSQL for persistent storage of users and encrypted messages.

## Current Phase

**Phase 2 (Core Email)**: Completed. The system now supports basic send/receive functionality with full encryption. Users can register, login, send encrypted emails, and retrieve them from their inbox via REST API.

## Upcoming Phases

- **Phase 3 (Advanced Features)**: Add support for attachments, email conversations/threading, and full-text search across messages.
- **Phase 4 (Web UI & Polish)**: Implement a simple web interface for email composition and inbox viewing, along with error handling and logging improvements.
- **Phase 5 (Testing & Demo)**: Add unit tests for encryption logic, API tests, and a basic security audit. Prepare for demo deployment.

## Prerequisites

- Go 1.19+
- PostgreSQL database

## Setup

1. Clone the repository:
   ```
   git clone <repo-url>
   cd secmail
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Set up PostgreSQL:
   - Create a database named `secmail`.
   - Update the DSN in `main.go` if needed (default: `host=localhost user=postgres password=postgres dbname=secmail port=5432 sslmode=disable`).

4. Run the server:
   ```
   go run main.go
   ```

5. The API will be available at `http://localhost:8080`.

## API Endpoints

### Public
- `POST /register`: Register a new user (email, password).
- `POST /login`: Login and receive JWT token.

### Protected (requires Authorization header with Bearer token)
- `POST /emails/send`: Send an email (recipients array, subject, body).
- `GET /emails/inbox`: Retrieve decrypted inbox messages.

## Security Notes

- Private keys are stored encrypted in the database (demo purposes only; in production, use secure key management).
- This is a prototype for educational purposes and not suitable for real-world use without additional security audits and features like key rotation, TLS, and compliance.

## Contributing

This project is in active development. Contributions are welcome for the upcoming phases.

## License

MIT License