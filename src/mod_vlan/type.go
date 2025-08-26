package mod_vlan

import (
	"net/netip"
)

type subnetMap [MaxVLAN]netip.Prefix
type subnetMaps struct {
	subnet map[netip.Addr]*subnetMap
}
