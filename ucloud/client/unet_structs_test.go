package client

import (
	"testing"
)

func TestSecurityGroupRule(t *testing.T) {
	var rule interface{} = SecurityGroupRule{
		ProtocolType: "TCP",
		DstPort:      "3306",
		SrcIP:        "0.0.0.0/0",
		RuleAction:   "DROP",
		Priority:     50,
	}

	p, ok := rule.(Parameterizable)
	if !ok {
		t.Fatal("SecurityGroupRule is not Parameterizable")
	}

	str, err := p.Parameterize()
	if err != nil {
		t.Fatal("Parameterize failed: ", err)
	}

	if str != "TCP|3306|0.0.0.0/0|DROP|50" {
		t.Error("Invalid parameterized string: ", str)
	}
}

func TestCreateSecurityGroupRequest(t *testing.T) {
	req := &CreateSecurityGroupRequest{
		Rule: []SecurityGroupRule{
			SecurityGroupRule{
				ProtocolType: "TCP",
				DstPort:      "3306",
				SrcIP:        "0.0.0.0/0",
				RuleAction:   "DROP",
				Priority:     50,
			},
			SecurityGroupRule{
				ProtocolType: "UDP",
				DstPort:      "53",
				SrcIP:        "0.0.0.0/0",
				RuleAction:   "ACCEPT",
				Priority:     50,
			},
		},
	}

	params, err := BuildParams(req)
	if err != nil {
		t.Fatal("Failed to build params: ", err)
	}

	cases := []struct{ Arg, Expectation string }{
		{"Rule.0", "TCP|3306|0.0.0.0/0|DROP|50"},
		{"Rule.1", "UDP|53|0.0.0.0/0|ACCEPT|50"},
	}

	for _, tc := range cases {
		real := params.Get(tc.Arg)
		if real != tc.Expectation {
			t.Errorf("Expect %s to be %s but got: %s", tc.Arg, tc.Expectation, real)
		}
	}
}
