package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	cli := kingpin.New("vault-init", "Automate the initialization and unsealing of HashiCorp Vault.")

	startCmdClause := cli.Command("start", "Start the Vault initialization and unsealing process.")
	startCmd := attachStartCommand(startCmdClause)

	showCmdClause := cli.Command("show", "Fetch and decrypt the root token and unseal keys generated during the Vault initialization process.")
	showCmd := attachShowCommand(showCmdClause)

	switch kingpin.MustParse(cli.Parse(os.Args[1:])) {
	case startCmdClause.FullCommand():
		cli.FatalIfError(startCmd.Run(), "start")
	case showCmdClause.FullCommand():
		cli.FatalIfError(showCmd.Execute(), "show")
	}
}
