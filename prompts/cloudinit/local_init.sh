#!/usr/bin/env bash

set -Eeuo pipefail

SWAP_DISK="${SWAP_DISK:-/dev/sdb}"
SSH_DROPIN_DIR="/etc/ssh/sshd_config.d"
SSH_ROOT_CONF="${SSH_DROPIN_DIR}/90-root-key-only.conf"
SSH_KEYS_ONLY_CONF="${SSH_DROPIN_DIR}/91-all-users-key-only.conf"
CRYPTTAB_FILE="/etc/crypttab"
FSTAB_FILE="/etc/fstab"
CRYPT_NAME="cryptswap"
MAPPER_PATH="/dev/mapper/${CRYPT_NAME}"
SWAP_LABEL="SWAP"
SWAP_PARTITION=""

AUTHORIZED_KEYS_MNG_CONTENT='type data comment'

AUTHORIZED_KEYS_ROOT_CONTENT='type data comment'

log() {
  printf '[+] %s\n' "$*" >&2
}

warn() {
  printf '[!] %s\n' "$*" >&2
}

die() {
  printf '[x] %s\n' "$*" >&2
  exit 1
}

require_root() {
  [[ ${EUID} -eq 0 ]] || die "run as root"
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "missing command: $1"
}

install_system_updates() {
  log "Installing system updates"
  export DEBIAN_FRONTEND=noninteractive
  apt-get update
  apt-get -y dist-upgrade
}

install_prerequisites() {
  log "Installing prerequisites"
  export DEBIAN_FRONTEND=noninteractive
  apt-get install -y \
    ca-certificates \
    curl \
    cryptsetup \
    gdisk \
    openssh-server \
    parted \
    sudo \
    util-linux
}

ensure_authorized_keys() {
  local user="$1"
  local content="$2"
  local home_dir ssh_dir target uid gid

  id "$user" >/dev/null 2>&1 || die "user not found: $user"
  [[ -n "$content" ]] || die "authorized_keys content is empty for user: $user"

  home_dir="$(getent passwd "$user" | cut -d: -f6)"
  [[ -n "$home_dir" ]] || die "failed to resolve home for user: $user"

  ssh_dir="${home_dir}/.ssh"
  target="${ssh_dir}/authorized_keys"
  uid="$(id -u "$user")"
  gid="$(id -g "$user")"

  install -d -m 700 -o "$uid" -g "$gid" "$ssh_dir"
  printf '%s\n' "$content" > "$target"
  chown "$uid:$gid" "$target"
  chmod 600 "$target"
  log "Installed embedded authorized_keys for ${user}"
}

configure_ssh_key_only() {
  log "Configuring sshd drop-ins"
  install -d -m 755 "$SSH_DROPIN_DIR"

  cat > "$SSH_ROOT_CONF" <<'CFG'
PermitRootLogin prohibit-password
CFG

  cat > "$SSH_KEYS_ONLY_CONF" <<'CFG'
PubkeyAuthentication yes
PasswordAuthentication no
KbdInteractiveAuthentication no
ChallengeResponseAuthentication no
AuthenticationMethods publickey
CFG

  chmod 644 "$SSH_ROOT_CONF" "$SSH_KEYS_ONLY_CONF"

  if sshd -t; then
    systemctl reload ssh || systemctl reload sshd || true
    log "sshd configuration validated and reloaded"
  else
    die "sshd configuration test failed"
  fi
}

list_current_swaps() {
  log "Active swap devices"
  swapon --show --noheadings --raw || true
  log "Swap-capable block devices"
  lsblk -o NAME,PATH,TYPE,FSTYPE,LABEL,PARTLABEL,UUID,MOUNTPOINTS | grep -E 'swap|NAME' || true
}

remove_old_swaps() {
  log "Disabling active swap devices"
  while read -r swap_name swap_type _; do
    [[ -n "${swap_name:-}" ]] || continue
    swapoff "$swap_name" || true
  done < <(swapon --show=NAME,TYPE --noheadings --raw || true)

  if [[ -f "$FSTAB_FILE" ]]; then
    log "Removing swap entries from ${FSTAB_FILE}"
    cp -a "$FSTAB_FILE" "${FSTAB_FILE}.bak.$(date +%s)"
    sed -i -E '/^[[:space:]]*[^#].*[[:space:]]swap[[:space:]]/d' "$FSTAB_FILE"
  fi

  if [[ -f "$CRYPTTAB_FILE" ]]; then
    log "Removing swap entries from ${CRYPTTAB_FILE}"
    cp -a "$CRYPTTAB_FILE" "${CRYPTTAB_FILE}.bak.$(date +%s)"
    sed -i -E '/cryptswap|swapfile|[[:space:]]swap([,[:space:]]|$)/d' "$CRYPTTAB_FILE"
  fi

  log "Deleting common swap files"
  rm -f /swap.img /swapfile

  log "Deleting active swap files from swapon output"
  while read -r swap_name swap_type _; do
    [[ "$swap_type" == "file" ]] || continue
    rm -f -- "$swap_name"
  done < <(swapon --show=NAME,TYPE --noheadings --raw || true)

  log "Wiping old swap signatures"
  while read -r path fstype; do
    [[ -n "${path:-}" ]] || continue
    swapoff "$path" || true
    wipefs -a "$path" || true
  done < <(lsblk -pnro PATH,FSTYPE | awk '$2=="swap"{print $1, $2}')

  if [[ -e "$MAPPER_PATH" ]]; then
    swapoff "$MAPPER_PATH" || true
    cryptsetup close "$CRYPT_NAME" || true
  fi
}

