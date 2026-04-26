#!/usr/bin/env bash

set -euo pipefail

DOMAIN_FILE="domain.txt"
ASN_FILE="asn.txt"
LOCAL_SUBNET_FILE="local-subnet.txt"
OUT_FILE="subnet.txt"
ASN_OUT_FILE="asn-all.txt"

TMP_IPS="$(mktemp)"
TMP_ASNS_FROM_DOMAINS="$(mktemp)"
TMP_ALL_ASNS="$(mktemp)"
TMP_SUBNETS="$(mktemp)"

cleanup() {
  rm -f "$TMP_IPS" "$TMP_ASNS_FROM_DOMAINS" "$TMP_ALL_ASNS" "$TMP_SUBNETS"
}
trap cleanup EXIT

log() {
  printf '[*] %s\n' "$*" >&2
}

warn() {
  printf '[!] %s\n' "$*" >&2
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    printf '[!] missing command: %s\n' "$1" >&2
    exit 1
  }
}

trim_line() {
  local s="$1"
  s="${s%%#*}"
  echo "$s" | xargs
}

normalize_asn() {
  local asn="$1"
  asn="$(echo "$asn" | tr '[:lower:]' '[:upper:]' | xargs)"

  if [[ "$asn" =~ ^AS[0-9]+$ ]]; then
    echo "$asn"
    return 0
  fi

  if [[ "$asn" =~ ^[0-9]+$ ]]; then
    echo "AS$asn"
    return 0
  fi

  return 1
}

normalize_subnet() {
  local subnet="$1"
  subnet="$(echo "$subnet" | xargs)"
  [[ "$subnet" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/[0-9]+$ ]] || return 1
  echo "$subnet"
}

for cmd in dig whois awk sed grep sort xargs tr; do
  need_cmd "$cmd"
done

resolve_domain_ipv4() {
  local domain="$1"
  dig +short A "$domain" 2>/dev/null \
    | grep -E '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$' \
    | sort -u || true
}

ip_to_asn() {
  local ip="$1"
  whois -h whois.cymru.com -- "-v $ip" 2>/dev/null \
    | awk -F'|' '
        NR > 1 {
          gsub(/^[ \t]+|[ \t]+$/, "", $1)
          if ($1 ~ /^[0-9]+$/) print "AS" $1
        }
      ' \
    | sort -u || true
}

asn_to_ipv4_routes() {
  local asn="$1"
  {
    whois -h whois.radb.net -- "-i origin $asn" 2>/dev/null || true
    whois -h whois.ripe.net -- "-rBGi origin $asn" 2>/dev/null || true
  } | awk '
        BEGIN { IGNORECASE=1 }
        /^route:[[:space:]]+/ { print $2 }
      ' \
    | grep -E '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/[0-9]+$' \
    | sort -u || true
}

if [[ ! -f "$DOMAIN_FILE" ]]; then
  warn "domain file not found: $DOMAIN_FILE"
fi

if [[ ! -f "$ASN_FILE" ]]; then
  warn "asn file not found: $ASN_FILE"
fi

if [[ ! -f "$LOCAL_SUBNET_FILE" ]]; then
  warn "local subnet file not found: $LOCAL_SUBNET_FILE"
fi

if [[ -f "$DOMAIN_FILE" ]]; then
  log "loading domains from $DOMAIN_FILE"

  while IFS= read -r raw_domain || [[ -n "$raw_domain" ]]; do
    domain="$(trim_line "$raw_domain")"
    [[ -z "$domain" ]] && continue

    log "resolving IPv4 for domain: $domain"

    while IFS= read -r ip; do
      [[ -z "$ip" ]] && continue
      echo "$ip" >> "$TMP_IPS"
      log "  found IP: $ip"

      while IFS= read -r asn; do
        [[ -z "$asn" ]] && continue
        echo "$asn" >> "$TMP_ASNS_FROM_DOMAINS"
        log "  mapped IP $ip -> $asn"
      done < <(ip_to_asn "$ip")
    done < <(resolve_domain_ipv4 "$domain")
  done < "$DOMAIN_FILE"
fi

if [[ -f "$ASN_FILE" ]]; then
  log "loading ASNs from $ASN_FILE"

  while IFS= read -r raw_asn || [[ -n "$raw_asn" ]]; do
    line="$(trim_line "$raw_asn")"
    [[ -z "$line" ]] && continue

    if normalized="$(normalize_asn "$line")"; then
      echo "$normalized" >> "$TMP_ALL_ASNS"
      log "  loaded ASN: $normalized"
    else
      warn "  invalid ASN skipped: $line"
    fi
  done < "$ASN_FILE"
fi

if [[ -s "$TMP_ASNS_FROM_DOMAINS" ]]; then
  log "merging ASNs discovered from domains"
  cat "$TMP_ASNS_FROM_DOMAINS" >> "$TMP_ALL_ASNS"
fi

sort -u "$TMP_ALL_ASNS" -o "$TMP_ALL_ASNS"

if [[ -s "$TMP_ALL_ASNS" ]]; then
  cp "$TMP_ALL_ASNS" "$ASN_OUT_FILE"
  log "wrote merged ASN list to $ASN_OUT_FILE"

  while IFS= read -r asn || [[ -n "$asn" ]]; do
    [[ -z "$asn" ]] && continue
    log "querying routes for $asn"

    while IFS= read -r subnet; do
      [[ -z "$subnet" ]] && continue
      echo "$subnet" >> "$TMP_SUBNETS"
      log "  found subnet: $subnet"
    done < <(asn_to_ipv4_routes "$asn")
  done < "$TMP_ALL_ASNS"
else
  warn "no ASNs collected"
  : > "$ASN_OUT_FILE"
fi

if [[ -f "$LOCAL_SUBNET_FILE" ]]; then
  log "loading local subnets from $LOCAL_SUBNET_FILE"

  while IFS= read -r raw_subnet || [[ -n "$raw_subnet" ]]; do
    line="$(trim_line "$raw_subnet")"
    [[ -z "$line" ]] && continue

    if subnet="$(normalize_subnet "$line")"; then
      echo "$subnet" >> "$TMP_SUBNETS"
      log "  added local subnet: $subnet"
    else
      warn "  invalid subnet skipped: $line"
    fi
  done < "$LOCAL_SUBNET_FILE"
fi

sort -u "$TMP_SUBNETS" > "$OUT_FILE"

log "done"
log "merged ASN list: $ASN_OUT_FILE"
log "IPv4 subnet list: $OUT_FILE"
