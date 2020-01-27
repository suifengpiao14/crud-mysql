package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	nlog "github.com/nuveo/log"
	"github.com/spf13/cobra"
	mysql "github.com/suifengpiao14/crud-mysql/adapter-mysql"
	"github.com/suifengpiao14/crud-mysql/config"
)

var (
	daemon bool
)

// startCmd represents the create command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start crud-mysql server",
	Run: func(cmd *cobra.Command, args []string) {
		if config.PrestConf.Adapter == nil {
			nlog.Warningln("adapter is not set. Using the default (mysql)")
			mysql.Load()
		}
		if config.PrestConf.SocketPath != "" {
			go func() {
				startSocketServer()
			}()
		}
		startServer()

	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().BoolVarP(&daemon, "daemon", "d", false, "run daemon")
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

// socket 服务
func startSocketServer() {
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
	l.Printf("listening on %s and serving on %s", config.PrestConf.SocketPath, config.PrestConf.ContextPath)
	if err := os.RemoveAll(config.PrestConf.SocketPath); err != nil {
		l.Fatal(err)
	}
	unixListener, err := net.Listen("unix", config.PrestConf.SocketPath)
	if err != nil {
		l.Fatal(err)
	}
	defer unixListener.Close()
	l.Fatal(http.Serve(unixListener, mux))

}
