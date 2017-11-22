/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helpers

import (
	"errors"
	"fmt"

	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// GatewayByName ...
func GatewayByName(vcloud *client.Client, name string) (*models.EdgeGateway, error) {
	filter := fmt.Sprintf("name==%s", name)

	records, err := vcloud.Queries.RecordsFilter(models.QueryEdgeGateway, filter, "1")
	if err != nil {
		return nil, err
	}

	if len(records.EdgeGatewayRecords) != 1 {
		return nil, errors.New("could not find edge gateway: " + name)
	}

	id := records.EdgeGatewayRecords[0].ID()

	return vcloud.Gateways.Get(id)
}
