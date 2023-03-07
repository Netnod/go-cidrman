// go test -v -run="TestRemoveCIDRs"

package cidrman

import (
	"reflect"
	"testing"
)

func TestRemoveCIDRs(t *testing.T) {
	type TestCase struct {
		Input  []string
		Remove []string
		Output []string
		Error  bool
	}

	testCases := []TestCase{
		{
			Input:  nil,
			Remove: nil,
			Output: nil,
			Error:  false,
		},
		{
			Input:  []string{},
			Remove: nil,
			Output: []string{},
			Error:  false,
		},
		{
			Input:  nil,
			Remove: []string{},
			Output: nil,
			Error:  false,
		},
		{
			Input:  []string{},
			Remove: []string{},
			Output: []string{},
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
			},
			Remove: []string{},
			Output: []string{
				"10.0.0.0/8",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
			},
			Remove: nil,
			Output: []string{
				"10.0.0.0/8",
			},
			Error:  false,
		},
		{
			Input:  nil,
			Remove: []string{
				"10.0.0.0/8",
			},
			Output: nil,
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
			},
			Remove: []string{
				"10.0.0.0/8",
			},
			Output: []string{},
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
				"0.0.0.0/0",
			},
			Remove: []string{},
			// With nothing to remove, we get back what we sent
			Output: []string{
				"10.0.0.0/8",
				"0.0.0.0/0",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
				"10.0.0.0/8",
			},
			Remove: []string{},
			// With nothing to remove, we get back what we sent
			Output: []string{
				"10.0.0.0/8",
				"10.0.0.0/8",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"192.0.128.0/24",
				"192.0.129.0/24",
			},
			Remove: []string{
				"10.0.0.0/8",
			},
			// RemoveIPNets will first do MergeIPNets() before processing the remove list
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
			Remove: []string{
				"192.0.128.0/23",
			},
			Output: []string{},
			Error:  false,
		},
		{
			Input:  []string{
				"192.0.128.0/24",
				"192.0.139.0/24",
			},
			Remove: []string{
				"192.0.128.0/23",
			},
			Output: []string{
				"192.0.139.0/24",
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
			Remove: []string{
				"172.16.8.0/22",
			},
			Output: []string{
				"172.16.12.0/23",
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
			Remove: []string{
				"172.16.12.0/22",
			},
			Output: []string{
				"172.16.10.0/23",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"172.16.8.0/21",
			},
			Remove: []string{
				"172.16.12.0/23",
			},
			Output: []string{
				"172.16.8.0/22",
				"172.16.14.0/23",
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
			Remove: []string{
				"172.16.8.0/21",
			},
			Output: []string{},
			Error:  false,
		},
		// IPv6 tests
		{
			Input: []string{
				"::/0",
			},
			Remove: []string{},
			Output: []string{
				"::/0",
			},
			Error: false,
		},
		{
			Input:  []string{
				"fd00::/8",
			},
			Remove: nil,
			Output: []string{
				"fd00::/8",
			},
			Error:  false,
		},
		{
			Input:  nil,
			Remove: []string{
				"fd00::/8",
			},
			Output: nil,
			Error:  false,
		},
		{
			Input:  []string{
				"fd00::/8",
			},
			Remove: []string{
				"fd00::/8",
			},
			Output: []string{},
			Error:  false,
		},
		{
			Input:  []string{
				"fd00::/8",
				"::/0",
			},
			Remove: []string{},
			// With nothing to remove, we get back what we sent
			Output: []string{
				"fd00::/8",
				"::/0",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"2001:db8:0:2::/64",
				"2001:db8:0:2::/64",
			},
			Remove: []string{},
			// With nothing to remove, we get back what we sent
			Output: []string{
				"2001:db8:0:2::/64",
				"2001:db8:0:2::/64",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"2001:db8:0:2::/64",
				"2001:db8:0:3::/64",
			},
			Remove: []string{
				"fd00::/8",
			},
			// RemoveIPNets will first do MergeIPNets() before processing the remove list
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
			Remove: []string{
				"2001:db8:0:2::/63",
			},
			Output: []string{},
			Error:  false,
		},
		{
			Input:  []string{
				"2001:db8:0:2::/64",
				"2001:db8:1:3::/64",
			},
			Remove: []string{
				"2001:db8::/48",
			},
			Output: []string{
				"2001:db8:1:3::/64",
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
			Remove: []string{
				"fd00:0:4711:8::/62",
			},
			Output: []string{
				"fd00:0:4711:c::/63",
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
			Remove: []string{
				"fd00:0:4711:c::/62",
			},
			Output: []string{
				"fd00:0:4711:a::/63",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"fd00:0:4711:8::/61",
			},
			Remove: []string{
				"fd00:0:4711:c::/63",
			},
			Output: []string{
				"fd00:0:4711:8::/62",
				"fd00:0:4711:e::/63",
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
			Remove: []string{
				"fd00:0:4711:8::/61",
			},
			Output: []string{},
			Error:  false,
		},
		// Mixed blocks
		{
			Input:  []string{
				"fd00::/8",
			},
			Remove: []string{
				"10.0.0.0/8",
			},
			Output: []string{
				"fd00::/8",
			},
			Error:  false,
		},
		{
			Input:  []string{
				"10.0.0.0/8",
			},
			Remove: []string{
				"fd00::/8",
			},
			Output: []string{
				"10.0.0.0/8",
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
			Remove: []string{
				"192.0.128.0/24",
				"2001:db8:0:3::/64",
			},
			Output: []string{
				"10.0.0.0/8",
				"192.0.129.0/24",
				"2001:db8:0:2::/64",
			},
			Error:  false,
		},
	}

	for _, testCase := range testCases {
		output, err := RemoveCIDRs(testCase.Input, testCase.Remove)
		if err != nil {
			if !testCase.Error {
				t.Errorf("RemoveCIDRs(%#v, %#v) failed: %s", testCase.Input, testCase.Remove, err.Error())
			}
		} else if testCase.Error {
			t.Errorf("RemoveCIDRs(%#v, %#v) didn't return with error as expected", testCase.Input, testCase.Remove)
		}

		if !reflect.DeepEqual(testCase.Output, output) {
			t.Errorf("RemoveCIDRs(%#v, %#v) expected: %#v, got: %#v", testCase.Input, testCase.Remove, testCase.Output, output)
		}
	}
}
