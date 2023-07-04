package main

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var (
	// Sqrx cluster network bridge
	// we expect that the network is already created
	sqrxNetwork types.NetworkResource
)

func initDockerSqrxNetwork() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	netw, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		panic("could not get network list: " + err.Error())
	}
	if len(netw) == 0 {
		panic("no network found")
	}

	found := false
	for _, n := range netw {
		if n.Name == "sqrx-network" {
			sqrxNetwork = n
			found = true
			break
		}
	}

	if !found {
		panic("sqrx-network not found")
	}
}
