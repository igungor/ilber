package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"net/http/pprof"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/ilber/command"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// flags
var (
	flagConfig = flag.String("c", "./ilber.conf", "configuration file path")
)

func usage() {
	fmt.Fprintf(os.Stderr, "ilber is a multi-purpose Telegram bot\n\n")
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "  ilber -c path-to-ilber.conf\n\n")
	fmt.Fprintf(os.Stderr, "flags:\n")
	flag.PrintDefaults()
}

func main() {
	logger := log.New(os.Stdout, "ilber: ", log.LstdFlags|log.Lshortfile)
	flag.Usage = usage
	flag.Parse()

	b, err := bot.New(*flagConfig, logger)
	if err != nil {
		log.Fatalf("Could not initialize the bot: %v\n", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", b.Handler())
	registerMetrics(mux)
	registerProfile(mux)

	go func() {
		addr := net.JoinHostPort(b.Config.Host, b.Config.Port)
		log.Fatal(http.ListenAndServe(addr, mux))
	}()

	ctx := context.Background()
	for msg := range b.Messages() {
		log.Printf("%v\n", msg)

		// react only to user sent messages
		if msg.IsService() {
			continue
		}
		// is message a bot command?
		cmdname := msg.Command()
		if cmdname == "" {
			continue
		}

		// is the command even registered?
		cmd := command.Lookup(cmdname)
		if cmd == nil {
			continue
		}

		// it is. cool, run it!
		go cmd.Run(ctx, b, msg)
	}
}

func registerMetrics(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}

func registerProfile(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
}
