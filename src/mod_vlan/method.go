package mod_vlan

import (
	"encoding/binary"
	"math"
	"net/netip"

	"rmm23/src/mod_errors"
)

func (r *subnetMaps) Subnets(baseIPAddr netip.Addr, targetVLANs ...int) (outbound []netip.Prefix, err error) {
	outbound = make([]netip.Prefix, len(targetVLANs), len(targetVLANs))
	switch {
	case (*r)[baseIPAddr] == nil:
		return nil, mod_errors.ENOTFOUND
	case len(targetVLANs) == 0:
		return
	}

	for a, targetVLAN := range targetVLANs {
		outbound[a] = (*r)[baseIPAddr][targetVLAN]
	}

	return
}
func (r *subnetMaps) SubnetsAll(baseIPAddr netip.Addr) (outbound []netip.Prefix, err error) {
	switch {
	case (*r)[baseIPAddr] == nil:
		return nil, mod_errors.ENOTFOUND
	}

	outbound = (*r)[baseIPAddr][:]

	return
}

func (r *subnetMaps) GenerateSubnets(baseIPAddr netip.Addr, subnetPrefixLen int) (err error) {
	switch {
	case (*r)[baseIPAddr] != nil:
		return
	case baseIPAddr.Is4():
		return r.generateSubnetsIPv4(baseIPAddr, subnetPrefixLen)
	case baseIPAddr.Is6():
		return r.generateSubnetsIPv6(baseIPAddr, subnetPrefixLen)
	default:
		return mod_errors.EUnwilling
	}
}

func (r *subnetMaps) generateSubnetsIPv4(baseIPAddr netip.Addr, subnetPrefixLen int) (err error) {
	(*r)[baseIPAddr] = new(subnetMap)

	var (
		baseIPAsInt       = int(binary.BigEndian.Uint32(baseIPAddr.AsSlice()[:]))
		vlanAddressOffset = 1 << (MaxIPv4Bits - subnetPrefixLen)
	)
	switch {
	case baseIPAsInt+((MaxVLAN-1)*vlanAddressOffset) > math.MaxUint32:
		return mod_errors.EUnwilling
	}

	for currentVLAN := 0; currentVLAN <= MaxVLAN-1; currentVLAN++ {
		var (
			vlanSubnetIPBytes [4]byte
		)
		binary.BigEndian.PutUint32(vlanSubnetIPBytes[:], uint32(baseIPAsInt+currentVLAN*vlanAddressOffset))
		(*r)[baseIPAddr][currentVLAN] = netip.PrefixFrom(netip.AddrFrom4(vlanSubnetIPBytes), subnetPrefixLen)
	}

	return
}

func (r *subnetMaps) generateSubnetsIPv6(baseIPAddr netip.Addr, subnetPrefixLen int) (err error) {
	return mod_errors.EUnwilling
}
