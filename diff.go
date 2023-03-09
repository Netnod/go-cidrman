package cidrman

import (
	"fmt"
	"net"
)

// DiffIPNets accepts two lists of mixed IP networks and return three lists of
// IPNets that only exists in first list, second lost and in both lists, together with a status number.
// The DiffIPNets() will return the smallest possible lists of IPNets.
// Example:
//    onlyInLeftNets, inBothNets, onlyInRightNets, status, err := DiffIPNets(leftListOfNets, rightListOfNets)
// status values:
//    0 = both lists with exact same nets after individual merge
//    1 = leftnets and rightnets differ with some overlap and some unique nets
//    2 = leftnets same as rightnets with some extra nets
//    3 = rightnets same as leftnets with some extra nets
//    4 = leftnets and rightnets had no overlap at all
//    5 = empty input nets
//
//   -1 = some error processing nets in sub functions, see err
//   -2 = missing input data, see err
func DiffIPNets(leftnets, rightnets []*net.IPNet) (leftunique, inboth, rightunique []*net.IPNet, status int, err error) {
	if leftnets == nil || rightnets == nil {
		return nil, nil, nil, -2, fmt.Errorf("Missing input data in DiffINets()")
	}
	if len(leftnets) == 0 && len(rightnets) == 0 {
		return make([]*net.IPNet, 0), make([]*net.IPNet, 0), make([]*net.IPNet, 0), 5, nil
	}
	if len(leftnets) == 0 {
		return make([]*net.IPNet, 0), make([]*net.IPNet, 0), rightnets, 3, nil
	}
	if len(rightnets) == 0 {
		return leftnets, make([]*net.IPNet, 0), make([]*net.IPNet, 0), 2, nil
	}

	// Merge nets and subsetnets individually to have the miminal set of largets networks
	leftnets, err = MergeIPNets(leftnets)
	if err != nil {
		return nil, nil, nil, -1, fmt.Errorf("Merge error leftnets: %w", err)
	}
	rightnets, err = MergeIPNets(rightnets)
	if err != nil {
		return nil, nil, nil, -1, fmt.Errorf("Merge error rightnets: %w", err)
	}

	// Split into IPv4 and IPv6 lists.
	// Handle the list separately.
	var left4s cidrBlock4s
	var left6s cidrBlock6s
	for _, net := range leftnets {
		ip4 := net.IP.To4()
		if ip4 != nil {
			left4s = append(left4s, newBlock4(ip4, net.Mask))
		} else {
			ip6 := net.IP.To16()
			left6s = append(left6s, newBlock6(ip6, net.Mask))
		}
	}
	var right4s cidrBlock4s
	var right6s cidrBlock6s
	for _, net := range rightnets {
		ip4 := net.IP.To4()
		if ip4 != nil {
			right4s = append(right4s, newBlock4(ip4, net.Mask))
		} else {
			ip6 := net.IP.To16()
			right6s = append(right6s, newBlock6(ip6, net.Mask))
		}
	}

	// Create leftunique
	var leftunique4s []*net.IPNet
	if len(left4s) > 0 {
		leftunique4s, err = remove4(copy4s(left4s), copy4s(right4s))
		if err != nil {
			return nil, nil, nil, -1, fmt.Errorf("Error in creating leftunique4s: %w", err)
		}
	}
	var leftunique6s []*net.IPNet
	if len(left6s) > 0 {
		leftunique6s, err = remove6(copy6s(left6s), copy6s(right6s))
		if err != nil {
			return nil, nil, nil, -1, fmt.Errorf("Error in creating leftunique6s: %w", err)
		}
	}
	leftunique = append(leftunique4s, leftunique6s...)

	// Create rightunique
	var rightunique4s []*net.IPNet
	if len(right4s) > 0 {
		rightunique4s, err = remove4(copy4s(right4s), copy4s(left4s))
		if err != nil {
			return nil, nil, nil, -1, fmt.Errorf("Error in creating rightunique4s: %w", err)
		}
	}
	var rightunique6s []*net.IPNet
	if len(right6s) > 0 {
		rightunique6s, err = remove6(copy6s(right6s), copy6s(left6s))
		if err != nil {
			return nil, nil, nil, -1, fmt.Errorf("Error in creating rightunique6s: %w", err)
		}
	}
	rightunique = append(rightunique4s, rightunique6s...)

	// Create inboth
	var inboth4s []*net.IPNet
	if len(left4s) > 0 {
		inboth4s, err = subset4(left4s, right4s)
		if err != nil {
			return nil, nil, nil, -1, fmt.Errorf("Error in creating inboth4s: %w", err)
		}
	}
	var inboth6s []*net.IPNet
	if len(left6s) > 0 {
		inboth6s, err = subset6(left6s, right6s)
		if err != nil {
			return nil, nil, nil, -1, fmt.Errorf("Error in creating inboth6s: %w", err)
		}
	}
	inboth = append(inboth4s, inboth6s...)

	if len(inboth) == 0 {
		status = 4
	} else if len(leftunique) == 0 && len(rightunique) == 0 {
		status = 0
	} else if len(leftunique) == 0 {
		status = 3
	} else if len(rightunique) == 0 {
		status = 2
	} else {
		status = 1
	}
	return leftunique, inboth, rightunique, status, nil
}

// DiffCIDRs accepts two lists of mixed CIDR blocks and return three lists of
// CIDRs that  only exists in first list, second lost and in both lists, together with a status number.
func DiffCIDRs(leftcidrs, rightcidrs []string) ([]string, []string, []string, int, error) {
	if leftcidrs == nil || rightcidrs == nil {
		return nil, nil, nil, -2, fmt.Errorf("Missing input data in DiffCIDRs()")
	}
	if len(leftcidrs) == 0 && len(rightcidrs) == 0 {
		return make([]string, 0), make([]string, 0), make([]string, 0), 5, nil
	}
	if len(leftcidrs) == 0 {
		return make([]string, 0), make([]string, 0), rightcidrs, 3, nil
	}
	if len(rightcidrs) == 0 {
		return leftcidrs, make([]string, 0), make([]string, 0), 2, nil
	}

	var leftnetworks []*net.IPNet
	for _, cidr := range leftcidrs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, nil, nil, -1, err
		}
		leftnetworks = append(leftnetworks, network)
	}
	var rightnetworks []*net.IPNet
	for _, cidr := range rightcidrs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, nil, nil, -1, err
		}
		rightnetworks = append(rightnetworks, network)
	}

	leftNets, bothNets, rightNets, status,  err := DiffIPNets(leftnetworks, rightnetworks)
	if err != nil {
		return nil, nil, nil, status, err
	}

	return ipNets(leftNets).toCIDRs(), ipNets(bothNets).toCIDRs(), ipNets(rightNets).toCIDRs(), status, nil
}
