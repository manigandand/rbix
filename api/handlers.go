package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/manigandand/adk/errors"
	"github.com/manigandand/adk/respond"
)

// creates a new rbix-rbi container. we assume that the `image` is already pulled
func newSqureXSessionHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	// NOTE: ignoring the auth cases to validate the client already authenticated
	// and or check if any session is already active

	ctx := r.Context()
	containerUniqeID := "dfv-" + uuid.New().String()
	containerInfo, err := orchestrator.StartRBIInstance(ctx, containerUniqeID)
	if err.NotNil() {
		return err
	}

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
	ctx := r.Context()
	token := chi.URLParam(r, "termination_token")
	if err := orchestrator.DestroyRBIInstance(ctx, token); err.NotNil() {
		return err
	}

	return respond.OK(w, respond.Msg("container deleted"))
}
