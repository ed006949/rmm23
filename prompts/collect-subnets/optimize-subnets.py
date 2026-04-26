#!/usr/bin/env python3

from ipaddress import ip_network, collapse_addresses
import csv, os

infile='subnet.txt'
with open(infile) as f:
    nets=[]
    for line in f:
        s=line.strip()
        if not s: 
            continue
        try:
            nets.append(ip_network(s, strict=False))
        except Exception:
            pass

unique=sorted(set(nets), key=lambda n:(int(n.network_address), n.prefixlen))
collapsed=list(collapse_addresses(unique))

with open('subnet_optimized.txt','w') as f:
    for n in collapsed:
        f.write(str(n)+'\n')

print({'original':len(nets),'unique':len(unique),'optimized':len(collapsed),'reduction_percent':round((len(unique)-len(collapsed))*100/len(unique),2) if unique else 0})