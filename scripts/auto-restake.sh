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
