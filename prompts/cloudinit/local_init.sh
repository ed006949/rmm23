#!/usr/bin/env bash
set -Eeuo pipefail

AUTHORIZED_KEYS_MNG="${AUTHORIZED_KEYS_MNG:-/root/.ssh/authorized_keys}"
AUTHORIZED_KEYS_ROOT="${AUTHORIZED_KEYS_ROOT:-/root/.ssh/authorized_keys}"
SWAP_DISK="${SWAP_DISK:-}"
SWAP_SIZE_MIB="${SWAP_SIZE_MIB:-0}"
SSH_DROPIN_DIR="/etc/ssh/sshd_config.d"
SSH_ROOT_CONF="${SSH_DROPIN_DIR}/90-root-key-only.conf"
SSH_KEYS_ONLY_CONF="${SSH_DROPIN_DIR}/91-all-users-key-only.conf"
CRYPTTAB_FILE="/etc/crypttab"
FSTAB_FILE="/etc/fstab"
CRYPT_NAME="cryptswap"
MAPPER_PATH="/dev/mapper/${CRYPT_NAME}"
SWAP_LABEL="SWAP"

log() {
  printf '[+] %s\n' "$*"
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
  local src="$2"
  local home_dir ssh_dir target uid gid

  id "$user" >/dev/null 2>&1 || die "user not found: $user"
  [[ -f "$src" ]] || die "authorized_keys source not found: $src"

  home_dir="$(getent passwd "$user" | cut -d: -f6)"
  [[ -n "$home_dir" ]] || die "failed to resolve home for user: $user"

  ssh_dir="${home_dir}/.ssh"
  target="${ssh_dir}/authorized_keys"
  uid="$(id -u "$user")"
  gid="$(id -g "$user")"

  install -d -m 700 -o "$uid" -g "$gid" "$ssh_dir"
  install -m 600 -o "$uid" -g "$gid" "$src" "$target"
  log "Installed authorized_keys for ${user} from ${src}"
}

configure_ssh_key_only() {
  log "Configuring sshd drop-ins"
  install -d -m 755 "$SSH_DROPIN_DIR"

  cat > "$SSH_ROOT_CONF" <<CFG
PermitRootLogin prohibit-password
CFG

  cat > "$SSH_KEYS_ONLY_CONF" <<CFG
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
  local swap_name swap_type

  log "Disabling active swap devices"
  while read -r swap_name swap_type _; do
    [[ -n "${swap_name:-}" ]] || continue
    swapoff "$swap_name" || true
  done < <(swapon --show=NAME,TYPE --noheadings --raw || true)

  log "Removing swap entries from ${FSTAB_FILE}"
  cp -a "$FSTAB_FILE" "${FSTAB_FILE}.bak.$(date +%s)"
  grep -Ev '^[[:space:]]*($|#)' "$FSTAB_FILE" | cat >/dev/null || true
  sed -i -E '/^[[:space:]]*[^#].*[[:space:]]swap[[:space:]]/d' "$FSTAB_FILE"
  sed -i -E '/cryptswap|swapfile|[[:space:]]swap[[:space:]]/d' "$CRYPTTAB_FILE" 2>/dev/null || true

  log "Deleting swap files referenced in current swapon output"
  while read -r swap_name swap_type _; do
    [[ "$swap_type" == "file" ]] || continue
    rm -f -- "$swap_name"
  done < <(swapon --show=NAME,TYPE --noheadings --raw || true)

  log "Wiping swap signatures on swap partitions"
  while read -r path fstype; do
    [[ "$fstype" == "swap" ]] || continue
    swapoff "$path" || true
    wipefs -a "$path" || true
  done < <(lsblk -pnro PATH,FSTYPE | awk '$2=="swap"{print $1, $2}')

  if [[ -e "$MAPPER_PATH" ]]; then
    swapoff "$MAPPER_PATH" || true
    cryptsetup close "$CRYPT_NAME" || true
  fi
}

resolve_swap_disk() {
  [[ -n "$SWAP_DISK" ]] && return 0

  SWAP_DISK="$(findmnt -n -o SOURCE / | sed -E 's/p?[0-9]+$//' | sed -E 's/[0-9]+$//')"
  [[ -b "$SWAP_DISK" ]] || die "failed to detect root disk automatically; set SWAP_DISK=/dev/sdX or /dev/nvme0n1"
}

