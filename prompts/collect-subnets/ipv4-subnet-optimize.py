#!/usr/bin/env python3

import sys
from ipaddress import IPv4Address, ip_network, summarize_address_range

ranges = []

for line in sys.stdin:
    line = line.split('#', 1)[0].replace(',', ' ').strip()

    for tok in line.split():
        orig = tok
        if '/' not in tok:
            tok += '/32'

        try:
            net = ip_network(tok, strict=False)
        except ValueError as err:
            if '/' in orig:
                print(f"skipping {orig}: {err}", file=sys.stderr)
            continue

        if net.version != 4:
            continue

        ranges.append((int(net.network_address), int(net.broadcast_address)))

ranges.sort()

merged = []
for start, end in ranges:
    if not merged or start > merged[-1][1] + 1:
        merged.append([start, end])
    elif end > merged[-1][1]:
        merged[-1][1] = end

for start, end in merged:
    for net in summarize_address_range(IPv4Address(start), IPv4Address(end)):
        print(net)
