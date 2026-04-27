#!/usr/bin/env python3

import sys, ipaddress
nets = []
for line in sys.stdin:
    s = line.strip()
    if not s:
        continue
    nets.append(ipaddress.ip_network(s, strict=True))
for n in ipaddress.collapse_addresses(sorted(set(nets), key=lambda x: (x.version, int(x.network_address), x.prefixlen))):
    print(n)
