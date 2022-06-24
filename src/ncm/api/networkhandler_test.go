// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package api

import (
	"encoding/json"
	netintents "gitlab.com/project-emco/core/emco-base/src/ncm/pkg/networkintents"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/validation"
	"strings"
	"testing"
)

type validationTest struct {
	valid bool
	json  string
}

func validateVirtualNetwork(t *testing.T, test validationTest) {
	var p netintents.Network
	err := json.NewDecoder(strings.NewReader(test.json)).Decode(&p)
	err, _ = validation.ValidateJsonSchemaData("../json-schemas/virtual-network.json", p)
	if test.valid && err != nil {
		t.Errorf("Valid JSON failed to validate against schema json: %s\nerr: %s", test.json, err)
	} else if !test.valid && err == nil {
		t.Errorf("Invalid JSON validated against schema %s", test.json)
	}
}

func TestVirtualNetworkValidation(t *testing.T) {

	t.Run("Require at least one of IPv4 and IPv6 subnet", func(t *testing.T) {
		tests := []validationTest{
			{
				valid: false,
				json: `
{
  "metadata": {
    "name": "name",
    "description": "description",
    "userData1": "user data 1",
    "userData2": "user data 2"
  },
  "spec": {
      "cniType": "ovn4nfv"
  }
}`,
			},
			{
				valid: true,
				json: `
{
  "metadata": {
    "name": "name",
    "description": "description",
    "userData1": "user data 1",
    "userData2": "user data 2"
  },
  "spec": {
      "cniType": "ovn4nfv",
      "ipv4Subnets": [
          {
              "subnet": "192.168.20.0/24",
              "name": "subnet4",
              "gateway":  "192.168.20.100"
          }
      ]
  }
}`,
			},
			{
				valid: true,
				json: `
{
  "metadata": {
    "name": "name",
    "description": "description",
    "userData1": "user data 1",
    "userData2": "user data 2"
  },
  "spec": {
      "cniType": "ovn4nfv",
      "ipv6Subnets": [
          {
              "subnet": "2001:db8:42:0::/56",
              "name": "subnet6",
              "gateway":  "2001:db8:42:0::1"
          }
      ]
  }
}`,
			},
			{
				valid: true,
				json: `
{
  "metadata": {
    "name": "name",
    "description": "description",
    "userData1": "user data 1",
    "userData2": "user data 2"
  },
  "spec": {
      "cniType": "ovn4nfv",
      "ipv4Subnets": [
          {
              "subnet": "192.168.20.0/24",
              "name": "subnet4",
              "gateway":  "192.168.20.100"
          }
      ],
      "ipv6Subnets": [
          {
              "subnet": "2001:db8:42:0::/56",
              "name": "subnet6",
              "gateway":  "2001:db8:42:0::1"
          }
      ]
  }
}`,
			},
		}
		for _, test := range tests {
			validateVirtualNetwork(t, test)
		}
	})

	t.Run("Require cniType", func(t *testing.T) {
		tests := []validationTest{
			{
				valid: false,
				json: `
{
  "metatest": {
    "name": "name",
  },
  "spec": {
      "ipv4Subnets": [
          {
              "subnet": "192.168.20.0/24",
              "name": "subnet4",
              "gateway":  "192.168.20.100"
          }
      ]
  }
}`,
			},
			{
				valid: false,
				json: `
{
  "metadata": {
    "name": "name",
  },
  "spec": {
      "ipv6Subnets": [
          {
              "subnet": "2001:db8:42:0::/56",
              "name": "subnet6",
              "gateway":  "2001:db8:42:0::1"
          }
      ]
  }
}`,
			},
		}
		for _, test := range tests {
			validateVirtualNetwork(t, test)
		}
	})
}
