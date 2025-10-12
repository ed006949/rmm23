package mod_net

import (
	"encoding/binary"
	"math"
	"net/netip"
	"sync"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_reflect"
	"rmm23/src/mod_slices"
)

var Subnets = new(subnetsStruct)

type subnetsStruct struct {
	// mu provides thread-safe access to the subnets map.
	// It's locked before any read/write operation to prevent race conditions
	// when multiple goroutines attempt to generate or retrieve subnets simultaneously.
	mu sync.Mutex

	// subnets is a three-level nested map that caches generated subnet divisions.
	//
	// Structure: map[basePrefix]map[subnetPrefixLen][]subnet
	//
	// Level 1 - map[netip.Prefix]: The base network prefix (e.g., "192.168.0.0/16")
	//   - Key: The original network that will be subdivided
	//   - Value: A map of all subnet configurations for this base network
	//
	// Level 2 - map[int]: The desired subnet prefix length (e.g., 24 for /24 subnets)
	//   - Key: The prefix length of the subdivided subnets (must be >= base prefix length)
	//   - Value: A slice containing all possible subnets of this size
	//
	// Level 3 - []netip.Prefix: The actual list of generated subnets
	//   - Index: The subnet ID (0-based sequential identifier)
	//   - Value: The subnet prefix with its network address and mask
	//
	// Example:
	//   Base network: 10.0.0.0/16
	//   Desired subnet size: /24
	//   Result: subnets[10.0.0.0/16][24] = [10.0.0.0/24, 10.0.1.0/24, ..., 10.0.255.0/24]
	//            Total of 256 subnets (2^(24-16) = 2^8 = 256)
	//
	// This structure allows:
	//   - Lazy generation: Subnets are only computed when first requested
	//   - Efficient reuse: Once generated, subnets are cached for future lookups
	//   - Multiple configurations: Different base networks and subnet sizes can coexist
	//   - Fast retrieval: O(1) lookup time for any specific subnet by ID
	subnets map[netip.Prefix]map[int][]netip.Prefix

	index map[netip.Prefix]map[int]*subnetsIndex
}
type subnetsIndex struct {
	index map[netip.Prefix]int
	used  map[int]struct{}
}

func (r *subnetsStruct) PrefixIsFree(basePrefix netip.Prefix, subnetPrefixLen int, prefix netip.Prefix) (flag bool, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch err = r.validate(basePrefix, subnetPrefixLen); {
	case err != nil:
		return
	}

	var (
		index int
	)
	switch index, err = r.indexByPrefix(basePrefix, subnetPrefixLen, prefix); {
	case err != nil:
		return
	}

	return r.indexIsFree(basePrefix, subnetPrefixLen, index), nil
}

func (r *subnetsStruct) PrefixFree(basePrefix netip.Prefix, subnetPrefixLen int, prefix netip.Prefix) (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch err = r.validate(basePrefix, subnetPrefixLen); {
	case err != nil:
		return
	}

	var (
		index int
	)
	switch index, err = r.indexByPrefix(basePrefix, subnetPrefixLen, prefix); {
	case err != nil:
		return
	}

	r.indexFree(basePrefix, subnetPrefixLen, index)

	return
}
func (r *subnetsStruct) PrefixUse(basePrefix netip.Prefix, subnetPrefixLen int, prefix netip.Prefix) (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch err = r.validate(basePrefix, subnetPrefixLen); {
	case err != nil:
		return
	}

	var (
		index int
	)
	switch index, err = r.indexByPrefix(basePrefix, subnetPrefixLen, prefix); {
	case err != nil:
		return
	}

	r.indexUse(basePrefix, subnetPrefixLen, index)

	return
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
		r.index = make(map[netip.Prefix]map[int]*subnetsIndex)

		fallthrough
	case r.subnets[basePrefix] == nil:
		r.subnets[basePrefix] = make(map[int][]netip.Prefix)
		r.index[basePrefix] = make(map[int]*subnetsIndex)

		fallthrough
	case r.subnets[basePrefix][subnetPrefixLen] == nil:
		var (
			totalIDs = 1 << (subnetPrefixLen - basePrefix.Bits())
		)

		r.subnets[basePrefix][subnetPrefixLen] = make([]netip.Prefix, totalIDs)
		// r.index[basePrefix][subnetPrefixLen] = new(subnetsIndex)
		r.index[basePrefix][subnetPrefixLen] = &subnetsIndex{
			index: make(map[netip.Prefix]int, totalIDs),
			used:  make(map[int]struct{}, totalIDs),
		}

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

	r.index[basePrefix][subnetPrefixLen].index = mod_slices.Index(r.subnets[basePrefix][subnetPrefixLen])
}

func (r *subnetsStruct) generateIPv6(basePrefix netip.Prefix, subnetPrefixLen int, totalIDs int) {
	var (
		baseAddrBytes = basePrefix.Addr().AsSlice()[:] // 16 bytes
		// We add offset in the subnet bit position range [basePrefix.Bits(), subnetPrefixLen)
		shift = uint(MaxIPv6Bits - subnetPrefixLen)
		// Prepare a working 16-byte array for arithmetic
		base [16]byte
	)
	copy(base[:], baseAddrBytes)

	for currentID := 0; currentID < totalIDs; currentID++ {
		var (
			// offset = currentID << shift
			offset [16]byte
		)

		switch {
		case shift >= MaxIPv6Bits/2:
			binary.BigEndian.PutUint64(offset[0:8], uint64(currentID)<<(shift-MaxIPv6Bits/2))
		default:
			var (
				ohi, olo uint64 = 0, uint64(currentID) << shift
			)
			switch {
			case shift > 0:
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
			var (
				sum = uint16(base[i]) + uint16(offset[i]) + carry
			)

			cur[i] = byte(sum & math.MaxUint8)

			carry = sum >> mod_reflect.BitsPerByte

			switch {
			case i == 0:
				break
			}
		}

		r.subnets[basePrefix][subnetPrefixLen][currentID] = netip.PrefixFrom(netip.AddrFrom16(cur), subnetPrefixLen)
	}

	r.index[basePrefix][subnetPrefixLen].index = mod_slices.Index(r.subnets[basePrefix][subnetPrefixLen])
}

func (r *subnetsStruct) indexByPrefix(basePrefix netip.Prefix, subnetPrefixLen int, prefix netip.Prefix) (index int, err error) {
	var (
		ok bool
	)
	switch index, ok = r.index[basePrefix][subnetPrefixLen].index[prefix]; {
	case !ok:
		return 0, mod_errors.ENOTFOUND
	}

	return index, err
}

func (r *subnetsStruct) indexIsFree(basePrefix netip.Prefix, subnetPrefixLen int, index int) (flag bool) {
	_, flag = r.index[basePrefix][subnetPrefixLen].used[index]

	return
}

func (r *subnetsStruct) indexFree(basePrefix netip.Prefix, subnetPrefixLen int, index int) {
	delete(r.index[basePrefix][subnetPrefixLen].used, index)
}

func (r *subnetsStruct) indexUse(basePrefix netip.Prefix, subnetPrefixLen int, index int) {
	r.index[basePrefix][subnetPrefixLen].used[index] = struct{}{}
}
