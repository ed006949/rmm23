package mod_net

import (
	"encoding/binary"
	"math"
	"net/netip"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_reflect"
)

type subnetMaps struct {
	subnet map[netip.Prefix]subnetMap
}

type subnetMap []netip.Prefix

func NewSubnets() (outbound *subnetMaps) { return new(subnetMaps) }

func (r *subnetMaps) SubnetList(basePrefix netip.Prefix, subnetIDs ...int) (outbound subnetMap, err error) {
	switch value, ok := r.subnet[basePrefix]; {
	case !ok:
		return nil, mod_errors.ENOTFOUND
	case len(subnetIDs) == 0:
		return
	default:
		mod_reflect.MakeSliceIfNil(&outbound, len(subnetIDs), len(subnetIDs))

		for _, subnetID := range subnetIDs {
			outbound[subnetID] = value[subnetID]
		}

		return
	}
}

func (r *subnetMaps) Subnet(basePrefix netip.Prefix, subnetID int) (outbound netip.Prefix, err error) {
	switch _, ok := r.subnet[basePrefix]; {
	case !ok:
		return outbound, mod_errors.ENOTFOUND
	default:
		return r.subnet[basePrefix][subnetID], nil
	}
}

func (r *subnetMaps) Subnets(basePrefix netip.Prefix) (outbound subnetMap, err error) {
	switch value, ok := r.subnet[basePrefix]; {
	case !ok:
		return nil, mod_errors.ENOTFOUND
	default:
		return value, nil
	}
}

func (r *subnetMaps) GenerateSubnets(basePrefix netip.Prefix, subnetPrefixLen int) (err error) {
	switch _, ok := r.subnet[basePrefix]; {
	case ok:
		return
	}

	mod_reflect.MakeMapIfNil(&r.subnet)

	switch {
	case basePrefix.Addr().Is4():
		return r.generateSubnetsIPv4(basePrefix, subnetPrefixLen)
	case basePrefix.Addr().Is6():
		return r.generateSubnetsIPv6(basePrefix, subnetPrefixLen)
	default:
		return mod_errors.EUnwilling
	}
}

func (r *subnetMaps) generateSubnetsIPv4(basePrefix netip.Prefix, subnetPrefixLen int) (err error) {
	var (
		totalIDs = 1 << (subnetPrefixLen - basePrefix.Bits())
	)

	r.subnet[basePrefix] = make(subnetMap, totalIDs, totalIDs)

	var (
		baseAddrAsInt  = int(binary.BigEndian.Uint32(basePrefix.Addr().AsSlice()[:]))
		baseAddrOffset = 1 << (MaxIPv4Bits - subnetPrefixLen)
	)
	switch {
	case baseAddrAsInt+((totalIDs-1)*baseAddrOffset) > math.MaxUint32:
		return mod_errors.EUnwilling
	}

	for currentID := 0; currentID <= totalIDs-1; currentID++ {
		var (
			currentAddrBytes [4]byte
		)
		binary.BigEndian.PutUint32(currentAddrBytes[:], uint32(baseAddrAsInt+currentID*baseAddrOffset))
		r.subnet[basePrefix][currentID] = netip.PrefixFrom(netip.AddrFrom4(currentAddrBytes), subnetPrefixLen)
	}

	return
}

func (r *subnetMaps) generateSubnetsIPv6(basePrefix netip.Prefix, subnetPrefixLen int) (err error) {
	return mod_errors.EUnwilling
}
