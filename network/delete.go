/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package network

import (
	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/config"
)

// Delete : delete an org vdc network
func (n *Network) Delete() error {
	cfg := config.New(n.Credentials.VCloudURL, "27.0").WithCredentials(n.Credentials.Username, n.Credentials.GetPassword())
	vcloud := client.New(cfg)

	err := vcloud.Authenticate()
	if err != nil {
		return err
	}

	task, err := vcloud.Networks.Delete(n.ID)
	if err != nil {
		return err
	}

	return vcloud.Tasks.Wait(task)
}
