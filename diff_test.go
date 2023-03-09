// go test -v -run="TestDiffCIDRs"

package cidrman

import (
	"reflect"
	"testing"
)

func TestDiffCIDRs(t *testing.T) {
	type TestCase struct {
		TestName string
		LeftIn   []string
		RightIn  []string
		LeftOut  []string
		BothOut  []string
		RightOut []string
		Status   int
		Error    bool
	}

	testCases := []TestCase{
		{
			TestName: "#1",
			LeftIn:   nil,
			RightIn:  nil,
			LeftOut:  nil,
			BothOut:  nil,
			RightOut: nil,
			Status:   -2,
			Error:    true,
		},
		{
			TestName: "#2",
			LeftIn:   []string{},
			RightIn:  nil,
			LeftOut:  nil,
			BothOut:  nil,
			RightOut: nil,
			Status:   -2,
			Error:    true,
		},
		{
			TestName: "#3",
			LeftIn:   nil,
			RightIn:  []string{},
			LeftOut:  nil,
			BothOut:  nil,
			RightOut: nil,
			Status:   -2,
			Error:    true,
		},
		{
			TestName: "#4",
			LeftIn:   []string{},
			RightIn:  []string{},
			LeftOut:  []string{},
			BothOut:  []string{},
			RightOut: []string{},
			Status:   5,
			Error:    false,
		},
		{
			TestName: "#5",
			LeftIn:   []string{
				"10.0.0.0/8",
			},
			RightIn:  []string{},
			LeftOut:  []string{
				"10.0.0.0/8",
			},
			BothOut:  []string{},
			RightOut: []string{},
			Status:   2,
			Error:    false,
		},
		{
			TestName: "#6",
			LeftIn:   []string{},
			RightIn:  []string{
				"10.0.0.0/8",
			},
			LeftOut:  []string{},
			BothOut:  []string{},
			RightOut: []string{
				"10.0.0.0/8",
			},
			Status:   3,
			Error:    false,
		},
		{
			TestName: "#7",
			LeftIn:   []string{
				"10.0.0.0/8",
				"11.0.0.0/8",
			},
			RightIn:  []string{
				"11.0.0.0/8",
				"10.0.0.0/8",
			},
			LeftOut:  []string{},
			BothOut:  []string{
				"10.0.0.0/7",
			},
			RightOut: []string{},
			Status:   0,
			Error:    false,
		},
		{
			TestName: "#8",
			LeftIn:   []string{
				"10.0.0.0/8",
				"11.0.0.0/8",
			},
			RightIn:  []string{
				"10.0.0.0/8",
			},
			LeftOut:  []string{
				"11.0.0.0/8",
			},
			BothOut:  []string{
				"10.0.0.0/8",
			},
			RightOut: []string{},
			Status:   2,
			Error:    false,
		},
		{
			TestName: "#9",
			LeftIn:   []string{
				"10.0.0.0/8",
				"172.16.0.0/16",
			},
			RightIn:  []string{
				"172.16.16.0/20",
				"192.168.16.0/24",
			},
			LeftOut:  []string{
				"10.0.0.0/8",
				"172.16.0.0/20",
				"172.16.32.0/19",
				"172.16.64.0/18",
				"172.16.128.0/17",
			},
			BothOut:  []string{
				"172.16.16.0/20",
			},
			RightOut: []string{
				"192.168.16.0/24",
			},
			Status:   1,
			Error:    false,
		},
		{
			TestName: "#10",
			LeftIn:   []string{
				"192.0.128.128/25",
				"192.0.129.0/25",
			},
			RightIn:  []string{
				"192.0.128.0/24",
				"192.0.129.0/24",
			},
			LeftOut:  []string{},
			BothOut:  []string{
				"192.0.128.128/25",
				"192.0.129.0/25",
			},
			RightOut: []string{
				"192.0.128.0/25",
				"192.0.129.128/25",
			},
			Status:   3,
			Error:    false,
		},
		{
			TestName: "#11",
			LeftIn:   []string{
				"172.16.8.0/21",
			},
			RightIn:  []string{
				"172.16.10.0/24",
				"172.16.11.0/24",
				"172.16.12.0/24",
				"172.16.13.0/24",
				"172.16.14.0/24",
				"172.16.9.0/24",
				"172.16.8.0/24",
				"172.16.15.0/24",
			},
			LeftOut:  []string{},
			BothOut:  []string{
				"172.16.8.0/21",
			},
			RightOut: []string{},
			Status:   0,
			Error:    false,
		},
		// IPv6 tests
		{
			TestName: "#12",
			LeftIn:   []string{
				"fd00::/8",
			},
			RightIn:  []string{},
			LeftOut:  []string{
				"fd00::/8",
			},
			BothOut:  []string{},
			RightOut: []string{},
			Status:   2,
			Error:    false,
		},
		{
			TestName: "#13",
			LeftIn:   []string{},
			RightIn:  []string{
				"fd00::/8",
			},
			LeftOut:  []string{},
			BothOut:  []string{},
			RightOut: []string{
				"fd00::/8",
			},
			Status:   3,
			Error:    false,
		},
		{
			TestName: "#14",
			LeftIn:   []string{
				"2001:db8:0:2::/64",
				"2001:db8:0:3::/64",
			},
			RightIn:  []string{
				"2001:db8:0:3::/64",
				"2001:db8:0:2::/64",
			},
			LeftOut:  []string{},
			BothOut:  []string{
				"2001:db8:0:2::/63",
			},
			RightOut: []string{},
			Status:   0,
			Error:    false,
		},
		{
			TestName: "#15",
			LeftIn:   []string{
				"2001:db8:0:2::/64",
				"2001:db8:0:3::/64",
			},
			RightIn:  []string{
				"2001:db8:0:2::/64",
			},
			LeftOut:  []string{
				"2001:db8:0:3::/64",
			},
			BothOut:  []string{
				"2001:db8:0:2::/64",
			},
			RightOut: []string{},
			Status:   2,
			Error:    false,
		},
		{
			TestName: "#16",
			LeftIn:   []string{
				"fc00::/8",
				"2001:db8:0:2::/64",
			},
			RightIn:  []string{
				"2001:db8:0:2:2::/80",
				"fd00:0:4711:a::/64",
			},
			LeftOut:  []string{
				"2001:db8:0:2::/79",
				"2001:db8:0:2:3::/80",
				"2001:db8:0:2:4::/78",
				"2001:db8:0:2:8::/77",
				"2001:db8:0:2:10::/76",
				"2001:db8:0:2:20::/75",
				"2001:db8:0:2:40::/74",
				"2001:db8:0:2:80::/73",
				"2001:db8:0:2:100::/72",
				"2001:db8:0:2:200::/71",
				"2001:db8:0:2:400::/70",
				"2001:db8:0:2:800::/69",
				"2001:db8:0:2:1000::/68",
				"2001:db8:0:2:2000::/67",
				"2001:db8:0:2:4000::/66",
				"2001:db8:0:2:8000::/65",
				"fc00::/8",
			},
			BothOut:  []string{
				"2001:db8:0:2:2::/80",
			},
			RightOut: []string{
				"fd00:0:4711:a::/64",
			},
			Status:   1,
			Error:    false,
		},
		{
			TestName: "#17",
			LeftIn:   []string{
				"2001:db8:0:4::/64",
				"2001:db8:0:3::/64",
			},
			RightIn:  []string{
				"2001:db8:0:2::/63",
				"2001:db8:0:4::/63",
			},
			LeftOut:  []string{},
			BothOut:  []string{
				"2001:db8:0:3::/64",
				"2001:db8:0:4::/64",
			},
			RightOut: []string{
				"2001:db8:0:2::/64",
				"2001:db8:0:5::/64",
			},
			Status:   3,
			Error:    false,
		},
		{
			TestName: "#18",
			LeftIn:   []string{
				"fd00:0:4711:8::/61",
			},
			RightIn:  []string{
				"fd00:0:4711:a::/64",
				"fd00:0:4711:b::/64",
				"fd00:0:4711:c::/64",
				"fd00:0:4711:d::/64",
				"fd00:0:4711:e::/64",
				"fd00:0:4711:8::/64",
				"fd00:0:4711:9::/64",
				"fd00:0:4711:f::/64",
			},
			LeftOut:  []string{},
			BothOut:  []string{
				"fd00:0:4711:8::/61",
			},
			RightOut: []string{},
			Status:   0,
			Error:    false,
		},
		// Mixed blocks
		{
			TestName: "#19",
			LeftIn:   []string{
				"fd00::/8",
			},
			RightIn:  []string{
				"10.0.0.0/8",
			},
			LeftOut:  []string{
				"fd00::/8",
			},
			BothOut:  []string{},
			RightOut: []string{
				"10.0.0.0/8",
			},
			Status:   4,
			Error:    false,
		},
		{
			TestName: "#20",
			LeftIn:   []string{
				"fd00::/8",
				"10.0.0.0/8",
			},
			RightIn:  []string{
				"10.0.0.0/9",
			},
			LeftOut:  []string{
				"10.128.0.0/9",
				"fd00::/8",
			},
			BothOut:  []string{
				"10.0.0.0/9",
			},
			RightOut: []string{},
			Status:   2,
			Error:    false,
		},
		{
			TestName: "#21",
			LeftIn:   []string{
				"10.0.0.0/8",
				"2001:db8:0:2::/64",
				"2001:db8:0:3::/64",
				"192.0.128.0/24",
				"192.0.129.0/24",
			},
			RightIn:  []string{
				"fd00::/8",
				"192.0.128.0/24",
				"2001:db8:0:3::/64",
			},
			LeftOut:  []string{
				"10.0.0.0/8",
				"192.0.129.0/24",
				"2001:db8:0:2::/64",
			},
			BothOut:  []string{
				"192.0.128.0/24",
				"2001:db8:0:3::/64",
			},
			RightOut: []string{
				"fd00::/8",
			},
			Status:   1,
			Error:    false,
		},
	}

	for _, testCase := range testCases {
		leftout, bothout, rightout, status, err := DiffCIDRs(testCase.LeftIn, testCase.RightIn)
		//t.Errorf("\nleftout: %#v\nbothout: %#v\nrightout: %#v\nstatus: %#v", leftout, bothout, rightout, status)
		if err != nil {
			if !testCase.Error {
				t.Errorf("DiffCIDRs - %s : (%#v, %#v) failed with unexpected error: %s", testCase.TestName, testCase.LeftIn, testCase.RightIn, err.Error())
			} else if testCase.Status != status {
				t.Errorf("DiffCIDRs - %s : (%#v, %#v) mismatching status: (expected: %v, got: %v), with error: %s", testCase.TestName, testCase.LeftIn, testCase.RightIn, testCase.Status, status, err.Error())
			}
		} else if testCase.Error {
			t.Errorf("DiffCIDRs - %s : (%#v, %#v) didn't return with error as expected", testCase.TestName, testCase.LeftIn, testCase.RightIn)
		}

		if !reflect.DeepEqual(testCase.LeftOut, leftout) {
			t.Errorf("DiffCIDRs - %s : (%#v, %#v)\nexpected left result: %#v, got: %#v", testCase.TestName, testCase.LeftIn, testCase.RightIn, testCase.LeftOut, leftout)
		}
		if !reflect.DeepEqual(testCase.BothOut, bothout) {
			t.Errorf("DiffCIDRs - %s : (%#v, %#v)\nexpected in both result: %#v, got: %#v", testCase.TestName, testCase.LeftIn, testCase.RightIn, testCase.BothOut, bothout)
		}
		if !reflect.DeepEqual(testCase.RightOut, rightout) {
			t.Errorf("DiffCIDRs - %s : (%#v, %#v)\nexpected right result: %#v, got: %#v", testCase.TestName, testCase.LeftIn, testCase.RightIn, testCase.RightOut, rightout)
		}
		if testCase.Status != status {
			t.Errorf("DiffCIDRs - %s : (%#v, %#v) mismatching status, expected: %v, got: %v", testCase.TestName, testCase.LeftIn, testCase.RightIn, testCase.Status, status)
		}
	}
}
