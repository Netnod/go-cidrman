// go test -v -run="TestSubsetCIDRs"

package cidrman

import (
	"reflect"
	"testing"
)

func TestSubsetCIDRs(t *testing.T) {
	type TestCase struct {
		Input  []string
		Subset []string
		Output []string
		Error  bool
	}

	testCases := []TestCase{
		{
			Input:  nil,
			Subset: nil,
			Output: nil,
			Error:  false,
		},
		{
			Input:  []string{},
			Subset: nil,
			Output: []string{},
			Error:  false,
		},
		{
			Input:  nil,
			Subset: []string{},
			Output: nil,
			Error:  false,
		},
		{
			Input:  []string{},
			Subset: []string{},
			Output: []string{},
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
			},
			Subset: []string{},
			// With nothing to keep, we get back an empty list
			Output: []string{},
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
			},
			Subset: nil,
			Output: []string{},
			Error:  false,
		},
		{
			Input:  nil,
			Subset: []string{
				"10.0.0.0/8",
			},
			Output: nil,
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
			},
			Subset: []string{
				"10.0.0.0/8",
			},
			Output: []string{
				"10.0.0.0/8",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
				"0.0.0.0/0",
			},
			Subset: []string{
				"127.0.0.0/8",
			},
			Output: []string{
				"127.0.0.0/8",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
				"10.0.0.0/8",
			},
			Subset: []string{
				"10.0.0.0/8",
			},
			Output: []string{
				"10.0.0.0/8",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"192.0.128.0/24",
				"192.0.129.0/24",
			},
			Subset: []string{
				"192.0.0.0/16",
			},
			// SubsetIPNets will first do MergeIPNets() before processing the subset list
			Output: []string{
				"192.0.128.0/23",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"192.0.128.0/24",
				"192.0.129.0/24",
			},
			Subset: []string{
				"192.0.128.128/25",
				"192.0.129.0/25",
			},
			Output: []string{
				"192.0.128.128/25",
				"192.0.129.0/25",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"192.0.128.0/24",
				"192.0.139.0/24",
			},
			Subset: []string{
				"192.0.128.0/23",
			},
			Output: []string{
				"192.0.128.0/24",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"172.16.10.0/24",
				"172.16.11.0/24",
				"172.16.12.0/24",
				"172.16.13.0/24",
				"172.16.14.0/24",
			},
			Subset: []string{
				"172.16.8.0/22",
			},
			Output: []string{
				"172.16.10.0/23",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"172.16.10.0/24",
				"172.16.11.0/24",
				"172.16.12.0/24",
				"172.16.13.0/24",
				"172.16.14.0/24",
			},
			Subset: []string{
				"172.16.12.0/22",
			},
			Output: []string{
				"172.16.12.0/23",
				"172.16.14.0/24",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"172.16.8.0/20",
			},
			Subset: []string{
				"172.16.12.0/24",
				"172.16.14.0/24",
			},
			Output: []string{
				"172.16.12.0/24",
				"172.16.14.0/24",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"172.16.10.0/24",
				"172.16.11.0/24",
				"172.16.12.0/24",
				"172.16.13.0/24",
				"172.16.14.0/24",
			},
			Subset: []string{
				"172.16.8.0/21",
			},
			Output: []string{
				"172.16.10.0/23",
				"172.16.12.0/23",
				"172.16.14.0/24",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
				"172.16.10.0/24",
				"172.16.11.0/24",
				"172.16.12.0/24",
				"172.16.13.0/24",
				"172.16.14.0/24",
				"192.0.128.0/23",
				"192.0.139.0/24",
			},
			Subset: []string{
				"172.16.8.0/22",
				"10.10.10.0/24",
				"172.16.13.0/24",
				"172.16.14.128/26",
			},
			Output: []string{
				"10.10.10.0/24",
				"172.16.10.0/23",
				"172.16.13.0/24",
				"172.16.14.128/26",
			},
			Error:  false,
		},
		// IPv6 tests
		{
			Input: []string{
				"::/0",
			},
			Subset: []string{},
			// With nothing to keep, we get back an empty list
			Output: []string{},
			Error: false,
		},
		{
			Input:  []string{
				"fd00::/8",
			},
			Subset: nil,
			Output: []string{},
			Error:  false,
		},
		{
			Input:  nil,
			Subset: []string{
				"fd00::/8",
			},
			Output: nil,
			Error:  false,
		},
		{
			Input:  []string{
				"fd00::/8",
			},
			Subset: []string{
				"fd00::/8",
			},
			Output: []string{
				"fd00::/8",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"fd00::/8",
				"::/0",
			},
			Subset: []string{
				"2001:db8:0:2::/64",
			},
			Output: []string{
				"2001:db8:0:2::/64",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"2001:db8:0:2::/64",
				"2001:db8:0:2::/64",
			},
			Subset: []string{
				"2001:db8:0:2::/64",
			},
			Output: []string{
				"2001:db8:0:2::/64",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"2001:db8:0:2::/64",
				"2001:db8:0:3::/64",
			},
			Subset: []string{
				"2001:db8::/32",
				"fd00::/8",
			},
			// SubsetIPNets will first do MergeIPNets() before processing the subset list
			Output: []string{
				"2001:db8:0:2::/63",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"2001:db8:0:2::/64",
				"2001:db8:0:3::/64",
			},
			Subset: []string{
				"2001:db8:0:2:ffff::/72",
				"2001:db8:0:3::/80",
			},
			Output: []string{
				"2001:db8:0:2:ff00::/72",
				"2001:db8:0:3::/80",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"2001:db8:0:2::/64",
				"2001:db8:1:3::/64",
			},
			Subset: []string{
				"2001:db8::/48",
			},
			Output: []string{
				"2001:db8:0:2::/64",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"fd00:0:4711:a::/64",
				"fd00:0:4711:b::/64",
				"fd00:0:4711:c::/64",
				"fd00:0:4711:d::/64",
				"fd00:0:4711:e::/64",
			},
			Subset: []string{
				"fd00:0:4711:8::/62",
			},
			Output: []string{
				"fd00:0:4711:a::/63",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"fd00:0:4711:a::/64",
				"fd00:0:4711:b::/64",
				"fd00:0:4711:c::/64",
				"fd00:0:4711:d::/64",
				"fd00:0:4711:e::/64",
			},
			Subset: []string{
				"fd00:0:4711:c::/62",
			},
			Output: []string{
				"fd00:0:4711:c::/63",
				"fd00:0:4711:e::/64",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"fd00:0:4711:8::/61",
			},
			Subset: []string{
				"fd00:0:4711:c::/64",
				"fd00:0:4711:e::/64",
			},
			Output: []string{
				"fd00:0:4711:c::/64",
				"fd00:0:4711:e::/64",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"fd00:0:4711:a::/64",
				"fd00:0:4711:b::/64",
				"fd00:0:4711:c::/64",
				"fd00:0:4711:d::/64",
				"fd00:0:4711:e::/64",
			},
			Subset: []string{
				"fd00:0:4711:8::/61",
			},
			Output: []string{
				"fd00:0:4711:a::/63",
				"fd00:0:4711:c::/63",
				"fd00:0:4711:e::/64",
			},
			Error:  false,
		},
		// Mixed blocks
		{
			Input:  []string{
				"fd00::/8",
			},
			Subset: []string{
				"10.0.0.0/8",
			},
			Output: []string{},
			Error:  false,
		},
		{
			Input:  []string{
				"fd00::/8",
				"0.0.0.0/4",
			},
			Subset: []string{
				"10.0.0.0/8",
			},
			Output: []string{
				"10.0.0.0/8",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
				"fc00::/7",
			},
			Subset: []string{
				"fd00::/8",
			},
			Output: []string{
				"fd00::/8",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
				"2001:db8:0:2::/64",
				"2001:db8:0:3::/64",
				"192.0.128.0/24",
				"192.0.129.0/24",
			},
			Subset: []string{
				"192.0.128.0/24",
				"2001:db8:0:3::/64",
			},
			Output: []string{
				"192.0.128.0/24",
				"2001:db8:0:3::/64",
			},
			Error:  false,
		},
	}

	for _, testCase := range testCases {
		output, err := SubsetCIDRs(testCase.Input, testCase.Subset)
		if err != nil {
			if !testCase.Error {
				t.Errorf("SubsetCIDRs(%#v, %#v) failed: %s", testCase.Input, testCase.Subset, err.Error())
			}
		} else if testCase.Error {
			t.Errorf("SubsetCIDRs(%#v, %#v) didn't return with error as expected", testCase.Input, testCase.Subset)
		}

		if !reflect.DeepEqual(testCase.Output, output) {
			t.Errorf("SubsetCIDRs(%#v, %#v) expected: %#v, got: %#v", testCase.Input, testCase.Subset, testCase.Output, output)
		}
	}
}
