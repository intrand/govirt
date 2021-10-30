package main

import (
	"errors"

	ovirt "github.com/ovirt/go-ovirt"
)

func getStorageDomain(conn *ovirt.Connection, domainName string) (ovirt.StorageDomain, error) {
	domain := ovirt.StorageDomain{}

	sdsService := conn.SystemService().StorageDomainsService()

	sdsResp, err := sdsService.List().Search("name=" + domainName).Send()
	if err != nil {
		return domain, err
	}

	sdSlice, ok := sdsResp.StorageDomains()
	if !ok {
		return domain, errors.New("couldn't get list of storage domains")
	}

	for _, sd := range sdSlice.Slice() {
		name, ok := sd.Name()
		if !ok {
			return domain, errors.New("couldn't set name of storage domain")
		}
		if name == domainName {
			domain = *sd
		}
	}
	return domain, nil
}
