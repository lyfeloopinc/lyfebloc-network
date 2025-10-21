#!/usr/bin/env bash
set -euo pipefail

CHAIN_ID="lyfebloc-devnet"
HOME_DIR="$(pwd)/data/devnet"
DENOM="ulbt"
MONIKER="lyfebloc-devnet"
KEY_NAME="alice"
KEYRING="test"
MIN_GAS="0.025${DENOM}"

cat <<'SCRIPT' > scripts/auto-restake.sh
#!/usr/bin/env bash
set -euo pipefail

HOME_DIR="$(pwd)/data/devnet"
CHAIN_ID="lyfebloc-devnet"
KEY_NAME="alice"
KEYRING="test"

VALOPER="$(./lyfebloc-networkd keys show $KEY_NAME --bech val --address --home "$HOME_DIR" --keyring-backend "$KEYRING")"

while true; do
  ./lyfebloc-networkd tx distribution withdraw-rewards "$VALOPER" \
    --chain-id "$CHAIN_ID" \
    --from "$KEY_NAME" \
    --keyring-backend "$KEYRING" \
    --yes \
    --home "$HOME_DIR" \
    --gas auto \
    --gas-adjustment 1.5 \
    --fees 2500ulbt \
    >/dev/null 2>&1 || true
  sleep 20
done
SCRIPT
chmod +x scripts/auto-restake.sh

echo "🚀 Building Lyfebloc binary..."
GOCACHE=$(pwd)/.gocache go build -o lyfebloc-networkd ./cmd/lyfebloc-networkd

echo "🧹 Resetting old devnet..."
rm -rf "$HOME_DIR"

mkdir -p "$HOME_DIR"

trap 'echo "Shutting down..."' EXIT

echo "🌱 Initializing new devnet..."
./lyfebloc-networkd init $MONIKER --chain-id $CHAIN_ID --home "$HOME_DIR"

echo "🔑 Creating validator key..."
./lyfebloc-networkd keys add $KEY_NAME --home "$HOME_DIR" --keyring-backend $KEYRING || true

ADDRESS=$(./lyfebloc-networkd keys show $KEY_NAME -a --home "$HOME_DIR" --keyring-backend $KEYRING)
VALOPER=$(./lyfebloc-networkd keys show $KEY_NAME --bech val --address --home "$HOME_DIR" --keyring-backend $KEYRING)

echo "💰 Funding account..."
./lyfebloc-networkd genesis add-genesis-account "$ADDRESS" 1000000000${DENOM} --home "$HOME_DIR"

echo "🧾 Creating gentx..."
./lyfebloc-networkd genesis gentx $KEY_NAME 100000000${DENOM} \
  --chain-id $CHAIN_ID \
  --home "$HOME_DIR" \
  --keyring-backend $KEYRING

echo "🧩 Collecting gentxs..."
./lyfebloc-networkd genesis collect-gentxs --home "$HOME_DIR"

echo "⚙️ Setting minimum gas prices..."
python3 - <<'PY'
from pathlib import Path
path = Path('data/devnet/config/app.toml')
text = path.read_text().splitlines()
with path.open('w') as f:
    for line in text:
        if line.strip().startswith('minimum-gas-prices ='):
            f.write('minimum-gas-prices = "0.025ulbt"\n')
        else:
            f.write(line + '\n')
PY

echo "🔥 Starting chain..."
./lyfebloc-networkd start \
  --home "$HOME_DIR" \
  --minimum-gas-prices "$MIN_GAS" \
  --rpc.laddr tcp://0.0.0.0:26657 \
  --grpc.address 0.0.0.0:9090 \
  | tee "$HOME_DIR/lyfebloc.log"
