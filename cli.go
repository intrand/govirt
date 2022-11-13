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

	cmd_version = app.Command("version", "prints version and exits")
	cmd_update  = app.Command("update", "updates to latest version and exits")

	// init
	cmd_cloud_init             = app.Command("cloud-init", " ")
	cmd_cloud_init_create      = cmd_cloud_init.Command("create", "generate cloud-init configuration")
	cmd_cloud_init_create_fqdn = cmd_cloud_init_create.Flag("fqdn", "fqdn").Short('f').Envar("govirt_cloud_init_create_fqdn").Required().String()
	// cmd_cloud_init_create_output        = cmd_cloud_init_create.Flag("output", "output file").Envar("govirt_cloud_init_create_output").Envar("govirt_cloud_init_create_output").Default("script.yml").String()
	cmd_cloud_init_create_cidr          = cmd_cloud_init_create.Flag("cidr", "IPv4 address in CIDR notation to give vm on startup (eg, 192.168.0.99/24)").Short('c').Envar("govirt_cloud_init_create_cidr").String()
	cmd_cloud_init_create_gateway4      = cmd_cloud_init_create.Flag("gateway4", "IPv4 address to give vm on startup (eg, 192.168.0.1)").Short('g').Envar("govirt_cloud_init_create_gateway").String()
	cmd_cloud_init_create_dns_addresses = cmd_cloud_init_create.Flag("dns-addresses", "comma-separated list of IPv4 addresses to give vm on startup").Short('d').Envar("govirt_cloud_init_create_dns_addresses").Default("1.1.1.1,1.0.0.1,8.8.8.8").String()
	cmd_cloud_init_create_user          = cmd_cloud_init_create.Flag("user", "user to customize").Short('u').Envar("govirt_cloud_init_create_user").Default("root").String()
	cmd_cloud_init_create_ssh_key       = cmd_cloud_init_create.Flag("ssh-key", "public SSH key to go to the specified user").Short('s').Envar("govirt_cloud_init_create_ssh_key").String()
	cmd_cloud_init_create_nic           = cmd_cloud_init_create.Flag("nic", "name of NIC to configure").Short('n').Envar("govirt_cloud_init_create_nic").Default("enp1s0").String()
	cmd_cloud_init_create_output        = cmd_cloud_init_create.Flag("output", "path to output file").Short('o').Envar("govirt_cloud_init_create_output").String()

	// vm
	cmd_vm = app.Command("vm", "virtual machines")

	cmd_vm_create                  = cmd_vm.Command("create", "creates VM from template")
	cmd_vm_create_name             = cmd_vm_create.Flag("name", "name of vm to create").Envar("govirt_vm_create_name").Required().String()
	cmd_vm_create_cluster          = cmd_vm_create.Flag("cluster", "cluster in which to create vm").Envar("govirt_vm_create_cluster").Default("Default").String()
	cmd_vm_create_template         = cmd_vm_create.Flag("template", "template from which to create vm").Envar("govirt_vm_create_template").Required().String()
	cmd_vm_create_template_version = cmd_vm_create.Flag("template-version", "version of template from which to create vm").Envar("govirt_vm_create_template_version").Required().Int64()
	cmd_vm_create_datastore        = cmd_vm_create.Flag("datastore", "storage domain in which to create vm").Envar("govirt_vm_create_datastore").Required().String()
	cmd_vm_create_cpu              = cmd_vm_create.Flag("cpu", "cpu socket count (1 core/socket)").Envar("govirt_vm_create_cpu").Default("1").Int64()
	cmd_vm_create_memory           = cmd_vm_create.Flag("memory", "memory (RAM) in GB").Envar("govirt_vm_create_memory").Default("1").Int64()

	cmd_vm_get      = cmd_vm.Command("get", "get details of a vm")
	cmd_vm_get_name = cmd_vm_get.Flag("name", "name of vm to get info of").Envar("govirt_vm_get_name").Required().String()

	cmd_vm_rm      = cmd_vm.Command("rm", "remove a vm")
	cmd_vm_rm_name = cmd_vm_rm.Flag("name", "name of vm to remove").Envar("govirt_vm_rm_name").Required().String()
	cmd_vm_rm_yes  = cmd_vm_rm.Flag("yes", "remove vm from ovirt").Envar("govirt_vm_rm_yes").Default("false").Bool()

	cmd_vm_start        = cmd_vm.Command("start", "start virtual machine with cloud-init")
	cmd_vm_start_name   = cmd_vm_start.Flag("name", "name of vm to start").Envar("govirt_vm_start_name").Required().String()
	cmd_vm_start_init   = cmd_vm_start.Flag("init", "start with cloud-init enabled").Envar("govirt_vm_start_init").Default("false").Bool()
	cmd_vm_start_script = cmd_vm_start.Flag("script", "/path/to/cloud-init-script.yml (eg, /etc/govirt/cloud-init-script.yml)").Envar("govirt_vm_start_script").Default("cloud-init-script.yml").String()

	cmd_vm_stop       = cmd_vm.Command("stop", "shutdown virtual machine gracefully")
	cmd_vm_stop_name  = cmd_vm_stop.Flag("name", "name of vm to stop").Envar("govirt_vm_stop_name").Required().String()
	cmd_vm_stop_force = cmd_vm_stop.Flag("force", "don't wait for shutdown; pull virtual plug").Envar("govirt_vm_stop_force").Default("false").Bool()

	cmd_vm_summary = cmd_vm.Command("summary", "output summary of virtual machines")

	// storage
	cmd_storage         = app.Command("storage", "storage domains")
	cmd_storage_summary = cmd_storage.Command("summary", "output summary of storage")

	// testing
	// cmd_test = app.Command("test", "I wouldn't run this if I were you. You've been warned. No, really, it might delete your entire datacenter without prompting.")
)
