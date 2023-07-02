package main

import (
	"net/http"

	"github.com/manigandand/adk/errors"
	"github.com/manigandand/adk/respond"
)

func newSqureXSessionHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	ctx := r.Context()
	return respond.OK(w, nil)
}

func getContainerStatHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	ctx := r.Context()
	return respond.OK(w, nil)
}

func stopContainerHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	ctx := r.Context()
	return respond.OK(w, nil)
}
