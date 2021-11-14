package main

import "net"

type storageDomain struct {
	Name          string
	Total         string
	Used          string
	Overcommitted string
}

type cloudInit struct {
	User         string
	SSHPubKey    string
	NICName      string
	FQDN         string
	Hostname     string
	Domain       string
	CIDR         string
	IPAddress    net.IP
	Network      net.IPNet
	Gateway4     net.IP
	DNSAddresses []string
}
