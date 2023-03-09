package cidrman

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"sort"
)

const widthUInt32 = 32

const maxUInt32 = math.MaxUint32

// ipv4ToUInt32 converts an IPv4 address to an unsigned 32-bit integer.
func ipv4ToUInt32(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip)
}

// uint32ToIPV4 converts an unsigned 32-bit integer to an IPv4 address.
func uint32ToIPV4(addr uint32) net.IP {
	ip := make([]byte, net.IPv4len)
	binary.BigEndian.PutUint32(ip, addr)
	return ip
}

// The following functions are inspired by http://www.cs.colostate.edu/~somlo/iprange.c.

// setBit sets the specified bit in an address to 0 or 1.
func setBit(addr uint32, bit uint, val uint) uint32 {
	if bit < 0 {
		panic("negative bit index")
	}

	if val == 0 {
		return addr & ^(1 << (widthUInt32 - bit))
	} else if val == 1 {
		return addr | (1 << (widthUInt32 - bit))
	} else {
		panic("set bit is not 0 or 1")
	}
}

// hostmask4 returns the hostmask for the specified prefix.
func hostmask4(prefix uint) uint32 {
	return (1 << (widthUInt32 - prefix)) - 1
}

// netmask4 returns the netmask for the specified prefix.
func netmask4(prefix uint) uint32 {
	if prefix == 0 {
		return 0
	}
	return maxUInt32 ^ hostmask4(prefix)
}

// broadcast4 returns the broadcast address for the given address and prefix.
func broadcast4(addr uint32, prefix uint) uint32 {
	return addr | hostmask4(prefix)
}

// network4 returns the network address for the given address and prefix.
func network4(addr uint32, prefix uint) uint32 {
	return addr & netmask4(prefix)
}

// splitRange4 recursively computes the CIDR blocks to cover the range lo to hi.
func splitRange4(addr uint32, prefix uint, lo, hi uint32, cidrs *[]*net.IPNet) error {
	if prefix > widthUInt32 {
		return fmt.Errorf("Invalid mask size: %d", prefix)
	}

	bc := broadcast4(addr, prefix)
	if (lo < addr) || (hi > bc) {
		return fmt.Errorf("%d, %d out of range for network %d/%d, broadcast %d", lo, hi, addr, prefix, bc)
	}

	if (lo == addr) && (hi == bc) {
		cidr := net.IPNet{IP: uint32ToIPV4(addr), Mask: net.CIDRMask(int(prefix), 8*net.IPv4len)}
		*cidrs = append(*cidrs, &cidr)
		return nil
	}

	prefix++
	lowerHalf := addr
	upperHalf := setBit(addr, prefix, 1)
	if hi < upperHalf {
		return splitRange4(lowerHalf, prefix, lo, hi, cidrs)
	} else if lo >= upperHalf {
		return splitRange4(upperHalf, prefix, lo, hi, cidrs)
	} else {
		err := splitRange4(lowerHalf, prefix, lo, broadcast4(lowerHalf, prefix), cidrs)
		if err != nil {
			return err
		}
		return splitRange4(upperHalf, prefix, upperHalf, hi, cidrs)
	}
}

// IPv4 CIDR block.

type cidrBlock4 struct {
	first uint32
	last  uint32
}

type cidrBlock4s []*cidrBlock4

// newBlock4 returns a new IPv4 CIDR block.
func newBlock4(ip net.IP, mask net.IPMask) *cidrBlock4 {
	var block cidrBlock4

	block.first = ipv4ToUInt32(ip)
	prefix, _ := mask.Size()
	block.last = broadcast4(block.first, uint(prefix))

	return &block
}

// Sort interface.

func (c cidrBlock4s) Len() int {
	return len(c)
}

func (c cidrBlock4s) Less(i, j int) bool {
	lhs := c[i]
	rhs := c[j]

	// By last IP in the range.
	if lhs.last < rhs.last {
		return true
	} else if lhs.last > rhs.last {
		return false
	}

	// Then by first IP in the range.
	if lhs.first < rhs.first {
		return true
	} else if lhs.first > rhs.first {
		return false
	}

	return false
}

func (c cidrBlock4s) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// make a copy of list of CIDRs as most functions below will will destroy the original list(s)
// can be used like this if 'blocks' and 'removes' has to be re-used later
//    		result, err = remove4(copy4s(blocks), copy4s(removes))
func copy4s(inblocks cidrBlock4s) cidrBlock4s {
	outblocks := make(cidrBlock4s, len(inblocks))
	for i, ptr := range inblocks {
		if ptr == nil {
			// Skip to next for nil source pointer
			continue
		}
		// Create shallow copy of source cidrBlock
		newblock := *ptr

		outblocks[i] = &newblock
	}
	return outblocks
}

