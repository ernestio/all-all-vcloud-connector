/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package instance

import (
	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/config"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// Update : update a vapp/vm instance
func (i *Instance) Update() error {
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

	i.UpdateProviderType(vapp)

	vm := vapp.Children.Vms[0]

	// check power state!
	if vm.Status != "8" {
		task, err = vcloud.Vms.PowerOff(vm.GetID())
		if err != nil {
			return err
		}

		err = vcloud.Tasks.Wait(task)
		if err != nil {
			return err
		}
	}

	task, err = vcloud.Vms.Update(vm)
	if err != nil {
		return err
	}

	err = vcloud.Tasks.Wait(task)
	if err != nil {
		return err
	}

	metadata, err := vcloud.VApps.GetMetadata(vapp.GetID())
	if err != nil {
		return err
	}

	for k, v := range i.Tags {
		metadata.Remove(k)
		metadata.Add(k, v)
	}

	metadata.RemoveDefaults()

	task, err = vcloud.VApps.UpdateMetadata(vapp.GetID(), metadata)
	if err != nil {
		return err
	}

	err = vcloud.Tasks.Wait(task)
	if err != nil {
		return err
	}

	task, err = vcloud.Vms.PowerOn(vm.GetID())
	if err != nil {
		return err
	}

	err = vcloud.Tasks.Wait(task)
	if err != nil {
		return err
	}

	i.ConvertProviderType(vapp)

	return nil
}
