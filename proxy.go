// Copyright 2015 Jeff Martinez. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file
// or at http://opensource.org/licenses/MIT

/*
See README.md for full description and usage info.
*/

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jeffbmartinez/cleanexit"
	"github.com/jeffbmartinez/config"
	"github.com/jeffbmartinez/delay"
	"github.com/jeffbmartinez/log"
	"github.com/jeffbmartinez/stdoutlog"

	"github.com/jeffbmartinez/proxy/handler"
)

const EXIT_SUCCESS = 0
const EXIT_FAILURE = 1
const EXIT_USAGE_FAILURE = 2 // Same as golang's flag module uses, hardcoded at https://github.com/golang/go/blob/release-branch.go1.4/src/flag/flag.go#L812

const PROJECT_NAME = "proxy"

type ProgramConfig struct {
	Default   string
	Overrides []Override
}

type Override struct {
	From string
	To   string
}

func main() {
	cleanexit.SetUpExitOnCtrlC(showNiceExitMessage)

	allowAnyHostToConnect, listenPort, configFilename := getCommandLineArgs()
	if configFilename == "" {
		log.Fatalf("Didn't supply a configuration file. Use -c filename.json to do so.")
	}

	settings := &ProgramConfig{}

	if err := config.ReadSpecific(configFilename, &settings); err != nil {
		log.Fatalf("Problem reading from config file: %v", err)
	}

	router := mux.NewRouter()

	for _, override := range settings.Overrides {
		router.HandleFunc(override.From, handler.ForwardTo(override.To))
	}

	router.HandleFunc("/{pathname:.*}", handler.Forward(settings.Default))

	n := negroni.New()
	n.Use(delay.Middleware{})
	n.Use(stdoutlog.Middleware{})
	n.UseHandler(router)

	listenHost := "localhost"
	if allowAnyHostToConnect {
		listenHost = ""
	}

	displayServerInfo(listenHost, listenPort, configFilename, settings)

	listenAddress := fmt.Sprintf("%v:%v", listenHost, listenPort)
	n.Run(listenAddress)
}

func showNiceExitMessage() {
	/* \b is the equivalent of hitting the back arrow. With the two following
	   space characters they serve to hide the "^C" that is printed when
	   ctrl-c is typed.
	*/
	fmt.Printf("\b\b  \n[ctrl-c] %v is shutting down\n", PROJECT_NAME)
}

func getCommandLineArgs() (allowAnyHostToConnect bool, port int, config string) {
	const DEFAULT_PORT = 8000
	const DEFAULT_CONFIG = ""

	flag.BoolVar(&allowAnyHostToConnect, "a", false, "Use to allow any ip address (any host) to connect. Default allows ony localhost.")
	flag.IntVar(&port, "port", DEFAULT_PORT, "Port on which to listen for connections.")
	flag.StringVar(&config, "c", DEFAULT_CONFIG, "Config file with overrides.")

	flag.Parse()

	/* Don't accept any positional command line arguments. flag.NArgs()
	   counts only non-flag arguments. */
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(EXIT_USAGE_FAILURE)
	}

	return
}

func displayServerInfo(listenHost string, listenPort int, configFilename string, settings *ProgramConfig) {
	visibleTo := listenHost
	if visibleTo == "" {
		visibleTo = "All ip addresses"
	}

	fmt.Printf("%v is running.\n\n", PROJECT_NAME)
	fmt.Printf("Visible to: %v\n", visibleTo)
	fmt.Printf("Port: %v\n", listenPort)
	fmt.Printf("Configuration file: %v\n\n", configFilename)
	fmt.Printf("Proxy default: %v\n", settings.Default)
	fmt.Println("Proxy overrides:")
	for _, override := range settings.Overrides {
		fmt.Printf("  %v -> %v\n", override.From, override.To)
	}
	fmt.Println()
	fmt.Printf("Hit [ctrl-c] to quit\n")
}
