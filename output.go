package main

import "github.com/tatsushid/go-prettytable"

func tableOut(storageDomains []storageDomain) {
	table, err := prettytable.NewTable(
		prettytable.Column{Header: "Name"},
		prettytable.Column{Header: "Size"},
		prettytable.Column{Header: "Used"},
		prettytable.Column{Header: "Overcommitted"},
	)
	if err != nil {
		panic(err)
	}
	table.Separator = " | "
	for _, sd := range storageDomains {
		table.AddRow(
			sd.Name,
			sd.Total,
			sd.Used,
			sd.Overcommitted,
		)
	}
	table.Print()
}
