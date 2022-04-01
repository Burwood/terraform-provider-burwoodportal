package main

import (
	"burwoodportal/burwoodportal"
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"log"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: burwoodportal.Provider}

	if debugMode {
		err := plugin.Debug(context.Background(), "burwood.com/portal/burwoodportal", opts)
		if err != nil {
			log.Print("Error with debugger")
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
