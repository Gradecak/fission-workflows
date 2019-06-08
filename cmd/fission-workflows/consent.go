package main

import (
	// "context"
	"github.com/sirupsen/logrus"
	// "github.com/gradecak/fission-workflows/pkg/types"
	"github.com/urfave/cli"
	"strings"
)

var cmdConsent = cli.Command{
	Name:        "consent",
	Usage:       "consent <id> <status>",
	Description: "Inject consent status message into the workflow engine",
	Action: commandContext(func(ctx Context) error {

		cId := ctx.Args().Get(0)
		status := ctx.Args().Get(1)
		status = strings.ToUpper(status)
		client := getClient(ctx)
		// msg := &types.ConsentMessage{}
		err := client.Consent.UpdateFromString(ctx, cId, status)
		if err != nil {
			logrus.Fatalf("%v", err)
		}
		return err
	}),
}
