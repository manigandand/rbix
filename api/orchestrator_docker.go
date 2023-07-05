package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/manigandand/adk/errors"
)

type docker struct {
	// Sqrx cluster network bridge
	// we expect that the network is already created
	sqrxNetwork types.NetworkResource
}

// NewDockerOrchestrator - create new docker orchestrator
func NewDockerOrchestrator() (Orchestrator, *errors.AppError) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.InternalServer("could not create docker client: " + err.Error())
	}
	defer cli.Close()

	// fetch sqrx-network
	netw, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return nil, errors.InternalServer("could not get network list: " + err.Error())
	}
	if len(netw) == 0 {
		return nil, errors.InternalServer("no network found")
	}

	doc := &docker{}

	found := false
	for _, n := range netw {
		if n.Name == "sqrx-network" {
			doc.sqrxNetwork = n
			found = true
			break
		}
	}

	if !found {
		return nil, errors.InternalServer("sqrx-network not found")
	}

	return doc, nil
}

// StartRBIInstance - start new sqrx-rbi instance
func (d *docker) StartRBIInstance(ctx context.Context, containerUniqeID string,
) (*ContainerInfo, *errors.AppError) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.InternalServer("could not create docker client: " + err.Error())
	}
	defer cli.Close()

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
		return nil, errors.InternalServer("could not create container: " + err.Error())
	}
	log.Println("container created: ", resp)
	containerInfo.CreatedAt = time.Now()

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, errors.InternalServer("could not start container: " + err.Error())
	}
	containerInfo.StartedAt = time.Now()

	// connect container to the network
	if err := cli.NetworkConnect(ctx, d.sqrxNetwork.ID, resp.ID, nil); err != nil {
		return nil, errors.InternalServer("could not connect container to network: " + err.Error())
	}
	containerInfo.DocContainer = resp
	containerInfo.ValidTill = time.Now().Add(10 * time.Minute)
	containerInfo.Session = fmt.Sprintf("%s/%s/ws", SqrxWSLoadbalncerHost, containerUniqeID)
	containerInfo.TerminationToken = uuid.New().String()

	// save to db
	db.SaveContainerInfo(containerUniqeID, containerInfo)
	db.SaveTerminationTokenInfo(containerInfo.TerminationToken, containerUniqeID)

	// spinup background process to delete the container after 10 minutes
	go func() {
		<-time.After(10 * time.Minute)
		if err := d.deleteContainer(context.Background(), containerInfo); err.NotNil() {
			log.Println("couldn't able to delete container: ", err.Error())
		}
	}()

	return containerInfo, nil
}

// DestroyRBIInstance - destroy sqrx-rbi instance
func (d *docker) DestroyRBIInstance(ctx context.Context, terminationToken string) *errors.AppError {
	cInfo, err := db.GetContainerInfoByTermToken(terminationToken)
	if err.NotNil() {
		return err
	}

	return d.deleteContainer(ctx, cInfo)
}

// Stop - stop orchestrator
func (d *docker) Stop() {}

func (d *docker) deleteContainer(ctx context.Context, cInfo *ContainerInfo) *errors.AppError {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.InternalServer("could not create docker client: " + err.Error())
	}
	defer cli.Close()

	// stop container, ingore error since we try force remove
	cli.ContainerStop(ctx, cInfo.DocContainer.ID, container.StopOptions{})

	if err := cli.ContainerRemove(ctx, cInfo.DocContainer.ID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return errors.InternalServer("could not remove container: " + err.Error())
	}

	db.DeleteContainer(cInfo.ContainerID)

	return nil
}
