/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package base

// DefaultFields ...
type DefaultFields struct {
	ProviderType  string       `json:"_provider"`
	ComponentType string       `json:"_component"`
	ComponentID   string       `json:"_component_id"`
	State         string       `json:"_state"`
	Action        string       `json:"_action"`
	Credentials   *Credentials `json:"_credentials"`
	ErrorMessage  string       `json:"error,omitempty"`
	Service       string       `json:"service"`
}
