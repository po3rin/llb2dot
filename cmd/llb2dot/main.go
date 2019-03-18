package main

import (
	"log"
	"os"

	"github.com/po3rin/llb2dot"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "llb2dot"
	app.Usage = "convert buildkit LLB DAG graph to dot language"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "llb, l",
			Usage: "generate dot using LLB directory",
		},
		cli.StringFlag{
			Name:  "file, f",
			Value: "Dockerfile",
			Usage: "Dockerfile path",
		},
	}
	app.Action = action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) error {
	var ops llb2dot.LLBOps
	var err error

	if c.Bool("l") {
		ops, err = llb2dot.LoadLLB(os.Stdin)
		if err != nil {
			return err
		}
	} else {
		f, err := os.Open(c.String("f"))
		if err != nil {
			return err
		}
		ops, err = llb2dot.LoadDockerfile(f)
		if err != nil {
			return err
		}
	}

	g, err := llb2dot.LLB2Graph(ops)
	if err != nil {
		return err
	}
	llb2dot.WriteDOT(os.Stdout, g)
	return nil
}
