/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package instance

import (
	"fmt"

	"github.com/r3labs/vcloud-go-sdk/client"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// Template : template details
type Template struct {
	ID      string
	Name    string
	Catalog string
}

// TemplateCache : stores all found vapp templates for later lookups
type TemplateCache map[string]*Template

// NewTemplateCache : creates a new template cache
func NewTemplateCache() *TemplateCache {
	tc := make(TemplateCache)
	return &tc
}

// Get : get a templates information
func (t *TemplateCache) Get(vcloud *client.Client, imageID string) (*Template, error) {
	fmt.Println(t)
	if t == nil {
		fmt.Println("map is nil!")
		tc := make(TemplateCache)
		t = &tc
	}

	if (*t)[imageID] == nil {
		q, err := vcloud.Queries.RecordsFilter(models.QueryVAppTemplate, "id=="+imageID, "1")
		if err != nil {
			return nil, err
		}

		if len(q.VAppTemplates) < 1 {
			return nil, fmt.Errorf("could not find vapp template image: (%s)", imageID)
		}

		result := q.VAppTemplates[0]

		(*t)[imageID] = &Template{
			ID:      imageID,
			Name:    result.Name,
			Catalog: result.Catalog,
		}
	}

	return (*t)[imageID], nil
}

// GetDetails : returns the name and catalog of a vapp template
func (t *TemplateCache) GetDetails(vcloud *client.Client, imageID string) (string, string, error) {
	c, err := t.Get(vcloud, imageID)
	if err != nil {
		return "", "", err
	}

	return c.Name, c.Catalog, nil
}
