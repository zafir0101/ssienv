# ssienv

**ssienv** is a PoC of identity environment that implements the basic concepts of **Self-Sovereign Identity (SSI)**.

This tool allows you to generate and manage **DIDs** (Decentralized Identifiers), **Verifiable Credentials**, **DIDComm** connections, schemas, and **Proof Requests**, interacting with Cloud and Edge agents.

## Features

| Resource        | Description                                                              |
|-----------------|--------------------------------------------------------------------------|
| `did`           | Decentralized Identifier – unique, user-controlled digital address       |
| `connection`    | Peer-to-peer DIDComm connection                                          |
| `credential`    | Digital verifiable credential                                            |
| `proof-request` | Cryptographic query sent by a verifier to a holder                       |
| `schema`        | Digital blueprint that defines structure and attributes of a credential  |

## Installation

```bash
git clone https://github.com/zafir0101/ssienv.git
cd ssienv
go build -o ssienv .
```

Or install directly:
```bash
go install github.com/zafir0101/ssienv@latest
```
## Starting the Agents

Before using `ssienv`, you need to start the local Cloud Agents. Go into the infrastructure folder and bring the stack up:

```bash
cd ssienv/infra/st-multi
docker compose up
```
### Configuration

| Exposed Service             | Description                             |
| --------------------------- | --------------------------------------- |
| `host.docker.internal:8080` | Single-tenant Cloud Agent #1 (issuer)   |
| `host.docker.internal:8081` | Single-tenant Cloud Agent #2 (holder)   |
| `host.docker.internal:8082` | Single-tenant Cloud Agent #3 (verifier) |

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
