package claim

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/ravenix/hcloud-route-claimer/config"
	"net"
)

type Claimer struct {
	client *hcloud.Client
	ctx context.Context
}

func NewClaimer(client *hcloud.Client, ctx context.Context) *Claimer {
	return &Claimer{
		client,
		ctx,
	}
}

func (c *Claimer) Assign(claim *config.Claim) error {
	network, err := c.getNetwork(claim.Network)
	if err != nil {
		return err
	}

	toBeDeleted := claim.GetMatchingRoutes(network.Routes)
	for _, route := range toBeDeleted {
		_, _, err := c.client.Network.DeleteRoute(c.ctx, network, hcloud.NetworkDeleteRouteOpts{
			Route: route,
		})

		if err != nil {
			return fmt.Errorf("could not delete route: %v", err)
		}
	}

	for _, route := range claim.Routes {
		_, routeNet, _ := net.ParseCIDR(route)
		_, _, err := c.client.Network.AddRoute(c.ctx, network, hcloud.NetworkAddRouteOpts{
			Route: hcloud.NetworkRoute{
				Destination: routeNet,
				Gateway: claim.Gateway,
			},
		})

		if err != nil {
			return fmt.Errorf("could not create route: %v", err)
		}
	}

	return nil
}

func (c *Claimer) getNetwork(nameOrId string) (*hcloud.Network, error) {
	network, _, err := c.client.Network.Get(c.ctx, nameOrId)
	if err != nil {
		return nil, err
	}

	return network, nil
}
