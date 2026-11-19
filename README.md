# skdesk

[![CI](https://github.com/domehahn/skdesk/actions/workflows/ci.yml/badge.svg)](https://github.com/domehahn/skdesk/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/domehahn/skdesk)](https://goreportcard.com/report/github.com/domehahn/skdesk)
[![Go Reference](https://pkg.go.dev/badge/github.com/domehahn/skdesk.svg)](https://pkg.go.dev/github.com/domehahn/skdesk)
[![Go Version](https://img.shields.io/github/go-mod/go-version/domehahn/skdesk)](https://github.com/domehahn/skdesk/blob/main/go.mod)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/domehahn/skdesk/blob/main/LICENSE)
[![Status](https://img.shields.io/badge/Status-Production--Ready-brightgreen.svg)](https://github.com/domehahn/skdesk)

`skdesk` (**Human-in-the-Loop Review & Attestation Portal**) is a governance dashboard and review management plane for evaluating `skgate` security escalations and signing cryptographically bound approval receipts.

## Features

- **Review Ticket Lifecycle**: Manages `PENDING`, `APPROVED`, and `REJECTED` tickets for skills requiring elevated permissions.
- **Cryptographic Ed25519 Sign-off**: Produces verifiable approval receipts ingestible by `skgate`.
- **Embedded Web Dashboard**: Self-contained HTTP server and responsive UI requiring zero external dependencies.

## Installation

```bash
go install github.com/domehahn/skdesk/cmd/skdesk@latest
```

## Quick Start

```bash
# Generate reviewer Ed25519 keypair
skdesk keygen

# Start review server
skdesk serve --port 8080
```
