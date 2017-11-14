/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package network

import (
	"github.com/ernestio/all-all-vcloud-connector/base"
	"github.com/ernestio/all-all-vcloud-connector/helpers"
	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/config"
)

// Collection ...
type Collection struct {
	base.DefaultFields
	Components []*Network `json:"components,omitempty"`
}

// Find : finds all networks related to a vdc
func (c *Collection) Find() error {
	cfg := config.New(c.Credentials.VCloudURL, "27.0").WithCredentials(c.Credentials.Username, c.Credentials.Password)
	vcloud := client.New(cfg)

	err := vcloud.Authenticate()
	if err != nil {
		return err
	}

	vdc, err := helpers.VdcByName(vcloud, cfg.Org(), c.Credentials.Vdc)
	if err != nil {
		return err
	}

	for _, nr := range vdc.NetworkRefs() {
		var n Network

		nw, err := vcloud.Networks.Get(nr.ID())
		if err != nil {
			return err
		}

		n.ConvertProviderType(nw)

		c.Components = append(c.Components, &n)
	}

	return nil
}
