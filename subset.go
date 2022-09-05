package cidrman

import (
	"net"
)

// SubsetIPNets accepts two lists of mixed IP networks and return a new list of IPNets that exsists/overlaps in both lists.
// The SubsetIPNets() will return the smallest possible list of IPNets.
// Example:
//     internalNets, err := SubsetIPNets(mixedListOfNets, rfc1918nets)
func SubsetIPNets(nets, subsetnets []*net.IPNet) ([]*net.IPNet, error) {
	if nets == nil {
		return nil, nil
	}
	if len(nets) == 0 {
		return make([]*net.IPNet, 0), nil
	}
	if subsetnets == nil {
		return nets, nil
	}
	if len(subsetnets) == 0 {
		return nets, nil
	}

	// Merge nets and subsetnets individually to have the miminal set of largets networks
	nets, err := MergeIPNets(nets)
	if err != nil {
		return nil, err
	}
	subsetnets, err = MergeIPNets(subsetnets)
	if err != nil {
		return nil, err
	}

	// Split into IPv4 and IPv6 lists.
	// Handle the list separately and then combine.
	var block4s cidrBlock4s
	var block6s cidrBlock6s
	for _, net := range nets {
		ip4 := net.IP.To4()
		if ip4 != nil {
			block4s = append(block4s, newBlock4(ip4, net.Mask))
		} else {
			ip6 := net.IP.To16()
			block6s = append(block6s, newBlock6(ip6, net.Mask))
		}
	}
	var subset4s cidrBlock4s
	var subset6s cidrBlock6s
	for _, net := range subsetnets {
		ip4 := net.IP.To4()
		if ip4 != nil {
			subset4s = append(subset4s, newBlock4(ip4, net.Mask))
		} else {
			ip6 := net.IP.To16()
			subset6s = append(subset6s, newBlock6(ip6, net.Mask))
		}
	}

	new4s, err := subset4(block4s, subset4s)
	if err != nil {
		return nil, err
	}

	new6s, err := subset6(block6s, subset6s)
	if err != nil {
		return nil, err
	}

	merged := append(new4s, new6s...)
	return merged, nil
}

// SubsetCIDRs accepts two lists of mixed CIDR blocks and return a new list of CIDRs that exsists/overlaps in both lists.
func SubsetCIDRs(cidrs, subsets []string) ([]string, error) {
	if cidrs == nil {
		return nil, nil
	}
	if len(cidrs) == 0 {
		return make([]string, 0), nil
	}
	if subsets == nil {
		return cidrs, nil
	}
	if len(subsets) == 0 {
		return cidrs, nil
	}

	var networks []*net.IPNet
	for _, cidr := range cidrs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		networks = append(networks, network)
	}
	var subsetnets []*net.IPNet
	for _, cidr := range subsets {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		subsetnets = append(subsetnets, network)
	}

	newNets, err := SubsetIPNets(networks, subsetnets)
	if err != nil {
		return nil, err
	}
	// Handle the situation where no cidrs overlapped
	if len(newNets) == 0 {
		return make([]string, 0), nil
	}

	return ipNets(newNets).toCIDRs(), nil
}
