#!/usr/bin/env bash
set -euo pipefail

CA_NAME="CA"
CA_KEY="${CA_NAME}.key.pem"
CA_CRT="${CA_NAME}.crt.pem"
CA_CRL="${CA_NAME}.crl.pem"
INDEX="index.txt"
SERIAL="serial"
CRLNUM="crlnumber"
KEY_BITS=3072
DAYS=825      # non-CA cert validity
CA_DAYS=3650  # CA cert validity
VERBOSE=0

log() { ((VERBOSE)) && echo "[*] $*"; }
err() {
	echo    "ERROR: $*" >&2
	exit                          1
}

ca_conf() {
	cat <<EOF
[ ca ]
default_ca = ca_local
[ ca_local ]
dir               = .
database          = $INDEX
certs             = .
new_certs_dir     = newcerts
certificate       = $CA_CRT
private_key       = $CA_KEY
serial            = $SERIAL
crlnumber         = $CRLNUM
default_md        = sha256
policy            = policy_any
default_days      = $DAYS
unique_subject    = no
default_crl_days  = 365
x509_extensions   = v3_end

[ policy_any ]
commonName        = supplied

[ v3_end ]
basicConstraints  = CA:false
keyUsage          = digitalSignature, keyEncipherment
extendedKeyUsage  = serverAuth, clientAuth
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer

[ v3_ca ]
basicConstraints       = critical,CA:true
keyUsage               = critical,keyCertSign,cRLSign
subjectKeyIdentifier   = hash
authorityKeyIdentifier = keyid:always,issuer
EOF
}

ensure_db() {
	mkdir -p newcerts
	[[ -f "$INDEX" ]] || : >"$INDEX"
	[[ -f "$SERIAL" ]] || echo 1000 >"$SERIAL"
	[[ -f "$CRLNUM" ]] || echo 1000 >"$CRLNUM"
}

ensure_ca() {
	ensure_db
	if [[ -f "$CA_KEY" && -f "$CA_CRT" ]]; then
		log   "CA already exists"
		return
	fi
	log "Creating self-signed CA..."
	openssl genrsa -out "$CA_KEY" "$KEY_BITS"
	openssl req -x509 -new -key "$CA_KEY" -days "$CA_DAYS" -sha256 \
		-subj   "/CN=$CA_NAME" -out "$CA_CRT" -extensions v3_ca -config <(ca_conf)
	chmod 600 "$CA_KEY"
	chmod 644 "$CA_CRT"
}

op() {
	local desc="$1"
	shift
	if ((VERBOSE)); then echo "[op] $desc"; fi
	if "$@" >/tmp/pki.out 2>&1; then
		((VERBOSE))   && echo "[res] ok"
	else
		((VERBOSE))   && {
			echo                   "[res] fail"
			cat                                      /tmp/pki.out
		}
		rm   -f /tmp/pki.out
		return   1
	fi
	rm -f /tmp/pki.out
}

verify_key_cert_match() {
	[[ "$(openssl pkey -in "$1" -pubout)" == "$(openssl x509 -in "$2" -pubkey -noout)" ]]
}

verify_cert_signed_by_ca() {
	openssl verify -CAfile "$CA_CRT" "$1" >/dev/null
}

cmd_verify() {
	local name="$1"
	if [[ "$name" == "$CA_NAME" ]]; then
		[[ -f "$CA_KEY" && -f "$CA_CRT"   ]] || err "CA key or cert missing"
		op   "CA key matches cert" verify_key_cert_match "$CA_KEY" "$CA_CRT" || err "CA key and cert mismatch"
		echo   "CA verification OK"
	else
		ensure_ca
		local   crt="${name}.crt.pem" key="${name}.key.pem"
		[[ -f "$crt" && -f "$key"   ]] || err "Missing $crt or $key"
		op   "key matches cert" verify_key_cert_match "$key" "$crt" || err "Mismatch key/cert"
		op   "cert signed by CA" verify_cert_signed_by_ca "$crt" || err "Not signed by CA"
		echo   "Verification OK for $name"
	fi
}

