package main

import (
	"errors"
	"net"
	"os"
	"strings"
	"text/template"
)

func incrementIP(network net.IPNet) net.IP {
	guessedGateway := net.ParseIP(network.IP.String())
	for i := len(guessedGateway) - 1; i >= 0; i-- {
		guessedGateway[i]++
		if guessedGateway[i] != 0 {
			break
		}
	}

	return guessedGateway
}

func intakeCloudInitInput() (cloudInit, error) {
	var conf cloudInit
	var haveCidr bool = false
	var haveGateway4 bool = false

	conf.NICName = *cmd_cloud_init_create_nic

	if len(*cmd_cloud_init_create_ssh_key) > 0 {
		conf.SSHPubKey = *cmd_cloud_init_create_ssh_key
	}

	if len(*cmd_cloud_init_create_user) < 1 {
		return conf, errors.New("you specified an empty username; cannot use")
	}

	if len(*cmd_cloud_init_create_user) > 0 {
		conf.User = *cmd_cloud_init_create_user
	}

	if len(*cmd_cloud_init_create_fqdn) < 1 {
		return conf, errors.New("--fqdn is mandatory")
	}
	conf.FQDN = *cmd_cloud_init_create_fqdn

	if len(*cmd_cloud_init_create_cidr) > 0 {
		conf.CIDR = *cmd_cloud_init_create_cidr
		// fmt.Println("parse cidr")
		ip, ipNet, err := net.ParseCIDR(*cmd_cloud_init_create_cidr)
		conf.IPAddress = ip
		conf.Network = *ipNet
		if err != nil {
			return conf, err
		}
		haveCidr = true
	}

	if len(*cmd_cloud_init_create_gateway4) > 0 {
		// fmt.Println("parse gw")
		ip := net.ParseIP(*cmd_cloud_init_create_gateway4)
		conf.Gateway4 = ip
		if len(ip) > 0 {
			haveGateway4 = true
		}
	}

	if haveCidr && !haveGateway4 {
		// fmt.Println("network from CIDR, but no gateway; guessing .1")
		guessedGateway := incrementIP(conf.Network)
		conf.Gateway4 = guessedGateway
	}

	conf.DNSAddresses = strings.Split(*cmd_cloud_init_create_dns_addresses, ",")

	return conf, nil
}

func genCloudInitScript(config cloudInit) error {
	splitFQDN := strings.Split(config.FQDN, ".")
	config.Hostname = splitFQDN[0]
	if len(splitFQDN) < 2 {
		config.Domain = ""
	} else {
		config.Domain = strings.Join(splitFQDN[1:], ".")
	}
	templt := "fqdn: {{ .FQDN }}\n" +
		"growpart:\n" +
		"  mode: \"off\"\n" +
		"disable_root: true\n" +
		"ssh_deletekeys: false\n" +
		"preserve_hostname: true\n" +
		"resize_rootfs: false\n" +
		"packages_update: true\n" +
		"packages_upgrade: true\n" +
		"packages:\n" +
		"  - openssh-server\n" +
		"write_files:\n" +
		"  - path: /etc/netplan/00-installer-config.yaml\n" +
		"    permissions: 0o644\n" +
		"    content: |\n" +
		"      ---\n" +
		"      network: {}\n" +
		"  - path: /etc/netplan/50-cloud-init.yaml\n" +
		"    permissions: 0o644\n" +
		"    content: |\n" +
		"      ---\n" +
		"      network: {}\n" +
		"  - path: /etc/hosts\n" +
		"    permissions: 0o644\n" +
		"    content: |\n" +
		"      127.0.0.1 localhost\n" +
		"      127.0.1.1 {{ .FQDN }} {{ .Hostname }}\n" +
		"      \n" +
		"      # The following lines are desirable for IPv6 capable hosts\n" +
		"      ::1     ip6-localhost ip6-loopback\n" +
		"      fe00::0 ip6-localnet\n" +
		"      ff00::0 ip6-mcastprefix\n" +
		"      ff02::1 ip6-allnodes\n" +
		"      ff02::2 ip6-allrouters\n" +
		"  - path: /etc/cloud/cloud.cfg.d/no_datasource.cfg\n" +
		"    permissions: 0o644\n" +
		"    content: |\n" +
		"      datasource:\n" +
		"        None: {}\n" +
		"      datasource_list:\n" +
		"        - None\n" +
		"  - path: /etc/cloud/cloud.cfg.d/99-custom-networking.cfg\n" +
		"    permissions: 0o644\n" +
		"    content: |\n" +
		"      network: {config: disabled}\n" +
		"  - path: /etc/netplan/config.yaml\n" +
		"    permissions: 0o644\n" +
		"    content: |\n" +
		"      ---\n" +
		"      network:\n" +
		"        version: 2\n" +
		"        renderer: networkd\n" +
		"        ethernets:\n" +
		"          {{ .NICName }}:\n" +
		"            optional: yes\n"
	if len(config.CIDR) > 0 {
		templt = templt +
			"            dhcp4: no\n" +
			"            dhcp6: no\n" +
			"            addresses:\n" +
			"              - {{ .CIDR }}\n" +
			"            routes:\n" +
			"              - to: default\n" +
			"                via: {{ .Gateway4 }}\n"
		if len(config.DNSAddresses) > 0 {
			templt = templt +
				"            nameservers:\n" +
				"              addresses:\n"
			for _, addr := range config.DNSAddresses {
				templt = templt +
					"                - " + addr + "\n"
			}
		}
	} else {
		templt = templt +
			"            dhcp4: yes\n" +
			"            dhcp6: no\n"
	}
	templt = templt +
		"runcmd:\n" +
		"  - date > /opt/.creation\n" +
		"  - netplan apply\n" +
		"users:\n" +
		"  - name: {{ .User }}\n"
	if len(config.SSHPubKey) > 0 {
		templt = templt +
			"    ssh_authorized_keys:\n" +
			"      - \"{{ .SSHPubKey }}\"\n"
	}

	conf, err := template.New("conf").Parse(templt)
	if err != nil {
		return err
	}

	out := os.Stdout
	if len(*cmd_cloud_init_create_output) > 0 {
		out, err = os.Create(*cmd_cloud_init_create_output)
		if err != nil {
			return err
		}
	}

	err = conf.Execute(out, config)
	if err != nil {
		return err
	}

	out.Close()

	return nil
}
