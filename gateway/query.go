/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package gateway

import (
	"errors"

	"github.com/ernestio/all-all-vcloud-connector/base"
	"github.com/ernestio/all-all-vcloud-connector/helpers"
	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/config"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// Collection ...
type Collection struct {
	base.DefaultFields
	Components []*Gateway `json:"components,omitempty"`
}

// SetState : sets the collections state
func (c *Collection) SetState(state string) {
	c.State = state
}

// SetError : sets the collections error message
func (c *Collection) SetError(err error) {
	c.ErrorMessage = err.Error()
}

// GetCredentials ...
func (c *Collection) GetCredentials() *base.Credentials {
	return c.Credentials
}

// Create ...
func (c *Collection) Create() error {
	return errors.New("not implemented")
}

// Update ...
func (c *Collection) Update() error {
	return errors.New("not implemented")
}

// Delete ...
func (c *Collection) Delete() error {
	return errors.New("not implemented")
}

// Find : finds all edge gateways related to a vdc
func (c *Collection) Find() error {
	cfg := config.New(c.Credentials.VCloudURL, "27.0").WithCredentials(c.Credentials.Username, c.Credentials.GetPassword())
	vcloud := client.New(cfg)

	err := vcloud.Authenticate()
	if err != nil {
		return err
	}

	vdc, err := helpers.VdcByName(vcloud, cfg.Org(), c.Credentials.Vdc)
	if err != nil {
		return err
	}

	records, err := vcloud.Queries.RecordsFilter(models.QueryEdgeGateway, "vdc=="+vdc.Href, "1")
	if err != nil {
		return err
	}

	for _, gr := range records.EdgeGatewayRecords {
		var g Gateway

		gw, err := vcloud.Gateways.Get(gr.ID())
		if err != nil {
			return err
		}

		g.ConvertProviderType(gw)

		c.Components = append(c.Components, &g)
	}

	return nil
}
