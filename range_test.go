// go test -v -run="TestIPRangeToCIDRs"

package cidrman

import (
	"reflect"
	"testing"
)

func TestIPRangeToCIDRs(t *testing.T) {
	type TestCase struct {
		Lo     string
		Hi     string
		Output []string
		Error  bool
	}

	testCases := []TestCase{
		{
			Lo:     "abcdefgh",
			Hi:     "",
			Output: nil,
			Error:  true,
		},
		{
			Lo:     "192.168.1.12",
			Hi:     "192.168.1.11",
			Output: nil,
			Error:  true,
		},
		{
			Lo:     "0.0.0.1",
			Hi:     "2001:0db8:0000:0000:0000:ff00:0042:8329",
			Output: nil,
			Error:  true,
		},
		{
			Lo: "192.168.1.1",
			Hi: "192.168.1.1",
			Output: []string{
				"192.168.1.1/32",
			},
			Error: false,
		},
		{
			Lo: "192.168.1.1",
			Hi: "192.168.1.12",
			Output: []string{
				"192.168.1.1/32",
				"192.168.1.2/31",
				"192.168.1.4/30",
				"192.168.1.8/30",
				"192.168.1.12/32",
			},
			Error: false,
		},
		// Worst case.
		{
			Lo: "0.0.0.1",
			Hi: "255.255.255.254",
			Output: []string{
				"0.0.0.1/32",
				"0.0.0.2/31",
				"0.0.0.4/30",
				"0.0.0.8/29",
				"0.0.0.16/28",
				"0.0.0.32/27",
				"0.0.0.64/26",
				"0.0.0.128/25",
				"0.0.1.0/24",
				"0.0.2.0/23",
				"0.0.4.0/22",
				"0.0.8.0/21",
				"0.0.16.0/20",
				"0.0.32.0/19",
				"0.0.64.0/18",
				"0.0.128.0/17",
				"0.1.0.0/16",
				"0.2.0.0/15",
				"0.4.0.0/14",
				"0.8.0.0/13",
				"0.16.0.0/12",
				"0.32.0.0/11",
				"0.64.0.0/10",
				"0.128.0.0/9",
				"1.0.0.0/8",
				"2.0.0.0/7",
				"4.0.0.0/6",
				"8.0.0.0/5",
				"16.0.0.0/4",
				"32.0.0.0/3",
				"64.0.0.0/2",
				"128.0.0.0/2",
				"192.0.0.0/3",
				"224.0.0.0/4",
				"240.0.0.0/5",
				"248.0.0.0/6",
				"252.0.0.0/7",
				"254.0.0.0/8",
				"255.0.0.0/9",
				"255.128.0.0/10",
				"255.192.0.0/11",
				"255.224.0.0/12",
				"255.240.0.0/13",
				"255.248.0.0/14",
				"255.252.0.0/15",
				"255.254.0.0/16",
				"255.255.0.0/17",
				"255.255.128.0/18",
				"255.255.192.0/19",
				"255.255.224.0/20",
				"255.255.240.0/21",
				"255.255.248.0/22",
				"255.255.252.0/23",
				"255.255.254.0/24",
				"255.255.255.0/25",
				"255.255.255.128/26",
				"255.255.255.192/27",
				"255.255.255.224/28",
				"255.255.255.240/29",
				"255.255.255.248/30",
				"255.255.255.252/31",
				"255.255.255.254/32",
			},
			Error: false,
		},
		{
			Lo: "0.0.0.0",
			Hi: "255.255.255.255",
			Output: []string{
				"0.0.0.0/0",
			},
			Error: false,
		},
		{
			Lo: "0000:0000:0000:0000:0000:0000:0000:0000",
			Hi: "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			Output: []string{
				"::/0",
			},
			Error: false,
		},
		{
			Lo:     "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			Hi:     "ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffd",
			Output: nil,
			Error:  true,
		},
		{
			Lo: "ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffd",
			Hi: "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			Output: []string{
				"ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffd/128",
				"ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffe/127",
			},
			Error: false,
		},
		{
			Lo: "2001:0db8:0000:0000:0000:ff00:0042:8328",
			Hi: "2001:0db8:0000:0000:0000:ff00:0042:8328",
			Output: []string{
				"2001:db8::ff00:42:8328/128",
			},
			Error: false,
		},
	}

	for _, testCase := range testCases {
		output, err := IPRangeToCIDRs(testCase.Lo, testCase.Hi)
		if err != nil {
			if !testCase.Error {
				t.Errorf("IPRangeToCIDRs(%s, %s) failed: %s", testCase.Lo, testCase.Hi, err.Error())
			}
			continue
		}
		if !reflect.DeepEqual(testCase.Output, output) {
			t.Errorf("IPRangeToCIDRs(%s, %s) expected: %#v, got: %#v", testCase.Lo, testCase.Hi, testCase.Output, output)
		}
	}
}
