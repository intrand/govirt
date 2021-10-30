package main

import (
	"time"

	ovirt "github.com/ovirt/go-ovirt"
)

func connect() (*ovirt.Connection, error) {
	conn, err := ovirt.NewConnectionBuilder().
		URL(*cmd_url).
		Username(*cmd_username).
		Password(*cmd_password).
		Insecure(true).
		Compress(true).
		Timeout(time.Second * 10).
		Build()
	if err != nil {
		return nil, err
		// log.Fatalf("Make connection failed, reason: %s", err.Error())
	}

	// Never forget to close connection
	defer conn.Close()

	return conn, nil
}
