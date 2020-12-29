package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	cli := kingpin.New("vault-init", "Automate the initialization and unsealing of HashiCorp Vault.")

	startCmdClause := cli.Command("start", "Start the Vault initialization and unsealing process.")
	startCmd := attachStartCommand(startCmdClause)

	switch kingpin.MustParse(cli.Parse(os.Args[1:])) {
	case startCmdClause.FullCommand():
		cli.FatalIfError(startCmd.Run(), "start")
	}
}
