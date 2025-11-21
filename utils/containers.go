package utils

import (
	"CipherOps/config"
)

// Container config
type ContainerConfig struct {
	ID string
	ImageName string
	Port string
	EnvVars map[string]string
}

var ContainerRegistry = map[string]ContainerConfig {
	"postgres": {},
	"mysql": {},
	"web": {}
}

// Define Docker actions 
type DockerManager interface {
	Create(service string) (string, error)
	Remove(containerID string) error
	RemoveImage(imageName string) error
	Start(containerID string) error
	Stop(ContainerID string) error
	ListContainers(all bool) ([]ContainerInfo, error)
    ListImages() ([]ImageInfo, error)
}

func (*config.Config) Create(service string) (string, error) {
	
}