/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package network

import (
	"github.com/ernestio/all-all-vcloud-connector/helpers"
	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/config"
)

// Create : create an org vdc network
func (n *Network) Create() error {
	cfg := config.New(n.Credentials.VCloudURL, "27.0").WithCredentials(n.Credentials.Username, n.Credentials.Password)
	vcloud := client.New(cfg)

	err := vcloud.Authenticate()
	if err != nil {
		return err
	}

	vdc, err := helpers.VdcByName(vcloud, cfg.Org(), n.Credentials.Vdc)
	if err != nil {
		return err
	}

	gateway, err := helpers.GatewayByName(vcloud, n.EdgeGateway)
	if err != nil {
		return err
	}

	nw := n.ConvertErnestType()
	nw.SetEdgeGateway(gateway.Href, gateway.Name)

	err = vcloud.Networks.Create(vdc.GetID(), nw)
	if err != nil {
		return err
	}

	for _, task := range nw.GetTasks() {
		err = vcloud.Tasks.Wait(&task)
		if err != nil {
			return err
		}
	}

	n.ConvertProviderType(nw)

	return nil
}
