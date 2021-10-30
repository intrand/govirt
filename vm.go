package main

import (
	"errors"

	ovirt "github.com/ovirt/go-ovirt"
)

func createVmFromTemplate(conn *ovirt.Connection, template ovirt.Template, domain ovirt.StorageDomain, vmName string, clusterName string) error {
	disksSlice, ok := template.DiskAttachments()
	if !ok {
		return errors.New("couldn't get list of disk attachments")
	}

	_disks, err := conn.FollowLink(disksSlice)
	if err != nil {
		return err
	}

	disks := _disks.(*ovirt.DiskAttachmentSlice).Slice()
	if len(disks) < 1 {
		return errors.New("couldn't find disks attached to the template")
	}

	disk, ok := disks[0].Disk()
	if !ok {
		return errors.New("couldn't get first disk")
	}

	templateId, ok := template.Id()
	if !ok {
		return errors.New("couldn't get template id")
	}

	newTemplate, err := ovirt.NewTemplateBuilder().Id(templateId).Build()
	if err != nil {
		return err
	}

	newCluster, err := ovirt.NewClusterBuilder().Name(clusterName).Build()
	if err != nil {
		return err
	}

	domainId, ok := domain.Id()
	if !ok {
		return errors.New("couldn't get storage domain id")
	}

	newDomain, err := ovirt.NewStorageDomainBuilder().Id(domainId).Build()
	if err != nil {
		return err
	}

	diskId, ok := disk.Id()
	if !ok {
		return errors.New("couldn't get disk id")
	}

	newDisk, err := ovirt.NewDiskBuilder().Id(diskId).Format(ovirt.DISKFORMAT_COW).StorageDomainsOfAny(newDomain).Build()
	if err != nil {
		return err
	}

	newAttachment, err := ovirt.NewDiskAttachmentBuilder().Disk(newDisk).Build()
	if err != nil {
		return err
	}

	newVm, err := ovirt.NewVmBuilder().Name(vmName).Cluster(newCluster).Template(newTemplate).DiskAttachmentsOfAny(newAttachment).Build()
	if err != nil {
		return err
	}

	vmsService := conn.SystemService().VmsService()
	_, err = vmsService.Add().Vm(newVm).Send()
	if err != nil {
		return err
	}
	return nil
}

func getVm(conn *ovirt.Connection, vmName string) (ovirt.Vm, error) {
	vm := ovirt.Vm{}
	vmsService := conn.SystemService().VmsService()

	res, err := vmsService.List().Search("name=" + vmName).Send()
	if err != nil {
		return vm, err
	}

	vms, ok := res.Vms()
	if !ok {
		return vm, errors.New("couldn't get virtual machines")
	}

	if len(vms.Slice()) == 1 {
		return *vms.Slice()[0], nil
	} else {
		return vm, errors.New("requested virtual machine not found")
	}
}

func deleteVm(conn *ovirt.Connection, vm ovirt.Vm) error {
	vmsService := conn.SystemService().VmsService()

	vmId, ok := vm.Id()
	if !ok {
		return errors.New("couldn't get virtual machine id")
	}

	vmService := vmsService.VmService(vmId)

	_, err := vmService.Remove().Send()
	if err != nil {
		return err
	}

	return nil
}
