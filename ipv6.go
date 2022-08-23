package cidrman

import (
	"fmt"
	"math/big"
	"net"
	"sort"
)

const widthUInt128 = 128

var maxUInt128 = big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(widthUInt128), nil), big.NewInt(1))

// ipv6ToUInt128 converts an IPv6 address to an unsigned 128-bit integer.
func ipv6ToUInt128(ip net.IP) *big.Int {
	return big.NewInt(0).SetBytes(ip)
}

// uint128ToIPV6 converts an unsigned 128-bit integer to an IPv6 address.
func uint128ToIPV6(addr *big.Int) net.IP {
	ip := make([]byte, net.IPv6len)
	ab := addr.Bytes()
	copy(ip[len(ip)-len(ab):], ab)
	return ip
}

// copyUInt128 copies an unsigned 128-bit integer.
func copyUInt128(x *big.Int) *big.Int {
	return big.NewInt(0).Set(x)
}

// hostmask6 returns the hostmask for the specified prefix.
func hostmask6(prefix uint) *big.Int {
	z := big.NewInt(0)

	z.Lsh(big.NewInt(1), widthUInt128-prefix)
	z.Sub(z, big.NewInt(1))

	return z
}

// netmask6 returns the netmask for the specified prefix.
func netmask6(prefix uint) *big.Int {
	z := big.NewInt(0)

	if prefix == 0 {
		return z
	}

	z.Xor(maxUInt128, hostmask6(prefix))

	return z
}

// broadcast6 returns the broadcast address for the given address and prefix.
func broadcast6(addr *big.Int, prefix uint) *big.Int {
	z := big.NewInt(0)

	z.Or(addr, hostmask6(prefix))

	return z
}

// network6 returns the network address for the given address and prefix.
func network6(addr *big.Int, prefix uint) *big.Int {
	z := big.NewInt(0)

	z.And(addr, netmask6(prefix))

	return z
}

// splitRange6 recursively computes the CIDR blocks to cover the range lo to hi.
func splitRange6(addr *big.Int, prefix uint, lo, hi *big.Int, cidrs *[]*net.IPNet) error {
	if prefix > widthUInt128 {
		return fmt.Errorf("Invalid mask size: %d", prefix)
	}

	bc := broadcast6(addr, prefix)
	//	fmt.Printf("%v/%v, %v-%v, %v\n", uint128ToIPV6(addr), prefix, uint128ToIPV6(lo), uint128ToIPV6(hi), uint128ToIPV6(bc))
	if (lo.Cmp(addr) < 0) || (hi.Cmp(bc) > 0) {
		return fmt.Errorf("%v, %v out of range for network %v/%d, broadcast %v", uint128ToIPV6(lo), uint128ToIPV6(hi), uint128ToIPV6(addr), prefix, uint128ToIPV6(bc))
	}

	if (lo.Cmp(addr) == 0) && (hi.Cmp(bc) == 0) {
		cidr := net.IPNet{IP: uint128ToIPV6(addr), Mask: net.CIDRMask(int(prefix), 8*net.IPv6len)}
		*cidrs = append(*cidrs, &cidr)
		return nil
	}

	prefix++
	lowerHalf := copyUInt128(addr)
	upperHalf := copyUInt128(addr)
	upperHalf.SetBit(upperHalf, int(widthUInt128-prefix), 1)
	if hi.Cmp(upperHalf) < 0 {
		return splitRange6(lowerHalf, prefix, lo, hi, cidrs)
	} else if lo.Cmp(upperHalf) >= 0 {
		return splitRange6(upperHalf, prefix, lo, hi, cidrs)
	} else {
		err := splitRange6(lowerHalf, prefix, lo, broadcast6(lowerHalf, prefix), cidrs)
		if err != nil {
			return err
		}
		return splitRange6(upperHalf, prefix, upperHalf, hi, cidrs)
	}
}

// IPv6 CIDR block.

type cidrBlock6 struct {
	first *big.Int
	last  *big.Int
}

type cidrBlock6s []*cidrBlock6

