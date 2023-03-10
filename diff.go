package cidrman

import (
	"fmt"
	"net"
)

// DiffIPNets accepts two lists of mixed IP networks and return three lists of
// IPNets that only exists in first list, second list and in both lists, together with a status number.
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

	leftunique, err = RemoveIPNets(leftnets, rightnets)
	if err != nil {
		return nil, nil, nil, -1, fmt.Errorf("Error in creating leftunique: %w", err)
	}
	rightunique, err = RemoveIPNets(rightnets, leftnets)
	if err != nil {
		return nil, nil, nil, -1, fmt.Errorf("Error in creating rightunique: %w", err)
	}
	inboth, err = SubsetIPNets(leftnets, rightnets)
	if err != nil {
		return nil, nil, nil, -1, fmt.Errorf("Error in creating inboth: %w", err)
	}

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
func DiffCIDRs(leftcidrs, rightcidrs []string) (leftreturn []string, bothreturn []string, rightreturn []string, status int, err error) {
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

	// Handle the situation empty data
	if len(leftNets) == 0 {
		leftreturn = make([]string, 0)
	} else {
		leftreturn = ipNets(leftNets).toCIDRs()
	}
	if len(bothNets) == 0 {
		bothreturn = make([]string, 0)
	} else {
		bothreturn = ipNets(bothNets).toCIDRs()
	}
	if len(rightNets) == 0 {
		rightreturn = make([]string, 0)
	} else {
		rightreturn = ipNets(rightNets).toCIDRs()
	}

	return leftreturn, bothreturn, rightreturn, status, nil
}
