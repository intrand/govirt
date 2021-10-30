package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// global
	app          = kingpin.New("govirt", "Manages ovirt").Author("intrand")
	cmd_username = app.Flag("username", "Username").Envar("govirt_username").Default("admin@internal").String()
	cmd_password = app.Flag("password", "Password").Envar("govirt_password").Required().String()
	cmd_url      = app.Flag("url", "URL to API (eg, https://1.2.3.4/ovirt-engine/api)").Envar("govirt_url").Required().String()

	// vm
	cmd_vm = app.Command("vm", "virtual machines")

	cmd_vm_create                  = cmd_vm.Command("create", "creates VM from template")
	cmd_vm_create_name             = cmd_vm_create.Flag("name", "name of vm to create").Envar("govirt_vm_create_name").Required().String()
	cmd_vm_create_cluster          = cmd_vm_create.Flag("cluster", "cluster in which to create vm").Envar("govirt_vm_create_cluster").Default("Default").String()
	cmd_vm_create_template         = cmd_vm_create.Flag("template", "template from which to create vm").Envar("govirt_vm_create_template").Required().String()
	cmd_vm_create_template_version = cmd_vm_create.Flag("template-version", "version of template from which to create vm").Envar("govirt_vm_create_template_version").Required().Int64()
	cmd_vm_create_datastore        = cmd_vm_create.Flag("datastore", "storage domain in which to create vm").Envar("govirt_vm_create_datastore").Required().String()

	cmd_vm_rm      = cmd_vm.Command("rm", "remove a vm")
	cmd_vm_rm_name = cmd_vm_rm.Flag("name", "name of vm to remove").Envar("govirt_vm_rm_name").Required().String()
	cmd_vm_rm_yes  = cmd_vm_rm.Flag("yes", "remove vm from ovirt").Envar("govirt_vm_rm_yes").Default("false").Bool()

	cmd_vm_summary = cmd_vm.Command("summary", "output summary of virtual machines")

	// storage
	cmd_storage         = app.Command("storage", "storage domains")
	cmd_storage_summary = cmd_storage.Command("summary", "output summary of storage")
)
