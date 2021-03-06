package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	nlog "github.com/nuveo/log"

	"github.com/spf13/cobra"
	mysql "github.com/suifengpiao14/crud-mysql/adapters/mysql"
	"github.com/suifengpiao14/crud-mysql/config"
	"github.com/suifengpiao14/crud-mysql/config/router"
	"github.com/suifengpiao14/crud-mysql/controllers"
	"github.com/suifengpiao14/crud-mysql/middlewares"
	"github.com/urfave/negroni"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "start",
	Short: "Serve a RESTful API from mysql database",
	Long:  `Serve a RESTful API from mysql database, start HTTP server`,
	Run: func(cmd *cobra.Command, args []string) {
		mysql.Load()
		nlog.Warningln("start mysql restful server")
		startServer()

	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	migrateCmd.AddCommand(createCmd)
	migrateCmd.AddCommand(downCmd)
	migrateCmd.AddCommand(gotoCmd)
	migrateCmd.AddCommand(mversionCmd)
	migrateCmd.AddCommand(nextCmd)
	migrateCmd.AddCommand(redoCmd)
	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(resetCmd)
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(migrateCmd)
	migrateCmd.PersistentFlags().StringVar(&urlConn, "url", driverURL(), "Database driver url")
	migrateCmd.PersistentFlags().StringVar(&path, "path", config.PrestConf.MigrationsPath, "Migrations directory")
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// MakeHandler reagister all routes
func MakeHandler() http.Handler {
	n := middlewares.GetApp()
	r := router.Get()
	r.HandleFunc("/databases", controllers.GetDatabases).Methods("GET")
	r.HandleFunc("/schemas", controllers.GetSchemas).Methods("GET")
	r.HandleFunc("/tables", controllers.GetTables).Methods("GET")
	r.HandleFunc("/_QUERIES/{queriesLocation}/{script}", controllers.ExecuteFromScripts)
	r.HandleFunc("/tpl", controllers.ExecuteTplQuery)
	r.HandleFunc("/{database}/{schema}", controllers.GetTablesByDatabaseAndSchema).Methods("GET")
	crudRoutes := mux.NewRouter().PathPrefix("/").Subrouter().StrictSlash(true)
	crudRoutes.HandleFunc("/{database}/{schema}/{table}", controllers.SelectFromTables).Methods("GET")
	crudRoutes.HandleFunc("/{database}/{schema}/{table}", controllers.InsertInTables).Methods("POST")
	crudRoutes.HandleFunc("/async/{database}/{schema}/{table}", controllers.AsyncInsertInTables).Methods("POST")
	crudRoutes.HandleFunc("/batch/{database}/{schema}/{table}", controllers.BatchInsertInTables).Methods("POST")
	crudRoutes.HandleFunc("/{database}/{schema}/{table}", controllers.DeleteFromTable).Methods("DELETE")
	crudRoutes.HandleFunc("/{database}/{schema}/{table}", controllers.UpdateTable).Methods("PUT", "PATCH")
	r.PathPrefix("/").Handler(negroni.New(
		middlewares.AccessControl(),
		negroni.Wrap(crudRoutes),
	))
	n.UseHandler(r)
	return n
}

func startServer() {

	mux := http.NewServeMux()
	mux.Handle(config.PrestConf.ContextPath, MakeHandler())
	l := log.New(os.Stdout, "[prest] ", 0)

	if !config.PrestConf.AccessConf.Restrict {
		nlog.Warningln("You are running pREST in public mode.")
	}

	if config.PrestConf.Debug {
		nlog.DebugMode = config.PrestConf.Debug
		nlog.Warningln("You are running pREST in debug mode.")
	}
	addr := fmt.Sprintf("%s:%d", config.PrestConf.HTTPHost, config.PrestConf.HTTPPort)
	l.Printf("listening on %s and serving on %s", addr, config.PrestConf.ContextPath)
	if config.PrestConf.HTTPSMode {
		l.Fatal(http.ListenAndServeTLS(addr, config.PrestConf.HTTPSCert, config.PrestConf.HTTPSKey, mux))
	}
	l.Fatal(http.ListenAndServe(addr, mux))
}
