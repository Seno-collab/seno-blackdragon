# Seno Blackdragon

This project demonstrates a simple Go service following a modular Clean Architecture.
Domains are separated into packages with clear responsibility boundaries:

- **User**
- **Token**
- **Dragon**
- **Skill**
- **Wallet**

Each domain exposes models, repositories, services and HTTP handlers under `internal/<domain>`.
The entry point is in `cmd/server` and uses in-memory repositories.
