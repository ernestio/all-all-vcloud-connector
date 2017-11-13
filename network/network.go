/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package network

import (
	"github.com/ernestio/all-all-vcloud-connector/credentials"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// Network : Mapping of a network component
type Network struct {
	ProviderType  string                   `json:"_provider"`
	ComponentType string                   `json:"_component"`
	ComponentID   string                   `json:"_component_id"`
	State         string                   `json:"_state"`
	Action        string                   `json:"_action"`
	ID            string                   `json:"id"`
	Name          string                   `json:"name"`
	Subnet        string                   `json:"range"`
	Netmask       string                   `json:"netmask"`
	StartAddress  string                   `json:"start_address"`
	EndAddress    string                   `json:"end_address"`
	Gateway       string                   `json:"gateway"`
	DNS           []string                 `json:"dns"`
	EdgeGateway   string                   `json:"edge_gateway"`
	EdgeGatewayID string                   `json:"edge_gateway_id"`
	Credentials   *credentials.Credentials `json:"credentials"`
	Service       string                   `json:"service"`
}

// ToVcloudType : converts the ernest network to an org vdc network type
func (n *Network) ToVcloudType() *models.Network {
	nw := models.Network{
		Name: n.Name,
	}

	nw.SetNetmask(n.Netmask)
	nw.SetStartAddress(n.StartAddress)
	nw.SetEndAddress(n.EndAddress)
	nw.SetGateway(n.Gateway)
	nw.SetFenceMode("natRouted")
	nw.SetIsEnabled(true)

	if len(n.DNS) > 1 {
		nw.SetDNS1(n.DNS[0])
	}

	if len(n.DNS) > 2 {
		nw.SetDNS2(n.DNS[1])
	}

	return &nw
}

// FromVcloudType : converts the org vdc network to a ernest network
func (n *Network) FromVcloudType(nw *models.Network) {
	n.ID = nw.ID
	n.Name = nw.Name
	n.EdgeGateway = nw.EdgeGatewayName()
	n.EdgeGatewayID = nw.EdgeGatewayID()
	n.Netmask = nw.Netmask()
	n.Gateway = nw.Gateway()
	n.StartAddress = nw.StartAddress()
	n.EndAddress = nw.EndAddress()

	n.DNS = []string{
		nw.DNS1(),
		nw.DNS2(),
	}
}
