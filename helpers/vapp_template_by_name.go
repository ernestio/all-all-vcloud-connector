/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helpers

import (
	"errors"

	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// CatalogItemByName ...
func CatalogItemByName(vcloud *client.Client, org, catalog, name string) (*models.CatalogItem, error) {
	o, err := OrgByName(vcloud, org)
	if err != nil {
		return nil, err
	}

	cref := o.CatalogRefs().Get(catalog)
	if cref == nil {
		return nil, errors.New("could not find template catalog: " + catalog)
	}

	c, err := vcloud.Catalogs.Get(cref.ID())
	if err != nil {
		return nil, err
	}

	item := c.Items().ByName(name)

	return vcloud.Catalogs.GetItem(item.GetID())
}
