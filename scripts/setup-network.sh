#!/usr/bin/env bash
set -euo pipefail

CHAIN_TYPE="${1:-devnet}"
HOME_ROOT="$(pwd)/data"
DENOM="ulbt"
BINARY="./lyfebloc-networkd"
KEYRING="test"

case "$CHAIN_TYPE" in
  devnet)
    CHAIN_ID="lyfebloc-devnet"
    VALIDATORS=2
    GAS="0.025${DENOM}"
    PARAMS_FILE=""
    ;;
  testnet)
    CHAIN_ID="lyfebloc-testnet"
    VALIDATORS=10
    GAS="0.020${DENOM}"
    PARAMS_FILE="config/testnet-params.json"
    ;;
  mainnet)
    CHAIN_ID="lyfebloc-mainnet"
    VALIDATORS=30
    GAS="0.015${DENOM}"
    PARAMS_FILE="config/mainnet-params.json"
    ;;
  *)
    echo "Usage: $0 {devnet|testnet|mainnet}"
    exit 1
    ;;
esac

HOME_DIR="${HOME_ROOT}/${CHAIN_TYPE}"
MONIKER="${CHAIN_TYPE}-node"

if [ ! -x "$BINARY" ]; then
  echo "⚙️  Building Lyfebloc binary..."
  GOCACHE=$(pwd)/.gocache go build -o lyfebloc-networkd ./cmd/lyfebloc-networkd
fi

echo "🧹 Resetting $CHAIN_TYPE at $HOME_DIR"
rm -rf "$HOME_DIR"
mkdir -p "$HOME_DIR"

trap 'echo "Shutting down setup"' EXIT

echo "🌱 Initializing ${CHAIN_ID}..."
$BINARY init "$MONIKER" --chain-id "$CHAIN_ID" --home "$HOME_DIR"

mkdir -p "$HOME_DIR/config/gentx"

# bootstrap primary key
if ! $BINARY keys show alice --home "$HOME_DIR" --keyring-backend "$KEYRING" >/dev/null 2>&1; then
  $BINARY keys add alice --home "$HOME_DIR" --keyring-backend "$KEYRING" >/dev/null 2>&1
fi
ALICE_ADDR=$($BINARY keys show alice -a --home "$HOME_DIR" --keyring-backend "$KEYRING")

# fund alice heavily
echo "💰 Funding alice..."
$BINARY genesis add-genesis-account "$ALICE_ADDR" 1000000000${DENOM} --home "$HOME_DIR"

# create validator accounts
for i in $(seq 1 $VALIDATORS); do
  VAL_KEY="validator${i}"
  if ! $BINARY keys show "$VAL_KEY" --home "$HOME_DIR" --keyring-backend "$KEYRING" >/dev/null 2>&1; then
    $BINARY keys add "$VAL_KEY" --home "$HOME_DIR" --keyring-backend "$KEYRING" >/dev/null 2>&1
  fi
  VAL_ADDR=$($BINARY keys show "$VAL_KEY" -a --home "$HOME_DIR" --keyring-backend "$KEYRING")
  $BINARY genesis add-genesis-account "$VAL_ADDR" 200000000${DENOM} --home "$HOME_DIR"
  GENTX_FILE="$HOME_DIR/config/gentx/gentx-$VAL_KEY.json"
  rm -f "$GENTX_FILE"
  $BINARY genesis gentx "$VAL_KEY" 100000000${DENOM} \
    --chain-id "$CHAIN_ID" \
    --home "$HOME_DIR" \
    --keyring-backend "$KEYRING" \
    --output-document "$GENTX_FILE"
done

echo "🧩 Collecting gentxs..."
$BINARY genesis collect-gentxs --home "$HOME_DIR"

# apply custom params if provided
if [ -n "$PARAMS_FILE" ] && [ -f "$PARAMS_FILE" ]; then
  python3 - "$PARAMS_FILE" "$HOME_DIR/config/genesis.json" <<'PY'