find_or_create_swap_partition() {
  local part_path part_num sectors mib start_mib end_mib free_line free_start free_end

  resolve_swap_disk
  log "Using disk: ${SWAP_DISK}"

  part_path="$(lsblk -pnro PATH,PARTLABEL "$SWAP_DISK" | awk '$2=="'"$SWAP_LABEL"'"{print $1; exit}')"
  if [[ -n "$part_path" ]]; then
    printf '%s\n' "$part_path"
    return 0
  fi

  [[ "$SWAP_SIZE_MIB" =~ ^[0-9]+$ ]] || die "SWAP_SIZE_MIB must be an integer MiB value"
  (( SWAP_SIZE_MIB > 0 )) || die "set SWAP_SIZE_MIB to the desired swap partition size in MiB"

  free_line="$(parted -m "$SWAP_DISK" unit MiB print free | awk -F: '$5=="free"{last=$0} END{print last}')"
  [[ -n "$free_line" ]] || die "no free space found on ${SWAP_DISK}"

  free_start="$(awk -F: '{gsub("MiB", "", $2); print int($2)}' <<<"$free_line")"
  free_end="$(awk -F: '{gsub("MiB", "", $3); print int($3)}' <<<"$free_line")"
  start_mib="$free_start"
  end_mib="$(( start_mib + SWAP_SIZE_MIB ))"
  (( end_mib <= free_end )) || die "not enough free space on ${SWAP_DISK}: need ${SWAP_SIZE_MIB} MiB"

  log "Creating GPT swap partition ${SWAP_LABEL} from ${start_mib}MiB to ${end_mib}MiB"
  parted -s "$SWAP_DISK" mkpart "$SWAP_LABEL" linux-swap "${start_mib}MiB" "${end_mib}MiB"
  partprobe "$SWAP_DISK"
  udevadm settle

  part_path="$(lsblk -pnro PATH,PARTLABEL "$SWAP_DISK" | awk '$2=="'"$SWAP_LABEL"'"{print $1; exit}')"
  [[ -n "$part_path" ]] || die "failed to discover created swap partition"

  part_num="$(lsblk -no PARTN "$part_path")"
  [[ -n "$part_num" ]] || die "failed to resolve partition number for ${part_path}"
  sgdisk --change-name="${part_num}:${SWAP_LABEL}" "$SWAP_DISK"

  printf '%s\n' "$part_path"
}

configure_encrypted_swap() {
  local swap_part partuuid

  swap_part="$(find_or_create_swap_partition)"
  [[ -b "$swap_part" ]] || die "swap partition not found: ${swap_part}"

  log "Preparing encrypted swap on ${swap_part}"
  swapoff "$swap_part" || true
  wipefs -a "$swap_part" || true

  partuuid="$(blkid -s PARTUUID -o value "$swap_part")"
  [[ -n "$partuuid" ]] || die "failed to read PARTUUID for ${swap_part}"

  grep -vE "^[[:space:]]*${CRYPT_NAME}[[:space:]]" "$CRYPTTAB_FILE" 2>/dev/null > "${CRYPTTAB_FILE}.tmp" || true
  mv "${CRYPTTAB_FILE}.tmp" "$CRYPTTAB_FILE" 2>/dev/null || true
  printf '%s\n' "${CRYPT_NAME} PARTUUID=${partuuid} /dev/urandom swap,cipher=aes-xts-plain64,size=256,discard" >> "$CRYPTTAB_FILE"

  cryptsetup open --type plain --key-file /dev/urandom "$swap_part" "$CRYPT_NAME"
  mkswap -f "$MAPPER_PATH"
  swapon "$MAPPER_PATH"

  cp -a "$FSTAB_FILE" "${FSTAB_FILE}.bak.swap.$(date +%s)"
  sed -i -E '/^[[:space:]]*[^#].*[[:space:]]swap[[:space:]]/d' "$FSTAB_FILE"
  printf '%s\n' "${MAPPER_PATH} none swap sw 0 0" >> "$FSTAB_FILE"

  log "Encrypted swap configured on ${swap_part} via ${MAPPER_PATH}"
}

main() {
  require_root

  require_cmd apt-get
  require_cmd swapon
  require_cmd lsblk
  require_cmd sshd

  # do NOT require cryptsetup/parted/sgdisk/etc before installing prerequisites
  install_system_updates
  install_prerequisites

  # now we can safely insist on them
  require_cmd cryptsetup
  require_cmd parted
  require_cmd sgdisk
  require_cmd findmnt
  require_cmd blkid

  ensure_authorized_keys mng "$AUTHORIZED_KEYS_MNG"
  ensure_authorized_keys root "$AUTHORIZED_KEYS_ROOT"
  configure_ssh_key_only
  list_current_swaps
  remove_old_swaps
  configure_encrypted_swap
  list_current_swaps
  log "Done"
}

main "$@"
