/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"

	"github.com/ernestio/all-all-vcloud-connector/gateway"
	"github.com/ernestio/all-all-vcloud-connector/instance"
	"github.com/ernestio/all-all-vcloud-connector/network"
	"github.com/mitchellh/mapstructure"
)

// Event : defines the interface that all events will conform to
type Event interface {
	SetError(error)
	SetState(string)
	Create() error
	Update() error
	Delete() error
	Find() error
}

func event(subject string, data []byte) (Event, error) {
	var e Event
	var m map[string]interface{}

	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	switch subject {
	case "router.find.vcloud":
		e = &gateway.Collection{}
	case "router.create.vcloud", "router.update.vcloud":
		e = &gateway.Gateway{}
	case "network.find.vcloud":
		e = &network.Collection{}
	case "network.create.vcloud", "network.update.vcloud", "network.delete.vcloud":
		e = &network.Network{}
	case "instance.find.vcloud":
		e = &instance.Collection{}
	case "instance.create.vcloud", "instance.update.vcloud", "instance.delete.vcloud":
		e = &instance.Instance{}
	}

	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   e,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(m)
	if err != nil {
		return nil, err
	}

	return e, nil
}
