# dsus

**D**arn **S**imple **U**pdate **S**erver

A lightweight, secure update server written in Go. Upload signed binaries and serve them to clients with cryptographic verification.

## Features

- RSA signature verification for uploaded files
- Optional HTTP Basic Auth
- SHA256 versioning
- 200MB file upload limit
- Zero-config debug mode
- Simple REST API

## Quick Start

```bash
# Debug build
just build-debug

# Run
./dsus
```

Server starts on `http://localhost:8080`

## API

### Upload File

```bash
curl -X POST http://localhost:8080/upload \
  -F "executable=@your-app" \
  -F "signature=@signature.sig"
```

**Response:** `OK` on success, `422` on invalid signature

### Download Files

| Endpoint | Description |
|----------|-------------|
| `GET /latest` | Latest uploaded executable |
| `GET /signature` | Signature file |
| `GET /version` | SHA256 hash of latest (file) or server version (build) |

## Authentication

Set environment variables to enable Basic Auth:

```bash
export DSUS_USER=admin
export DSUS_PASS=secret
./dsus
```

## Build

```bash
# Debug build
just build-debug

# Release build (production)
just build-release

# Clean
just clean
```

## Configuration

### Debug Mode (default)

Files served from `./files`, certs from `./certs/`

### Production Mode

- Files: `/var/lib/dsus/files/`
- Certs: `/etc/dsus/certs/`

## Certificate Setup

Place in `./certs/` (debug) or `/etc/dsus/certs/` (production):

- `publickey.pub` - RSA public key for signature verification

## Generating Keys

```bash
# Private key
openssl genrsa -out private.pem 2048

# Public key
openssl rsa -in private.pem -pubout -out publickey.pub

# Sign a file
openssl dgst -sha256 -sign private.pem -out signature.sig your-app
```

## Tech Stack

- [Go](https://golang.org/) 1.25+
- [Fiber](https://gofiber.io/) v3 web framework

## License

MIT
