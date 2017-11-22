/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helpers

import (
	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// OrgByName ...
func OrgByName(vcloud *client.Client, name string) (*models.Org, error) {
	orgs, err := vcloud.Orgs.List()
	if err != nil {
		return nil, err
	}

	orgID := orgs.ByName(name)[0].ID()

	return vcloud.Orgs.Get(orgID)
}
