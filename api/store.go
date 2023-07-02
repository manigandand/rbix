package main

import (
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/manigandand/adk/errors"
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
	mx               sync.RWMutex
}

// SaveContainerInfo - save container info
func (s *Store) SaveContainerInfo(cID string, cInfo *ContainerInfo) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	s.containers[cID] = cInfo

	return
}

// SaveTerminationTokenInfo - save container info
func (s *Store) SaveTerminationTokenInfo(terminationToken, cID string) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	s.terminationToken[terminationToken] = cID

	return
}

// GetContainerInfo - get container info
func (s *Store) GetContainerInfo(cID string) (*ContainerInfo, *errors.AppError) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	c, ok := s.containers[cID]
	if !ok {
		return nil, errors.NotFound("container not found")
	}

	return c, nil
}

// GetContainerInfoByTermToken - get container info by termination token
func (s *Store) GetContainerInfoByTermToken(token string) (*ContainerInfo, *errors.AppError) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	cID, ok := s.terminationToken[token]
	if !ok {
		return nil, errors.NotFound("invalid termination token")
	}

	c, ok := s.containers[cID]
	if !ok {
		return nil, errors.NotFound("session not found")
	}

	return c, nil
}

// DeleteContainer - delete container
func (s *Store) DeleteContainer(cID string) *errors.AppError {
	s.mx.RLock()
	defer s.mx.RUnlock()

	c, ok := s.containers[cID]
	if !ok {
		return errors.NotFound("session not found")
	}

	delete(s.containers, cID)
	delete(s.terminationToken, c.TerminationToken)

	return nil
}
