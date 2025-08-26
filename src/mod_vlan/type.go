package mod_vlan

import (
	"net/netip"
)

type subnetMap [MaxVLAN]netip.Prefix
type subnetMaps map[netip.Addr]*subnetMap