import json
from pathlib import Path
import sys
params_path = Path(sys.argv[1])
if not params_path.exists():
    sys.exit(0)
override = json.loads(params_path.read_text())
path = Path(sys.argv[2])
data = json.loads(path.read_text())

# simple deep merge for expected sections
for section, value in override.items():
    if isinstance(value, dict):
        data.setdefault(section, {})
        for k, v in value.items():
            if isinstance(v, dict):
                data[section].setdefault(k, {})
                data[section][k].update(v)
            else:
                data[section][k] = v
    else:
        data[section] = value

path.write_text(json.dumps(data, indent=2))
PY
fi

# gas price tweak (BSD-compatible sed)
sed -i.bak "s/^minimum-gas-prices *=.*/minimum-gas-prices = \"$GAS\"/" "$HOME_DIR/config/app.toml"
rm -f "$HOME_DIR/config/app.toml.bak"

python3 - "$HOME_DIR/config/app.toml" "$HOME_DIR/config/config.toml" <<'PY'
from pathlib import Path
import sys

app_path = Path(sys.argv[1])
config_path = Path(sys.argv[2])

def with_indent(src: str, replacement: str) -> str:
    stripped = src.lstrip()
    prefix = src[: len(src) - len(stripped)]
    return f"{prefix}{replacement}"

def update_app(path: Path):
    if not path.exists():
        return
    lines = path.read_text().splitlines()
    out = []
    section = None
    for line in lines:
        stripped = line.strip()
        if stripped.startswith("[") and stripped.endswith("]"):
            section = stripped[1:-1]
        if section == "api":
            if stripped.startswith("address ="):
                line = with_indent(line, 'address = "tcp://0.0.0.0:1317"')
            elif stripped.startswith("enable ="):
                line = with_indent(line, "enable = true")
            elif stripped.startswith("swagger ="):
                line = with_indent(line, "swagger = true")
            elif stripped.startswith("enabled-unsafe-cors ="):
                line = with_indent(line, "enabled-unsafe-cors = true")
        elif section == "grpc":
            if stripped.startswith("address ="):
                line = with_indent(line, 'address = "0.0.0.0:9090"')
            elif stripped.startswith("enable ="):
                line = with_indent(line, "enable = true")
        elif section == "grpc-web":
            if stripped.startswith("address ="):
                line = with_indent(line, 'address = "0.0.0.0:9091"')
            elif stripped.startswith("enable ="):
                line = with_indent(line, "enable = true")
        out.append(line)
    path.write_text("\n".join(out) + "\n")

def update_config(path: Path):
    if not path.exists():
        return
    lines = path.read_text().splitlines()
    out = []
    section = None
    for line in lines:
        stripped = line.strip()
        if stripped.startswith("[") and stripped.endswith("]"):
            section = stripped[1:-1]
        if section == "rpc":
            if stripped.startswith("laddr ="):
                line = with_indent(line, 'laddr = "tcp://0.0.0.0:26657"')
            elif stripped.startswith("pprof_laddr ="):
                line = with_indent(line, 'pprof_laddr = "0.0.0.0:6060"')
        elif section == "p2p":
            if stripped.startswith("laddr ="):
                line = with_indent(line, 'laddr = "tcp://0.0.0.0:26656"')
        elif section == "instrumentation":
            if stripped.startswith("prometheus ="):
                line = with_indent(line, "prometheus = true")
            elif stripped.startswith("prometheus_listen_addr ="):
                line = with_indent(line, 'prometheus_listen_addr = ":26660"')
        out.append(line)
    path.write_text("\n".join(out) + "\n")

update_app(app_path)
update_config(config_path)
PY

echo "✅ $CHAIN_TYPE genesis ready at $HOME_DIR/config/genesis.json"
echo "🔹 Start with:"
echo "$BINARY start --home $HOME_DIR --minimum-gas-prices $GAS"
