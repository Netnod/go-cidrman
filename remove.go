package cidrman

import (
	"net"
)

// RemoveIPNets accepts two lists of mixed IP networks and removes the second list from the first and return a new list of IPNets.
// The remove will return the smallest possible list of IPNets.
// Example:
//     routableNets, err := RemoveIPNets(mixedListOfNets, rfc1918nets)
func RemoveIPNets(nets, rmnets []*net.IPNet) ([]*net.IPNet, error) {
	if nets == nil {
		return nil, nil
	}
	if len(nets) == 0 {
		return make([]*net.IPNet, 0), nil
	}
	if rmnets == nil {
		return nets, nil
	}
	if len(rmnets) == 0 {
		return nets, nil
	}

	// Merge nets and rmnet individually to have the miminal set of largets networks
	nets, err := MergeIPNets(nets)
	if err != nil {
		return nil, err
	}
	rmnets, err = MergeIPNets(rmnets)
	if err != nil {
		return nil, err
	}

	// Split into IPv4 and IPv6 lists.
	// Handle the lists separately and then combine.
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
	var remove4s cidrBlock4s
	var remove6s cidrBlock6s
	for _, net := range rmnets {
		ip4 := net.IP.To4()
		if ip4 != nil {
			remove4s = append(remove4s, newBlock4(ip4, net.Mask))
		} else {
			ip6 := net.IP.To16()
			remove6s = append(remove6s, newBlock6(ip6, net.Mask))
		}
	}

	new4s, err := remove4(block4s, remove4s)
	if err != nil {
		return nil, err
	}

	new6s, err := remove6(block6s, remove6s)
	if err != nil {
		return nil, err
	}

	merged := append(new4s, new6s...)
	return merged, nil
}

// RemoveCIDRs accepts two lists of mixed CIDR blocks and removes the second list from the first and return new a list of CIDRs.
func RemoveCIDRs(cidrs, removes []string) ([]string, error) {
	if cidrs == nil {
		return nil, nil
	}
	if len(cidrs) == 0 {
		return make([]string, 0), nil
	}
	if removes == nil {
		return cidrs, nil
	}
	if len(removes) == 0 {
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
	var rmnets []*net.IPNet
	for _, cidr := range removes {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		rmnets = append(rmnets, network)
	}

	newNets, err := RemoveIPNets(networks, rmnets)
	if err != nil {
		return nil, err
	}
	// Handle the situation where all cidrs were removed
	if len(newNets) == 0 {
		return make([]string, 0), nil
	}

	return ipNets(newNets).toCIDRs(), nil
}
