/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package instance

import (
	"errors"

	"github.com/ernestio/all-all-vcloud-connector/helpers"
	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/config"
)

// Create : create a vapp/vm instance
func (i *Instance) Create() error {
	cfg := config.New(i.Credentials.VCloudURL, "27.0").WithCredentials(i.Credentials.Username, i.Credentials.Password)
	vcloud := client.New(cfg)

	err := vcloud.Authenticate()
	if err != nil {
		return err
	}

	vdc, err := helpers.VdcByName(vcloud, cfg.Org(), i.Credentials.Vdc)
	if err != nil {
		return err
	}

	nwref := vdc.NetworkRefs().Get(i.Network)
	if nwref == nil {
		return errors.New("could not find network: " + i.Network)
	}

	template, err := helpers.CatalogItemByName(vcloud, cfg.Org(), i.Catalog, i.Image)
	if err != nil {
		return err
	}

	vr := i.CreateProviderRequest()
	vr.SetNetwork(nwref.Name, "bridged", nwref.Href)
	vr.SetSource(i.Image, template.Entity.Href)

	vapp, err := vcloud.VApps.Create(vdc.GetID(), vr)
	if err != nil {
		return err
	}

	for _, task := range vapp.GetTasks() {
		err = vcloud.Tasks.Wait(&task)
		if err != nil {
			return err
		}
	}

	i.ConvertProviderType(vapp)

	return i.Update()
}
