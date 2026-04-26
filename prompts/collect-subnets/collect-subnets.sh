#!/usr/bin/env bash

set -euo pipefail

###############################################################################
# Files
###############################################################################
DOMAIN_FILE="domain.txt"
ASN_FILE="asn.txt"
LOCAL_SUBNET_FILE="local-subnet.txt"
OUT_FILE="subnet.txt"
ASN_OUT_FILE="asn-all.txt"
OPTIMIZER_SCRIPT="optimize-subnets.py"
OPTIMIZED_OUT_FILE="subnet_optimized.txt"

###############################################################################
# Temp files
###############################################################################
TMP_IPS="$(mktemp)"
TMP_ASNS_FROM_DOMAINS="$(mktemp)"
TMP_ALL_ASNS="$(mktemp)"
TMP_SUBNETS="$(mktemp)"

cleanup() {
  rm -f "$TMP_IPS" "$TMP_ASNS_FROM_DOMAINS" "$TMP_ALL_ASNS" "$TMP_SUBNETS"
}
trap cleanup EXIT

###############################################################################
# Helpers
###############################################################################
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

  if [[ "$subnet" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/[0-9]+$ ]]; then
    echo "$subnet"
    return 0
  fi

  return 1
}

run_optimizer() {
  if [[ ! -f "$OPTIMIZER_SCRIPT" ]]; then
    warn "optimizer script not found: $OPTIMIZER_SCRIPT"
    return 0
  fi

  if command -v python3 >/dev/null 2>&1; then
    log "running optimizer with python3: $OPTIMIZER_SCRIPT"
    python3 "$OPTIMIZER_SCRIPT"
  elif command -v python >/dev/null 2>&1; then
    log "running optimizer with python: $OPTIMIZER_SCRIPT"
    python "$OPTIMIZER_SCRIPT"
  else
    warn "python/python3 not found, skipping optimizer"
    return 0
  fi

  if [[ -f "$OPTIMIZED_OUT_FILE" ]]; then
    log "optimized subnet list written: $OPTIMIZED_OUT_FILE"
  else
    warn "optimizer finished but output not found: $OPTIMIZED_OUT_FILE"
  fi
}

###############################################################################
# Check commands
###############################################################################
for cmd in dig whois awk sed grep sort xargs tr; do
  need_cmd "$cmd"
done

###############################################################################
# DNS: domain -> IPv4 addresses
###############################################################################
resolve_domain_ipv4() {
  local domain="$1"

  dig +short A "$domain" 2>/dev/null \
    | grep -E '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$' \
    | sort -u || true
}

###############################################################################
# Whois: IP -> origin ASN
###############################################################################
ip_to_asn() {
  local ip="$1"

  whois -h whois.cymru.com -- "-v $ip" 2>/dev/null \
    | awk -F'|' '
        NR > 1 {
          gsub(/^[ \t]+|[ \t]+$/, "", $1)
          if ($1 ~ /^[0-9]+$/) {
            print "AS" $1
          }
        }
      ' \
    | sort -u || true
}

###############################################################################
# IRR/Whois: ASN -> IPv4 route objects
###############################################################################
asn_to_ipv4_routes() {
  local asn="$1"

  {
    whois -h whois.radb.net -- "-i origin $asn" 2>/dev/null || true
    whois -h whois.ripe.net -- "-rBGi origin $asn" 2>/dev/null || true
  } | awk '
        BEGIN { IGNORECASE=1 }
        /^route:[[:space:]]+/ {
          print $2
        }
      ' \
    | grep -E '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/[0-9]+$' \
    | sort -u || true
}

###############################################################################
# Validate inputs
###############################################################################
if [[ ! -f "$DOMAIN_FILE" ]]; then
  warn "domain file not found: $DOMAIN_FILE"
fi

if [[ ! -f "$ASN_FILE" ]]; then
  warn "asn file not found: $ASN_FILE"
fi

if [[ ! -f "$LOCAL_SUBNET_FILE" ]]; then
  warn "local subnet file not found: $LOCAL_SUBNET_FILE"
fi

###############################################################################
# Step 1: Read domains, resolve IPv4, discover ASNs from resolved IPs
###############################################################################
if [[ -f "$DOMAIN_FILE" ]]; then
  log "loading domains from $DOMAIN_FILE"

  while IFS= read -r raw_domain || [[ -n "$raw_domain" ]]; do
    domain="$(trim_line "$raw_domain")"
    [[ -z "$domain" ]] && continue

    log "resolving IPv4 for domain: $domain"

    resolved_any=0
    while IFS= read -r ip; do
      [[ -z "$ip" ]] && continue
      resolved_any=1

      echo "$ip" >> "$TMP_IPS"
      log "  found IP: $ip"

      while IFS= read -r asn; do
        [[ -z "$asn" ]] && continue
        echo "$asn" >> "$TMP_ASNS_FROM_DOMAINS"
        log "  mapped IP $ip -> $asn"
      done < <(ip_to_asn "$ip")
    done < <(resolve_domain_ipv4 "$domain")

    if [[ "$resolved_any" -eq 0 ]]; then
      warn "  no IPv4 records found for: $domain"
    fi
  done < "$DOMAIN_FILE"
fi

###############################################################################
# Step 2: Read ASNs from asn.txt and normalize them
###############################################################################
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

###############################################################################
# Step 3: Merge ASNs from domains + ASNs from file
###############################################################################
if [[ -s "$TMP_ASNS_FROM_DOMAINS" ]]; then
  log "merging ASNs discovered from domains"
  cat "$TMP_ASNS_FROM_DOMAINS" >> "$TMP_ALL_ASNS"
fi

sort -u "$TMP_ALL_ASNS" -o "$TMP_ALL_ASNS"

if [[ -s "$TMP_ALL_ASNS" ]]; then
  cp "$TMP_ALL_ASNS" "$ASN_OUT_FILE"
  log "wrote merged ASN list to $ASN_OUT_FILE"
else
  warn "no ASNs collected"
  : > "$ASN_OUT_FILE"
fi

###############################################################################
# Step 4: Query routes for each ASN
###############################################################################
if [[ -s "$TMP_ALL_ASNS" ]]; then
  log "querying IPv4 route objects for collected ASNs"

  while IFS= read -r asn || [[ -n "$asn" ]]; do
    [[ -z "$asn" ]] && continue
    log "querying routes for $asn"

    found_any=0
    while IFS= read -r subnet; do
      [[ -z "$subnet" ]] && continue
      found_any=1
      echo "$subnet" >> "$TMP_SUBNETS"
      log "  found subnet: $subnet"
    done < <(asn_to_ipv4_routes "$asn")

    if [[ "$found_any" -eq 0 ]]; then
      warn "  no IPv4 route objects found for $asn"
    fi
  done < "$TMP_ALL_ASNS"
fi

###############################################################################
# Step 5: Add local subnets from local-subnet.txt
###############################################################################
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

###############################################################################
# Step 6: Final output + optimizer
###############################################################################
sort -u "$TMP_SUBNETS" > "$OUT_FILE"

log "raw subnet list written: $OUT_FILE"
log "merged ASN list written: $ASN_OUT_FILE"

run_optimizer

log "done"