cmd_create() {
	local name="$1"
	if [[ "$name" == "$CA_NAME" ]]; then
		[[ -f "$CA_CRT"   ]] && err "CA already exists"
		op   "generate CA key" openssl genrsa -out "$CA_KEY" "$KEY_BITS"
		op   "generate CA cert" openssl req -x509 -new -key "$CA_KEY" \
			-days      "$CA_DAYS" -sha256 -subj "/CN=$CA_NAME" -out "$CA_CRT" \
			-extensions      v3_ca -config <(ca_conf)
		chmod   600 "$CA_KEY"
		chmod   644 "$CA_CRT"
		ensure_db
		echo   "CA created"
	else
		ensure_ca
		local   crt="${name}.crt.pem" key="${name}.key.pem"
		[[ -f "$crt" || -f "$key"   ]] && err "Target already exists"
		op   "generate key" openssl genrsa -out "$key" "$KEY_BITS"
		op   "generate CSR" openssl req -new -key "$key" -subj "/CN=$name" -out "${name}.csr.pem"
		op   "sign cert" openssl ca -batch -config <(ca_conf) \
			-in      "${name}.csr.pem" -out "$crt" -extensions v3_end
		rm   -f "${name}.csr.pem"
		chmod   600 "$key"
		chmod   644 "$crt"
		echo   "Created $key and $crt"
	fi
}

cmd_revoke() {
	local name="$1"
	[[ "$name" == "$CA_NAME" ]] && err "Cannot revoke CA"
	ensure_ca
	local crt="${name}.crt.pem"
	[[ -f "$crt" ]] || err "Missing $crt"
	verify_cert_signed_by_ca "$crt" || err "Not issued by this CA"
	log "Revoking $crt"
	op "revoke cert" openssl ca -batch -config <(ca_conf) -revoke "$crt"
	op "generate CRL" openssl ca -batch -config <(ca_conf) -gencrl -out "$CA_CRL"
	echo "Revoked $crt (CRL: $CA_CRL)"
}

cmd_delete() {
	local name="$1"
	[[ "$name" == "$CA_NAME" ]] && err "Cannot delete CA"
	cmd_revoke "$name" || true
	rm -f "${name}.key.pem" "${name}.crt.pem"
	echo "Deleted $name key and cert"
}

cmd_list() {
	ensure_ca
	echo "=== CA Database Contents ==="
	if [[ ! -s "$INDEX" ]]; then
		echo   "(no certificates issued)"
		return
	fi
	# Columns as per OpenSSL index.txt
	# V expiry revdate serial unknown /CN=alice
	# R expiry revdate serial unknown /CN=alice
	# E expiry ... etc
	printf "%-9s %-12s %-30s %-24s %-24s\n" "STATUS" "SERIAL" "CN" "EXPIRY/REVOCATION" "NOTES"
	awk -F'\t' '
    function fmt_ts(ts) {
        # OpenSSL gives YYMMDDHHMMSSZ; convert to ISO simple
        if (ts == "") return "-"
        if (match(ts, /../)) {
            y=substr(ts,1,2); m=substr(ts,3,2); d=substr(ts,5,2)
            H=substr(ts,7,2); M=substr(ts,9,2); S=substr(ts,11,2)
            return "20"y"-"m"-"d" "H":"M":"S"Z"
        } else return ts
    }
    {
        stat=$1; expiry=$2; revdate=$3; serial=$4; cn=$6;
        if (match(cn,/CN=([^\/]+)/,arr)) { cn=arr[1]; }
        note=""
        if (stat=="V") { status="VALID"; time=fmt_ts(expiry); }
        else if (stat=="E") { status="EXPIRED"; time=fmt_ts(expiry); }
        else if (stat=="R") { status="REVOKED"; time=fmt_ts(revdate); note="revoked on expiry="fmt_ts(expiry); }
        else { status=stat; time="-"; }
        printf "%-9s %-12s %-30s %-24s %-24s\n", status, serial, cn, time, note
    }
    ' "$INDEX"
}

usage() {
	cat <<EOF
Usage:
  $0 -verbose <command>
  $0 -verify <name>
  $0 -create <name>
  $0 -revoke <name>
  $0 -delete <name>
  $0 -list

Options:
  -verbose   Show each operation and its result.

Notes:
  Works only in current directory and maintains CA DB & CRL.
  -list shows status, serial, CN, expiry (for valid/expired), and revocation time.
EOF
}

main() {
	case "${1:-}" in
		-verbose)
			VERBOSE=1
			shift
			main                              "$@"
			;;
		-verify)
			shift
			cmd_verify                  "${1:-}"
			;;
		-create)
			shift
			cmd_create                  "${1:-}"
			;;
		-revoke)
			shift
			cmd_revoke                  "${1:-}"
			;;
		-delete)
			shift
			cmd_delete                  "${1:-}"
			;;
		-list)
			shift
			cmd_list
			;;
		*)   usage ;;
	esac
}

main "$@"
