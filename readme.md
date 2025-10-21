<p align="center">
  <img src="lb-net.png" alt="Lyfebloc Network Logo" width="300"/>
</p>

<h1 align="center" style="color:white;background-color:#0a0a0a;padding:16px;border-radius:8px;">
  🌐 Powering Universal Restaking Across Chains
</h1>

<p align="center">
  <b style="color:#2596be;">Cosmos EVM Layer 1 • Universal Restaking Layer • Cross-Chain Yield</b>
</p>

---

## 🪙 What Is Lyfebloc Network?

**Lyfebloc Network** is the world’s first **Universal Restaking Layer** built on the **Cosmos SDK + EVM** stack.  
It unifies staking economics across chains through a **dual-module restaking engine** that enables validators, delegators, and liquidity providers to compound and restake yield automatically — across multiple Cosmos and EVM networks.

> 💡 *“Restake once. Earn everywhere.”*

---

## ✨ Key Features

| Feature | Description |
|----------|--------------|
| **Dual Restaking Engine** | `BlocRestake` + `Restaking` modules combine auto-compounding with global liquidity restaking. |
| **Per-Validator Overrides** | Each validator can set a unique restake ratio for fine-tuned yield control. |
| **EVM Compatibility** | Fully integrated Solidity support for cross-chain restaking contracts. |
| **Auto-Restake Logic** | Validator rewards are automatically restaked every block or epoch. |
| **Cosmos IBC Native** | Designed to interconnect with Hub, Evmos, and Osmosis for multi-chain yield. |

---

## 🧠 Architecture Snapshot

├── app/ # Core app configuration (Cosmos SDK runtime)
├── x/blocrestake/ # BlocRestake module for IBC + staking integration
├── x/restaking/ # Chain-wide Restaking Engine
├── scripts/ # Local & Docker setup automation
├── docker/ # Docker + Prometheus + Grafana config
└── build/ # Compiled binaries

yaml
Copy code

---

## ⚙️ Build & Run Locally

### 1️⃣ Build Binary

```bash
make build
This compiles lyfebloc-networkd into the build/ folder.

2️⃣ Initialize Local Devnet
bash
Copy code
make init
Creates a local chain at data/devnet with funded validator accounts.

3️⃣ Launch Full Docker Stack
bash
Copy code
make docker-build
make docker-up
Access endpoints:

🌍 RPC: http://localhost:26657

🧠 API: http://localhost:1317

📊 Grafana Dashboard: http://localhost:3000 (user: admin / pass: admin)

📡 Prometheus: http://localhost:9091

4️⃣ View Logs
bash
Copy code
make docker-logs
🔄 Restaking Engine
Lyfebloc’s restaking engine supports chain-wide auto-compounding with validator-level overrides.

Global ratio (default: 25%) can be tuned via governance or CLI:

bash
Copy code
lyfebloc-networkd q restaking params
lyfebloc-networkd tx restaking set-ratio 0.35 --from alice --keyring-backend test
🧩 Developer Commands
Command	Description
make docker-down	Stop running containers
make docker-clean	Remove containers & volumes
make status	Display network & block status
make help	View available commands

🌍 Upcoming Testnet Beta (Q1 2026)
Lyfebloc Network Testnet Beta will introduce:

10 live validators with reward auto-restake

Cross-chain staking (Cosmos ↔ EVM)

Governance activation

Bridge integrations (IBC & ERC-20)



🪶 Vision
Lyfebloc is redefining the economics of validation.
By unifying staking, liquidity, and validator yield into one restaking layer, Lyfebloc creates a sustainable cross-chain economy where every chain can earn — not just stake.

Built for validators, liquidity providers, and Web3 builders who believe in the future of sovereign yield.

⚖️ License
Apache 2.0 — © 2025 Lyfeloop Inc.
All rights reserved.

