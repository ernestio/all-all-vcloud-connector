/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helpers

import (
	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// VdcByName ...
func VdcByName(vcloud *client.Client, org, name string) (*models.Vdc, error) {
	o, err := OrgByName(vcloud, org)
	if err != nil {
		return nil, err
	}

	vdcLink := o.Links.ByType(models.TypesVdc).ByName(name)

	return vcloud.Vdcs.Get(vdcLink[0].ID())
}
