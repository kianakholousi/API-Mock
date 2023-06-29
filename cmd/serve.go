package cmd

import (
	"flight-data-api/config"
	"flight-data-api/database"
	"flight-data-api/http/handler"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serveConfigPath string
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server with the specified configuration",
	Long: `Starts the server with the specified configuration.

You must specify a custom configuration file in YAML format using the --config flag. By default, this command will not run without a configuration file.

The server will listen for incoming requests on the port specified in the configuration file.

It is recommended to run this command to start the server when it is ready to accept incoming traffic.
	
Usage:
  mycommand serve --config [path]
	
Flags:
  -c, --config string   Path to the YAML configuration file (required)
  -h, --help            Help for serve`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&serveConfigPath, "config", "c", "", "Path to the YAML configuration file (required)")
	if err := serveCmd.MarkFlagRequired("config"); err != nil {
		panic(err)
	}
}

func serve() {
	cfg, err := config.Init(config.Params{FilePath: serveConfigPath, FileType: "yaml"})
	if err != nil {
		panic(err)
	}

	db, err := database.InitDB(cfg.Database) // should get db connection and use it in context
	if err != nil {
		panic(err)
	}

	e := echo.New()

	vldt := validator.New()

	flights := handler.Flights{DB: db, Validator: vldt}
	e.GET("/flights", flights.Get)

	if err := e.Start(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		panic(err)
	}
}
