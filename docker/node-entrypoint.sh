#!/usr/bin/env bash
set -euo pipefail

HOME_DIR="/home/lyfebloc/.lyfebloc-network"
SEED_DIR="/seed"
NETWORK="${NETWORK:-devnet}"
MONIKER="${MONIKER:-${NETWORK}-validator}"
PERSISTENT_PEERS="${PERSISTENT_PEERS:-}"

init_from_seed() {
  if [ ! -d "$HOME_DIR/config" ]; then
    mkdir -p "$HOME_DIR"
    if [ -d "$SEED_DIR/config" ]; then
      cp -R "$SEED_DIR/config" "$HOME_DIR/"
    fi
    if [ -d "$SEED_DIR/data" ]; then
      cp -R "$SEED_DIR/data" "$HOME_DIR/"
    fi
    if [ -d "$SEED_DIR/keyring-test" ]; then
      cp -R "$SEED_DIR/keyring-test" "$HOME_DIR/"
    fi
  fi
}

configure_node() {
  CONFIG_TOML="$HOME_DIR/config/config.toml"
  APP_TOML="$HOME_DIR/config/app.toml"

  if [ -f "$CONFIG_TOML" ]; then
    sed -i "s/^moniker = \".*\"/moniker = \"${MONIKER}\"/" "$CONFIG_TOML"
    sed -i "s/^laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/0.0.0.0:26657\"/" "$CONFIG_TOML"
    sed -i "s/^pprof_laddr = \"localhost:6060\"/pprof_laddr = \"0.0.0.0:6060\"/" "$CONFIG_TOML"
    sed -i "s/^prometheus = false/prometheus = true/" "$CONFIG_TOML"
    sed -i "s/^prometheus_listen_addr = \":26660\"/prometheus_listen_addr = \":26660\"/" "$CONFIG_TOML"
    if [ -n "$PERSISTENT_PEERS" ]; then
      sed -i "s/^persistent_peers = \".*\"/persistent_peers = \"${PERSISTENT_PEERS}\"/" "$CONFIG_TOML"
    fi
  fi

  if [ -f "$APP_TOML" ]; then
    sed -i "s/^address = \"tcp:\/\/127.0.0.1:1317\"/address = \"tcp:\/\/0.0.0.0:1317\"/" "$APP_TOML"
    sed -i "s/^address = \"0.0.0.0:9090\"/address = \"0.0.0.0:9090\"/" "$APP_TOML" || true
    sed -i "s/^address = \"localhost:9090\"/address = \"0.0.0.0:9090\"/" "$APP_TOML"
    sed -i "s/^address = \"0.0.0.0:9091\"/address = \"0.0.0.0:9091\"/" "$APP_TOML" || true
    sed -i "s/^address = \"localhost:9091\"/address = \"0.0.0.0:9091\"/" "$APP_TOML"
    sed -i "s/^enable = false/enable = true/" "$APP_TOML"
    sed -i "s/^enabled-unsafe-cors = false/enabled-unsafe-cors = true/" "$APP_TOML"
    sed -i "s/^minimum-gas-prices *=.*/minimum-gas-prices = \"0.025ulbt\"/" "$APP_TOML"
  fi
}

init_from_seed
configure_node

exec /usr/bin/lyfebloc-networkd "$@"
