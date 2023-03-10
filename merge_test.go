// go test -v -run="TestMergeCIDRs"

package cidrman

import (
	"reflect"
	"testing"
)

func TestMergeCIDRs(t *testing.T) {
	type TestCase struct {
		Input  []string
		Output []string
		Error  bool
	}

	testCases := []TestCase{
		{
			Input:  nil,
			Output: nil,
			Error:  false,
		},
		{
			Input:  []string{},
			Output: []string{},
			Error:  false,
		},
		{
			Input: []string{
				"10.0.0.0/8",
			},
			Output: []string{
				"10.0.0.0/8",
			},
			Error: false,
		},
		{
			Input: []string{
				"10.0.0.0/8",
				"0.0.0.0/0",
			},
			Output: []string{
				"0.0.0.0/0",
			},
			Error: false,
		},
		{
			Input: []string{
				"10.0.0.0/8",
				"10.0.0.0/8",
			},
			Output: []string{
				"10.0.0.0/8",
			},
			Error: false,
		},
		{
			Input: []string{
				"192.0.128.0/24",
				"192.0.129.0/24",
			},
			Output: []string{
				"192.0.128.0/23",
			},
			Error: false,
		},
		{
			Input: []string{
				"192.0.129.0/24",
				"192.0.130.0/24",
			},
			Output: []string{
				"192.0.129.0/24",
				"192.0.130.0/24",
			},
			Error: false,
		},
		{
			Input: []string{
				"192.0.2.112/30",
				"192.0.2.116/31",
				"192.0.2.118/31",
			},
			Output: []string{
				"192.0.2.112/29",
			},
			Error: false,
		},
		// The same as above out of order.
		{
			Input: []string{
				"192.0.2.116/31",
				"192.0.2.118/31",
				"192.0.2.112/30",
			},
			Output: []string{
				"192.0.2.112/29",
			},
			Error: false,
		},
		{
			Input: []string{
				"192.0.2.112/30",
				"192.0.2.116/32",
				"192.0.2.118/31",
			},
			Output: []string{
				"192.0.2.112/30",
				"192.0.2.116/32",
				"192.0.2.118/31",
			},
			Error: false,
		},
		{
			Input: []string{
				"192.0.2.112/31",
				"192.0.2.116/31",
				"192.0.2.118/31",
			},
			Output: []string{
				"192.0.2.112/31",
				"192.0.2.116/30",
			},
			Error: false,
		},
		{
			Input: []string{
				"192.0.1.254/31",
				"192.0.2.0/28",
				"192.0.2.16/28",
				"192.0.2.32/28",
				"192.0.2.48/28",
				"192.0.2.64/28",
				"192.0.2.80/28",
				"192.0.2.96/28",
				"192.0.2.112/28",
				"192.0.2.128/28",
				"192.0.2.144/28",
				"192.0.2.160/28",
				"192.0.2.176/28",
				"192.0.2.192/28",
				"192.0.2.208/28",
				"192.0.2.224/28",
				"192.0.2.240/28",
				"192.0.3.0/28",
			},
			Output: []string{
				"192.0.1.254/31",
				"192.0.2.0/24",
				"192.0.3.0/28",
			},
			Error: false,
		},
		{
			Input: []string{
				"1.2.3.4/16",
				"1.2.3.5/16",
			},
			Output: []string{
				"1.2.0.0/16",
			},
			Error: false,
		},
		{
			Input: []string{
				"1.2.3.4/32",
				"1.2.3.4/32",
			},
			Output: []string{
				"1.2.3.4/32",
			},
			Error: false,
		},
		{
			Input: []string{
				"1.2.3.4/32",
				"1.2.3.4/24",
			},
			Output: []string{
				"1.2.3.0/24",
			},
			Error: false,
		},
		// IPv6 tests
		{
			Input: []string{
				"::/0",
			},
			Output: []string{
				"::/0",
			},
			Error: false,
		},
		{
			Input: []string{
				"fd00::/8",
				"::/0",
			},
			Output: []string{
				"::/0",
			},
			Error: false,
		},
		{
			Input: []string{
				"fd00::/8",
				"fd00::/8",
			},
			Output: []string{
				"fd00::/8",
			},
			Error: false,
		},
		{
			Input: []string{
				"2001:db8:0:2::/64",
				"2001:db8:0:3::/64",
			},
			Output: []string{
				"2001:db8:0:2::/63",
			},
			Error: false,
		},
		{
			Input: []string{
				"2001:db8:0:1::/64",
				"2001:db8:0:2::/64",
			},
			Output: []string{
				"2001:db8:0:1::/64",
				"2001:db8:0:2::/64",
			},
			Error: false,
		},
		{
			Input: []string{
				"2001:db8:0:0::/64",
				"2001:db8:0:1::/64",
				"2001:db8:0:2::/63",
			},
			Output: []string{
				"2001:db8::/62",
			},
			Error: false,
		},
		// The same as above out of order.
		{
			Input: []string{
				"2001:db8:0:1::/64",
				"2001:db8:0:2::/63",
				"2001:db8:0:0::/64",
			},
			Output: []string{
				"2001:db8::/62",
			},
			Error: false,
		},
		{
			Input: []string{
				"2001:db8:0:1::/63",
				"2001:db8:0:2::/64",
				"2001:db8:0:4::/62",
			},
			Output: []string{
				"2001:db8::/63",
				"2001:db8:0:2::/64",
				"2001:db8:0:4::/62",
			},
			Error: false,
		},
		{
			Input: []string{
				"fd00:0:1:1::/63",
				"fd00:0:1:2::/64",
				"fd00:0:1:4::/62",
			},
			Output: []string{
				"fd00:0:1::/63",
				"fd00:0:1:2::/64",
				"fd00:0:1:4::/62",
			},
			Error: false,
		},
		{
			Input: []string{
				"2001:db8:0:1::/64",
				"2001:db8:0:2::/64",
				"2001:db8:0:3::/64",
			},
			Output: []string{
				"2001:db8:0:1::/64",
				"2001:db8:0:2::/63",
			},
			Error: false,
		},
		{
			Input: []string{
				"2001:db8:0:1::80/127",
				"2001:db8:0:2::/80",
				"2001:db8:0:2:1::/80",
				"2001:db8:0:2:2::/80",
				"2001:db8:0:2:3::/80",
				"2001:db8:0:2:4::/80",
				"2001:db8:0:2:5::/80",
				"2001:db8:0:2:6::/80",
				"2001:db8:0:2:7::/80",
				"2001:db8:0:2:8::/80",
				"2001:db8:0:2:9::/80",
				"2001:db8:0:2:a::/80",
				"2001:db8:0:2:b::/80",
				"2001:db8:0:2:c::/80",
				"2001:db8:0:2:d::/80",
				"2001:db8:0:2:e::/80",
				"2001:db8:0:2:f::/80",
				"2001:db8:0:3::/80",
			},
			Output: []string{
				"2001:db8:0:1::80/127",
				"2001:db8:0:2::/76",
				"2001:db8:0:3::/80",
			},
			Error: false,
		},
		{
			Input: []string{
				"fd00:0:1:2:12::/64",
				"fd00:0:1:2:24::/64",
			},
			Output: []string{
				"fd00:0:1:2::/64",
			},
			Error: false,
		},
		{
			Input: []string{
				"fd00:0:1:2:12::53/128",
				"fd00:0:1:2:12::53/128",
			},
			Output: []string{
				"fd00:0:1:2:12::53/128",
			},
			Error: false,
		},
		{
			Input: []string{
				"fd00:0:1:2:12::53/128",
				"fd00:0:1:2:12::53/80",
			},
			Output: []string{
				"fd00:0:1:2:12::/80",
			},
			Error: false,
		},
		// Mixed IPv4 and IPv6 tests
		{
			Input: []string{
				"2001:db8:0:2::/64",
				"2001:db8:0:3::/64",
				"192.0.128.0/24",
				"192.0.129.0/24",
			},
			Output: []string{
				"192.0.128.0/23",
				"2001:db8:0:2::/63",
			},
			Error: false,
		},
	}

	for _, testCase := range testCases {
		output, err := MergeCIDRs(testCase.Input)
		if err != nil {
			if !testCase.Error {
				t.Errorf("MergeCIDRS(%#v) failed: %s", testCase.Input, err.Error())
			}
		} else if testCase.Error {
			t.Errorf("MergeCIDRS(%#v) didn't return with error as expected", testCase.Input)
		}

		if !reflect.DeepEqual(testCase.Output, output) {
			t.Errorf("MergeCIDRS(%#v) expected: %#v, got: %#v", testCase.Input, testCase.Output, output)
		}
	}
}
