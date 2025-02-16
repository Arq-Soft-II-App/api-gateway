package services

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerService struct {
	cli *client.Client
}

// crea una nueva instancia del servicio Docker.
func NewDockerService() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerService{cli: cli}, nil
}

// lista todos los contenedores (incluso los detenidos).
func (ds *DockerService) ListContainers() ([]types.Container, error) {
	containers, err := ds.cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// inicia un contenedor dado su ID.
func (ds *DockerService) StartContainer(containerID string) error {
	return ds.cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
}

// detiene un contenedor dado su ID.
func (ds *DockerService) StopContainer(containerID string) error {
	return ds.cli.ContainerStop(context.Background(), containerID, nil)
}

// crea un nuevo contenedor con la imagen y configuración dada.
func (ds *DockerService) CreateContainer(image string, name string, exposedPort string) (string, error) {
	port, err := nat.NewPort("tcp", exposedPort)
	if err != nil {
		return "", err
	}

	networkingConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"backend-network": {},
		},
	}

	resp, err := ds.cli.ContainerCreate(context.Background(),
		&container.Config{
			Image: image,
			ExposedPorts: nat.PortSet{
				port: struct{}{},
			},
		},
		nil,
		networkingConfig,
		nil,
		name,
	)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (ds *DockerService) RemoveContainer(containerID string) error {
	return ds.cli.ContainerRemove(
		context.Background(),
		containerID,
		types.ContainerRemoveOptions{Force: true, RemoveVolumes: true},
	)
}
