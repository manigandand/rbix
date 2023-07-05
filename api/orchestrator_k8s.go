package main

import (
	"context"

	"github.com/manigandand/adk/errors"
)

type k8s struct {
}

// NewK8sOrchestrator - create new k8s orchestrator
func NewK8sOrchestrator() (Orchestrator, *errors.AppError) {
	return &k8s{}, nil
}

// StartRBIInstance - start new sqrx-rbi instance
func (k *k8s) StartRBIInstance(ctx context.Context, containerUniqeID string,
) (*ContainerInfo, *errors.AppError) {
	return nil, nil
}

// DestroyRBIInstance - destroy sqrx-rbi instance
func (k *k8s) DestroyRBIInstance(ctx context.Context, terminationToken string) *errors.AppError {
	return nil
}

// Stop - stop orchestrator
func (k *k8s) Stop() {}
