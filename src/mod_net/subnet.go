package mod_net

import (
	"encoding/binary"
	"net/netip"
	"sync"

	"rmm23/src/mod_errors"
)

var Subnets = new(subnetsStruct)

type subnetsStruct struct {
	mu      sync.Mutex
	subnets map[netip.Prefix]map[int][]netip.Prefix
}

func (r *subnetsStruct) Generate(basePrefix netip.Prefix, subnetPrefixLen int) (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.validate(basePrefix, subnetPrefixLen)
}

func (r *subnetsStruct) SubnetList(basePrefix netip.Prefix, subnetPrefixLen int, subnetIDs ...int) (outbound []netip.Prefix, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch err = r.validate(basePrefix, subnetPrefixLen); {
	case err != nil:
		return
	case len(r.subnets[basePrefix][subnetPrefixLen]) < len(subnetIDs):
		return nil, mod_errors.EUnwilling
	}

	outbound = make([]netip.Prefix, len(subnetIDs), len(subnetIDs))
	for a, subnetID := range subnetIDs {
		outbound[a] = r.subnets[basePrefix][subnetPrefixLen][subnetID]
	}

	return
}

func (r *subnetsStruct) Subnet(basePrefix netip.Prefix, subnetPrefixLen int, subnetID int) (outbound netip.Prefix, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch err = r.validate(basePrefix, subnetPrefixLen); {
	case err != nil:
		return
	}

	return r.subnets[basePrefix][subnetPrefixLen][subnetID], nil
}

func (r *subnetsStruct) Subnets(basePrefix netip.Prefix, subnetPrefixLen int) (outbound []netip.Prefix, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch err = r.validate(basePrefix, subnetPrefixLen); {
	case err != nil:
		return
	}

	return r.subnets[basePrefix][subnetPrefixLen], nil
}

func (r *subnetsStruct) validate(basePrefix netip.Prefix, subnetPrefixLen int) (err error) {
	switch {
	case subnetPrefixLen < basePrefix.Bits():
		return mod_errors.EUnwilling
	case r.subnets == nil:
		r.subnets = make(map[netip.Prefix]map[int][]netip.Prefix)

		fallthrough
	case r.subnets[basePrefix] == nil:
		r.subnets[basePrefix] = make(map[int][]netip.Prefix)

		fallthrough
	case r.subnets[basePrefix][subnetPrefixLen] == nil:
		var (
			totalIDs = 1 << (subnetPrefixLen - basePrefix.Bits())
		)

		r.subnets[basePrefix][subnetPrefixLen] = make([]netip.Prefix, totalIDs)

		return r.generate(basePrefix, subnetPrefixLen, totalIDs)
	default:
		return
	}
}

func (r *subnetsStruct) generate(basePrefix netip.Prefix, subnetPrefixLen int, totalIDs int) (err error) {
	switch {
	case basePrefix.Addr().Is4():
		return r.generateIPv4(basePrefix, subnetPrefixLen, totalIDs)
	case basePrefix.Addr().Is6():
		return r.generateIPv6(basePrefix, subnetPrefixLen, totalIDs)
	default:
		return mod_errors.EUnwilling
	}
}

func (r *subnetsStruct) generateIPv4(basePrefix netip.Prefix, subnetPrefixLen int, totalIDs int) (err error) {
	var (
		baseAddrAsInt  = int(binary.BigEndian.Uint32(basePrefix.Addr().AsSlice()[:]))
		baseAddrOffset = 1 << (MaxIPv4Bits - subnetPrefixLen)
	)
	for currentID, currentOffset := 0, baseAddrAsInt; currentID <= totalIDs-1; currentID, currentOffset = currentID+1, currentOffset+baseAddrOffset {
		var (
			currentAddrBytes [4]byte
		)
		binary.BigEndian.PutUint32(currentAddrBytes[:], uint32(currentOffset))
		r.subnets[basePrefix][subnetPrefixLen][currentID] = netip.PrefixFrom(netip.AddrFrom4(currentAddrBytes), subnetPrefixLen)
	}

	return
}

func (r *subnetsStruct) generateIPv6(basePrefix netip.Prefix, subnetPrefixLen int, totalIDs int) (err error) {
	return mod_errors.EUnwilling
}
