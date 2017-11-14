/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package instance

import (
	"log"
	"strconv"

	"github.com/ernestio/all-all-vcloud-connector/base"
	"github.com/r3labs/vcloud-go-sdk/models"
)

// Instance : Mapping of an instance component
type Instance struct {
	base.DefaultFields
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Hostname      string            `json:"hostname"`
	Catalog       string            `json:"reference_catalog"`
	Image         string            `json:"reference_image"`
	Cpus          int               `json:"cpus"`
	Memory        int               `json:"ram"`
	Network       string            `json:"network"`
	IP            string            `json:"ip"`
	Disks         []Disk            `json:"disks"`
	ShellCommands []string          `json:"shell_commands"`
	Tags          map[string]string `json:"tags"`
}

// CreateProviderRequest : converts the ernest instance into a vapp creation request
func (i *Instance) CreateProviderRequest() *models.InstantiateVAppParams {
	return &models.InstantiateVAppParams{
		Name:        i.Name,
		Description: i.Name,
		AcceptEULAs: true,
		Deploy:      true,
		PowerOn:     false,
	}
}

// UpdateProviderType : updates the provider type with values from the ernest instance
func (i *Instance) UpdateProviderType(vapp *models.VApp) {
	vm := vapp.Vms()[0]

	vm.Name = i.Name
	vhs := vm.VirtualHardwareSection
	con := vhs.GetDiskController()

	vhs.SetCPU(i.Cpus)
	vhs.SetRAM(i.Memory)

	if len(i.Disks) < 1 {
		return
	}

	for _, disk := range i.Disks {
		// Don't set the size of the boot drive
		vhs.RemoveDisk(con.InstanceID.Value, strconv.Itoa(disk.ID))
		vhs.AddDisk(con.InstanceID.Value, disk.Size)
	}
}

// ConvertProviderType : converts the org vdc network to an ernest network
func (i *Instance) ConvertProviderType(vapp *models.VApp) {
	vm := vapp.Vms()[0]
	vhs := vm.VirtualHardwareSection
	con := vhs.GetDiskController()

	i.ProviderType = "vcloud"
	i.ComponentType = "instance"
	i.ComponentID = "instance::" + vapp.Name

	i.ID = vapp.ID
	i.Name = vapp.Name
	i.Cpus = vhs.GetCPU()
	i.Memory = vhs.GetRAM()

	for _, disk := range vhs.Items.ByParent(con.InstanceID.Value) {
		var size int

		if disk.InstanceID.Value == "0" {
			continue
		}

		if disk.VirtualQuantityUnits.Value == "byte" {
			i64, err := strconv.ParseInt(disk.VirtualQuantity.Value, 10, 64)
			if err != nil {
				log.Println("could not get disk size (int64)")
				return
			}

			size = int(i64 / 1048576)
		} else {
			size, _ = strconv.Atoi(disk.VirtualQuantity.Value)
		}

		id, _ := strconv.Atoi(disk.InstanceID.Value)

		i.Disks = append(i.Disks, Disk{
			ID:   id,
			Size: size,
		})
	}
}
