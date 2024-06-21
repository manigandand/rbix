package main

import (
	"context"
	"log"

	"github.com/manigandand/adk/errors"
)

// orchestrator - global orchestrator instance
var orchestrator Orchestrator

// Orchestrator - orchestrator interface
type Orchestrator interface {
	// StartRBIInstance - start new rbix-rbi instance
	StartRBIInstance(ctx context.Context, containerUniqeID string) (*ContainerInfo, *errors.AppError)

	// DestroyRBIInstance - destroy rbix-rbi instance
	DestroyRBIInstance(ctx context.Context, terminationToken string) *errors.AppError

	// Stop - stop orchestrator
	Stop()
}

// InitOrchechestrator - init orchestrator
func InitOrchechestrator() {
	// init docker orchestrator
	if Env == EnvLoclDocker {
		docOrch, err := NewDockerOrchestrator()
		if err.NotNil() {
			log.Fatal("could not create docker orchestrator: " + err.Error())
		}

		orchestrator = docOrch
		log.Println("docker orchestrator initialized 👍")
		return
	}

	if Env == EnvLocalK8s {
		k8sOrch, err := NewK8sOrchestrator()
		if err.NotNil() {
			log.Fatal("could not create k8s orchestrator: " + err.Error())
		}

		orchestrator = k8sOrch
		log.Println("k8s orchestrator initialized 👍")
		return
	}

	log.Fatal("unsupported env ", Env)
}
