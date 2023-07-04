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
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/manigandand/adk/errors"
	"github.com/manigandand/adk/respond"
)

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
	db.SaveContainerInfo(containerUniqeID, containerInfo)
	db.SaveTerminationTokenInfo(containerInfo.TerminationToken, containerUniqeID)

	// spinup background process to delete the container after 10 minutes
	go func() {
		<-time.After(10 * time.Minute)
		if err := deleteContainer(containerInfo); err.NotNil() {
			log.Println("could not delete container: ", err.Error())
		}
	}()

	return respond.OK(w, containerInfo)
}

func getContainerStatHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	containerID := chi.URLParam(r, "container_id")

	cInfo, err := db.GetContainerInfo(containerID)
	if err.NotNil() {
		return err
	}

	return respond.OK(w, cInfo)
}

func stopContainerHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	token := chi.URLParam(r, "termination_token")
	cInfo, err := db.GetContainerInfoByTermToken(token)
	if err.NotNil() {
		return err
	}

	if err := deleteContainer(cInfo); err.NotNil() {
		return err
	}

	return respond.OK(w, respond.Msg("container deleted"))
}

func deleteContainer(cInfo *ContainerInfo) *errors.AppError {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.InternalServer("could not create docker client: " + err.Error())
	}
	defer cli.Close()

	// stop container, ingore error since we try force remove
	cli.ContainerStop(ctx, cInfo.Container.ID, container.StopOptions{})

	if err := cli.ContainerRemove(ctx, cInfo.Container.ID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return errors.InternalServer("could not remove container: " + err.Error())
	}

	db.DeleteContainer(cInfo.ContainerID)

	return nil
}
