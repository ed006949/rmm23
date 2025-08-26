package mod_vlan

import (
	"net/netip"
)

type subnetMap [MaxVLAN]netip.Prefix
type subnets map[netip.Addr]*subnetMap