// merge4 accepts a list of IPv4 networks and merges them into the smallest possible list of IPNets.
// It merges adjacent subnets where possible, those contained within others and removes any duplicates.
func merge4(blocks cidrBlock4s) ([]*net.IPNet, error) {
	sort.Sort(blocks)

	// Coalesce overlapping blocks.
	for i := len(blocks) - 1; i > 0; i-- {
		if blocks[i].first <= blocks[i-1].last+1 {
			blocks[i-1].last = blocks[i].last
			if blocks[i].first < blocks[i-1].first {
				blocks[i-1].first = blocks[i].first
			}
			blocks[i] = nil
		}
	}

	var merged []*net.IPNet
	for _, block := range blocks {
		if block == nil {
			continue
		}

		if err := splitRange4(0, 0, block.first, block.last, &merged); err != nil {
			return nil, err
		}
	}

	return merged, nil
}

// remove4 accepts two lists of IPv4 networks and removes the second list from the first and return a new list of IPNets.
// The remove will return the smallest possible list of IPNets.
func remove4(blocks, removes cidrBlock4s) ([]*net.IPNet, error) {
	sort.Sort(blocks)
	sort.Sort(removes)

	i := 0
	j := 0
	for i < len(blocks) {
		if j >= len(removes) {
			// No more remove blocks to compare with
			break
		}
		if removes[j].last < blocks[i].first {
			// Remove-block entirely before network-block, use next remove-block
			j++
		} else if blocks[i].last < removes[j].first {
			// Network-block entirely before remove-block, keep block and continue to next
			i++
		} else if blocks[i].first >= removes[j].first && blocks[i].last <= removes[j].last {
			// Network-block inside remove-block, remove that network-block
			blocks[i] = nil
			i++
		// From here on we have some sort of overlap
		} else if blocks[i].first >= removes[j].first {
			// Network-block starts inside remove-block, adjust start of network-block
			blocks[i].first = removes[j].last + 1
			j++
		} else if blocks[i].last <= removes[j].last {
			// Network-block ends inside remove-block, adjust end of network-block
			blocks[i].last = removes[j].first - 1
			i++
		} else {
			// Remove-block inside network-block, will split network-block into two new blocks
			//
			// Make room for new network block
			blocks = append(blocks, nil)
			copy(blocks[i+1:], blocks[i:])
			blocks[i] = new(cidrBlock4)
			// update first half of the network-block (new)
			blocks[i].first = blocks[i+1].first
			blocks[i].last = removes[j].first - 1
			// Update second half of the network-block (old)
			blocks[i+1].first = removes[j].last + 1
			i++
			j++
		}
	}

	var merged []*net.IPNet
	for _, block := range blocks {
		if block == nil {
			continue
		}

		if err := splitRange4(0, 0, block.first, block.last, &merged); err != nil {
			return nil, err
		}
	}

	return merged, nil
}

// subset4 accepts two lists of IPv4 networks and return a new list of IPNets that exsists/overlaps in both lists.
// The subset4() will return the smallest possible list of IPNets.
func subset4(blocks, subsets cidrBlock4s) ([]*net.IPNet, error) {
	sort.Sort(blocks)
	sort.Sort(subsets)

	i := 0
	j := 0
	for i < len(blocks) {
		if j >= len(subsets) || blocks[i].last < subsets[j].first {
			// No more subset blocks to compare with (then no more network blocks to keep), or
			// network-block entirely before subset-block, remove block and continue with next
			//
			// Clear network-block and continue to next
			blocks[i] = nil
			i++
		} else if subsets[j].last < blocks[i].first {
			// Subset-block entirely before network-block, use next subset-block
			j++
		} else if blocks[i].first >= subsets[j].first && blocks[i].last <= subsets[j].last {
			// Network-block inside subset-block, keep that network-block
			i++
		// From here on we have some sort of overlap
		} else if blocks[i].first >= subsets[j].first {
			// Network-block starts inside subset-block, adjust end of network-block
			blocks[i].last = subsets[j].last
			i++
			j++
		} else if blocks[i].last <= subsets[j].last {
			// Network-block ends inside subset-block, adjust start of network-block
			blocks[i].first = subsets[j].first
			i++
		} else {
			// Subset-block inside network-block, adjust both start and end of network-block
			blocks[i].first = subsets[j].first
			// Check if next subset-block exists and starts with network-block
			if j < len(subsets) - 1 && subsets[j+1].first < blocks[i].last {
				// Make room for new network block
				blocks = append(blocks, nil)
				copy(blocks[i+1:], blocks[i:])
				blocks[i] = new(cidrBlock4)
				// update first half of the network-block (new)
				blocks[i].first = subsets[j].first
				// Update second half of the network-block (old)
				blocks[i+1].first = subsets[j+1].first
			}
			// Finaly handle (first) network-block's last address
			blocks[i].last = subsets[j].last
			i++
			j++
		}
	}

	var overlaps []*net.IPNet
	for _, block := range blocks {
		if block == nil {
			continue
		}

		if err := splitRange4(0, 0, block.first, block.last, &overlaps); err != nil {
			return nil, err
		}
	}

	return overlaps, nil
}
