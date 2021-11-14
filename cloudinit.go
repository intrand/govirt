package main

import (
	"errors"
	"fmt"
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
		fmt.Println("network from CIDR, but no gateway; guessing .1")
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
	templt := "---\n" +
		"fqdn: {{ .FQDN }}\n" +
		"write_files:\n" +
		"- path: /etc/cloud/cloud.cfg.d/99-custom-networking.cfg\n" +
		"  permissions: '0644'\n" +
		"  content: |\n" +
		"    network: {config: disabled}\n" +
		"- path: /etc/netplan/config.yaml\n" +
		"  permissions: '0644'\n" +
		"  content: |\n" +
		"    ---\n" +
		"    network:\n" +
		"      version: 2\n" +
		"      renderer: networkd\n" +
		"      ethernets:\n" +
		"        {{ .NICName }}:\n" +
		"          optional: yes\n"
	if len(config.CIDR) > 0 {
		templt = templt +
			"          dhcp4: no\n" +
			"          dhcp6: no\n" +
			"          addresses:\n" +
			"            - {{ .CIDR }}\n" +
			"          gateway4: {{ .Gateway4 }}\n"
		if len(config.DNSAddresses) > 0 {
			templt = templt +
				"          nameservers:\n" +
				"            addresses:\n"
			for _, addr := range config.DNSAddresses {
				templt = templt +
					"              - " + addr + "\n"
			}
		}
	} else {
		templt = templt +
			"          dhcp4: yes\n" +
			"          dhcp6: no\n"
	}
	templt = templt +
		"runcmd:\n" +
		"- sed -i 's/template.domain template/{{ .FQDN }} {{ .Hostname }}/g' /etc/hosts  \n" +
		"- date > /opt/.creation\n" +
		"- rm /etc/netplan/50-cloud-init.yaml\n" +
		"- netplan generate\n" +
		"- netplan apply\n" +
		"- dpkg-reconfigure openssh-server\n" +
		"users:\n" +
		"- default\n" +
		"- name: {{ .User }}\n"
	if len(config.SSHPubKey) > 0 {
		templt = templt +
			"  ssh_authorized_keys:\n" +
			"  - \"{{ .SSHPubKey }}\"\n"
	}

	conf, err := template.New("conf").Parse(templt)
	if err != nil {
		return err
	}
	err = conf.Execute(os.Stdout, config)
	if err != nil {
		return err
	}

	return nil
}
