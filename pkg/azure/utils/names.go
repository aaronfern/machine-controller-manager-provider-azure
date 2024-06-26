// Copyright 2023 SAP SE or an SAP affiliate company
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"

	"github.com/gardener/machine-controller-manager-provider-azure/pkg/azure/api"
)

const (
	// NICSuffix is the suffix for NIC names.
	NICSuffix = "-nic"
	// OSDiskSuffix is the suffix for OSDisk names.
	OSDiskSuffix = "-os-disk"
	//DataDiskSuffix is the suffix for Data disk names.
	DataDiskSuffix = "-data-disk"
	// AzureCSIDriverName is the name of the CSI driver name for Azure provider
	AzureCSIDriverName = "disk.csi.azure.com"
)

// CreateNICName creates a NIC name given a VM name
func CreateNICName(vmName string) string {
	return fmt.Sprintf("%s%s", vmName, NICSuffix)
}

// ExtractVMNameFromNICName extracts VM Name from NIC name
func ExtractVMNameFromNICName(nicName string) string {
	return nicName[:len(nicName)-len(NICSuffix)]
}

// ExtractVMNameFromOSDiskName extracts VM name from OSDisk name
func ExtractVMNameFromOSDiskName(osDiskName string) string {
	return osDiskName[:len(osDiskName)-len(OSDiskSuffix)]
}

// CreateOSDiskName creates OSDisk name from VM name
func CreateOSDiskName(vmName string) string {
	return fmt.Sprintf("%s%s", vmName, OSDiskSuffix)
}

// CreateDataDiskName creates a name for a DataDisk using VM name and data disk name specified in the provider Spec
func CreateDataDiskName(vmName string, dataDisk api.AzureDataDisk) string {
	prefix := vmName
	suffix := GetDataDiskNameSuffix(dataDisk)
	return fmt.Sprintf("%s%s", prefix, suffix)
}

// GetDataDiskNameSuffix creates the suffix based on an optional data disk name and required lun fields.
func GetDataDiskNameSuffix(dataDisk api.AzureDataDisk) string {
	infix := getDataDiskInfix(dataDisk)
	return fmt.Sprintf("-%s%s", infix, DataDiskSuffix)
}

func getDataDiskInfix(dataDisk api.AzureDataDisk) string {
	name := dataDisk.Name
	if IsEmptyString(name) {
		return fmt.Sprintf("%d", *dataDisk.Lun)
	}
	return fmt.Sprintf("%s-%d", name, *dataDisk.Lun)
}
