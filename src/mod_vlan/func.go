package mod_vlan

import (
	"encoding/binary"
	"math"
	"net/netip"

	"rmm23/src/mod_errors"
)

func NewSubnets() (outbound subnets) {
	return make(subnets)
}

func (r *subnets) GetSubnets(baseIPAddr netip.Addr, targetVLANs ...int) (outbound []netip.Prefix, err error) {
	outbound = make([]netip.Prefix, len(targetVLANs), len(targetVLANs))
	switch {
	case (*r)[baseIPAddr] == nil:
		return nil, mod_errors.ENOTFOUND
	case len(targetVLANs) == 0:
		return
	}

	for a, b := range targetVLANs {
		outbound[a] = (*r)[baseIPAddr][b]
	}

	return
}

func (r *subnets) GenerateSubnet(baseIPAddr netip.Addr, subnetPrefixLen int) (err error) {
	switch {
	case (*r)[baseIPAddr] != nil:
		return
	case !baseIPAddr.Is4():
		return mod_errors.EUnwilling
	}

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
			vlanSubnetIPInt   = baseIPAsInt + currentVLAN*vlanAddressOffset
			vlanSubnetIPBytes [4]byte
		)
		binary.BigEndian.PutUint32(vlanSubnetIPBytes[:], uint32(vlanSubnetIPInt))

		var (
			vlanSubnetPrefix = netip.PrefixFrom(netip.AddrFrom4(vlanSubnetIPBytes), subnetPrefixLen)
		)

		(*r)[baseIPAddr][currentVLAN] = vlanSubnetPrefix
	}

	return
}
