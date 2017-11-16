/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package gateway

import (
	"errors"

	"github.com/ernestio/all-all-vcloud-connector/base"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// Gateway : mapping of a edge gateway component
type Gateway struct {
	base.DefaultFields
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	IP            string         `json:"ip"`
	NatRules      []NatRule      `json:"nat_rules"`
	FirewallRules []FirewallRule `json:"firewall_rules"`
}

// SetState : sets the edge gateways state
func (g *Gateway) SetState(state string) {
	g.State = state
}

// SetError : sets the edge gateways error message
func (g *Gateway) SetError(err error) {
	g.ErrorMessage = err.Error()
}

// Find : find an edge gateway
func (g *Gateway) Find() error {
	return errors.New("not implemented")
}

// UpdateProviderType : updates the provider type with values from the ernest gateway
func (g *Gateway) UpdateProviderType(gw *models.EdgeGateway) {
	gw.Nat().ClearRules()
	gw.Firewall().ClearRules()

	for _, r := range g.FirewallRules {
		var p models.FirewallProtocols

		gw.Firewall().AddRule(&models.FirewallRule{
			Description:          r.Name,
			SourceIP:             r.SourceIP,
			DestinationIP:        r.DestinationIP,
			SourcePortRange:      r.SourcePort,
			DestinationPortRange: r.DestinationPort,
			Protocols:            p.Set(r.Protocol),
			Policy:               "allow",
		})
	}

	for _, r := range g.NatRules {
		gw.Nat().AddRule(&models.NatRule{
			GatewayNatRule: &models.GatewayNatRule{
				Interface:      gw.DefaultInterface().Network,
				OriginalIP:     r.OriginIP,
				TranslatedIP:   r.TranslationIP,
				OriginalPort:   r.OriginPort,
				TranslatedPort: r.TranslationPort,
				Protocol:       r.Protocol,
			},
			Enabled: true,
		})
	}
}

// ConvertProviderType : converts the edge gateway to a ernest gateway
func (g *Gateway) ConvertProviderType(gw *models.EdgeGateway) {
	g.ProviderType = "vcloud"
	g.ComponentType = "router"
	g.ComponentID = "router::" + g.Name

	g.ID = gw.GetID()
	g.Name = gw.Name
	g.NatRules = make([]NatRule, 0)
	g.FirewallRules = make([]FirewallRule, 0)

	for _, r := range gw.Firewall().FirewallRules {
		g.FirewallRules = append(g.FirewallRules, FirewallRule{
			Name:            r.Description,
			SourceIP:        r.SourceIP,
			DestinationIP:   r.DestinationIP,
			SourcePort:      r.SourcePortRange,
			DestinationPort: r.DestinationPortRange,
			Protocol:        r.Protocols.Get(),
		})
	}

	for _, r := range gw.Nat().NatRules {
		g.NatRules = append(g.NatRules, NatRule{
			OriginIP:        r.GatewayNatRule.OriginalIP,
			TranslationIP:   r.GatewayNatRule.TranslatedIP,
			OriginPort:      r.GatewayNatRule.OriginalPort,
			TranslationPort: r.GatewayNatRule.TranslatedPort,
			Protocol:        r.GatewayNatRule.Protocol,
		})
	}
}