find_or_create_swap_partition() {
  local part_path part_num

  [[ -b "$SWAP_DISK" ]] || die "swap disk not found: ${SWAP_DISK}"
  log "Using disk: ${SWAP_DISK}"

  part_path="$(lsblk -pnro PATH,PARTLABEL "$SWAP_DISK" | awk '$2=="'"$SWAP_LABEL"'"{print $1; exit}')"
  if [[ -n "$part_path" && -b "$part_path" ]]; then
    SWAP_PARTITION="$part_path"
    log "Found existing swap partition: ${SWAP_PARTITION}"
    return 0
  fi

  log "Creating GPT partition table on ${SWAP_DISK}"
  parted -s "$SWAP_DISK" mklabel gpt

  log "Creating partition ${SWAP_LABEL} on ${SWAP_DISK}1"
  parted -s "$SWAP_DISK" mkpart primary linux-swap 1MiB 100%
  parted -s "$SWAP_DISK" name 1 "$SWAP_LABEL"
  partprobe "$SWAP_DISK"
  udevadm settle

  part_path="$(lsblk -pnro PATH,PARTLABEL "$SWAP_DISK" | awk '$2=="'"$SWAP_LABEL"'"{print $1; exit}')"
  [[ -n "$part_path" ]] || part_path="${SWAP_DISK}1"
  [[ -b "$part_path" ]] || die "failed to discover created swap partition"

  part_num="$(lsblk -no PARTN "$part_path")"
  [[ -n "$part_num" ]] || die "failed to resolve partition number for ${part_path}"
  sgdisk --change-name="${part_num}:${SWAP_LABEL}" "$SWAP_DISK" >/dev/null 2>&1 || true

  SWAP_PARTITION="$part_path"
  log "Created swap partition: ${SWAP_PARTITION}"
}

configure_encrypted_swap() {
  local partuuid crypttab_tmp

  find_or_create_swap_partition
  [[ -n "$SWAP_PARTITION" ]] || die "swap partition variable is empty"
  [[ -b "$SWAP_PARTITION" ]] || die "swap partition not found: ${SWAP_PARTITION}"

  log "Preparing encrypted swap on ${SWAP_PARTITION}"
  swapoff "$SWAP_PARTITION" || true
  wipefs -a "$SWAP_PARTITION" || true

  
#  partuuid="$(blkid -s PARTUUID -o value "$SWAP_PARTITION")"
#  [[ -n "$partuuid" ]] || die "failed to read PARTUUID for ${SWAP_PARTITION}"

  crypttab_tmp="$(mktemp)"
  if [[ -f "$CRYPTTAB_FILE" ]]; then
    grep -vE "^[[:space:]]*${CRYPT_NAME}[[:space:]]" "$CRYPTTAB_FILE" > "$crypttab_tmp" || true
  fi
#  printf '%s\n' "${CRYPT_NAME} PARTUUID=${partuuid} /dev/urandom swap,cipher=aes-xts-plain64,size=256,discard" >> "$crypttab_tmp"
  printf '%s\n' "${CRYPT_NAME} PARTLABEL=${SWAP_LABEL} /dev/urandom swap,cipher=aes-xts-plain64,size=256,discard" >> "$crypttab_tmp"
  install -m 644 "$crypttab_tmp" "$CRYPTTAB_FILE"
  rm -f "$crypttab_tmp"

  if cryptsetup status "$CRYPT_NAME" >/dev/null 2>&1; then
    cryptsetup close "$CRYPT_NAME" || true
  fi

  cryptsetup open --type plain --key-file /dev/urandom "$SWAP_PARTITION" "$CRYPT_NAME"
  mkswap -f "$MAPPER_PATH"
  swapon "$MAPPER_PATH"

  cp -a "$FSTAB_FILE" "${FSTAB_FILE}.bak.swap.$(date +%s)"
  sed -i -E '/^[[:space:]]*[^#].*[[:space:]]swap[[:space:]]/d' "$FSTAB_FILE"
  printf '%s\n' "${MAPPER_PATH} none swap sw 0 0" >> "$FSTAB_FILE"

  log "Encrypted swap configured on ${SWAP_PARTITION} via ${MAPPER_PATH}"
}

main() {
  require_root
  require_cmd apt-get
  require_cmd swapon
  require_cmd lsblk

  install_system_updates
  install_prerequisites

  require_cmd sshd
  require_cmd cryptsetup
  require_cmd parted
  require_cmd sgdisk
  require_cmd blkid
  require_cmd partprobe
  require_cmd udevadm

  ensure_authorized_keys mng "$AUTHORIZED_KEYS_MNG_CONTENT"
  ensure_authorized_keys root "$AUTHORIZED_KEYS_ROOT_CONTENT"
  configure_ssh_key_only
  list_current_swaps
  remove_old_swaps
  configure_encrypted_swap
  list_current_swaps
  log "Done"
}

main "$@"
