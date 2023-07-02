package main

import (
	"time"

	"github.com/docker/docker/api/types/container"
)

// ContainerInfo is the container information
type ContainerInfo struct {
	ContainerID      string                   `json:"-"`
	Container        container.CreateResponse `json:"-"`
	Session          string                   `json:"session"`
	TerminationToken string                   `json:"termination_token"`
	CreatedAt        time.Time                `json:"created_at"`
	StartedAt        time.Time                `json:"started_at"`
	ValidTill        time.Time                `json:"valid_till"`
}

var db *Store

// Store - local kv store
type Store struct {
	containers       map[string]*ContainerInfo
	terminationToken map[string]string
}
