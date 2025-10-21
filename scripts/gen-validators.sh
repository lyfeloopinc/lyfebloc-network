#!/usr/bin/env bash
set -euo pipefail

HOME_DIR="${1:-$(pwd)/data/testnet}"
COUNT="${2:-10}"
CHAIN_ID="${3:-lyfebloc-testnet}"
DENOM="${4:-ulbt}"
BINARY="./lyfebloc-networkd"
KEYRING="test"

for i in $(seq 1 "$COUNT"); do
  NAME="validator${i}"
  if ! $BINARY keys show "$NAME" --home "$HOME_DIR" --keyring-backend "$KEYRING" >/dev/null 2>&1; then
    $BINARY keys add "$NAME" --home "$HOME_DIR" --keyring-backend "$KEYRING" >/dev/null 2>&1
  fi
  ADDR=$($BINARY keys show "$NAME" -a --home "$HOME_DIR" --keyring-backend "$KEYRING")
  $BINARY genesis add-genesis-account "$ADDR" "200000000${DENOM}" --home "$HOME_DIR"
  $BINARY genesis gentx "$NAME" "100000000${DENOM}" \
    --chain-id "$CHAIN_ID" \
    --home "$HOME_DIR" \
    --keyring-backend "$KEYRING"
done

$BINARY genesis collect-gentxs --home "$HOME_DIR"

echo "✅ Generated $COUNT validator gentxs under $HOME_DIR"
