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

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
