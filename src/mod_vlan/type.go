package mod_vlan

import (
	"net/netip"
)

type subnetMap map[int]netip.Prefix
type subnetMaps struct {
	subnet map[netip.Addr]subnetMap
}
