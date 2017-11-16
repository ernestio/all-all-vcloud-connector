/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package instance

import (
	"errors"

	"github.com/ernestio/all-all-vcloud-connector/base"
	"github.com/ernestio/all-all-vcloud-connector/helpers"
	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/config"
)

// Collection ...
type Collection struct {
	base.DefaultFields
	Components []*Instance `json:"components,omitempty"`
}

// SetState : sets the collections state
func (c *Collection) SetState(state string) {
	c.State = state
}

// SetError : sets the collections error message
func (c *Collection) SetError(err error) {
	c.ErrorMessage = err.Error()
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

	for _, vappr := range vdc.VAppRefs() {
		var i Instance

		i.Tags = make(map[string]string)

		vapp, err := vcloud.VApps.Get(vappr.ID())
		if err != nil {
			return err
		}

		i.ConvertProviderType(vapp)

		metadata, err := vcloud.VApps.GetMetadata(vappr.ID())
		if err != nil {
			return err
		}

		for _, e := range metadata.Entries {
			i.Tags[e.Key] = e.TypedValue.Value
		}

		c.Components = append(c.Components, &i)
	}

	return nil
}
