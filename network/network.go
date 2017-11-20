/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package network

import (
	"errors"

	"github.com/ernestio/all-all-vcloud-connector/base"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// Network : Mapping of a network component
type Network struct {
	base.DefaultFields
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Subnet        string   `json:"range"`
	Netmask       string   `json:"netmask"`
	StartAddress  string   `json:"start_address"`
	EndAddress    string   `json:"end_address"`
	Gateway       string   `json:"gateway"`
	DNS           []string `json:"dns"`
	EdgeGateway   string   `json:"edge_gateway"`
	EdgeGatewayID string   `json:"edge_gateway_id"`
}

// SetState : sets the networks state
func (n *Network) SetState(state string) {
	n.State = state
}

// SetError : sets the networks error message
func (n *Network) SetError(err error) {
	n.ErrorMessage = err.Error()
}

// GetCredentials ...
func (n *Network) GetCredentials() *base.Credentials {
	return n.Credentials
}

// Find : find an org vdc network
func (n *Network) Find() error {
	return errors.New("not implemented")
}

// ConvertErnestType : converts the ernest network to an org vdc network type
func (n *Network) ConvertErnestType() *models.Network {
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

// ConvertProviderType : converts the org vdc network to an ernest network
func (n *Network) ConvertProviderType(nw *models.Network) {
	n.ProviderType = "vcloud"
	n.ComponentType = "network"
	n.ComponentID = "network::" + n.Name

	n.ID = nw.GetID()
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
