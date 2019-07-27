package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gradecak/fission-workflows/cmd/fission-workflows-bundle/bundle"
	"github.com/gradecak/fission-workflows/pkg/fes/backend/nats"
	"github.com/gradecak/fission-workflows/pkg/util"
	natsio "github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	ctx, cancelFn := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		for sig := range c {
			fmt.Println("Received signal: ", sig)
			go func() {
				time.Sleep(30 * time.Second)
				fmt.Println("Deadline exceeded; forcing shutdown.")
				os.Exit(0)
			}()
			cancelFn()
			break
		}
	}()

	cliApp := createCli()
	cliApp.Action = func(c *cli.Context) error {
		setupLogging(c)
		policy, err := bundle.ParseSchedulerConfig(c)
		if err != nil {
			logrus.Fatal("Error while initializing workflows: ", err)
		}

		proxyConfig, err := bundle.ParseFissionProxyConfig(c)
		if err != nil {
			logrus.Fatal("Error while parsing Fission Proxy: ", err)
		}

		return bundle.Run(ctx, &bundle.Options{
			NATS:                 parseNatsOptions(c),
			Fission:              parseFissionOptions(c),
			Scheduler:            policy,
			InternalRuntime:      c.Bool("internal"),
			InvocationController: c.Bool("controller") || c.Bool("invocation-controller"),
			WorkflowController:   c.Bool("controller") || c.Bool("workflow-con0troller"),
			Consent:              c.Bool("consent"),
			Provenance:           c.Bool("provenance"),
			AdminAPI:             c.Bool("api") || c.Bool("api-admin"),
			WorkflowAPI:          c.Bool("api") || c.Bool("api-workflow"),
			ConsentAPI:           c.Bool("api") || c.Bool("api-consent"),
			InvocationAPI:        c.Bool("api") || c.Bool("api-workflow-invocation"),
			HTTPGateway:          c.Bool("api") || c.Bool("api-http"),
			Metrics:              c.Bool("metrics"),
			Debug:                c.Bool("debug"),
			FissionProxy:         proxyConfig,
			UseNats:              c.Bool("nats"),
			ProvNats:             c.Bool("prov-nats"),
			ConsentNats:          c.Bool("consent-nats"),
			MaxParallel:          c.Int("max-parallel-invocations"),
		})
	}
	cliApp.Run(os.Args)
}

