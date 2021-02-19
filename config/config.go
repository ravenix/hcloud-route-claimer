package config

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
	"net"
)

type Config struct {
	HCloud HCloud `json:"hcloud"`
	Claims []*Claim `json:"claims"`
}

type Claim struct {
	Network string      `json:"network"`
	Gateway net.IP      `json:"gateway"`
	Routes  []string `json:"routes,omitempty"`
}

type HCloud struct {
	Token string `json:"token"`
}


func (c *Claim) GetMatchingRoutes(hcloudRoutes []hcloud.NetworkRoute) []hcloud.NetworkRoute {
	matchingRoutes := make([]hcloud.NetworkRoute, 0)

	for _, hcloudRoute := range hcloudRoutes {
		for _, route := range c.Routes {
			_, routeNet, _ := net.ParseCIDR(route)
			routeMask, _ := routeNet.Mask.Size()
			hcloudRouteMask, _ := hcloudRoute.Destination.Mask.Size()
			if routeMask <= hcloudRouteMask && routeNet.Contains(hcloudRoute.Destination.IP) {
				matchingRoutes = append(matchingRoutes, hcloudRoute)
			}
		}
	}

	return matchingRoutes
}