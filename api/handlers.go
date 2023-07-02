package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/manigandand/adk/errors"
	"github.com/manigandand/adk/respond"
)

var (
	// Sqrx cluster network bridge
	// we expect that the network is already created
	sqrxNetwork types.NetworkResource
)

func init() {
	// init db
	db = &Store{
		containers:       make(map[string]*ContainerInfo),
		terminationToken: make(map[string]string),
	}

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

// creates a new sqrx-rbi container. we assume that the `image` is already pulled
func newSqureXSessionHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	// NOTE: ignoring the auth cases to validate the client already authenticated
	// and or check if any session is already active

	ctx := r.Context()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.InternalServer("could not create docker client: " + err.Error())
	}
	defer cli.Close()

	containerUniqeID := "dfv-" + uuid.New().String()

	containerInfo := &ContainerInfo{
		ContainerID: containerUniqeID,
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:    SqrxRbiImage,
		Hostname: containerUniqeID,
		Env: []string{
			"CONTAINER_ID=" + containerUniqeID,
		},
		Tty: false,
	}, nil, nil, nil, containerUniqeID)
	if err != nil {
		return errors.InternalServer("could not create container: " + err.Error())
	}
	log.Println("container created: ", resp)
	containerInfo.CreatedAt = time.Now()

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return errors.InternalServer("could not start container: " + err.Error())
	}
	containerInfo.StartedAt = time.Now()

	// connect container to the network
	if err := cli.NetworkConnect(ctx, sqrxNetwork.ID, resp.ID, nil); err != nil {
		return errors.InternalServer("could not connect container to network: " + err.Error())
	}
	containerInfo.Container = resp
	containerInfo.ValidTill = time.Now().Add(10 * time.Minute)
	containerInfo.Session = fmt.Sprintf("%s/%s/ws", SqrxWSLoadbalncerHost, containerUniqeID)
	containerInfo.TerminationToken = uuid.New().String()

	// save to db
	db.containers[containerUniqeID] = containerInfo
	db.terminationToken[containerInfo.TerminationToken] = containerUniqeID

	return respond.OK(w, containerInfo)
}

func getContainerStatHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	// ctx := r.Context()
	return respond.OK(w, nil)
}

func stopContainerHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	// ctx := r.Context()
	return respond.OK(w, nil)
}
