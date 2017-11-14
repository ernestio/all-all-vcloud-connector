/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package instance

import (
	"fmt"

	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/config"
)

// Update : update a vapp/vm instance
func (i *Instance) Update() error {
	cfg := config.New(i.Credentials.VCloudURL, "27.0").WithCredentials(i.Credentials.Username, i.Credentials.Password)
	vcloud := client.New(cfg)

	err := vcloud.Authenticate()
	if err != nil {
		return err
	}

	vapp, err := vcloud.VApps.Get(i.ID)
	if err != nil {
		return err
	}

	fmt.Println(vapp)

	return nil
}
