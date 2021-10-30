// https://github.com/oVirt/go-ovirt/tree/master/examples
package main

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	// get input
	args := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch args {

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

	// vm rm
	case cmd_vm_rm.FullCommand():
		// get connection
		conn, err := connect()
		if err != nil {
			log.Fatalln(err)
		}

		// vm
		vm, err := getVm(conn, "test")
		if err != nil {
			log.Fatalln(err)
		}

		// warn
		log.Println("Requesting to PERMANENTLY DELETE VM " + *cmd_vm_rm_name + " WITH NO CHANCE FOR RECOVERY")

		// check for --yes
		if !*cmd_vm_rm_yes {
			log.Fatalln("Not approved. You MUST specify --yes as a safety precaution. Do not script this! You've been warned! Exiting.")
		}

		// delete
		err = deleteVm(conn, vm)
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

	// // Get the reference to the "clusters" service
	// clustersService := conn.SystemService().ClustersService()

	// // Use the "list" method of the "clusters" service to list all the clusters of the system
	// clustersResponse, err := clustersService.List().Send()
	// if err != nil {
	// 	fmt.Printf("Failed to get cluster list, reason: %v\n", err)
	// 	return
	// }

	// if clusters, ok := clustersResponse.Clusters(); ok {
	// 	// Print the datacenter names and identifiers
	// 	for _, cluster := range clusters.Slice() {
	// 		if clusterName, ok := cluster.Name(); ok {
	// 			fmt.Printf("Cluster name: %v\n", clusterName)
	// 		}
	// 		if clusterId, ok := cluster.Id(); ok {
	// 			fmt.Printf("Cluster id: %v\n", clusterId)
	// 		}
	// 	}
	// }
}
