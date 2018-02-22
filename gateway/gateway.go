/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package gateway

import (
	"errors"
	"strings"

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

// GetCredentials ...
func (g *Gateway) GetCredentials() *base.Credentials {
	return g.Credentials
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
			Policy:               r.Action,
			Enabled:              true,
		})
	}

	for _, r := range g.NatRules {
		gw.Nat().AddRule(&models.NatRule{
			RuleType: r.Type,
			Enabled:  true,
			GatewayNatRule: &models.GatewayNatRule{
				Interface:      gw.DefaultInterface().Network,
				OriginalIP:     r.OriginIP,
				TranslatedIP:   r.TranslationIP,
				OriginalPort:   r.OriginPort,
				TranslatedPort: r.TranslationPort,
				Protocol:       r.Protocol,
			},
		})
	}
}

// ConvertProviderType : converts the edge gateway to a ernest gateway
func (g *Gateway) ConvertProviderType(gw *models.EdgeGateway) {
	g.ProviderType = "vcloud"
	g.ComponentType = "router"
	g.ComponentID = "router::" + gw.Name

	g.ID = gw.GetID()
	g.Name = gw.Name
	g.NatRules = make([]NatRule, 0)
	g.FirewallRules = make([]FirewallRule, 0)

	if gw.Firewall() != nil {
		for _, r := range gw.Firewall().FirewallRules {
			g.FirewallRules = append(g.FirewallRules, FirewallRule{
				Name:            r.Description,
				SourceIP:        strings.ToLower(r.SourceIP),
				DestinationIP:   strings.ToLower(r.DestinationIP),
				SourcePort:      strings.ToLower(r.SourcePortRange),
				DestinationPort: strings.ToLower(r.DestinationPortRange),
				Protocol:        strings.ToLower(r.Protocols.Get()),
				Policy:          strings.ToLower(r.Action),
			})
		}
	}

	if gw.Nat() != nil {
		for _, r := range gw.Nat().NatRules {
			g.NatRules = append(g.NatRules, NatRule{
				Type:            r.RuleType,
				OriginIP:        strings.ToLower(r.GatewayNatRule.OriginalIP),
				TranslationIP:   strings.ToLower(r.GatewayNatRule.TranslatedIP),
				OriginPort:      r.GatewayNatRule.GetOriginalPort(),
				TranslationPort: r.GatewayNatRule.GetTranslatedPort(),
				Protocol:        r.GatewayNatRule.GetProtocol(),
			})
		}
	}
}
