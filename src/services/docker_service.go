package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

type DockerService struct {
	cli *client.Client
}

func NewDockerService() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerService{cli: cli}, nil
}

func (ds *DockerService) ListContainers() ([]types.Container, error) {
	containers, err := ds.cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (ds *DockerService) StartContainer(containerID string) error {
	return ds.cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
}

func (ds *DockerService) StopContainer(containerID string) error {
	return ds.cli.ContainerStop(context.Background(), containerID, nil)
}

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

func (ds *DockerService) GetContainerLogs(containerID, since, until string) (string, error) {
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "all",
		Since:      since,
		Until:      until,
	}
	reader, err := ds.cli.ContainerLogs(context.Background(), containerID, options)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var buf bytes.Buffer
	_, err = stdcopy.StdCopy(&buf, &buf, reader)
	if err != nil && err != io.EOF {
		return "", err
	}
	return buf.String(), nil
}

func (ds *DockerService) GetLogsForProject(projectName, since, until string) (map[string]string, error) {
	containers, err := ds.ListContainers()
	if err != nil {
		return nil, err
	}

	logsMap := make(map[string]string)
	for _, container := range containers {
		if container.Labels["com.docker.compose.project"] == projectName {
			log, err := ds.GetContainerLogs(container.ID, since, until)
			if err != nil {
				return nil, err
			}
			logsMap[container.ID] = log
		}
	}
	return logsMap, nil
}

func (ds *DockerService) GetLogsForService(projectName, serviceName, since, until string) (map[string]string, error) {
	containers, err := ds.ListContainers()
	if err != nil {
		return nil, err
	}

	logsMap := make(map[string]string)
	for _, container := range containers {
		if container.Labels["com.docker.compose.project"] == projectName &&
			container.Labels["com.docker.compose.service"] == serviceName {
			log, err := ds.GetContainerLogs(container.ID, since, until)
			if err != nil {
				return nil, err
			}
			logsMap[container.ID] = log
		}
	}
	return logsMap, nil
}

func (ds *DockerService) GetLogs(service, since, until string) (map[string]string, error) {
	if since != "" {
		if seconds, err := strconv.Atoi(since); err == nil {
			if seconds < 1000000000 {
				t := time.Now().Unix() - int64(seconds)
				since = strconv.FormatInt(t, 10)
			}
		}
	}
	if service == "" {
		return ds.GetLogsForProject("backend", since, until)
	}
	return ds.GetLogsForService("backend", service, since, until)
}

func (ds *DockerService) GetContainerStats(containerID string) (types.StatsJSON, error) {
	resp, err := ds.cli.ContainerStats(context.Background(), containerID, false) // false para obtener un snapshot único
	if err != nil {
		return types.StatsJSON{}, err
	}
	defer resp.Body.Close()

	var stats types.StatsJSON
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return types.StatsJSON{}, err
	}
	return stats, nil
}