// newBlock6 returns a new IPv6 CIDR block.
func newBlock6(ip net.IP, mask net.IPMask) *cidrBlock6 {
	var block cidrBlock6

	block.first = ipv6ToUInt128(ip)
	prefix, _ := mask.Size()
	block.last = broadcast6(block.first, uint(prefix))

	return &block
}

// Sort interface.

func (c cidrBlock6s) Len() int {
	return len(c)
}

func (c cidrBlock6s) Less(i, j int) bool {
	lhs := c[i]
	rhs := c[j]

	// By last IP in the range.
	if lhs.last.Cmp(rhs.last) < 0 {
		return true
	} else if lhs.last.Cmp(rhs.last) > 0 {
		return false
	}

	// Then by first IP in the range.
	if lhs.first.Cmp(rhs.first) < 0 {
		return true
	} else if lhs.first.Cmp(rhs.first) > 0 {
		return false
	}

	return false
}

func (c cidrBlock6s) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// merge6 accepts a list of IPv6 networks and merges them into the smallest possible list of IPNets.
// It merges adjacent subnets where possible, those contained within others and removes any duplicates.
func merge6(blocks cidrBlock6s) ([]*net.IPNet, error) {
	sort.Sort(blocks)

	// Coalesce overlapping blocks.
	for i := len(blocks) - 1; i > 0; i-- {
		cmp := big.NewInt(1)
		cmp.Add(cmp, blocks[i-1].last)
		if blocks[i].first.Cmp(cmp) <= 0 {
			blocks[i-1].last = blocks[i].last
			if blocks[i].first.Cmp(blocks[i-1].first) < 0 {
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

		if err := splitRange6(big.NewInt(0), 0, block.first, block.last, &merged); err != nil {
			return nil, err
		}
	}
	return merged, nil
}

// remove6 accepts two lists of IPv6 networks and removes the second list from the first and return a new list of IPNets.
// The remove will return the smallest possible list of IPNets.
func remove6(blocks, removes cidrBlock6s) ([]*net.IPNet, error) {
	sort.Sort(blocks)
	sort.Sort(removes)

	i := 0
	j := 0
	for i < len(blocks) {
		if j >= len(removes) {
			// No more remove blocks to compare with
			break
		}
		if removes[j].last.Cmp(blocks[i].first) < 0 {
			// Remove-block entirely before network-block, use next remove-block
			j++
		} else if blocks[i].last.Cmp(removes[j].first) < 0 {
			// Network-block entirely before remove-block, keep block and continue to next
			i++
		} else if blocks[i].first.Cmp(removes[j].first) >= 0 && blocks[i].last.Cmp(removes[j].last) <= 0 {
			// Network-block inside remove-block, remove that network-block
			blocks[i] = nil
			i++
		// From here on we have some sort of overlap
		} else if blocks[i].first.Cmp(removes[j].first) >= 0 {
			// Network-block starts inside remove-block, adjust start of network-block
			tmp := big.NewInt(1)
			blocks[i].first = tmp.Add(tmp, removes[j].last)
			j++
		} else if blocks[i].last.Cmp(removes[j].last) <= 0 {
			// Network-block ends inside remove-block, adjust end of network-block
			tmp := big.NewInt(0)
			blocks[i].last = tmp.Sub(tmp.Add(tmp, removes[j].first), big.NewInt(1))
			i++
		} else {
			// Remove-block inside network-block, will split network-block into two new blocks
			//
			// Make room for new network block
			blocks = append(blocks, nil)
			copy(blocks[i+1:], blocks[i:])
			blocks[i] = new(cidrBlock6)
			// update first half of the network-block (new)
			blocks[i].first = blocks[i+1].first
			tmp := big.NewInt(0)
			blocks[i].last = tmp.Sub(tmp.Add(tmp, removes[j].first), big.NewInt(1))
			// Update second half of the network-block (old)
			tmp = big.NewInt(1)
			tmp.Add(tmp, removes[j].last)
			blocks[i+1].first = tmp
			i++
			j++
		}
	}

	var merged []*net.IPNet
	for _, block := range blocks {
		if block == nil {
			continue
		}

		if err := splitRange6(big.NewInt(0), 0, block.first, block.last, &merged); err != nil {
			return nil, err
		}
	}

	return merged, nil
}
