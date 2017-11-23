/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package instance

import (
	"time"

	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/config"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// Delete : delete a vapp/vm instance
func (i *Instance) Delete() error {
	var task *models.Task

	cfg := config.New(i.Credentials.VCloudURL, "27.0").WithCredentials(i.Credentials.Username, i.Credentials.GetPassword())
	vcloud := client.New(cfg)

	err := vcloud.Authenticate()
	if err != nil {
		return err
	}

	vapp, err := vcloud.VApps.Get(i.ID)
	if err != nil {
		return err
	}

	if vapp.Deployed {
		task, err = vcloud.VApps.Undeploy(i.ID)
		if err != nil {
			return err
		}

		err = vcloud.Tasks.Wait(task)
		if err != nil {
			return err
		}
	}

	task, err = vcloud.VApps.Delete(i.ID)
	if err != nil {
		return err
	}

	err = vcloud.Tasks.Wait(task)
	if err != nil {
		return err
	}

	// vapp deletion does not register with edge gateway immendiately.
	time.Sleep(time.Second * 5)

	return nil
}
