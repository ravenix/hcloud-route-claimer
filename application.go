package main

import (
	"context"
	"flag"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/ravenix/hcloud-route-claimer/claim"
	"github.com/ravenix/hcloud-route-claimer/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	configPathPtr := flag.String("config", "", "path to config file")
	flag.Parse()

	if len(*configPathPtr) == 0 {
		log.Error("You must supply the path to a configuration file.")
		return
	}

	config, err := config.Load(*configPathPtr)
	if err != nil {
		log.Errorf("Failed to load config file: %v", err)
		return
	}

	client := hcloud.NewClient(hcloud.WithToken(config.HCloud.Token))
	ctx := context.Background()

	claimer := claim.NewClaimer(client, ctx)

	for _, claim := range config.Claims {
		err := claimer.Assign(claim)
		if err != nil {
			log.Warnf("Could not claim routes %v in network %s via %s: %v", claim.Routes, claim.Network, claim.Gateway, err)
		}
	}

	log.Info("successfully applied claims to hcloud")
}
