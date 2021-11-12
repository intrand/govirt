// https://github.com/oVirt/go-ovirt/tree/master/examples
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "" // to be filled in by goreleaser
	commit  = "" // to be filled in by goreleaser
	date    = "" // to be filled in by goreleaser
	builtBy = "" // to be filled in by goreleaser
	cmdname = filepath.Base(os.Args[0])
)

func main() {
	// get input
	args := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch args {

	case cmd_version.FullCommand():
		fmt.Println(
			"{\"version\":\"" + version + "\",\"commit\":\"" + commit + "\",\"date\":\"" + date + "\",\"built_by\":\"" + builtBy + "\"}")

	case cmd_update.FullCommand():
		// Handle updating to a new version
		log.Print("Attempting update of " + cmdname + "...")
		update_result, err := doSelfUpdate()
		if err != nil {
			log.Println("Couldn't update at this time. Please try again later. Exiting.")
			os.Exit(1)
		}
		if update_result {
			log.Println("Please run " + cmdname + " again.")
			os.Exit(0)
		}

	// vm create
	case cmd_vm_create.FullCommand():
		log.Println("Requesting to create VM " + *cmd_vm_create_name +
			" from template " + *cmd_vm_create_template + ":" + strconv.FormatInt(*cmd_vm_create_template_version, 20) +
			" in cluster " + *cmd_vm_create_cluster +
			" on datastore " + *cmd_vm_create_datastore)

		// get connection
		conn, err := connect()
		if err != nil {
			log.Fatalln(err)
		}

		// get template
		template, err := getTemplate(conn, *cmd_vm_create_template, *cmd_vm_create_template_version)
		if err != nil {
			log.Fatalln(err)
		}

		// get domain
		domain, err := getStorageDomain(conn, *cmd_vm_create_datastore)
		if err != nil {
			log.Fatalln(err)
		}

		// create vm
		err = createVmFromTemplate(conn, template, domain, *cmd_vm_create_name, *cmd_vm_create_cluster)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Successfully sent the request.")

	// vm get
	case cmd_vm_get.FullCommand():
		// get connection
		conn, err := connect()
		if err != nil {
			log.Fatalln(err)
		}

		// vm
		vm, err := getVm(conn, *cmd_vm_get_name)
		if err != nil {
			log.Fatalln(err)
		}

		vmId, ok := vm.Id()
		if !ok {
			log.Fatalln("couldn't get virtual machine id")
		}

		vmName, ok := vm.Name()
		if !ok {
			log.Fatalln("couldn't get virtual machine name")
		}

		vmState, ok := vm.Status()
		if !ok {
			log.Fatalln("couldn't get virtual machine status")
		}

		vmDescription, ok := vm.Description()
		if !ok {
			log.Fatalln("couldn't get virtual machine description")
		}

		fmt.Print(
			"name: " + vmName + "\n" +
				"\tid: " + vmId + "\n" +
				"\tdesc: " + vmDescription + "\n" +
				"\tstate: " + string(vmState) + "\n",
		)

	// vm rm
	case cmd_vm_rm.FullCommand():
		// get connection
		conn, err := connect()
		if err != nil {
			log.Fatalln(err)
		}

		// vm
		vm, err := getVm(conn, *cmd_vm_rm_name)
		if err != nil {
			log.Fatalln(err)
		}

		// warn
		log.Println("Requesting to PERMANENTLY DELETE VM " + *cmd_vm_rm_name + " WITH NO CHANCE FOR RECOVERY")

		// check for --yes
		if !*cmd_vm_rm_yes {
			log.Fatalln("Not approved. You MUST specify --yes as a safety precaution. Do not script this! You've been warned! Exiting.")
		}

		// shutdown
		err = shutdownVm(conn, vm, true) // force off
		if err != nil {
			log.Fatalln(err)
		}

		// delete
		err = deleteVm(conn, vm)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Successfully sent the request.")

	// vm start
	case cmd_vm_start.FullCommand():
		// if *cmd_vm_start_init {
		// 	if len(*cmd_vm_start_ip) < 1 {
		// 		log.Fatalln("You MUST specify --ip <address> when starting a VM with the --init argument.")
		// 	}
		// }
		// get connection
		conn, err := connect()
		if err != nil {
			log.Fatalln(err)
		}

		// vm
		vm, err := getVm(conn, *cmd_vm_start_name)
		if err != nil {
			log.Fatalln(err)
		}

		// warn
		if *cmd_vm_start_init {
			log.Println("Requesting to start virtual machine " + *cmd_vm_start_name + " with cloud-init config: " + *cmd_vm_start_script)
		} else {
			log.Println("Requesting to start virtual machine " + *cmd_vm_start_name)
		}

		err = startVm(conn, vm, *cmd_vm_start_init)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Successfully sent the request.")
		// } else {
		// 	log.Fatalln("Couldn't start virtual machine. Look above for errors.")
		// }

	// vm stop
	case cmd_vm_stop.FullCommand():
		// get connection
		conn, err := connect()
		if err != nil {
			log.Fatalln(err)
		}

		// vm
		vm, err := getVm(conn, *cmd_vm_stop_name)
		if err != nil {
			log.Fatalln(err)
		}

		// warn
		log.Println("Requesting to stop virtual machine " + *cmd_vm_stop_name)

		// shutdown
		err = shutdownVm(conn, vm, *cmd_vm_stop_force)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Successfully sent the request.")

	// vm summary
	case cmd_vm_summary.FullCommand():
		// get connection
		conn, err := connect()
		if err != nil {
			log.Fatalln(err)
		}

		vmsService := conn.SystemService().VmsService()
		vmsService.Add().Clone(true)
		vmResp, err := vmsService.List().Send()
		if err != nil {
			log.Fatalln(err)
		}

		vmSlice, ok := vmResp.Vms()
		if !ok {
			log.Fatalln(err)
		}

		allCores := []int64{}
		for _, whatever := range vmSlice.Slice() {
			cpu, _ := whatever.Cpu()
			top, _ := cpu.Topology()
			cores, _ := top.Cores()
			sockets, _ := top.Sockets()
			allCores = append(allCores, sockets*cores)
			// fmt.Println("sockets:", sockets, "; cores per socket:", cores, "; total vcpu:", sockets*cores)
		}

		if len(allCores) > 0 {
			var allCoresTotal int64
			for _, core := range allCores {
				allCoresTotal = allCoresTotal + core
			}

			avg := int(allCoresTotal) / len(allCores)
			log.Println("Average of " + strconv.Itoa(avg) + " cores per virtual machine")
		}

	// storage summary
	case cmd_storage_summary.FullCommand():
		// get connection
		conn, err := connect()
		if err != nil {
			log.Fatalln(err)
		}

		storageService := conn.SystemService().StorageDomainsService()
		storageResp, err := storageService.List().Send()
		if err != nil {
			log.Fatalln(err)
		}
		domains, ok := storageResp.StorageDomains()
		if !ok {
			log.Fatalln(err)
		}
		outputDomains := []storageDomain{}
		for _, domain := range domains.Slice() {
			name, _ := domain.Name()
			available, _ := domain.Available()
			used, _ := domain.Used()
			commited, _ := domain.Committed()

			total := available + used
			if total == 0 {
				continue
			}

			percent := float64(used) / float64(total)
			overcommit := float64(commited) / float64(total)
			outputDomains = append(outputDomains, storageDomain{
				Name:          name,
				Total:         strconv.Itoa(int(total/1024/1024/1024)) + "GB",
				Used:          strconv.FormatFloat(percent, 'f', 2, 64) + "%",
				Overcommitted: strconv.FormatFloat(overcommit, 'f', 2, 64) + "x",
			})
		}
		tableOut(outputDomains)
	}
}