func setupLogging(c *cli.Context) {
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func parseFissionOptions(c *cli.Context) *bundle.FissionOptions {
	if !c.Bool("fission") {
		return nil
	}

	return &bundle.FissionOptions{
		ExecutorAddress: c.String("fission-executor"),
		ControllerAddr:  c.String("fission-controller"),
		RouterAddr:      c.String("fission-router"),
	}
}

func parseNatsOptions(c *cli.Context) *nats.Config {
	// since dataflows utilises nats config, ensure opts are only nil if both
	// are disabled
	if !c.Bool("nats") && !c.Bool("prov-nats") && !c.Bool("consent") {
		return nil
	}

	client := c.String("nats-client")
	if client == "" {
		client = fmt.Sprintf("workflow-bundle-%s", util.UID())
	}

	return &nats.Config{
		URL:           c.String("nats-url"),
		Cluster:       c.String("nats-cluster"),
		Client:        client,
		AutoReconnect: true,
	}
}

func createCli() *cli.App {

	cliApp := cli.NewApp()

	cliApp.Flags = append([]cli.Flag{
		// Generic
		cli.BoolFlag{
			Name:   "d, debug",
			EnvVar: "WORKFLOW_DEBUG",
		},

		// NATS
		cli.StringFlag{
			Name:   "nats-url",
			Usage:  "URL to the data store used by the NATS event store.",
			Value:  natsio.DefaultURL, // http://nats-streaming.fission
			EnvVar: "ES_NATS_URL",
		},
		cli.StringFlag{
			Name:   "nats-cluster",
			Usage:  "Cluster name used for the NATS event store (if needed)",
			Value:  "test-cluster", // mqtrigger
			EnvVar: "ES_NATS_CLUSTER",
		},
		cli.StringFlag{
			Name:   "nats-client",
			Usage:  "Client name used for the NATS event store. By default it will generate a unique clientID.",
			EnvVar: "ES_NATS_CLIENT",
		},
		cli.BoolFlag{
			Name:  "nats",
			Usage: "Use NATS as the event store",
		},
		cli.BoolFlag{
			Name:  "prov-nats",
			Usage: "Use NATS for publishing provenance",
		},
		cli.BoolFlag{
			Name:  "consent-nats",
			Usage: "Use NATS for recieving provenance",
		},
		// Dataflow specifics
		cli.BoolFlag{
			Name:  "consent",
			Usage: "Enable consent verification on workflow execution",
		},
		cli.BoolFlag{
			Name:  "provenance",
			Usage: "Enable provenance generation on completed workflows",
		},
		// Fission Environment Proxy
		cli.BoolFlag{
			Name:  "fission.proxy, fission-proxy",
			Usage: "Enable Fission environment as a proxy",
		},
		cli.DurationFlag{
			Name:  "fission.proxy.timeout",
			Usage: "The default timeout assigned to workflow invocations coming from the Fission proxy",
			Value: 5 * time.Minute,
		},
		cli.StringFlag{
			Name:  "fission.proxy.addr",
			Usage: "The timeout assigned to workflow invocations coming from the Fission proxy",
			Value: ":8888",
		},

		// Fission Function Runtime
		cli.BoolFlag{
			Name:  "fission",
			Usage: "Use Fission as a function environment",
		},
		cli.StringFlag{
			Name:   "fission-executor",
			Usage:  "Address of the Fission executor to optimize executions",
			Value:  "http://executor.fission",
			EnvVar: "FNENV_FISSION_EXECUTOR",
		},
		cli.StringFlag{
			Name:   "fission-controller",
			Usage:  "Address of the Fission controller for resolving functions",
			Value:  "http://controller.fission",
			EnvVar: "FNENV_FISSION_CONTROLLER",
		},
		cli.StringFlag{
			Name:   "fission-router",
			Usage:  "Address of the Fission router for executing functions",
			Value:  "http://router.fission",
			EnvVar: "FNENV_FISSION_ROUTER",
		},

		// Components
		cli.BoolFlag{
			Name:  "internal",
			Usage: "Use internal function runtime",
		},
		cli.BoolFlag{
			Name:  "controller",
			Usage: "Run the controller with all components",
		},
		cli.IntFlag{
			Name:  "max-parallel-invocations",
			Usage: "Maximium number of parallel invocations to allow into the system",
			Value: 1500,
		},
		cli.BoolFlag{
			Name:  "workflow-controller",
			Usage: "Run the workflow controller",
		},
		cli.BoolFlag{
			Name:  "invocation-controller",
			Usage: "Run the invocation controller",
		},
		cli.BoolFlag{
			Name:  "api-http",
			Usage: "Serve the http apis of the apis",
		},
		cli.BoolFlag{
			Name:  "api-workflow-invocation",
			Usage: "Serve the workflow invocation gRPC api",
		},
		cli.BoolFlag{
			Name:  "api-workflow",
			Usage: "Serve the workflow gRPC api",
		},
		cli.BoolFlag{
			Name:  "api-consent",
			Usage: "Serve the consent gRPC api",
		},
		cli.BoolFlag{
			Name:  "api-admin",
			Usage: "Serve the admin gRPC api",
		},
		cli.BoolFlag{
			Name:  "metrics",
			Usage: "Serve prometheus metrics",
		},
		cli.BoolFlag{
			Name:  "api",
			Usage: "Shortcut for serving all APIs over both gRPC and HTTP",
		},

		// Scheduler
		cli.StringFlag{
			Name:  bundle.FlagSchedulerPolicy,
			Usage: "Policy to use for the scheduler (prewarm-all, prewarm-horizon, horizon)",
			Value: "horizon-mz",
		},
		cli.DurationFlag{
			Name:  bundle.FlagSchedulerColdStartDuration,
			Usage: "The static cold start duration to assume when using prewarm schedulers",
			Value: 1 * time.Second,
		},
	})

	return cliApp
}
