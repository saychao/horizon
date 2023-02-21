package main

import (
	"log"
	"runtime"

	"github.com/saychao/horizon/db2/history/schema"

	"github.com/saychao/horizon"
	"github.com/saychao/horizon/config"
	"github.com/spf13/cobra"
)

var app *horizon.App
var conf config.Config
var version string
var configFile string

var rootCmd *cobra.Command

func main() {
	if version != "" {
		horizon.SetVersion(version)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	rootCmd.Execute()
}

func init() {
	rootCmd = &cobra.Command{
		Use:   "horizon",
		Short: "client-facing api server for the TokenD network",
		Long:  "client-facing api server for the TokenD network",
		Run: func(cmd *cobra.Command, args []string) {
			initApp(cmd, args)
			app.Serve()
		},
	}
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "config file")

	rootCmd.AddCommand(dbCmd)
}

func initApp(cmd *cobra.Command, args []string) {
	err := initConfig(configFile)
	if err != nil {
		log.Println("Failed to init config")
		log.Fatal(err.Error())
	}

	if conf.MigrateUpOnStart {
		migrate("up", 0, schema.Migrate, conf.DatabaseURL)
		tryToEmptyDB()
	}

	app, err = horizon.NewApp(conf)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func initConfig(fn string) error {
	c := config.NewViperConfig(fn)
	if err := c.Init(); err != nil {
		return err
	}
	conf = c
	return nil
}
