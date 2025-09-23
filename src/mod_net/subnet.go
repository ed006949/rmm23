package mod_net

import (
	"encoding/binary"
	"math"
	"net/netip"
	"sync"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_slices"
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
	}

	outbound = make([]netip.Prefix, len(subnetIDs), len(subnetIDs))
	for a, subnetID := range subnetIDs {
		switch {
		case !mod_slices.HasIndex(r.subnets[basePrefix][subnetPrefixLen], subnetID):
			return nil, mod_errors.EUnwilling
		}

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
	}

	return
}

func (r *subnetsStruct) generate(basePrefix netip.Prefix, subnetPrefixLen int, totalIDs int) (err error) {
	switch {
	case basePrefix.Addr().Is4():
		r.generateIPv4(basePrefix, subnetPrefixLen, totalIDs)
	case basePrefix.Addr().Is6():
		r.generateIPv6(basePrefix, subnetPrefixLen, totalIDs)
	default:
		return mod_errors.EUnwilling
	}

	return
}

func (r *subnetsStruct) generateIPv4(basePrefix netip.Prefix, subnetPrefixLen int, totalIDs int) {
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
}

func (r *subnetsStruct) generateIPv6(basePrefix netip.Prefix, subnetPrefixLen int, totalIDs int) {
	var (
		baseAddrBytes = basePrefix.Addr().AsSlice()[:] // 16 bytes
		// We add offset in the subnet bit position range [basePrefix.Bits(), subnetPrefixLen)
		shift = uint(MaxIPv6Bits - subnetPrefixLen)
	)
	// Prepare a working 16-byte array for arithmetic
	var base [16]byte
	copy(base[:], baseAddrBytes)

	for currentID := 0; currentID < totalIDs; currentID++ {
		// offset = currentID << shift
		var offset [16]byte

		if shift >= MaxIPv6Bits/2 {
			// Only high 64 bits affected
			hi := uint64(currentID) << (shift - MaxIPv6Bits/2)
			binary.BigEndian.PutUint64(offset[0:8], hi)
		} else {
			// Both high and/or low 64 bits can be affected
			var ohi, olo uint64

			olo = uint64(currentID) << shift
			// carry to high if shift caused bits to overflow from low into high
			if shift > 0 {
				ohi = uint64(currentID) >> (MaxIPv6Bits/2 - shift)
			}

			binary.BigEndian.PutUint64(offset[0:8], ohi)
			binary.BigEndian.PutUint64(offset[8:16], olo)
		}

		// current = base + offset (big-endian 128-bit add)
		var (
			cur   [16]byte
			carry uint16
		)

		for i := 15; i >= 0; i-- {
			sum := uint16(base[i]) + uint16(offset[i]) + carry
			cur[i] = byte(sum & math.MaxUint8)
			// use named constant instead of magic number
			const byteBits = 8

			carry = sum >> byteBits

			if i == 0 {
				break
			}
		}

		addr := netip.AddrFrom16(cur)
		r.subnets[basePrefix][subnetPrefixLen][currentID] = netip.PrefixFrom(addr, subnetPrefixLen)
	}
}
