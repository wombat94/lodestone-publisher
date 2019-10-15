package main

import (
	"fmt"
	"github.com/analogj/go-util/utils"
	"github.com/analogj/lodestone-publisher/pkg/notify"
	"github.com/analogj/lodestone-publisher/pkg/version"
	"github.com/analogj/lodestone-publisher/pkg/watch"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"log"
	"os"
	"time"
)

var goos string
var goarch string

func main() {
	app := &cli.App{
		Name:     "lodestone-fs-publisher",
		Usage:    "Filesystem watcher & notifications for lodestone",
		Version:  version.VERSION,
		Compiled: time.Now(),
		Authors: []cli.Author{
			cli.Author{
				Name:  "Jason Kulatunga",
				Email: "jason@thesparktree.com",
			},
		},
		Before: func(c *cli.Context) error {

			capsuleUrl := "https://github.com/AnalogJ/lodestone-publisher"

			versionInfo := fmt.Sprintf("%s.%s-%s", goos, goarch, version.VERSION)

			subtitle := capsuleUrl + utils.LeftPad2Len(versionInfo, " ", 53-len(capsuleUrl))

			fmt.Fprintf(c.App.Writer, fmt.Sprintf(utils.StripIndent(
				`
			 __    _____  ____  ____  ___  ____  _____  _  _  ____ 
			(  )  (  _  )(  _ \( ___)/ __)(_  _)(  _  )( \( )( ___)
			 )(__  )(_)(  )(_) ))__) \__ \  )(   )(_)(  )  (  )__) 
			(____)(_____)(____/(____)(___/ (__) (_____)(_)\_)(____)
			%s
			`), subtitle))
			return nil
		},

		Commands: []cli.Command{
			{
				Name:  "start",
				Usage: "Start the Lodestone filesystem watcher",
				Action: func(c *cli.Context) error {

					var notifyClient notify.Interface

					notifyClient = new(notify.RedisNotify)
					notifyClient.Init(map[string]string{
						"addr":     fmt.Sprintf("%s:%d", c.String("redis-hostname"), c.Int("redis-port")),
						"password": c.String("redis-password"),
						"queue":    c.String("redis-queue"),
					})

					watcher := watch.FsWatcher{}
					watcher.Start(notifyClient, map[string]string{
						"dir":    c.String("dir"),
						"bucket": c.String("bucket"),
					})
					return nil
				},

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "dir",
						Usage: "The directory to watch for file changes.",
					},
					&cli.StringFlag{
						Name:  "bucket",
						Usage: "The name of the bucket",
					},

					&cli.StringFlag{
						Name:  "redis-hostname",
						Usage: "The redis server hostname",
						Value: "localhost",
					},
					&cli.IntFlag{
						Name:  "redis-port",
						Usage: "The redis server port",
						Value: 6379,
					},
					&cli.StringFlag{
						Name:  "redis-password",
						Usage: "The redis server password",
						Value: "",
					},
					&cli.StringFlag{
						Name:  "redis-queue",
						Usage: "The redis server queue",
						Value: "documentsevents",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(color.HiRedString("ERROR: %v", err))
	}
}