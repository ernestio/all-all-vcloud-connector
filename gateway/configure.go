/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package gateway

import (
	"github.com/ernestio/all-all-vcloud-connector/helpers"
	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/config"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// Configure : configure an edge gateway
func (g *Gateway) Configure() error {
	var gateway *models.EdgeGateway

	cfg := config.New(g.Credentials.VCloudURL, "27.0").WithCredentials(g.Credentials.Username, g.Credentials.Password)
	vcloud := client.New(cfg)

	err := vcloud.Authenticate()
	if err != nil {
		return err
	}

	if g.ID != "" {
		gateway, err = vcloud.Gateways.Get(g.ID)
	} else {
		gateway, err = helpers.GatewayByName(vcloud, g.Name)
	}

	if err != nil {
		return err
	}

	g.UpdateProviderType(gateway)

	task, err := vcloud.Gateways.Update(gateway)
	if err != nil {
		return err
	}

	err = vcloud.Tasks.Wait(task)
	if err != nil {
		return err
	}

	g.ConvertProviderType(gateway)

	return nil
}
