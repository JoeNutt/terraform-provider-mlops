package infrastructure

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	Client *client.Client
}

func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &DockerClient{
		Client: cli,
	}, nil
}

func (d *DockerClient) DeployAgent(ctx context.Context, name, imageRef, groqAPIKey string) error {
	// 1. Pull the specified image
	reader, err := d.Client.ImagePull(ctx, imageRef, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", imageRef, err)
	}
	defer reader.Close()
	// We should consume the reader to ensure pull is complete
	_, _ = io.Copy(io.Discard, reader)

	// 2. Create a container
	config := &container.Config{
		Image: imageRef,
		Env: []string{
			fmt.Sprintf("GROQ_API_KEY=%s", groqAPIKey),
		},
	}

	resp, err := d.Client.ContainerCreate(ctx, config, nil, nil, nil, name)
	if err != nil {
		return fmt.Errorf("failed to create container %s: %w", name, err)
	}

	// 3. Start the container
	if err := d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container %s (ID: %s): %w", name, resp.ID, err)
	}

	return nil
}

func (d *DockerClient) RemoveAgent(ctx context.Context, name string) error {
	// 1. Stop the container
	// Using a timeout of 10 seconds (standard default)
	stopOptions := container.StopOptions{}
	if err := d.Client.ContainerStop(ctx, name, stopOptions); err != nil {
		// If the container is already stopped, we ignore the error
		if !client.IsErrNotFound(err) {
			return fmt.Errorf("failed to stop container %s: %w", name, err)
		}
	}

	// 2. Remove the container
	removeOptions := container.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	}
	if err := d.Client.ContainerRemove(ctx, name, removeOptions); err != nil {
		if !client.IsErrNotFound(err) {
			return fmt.Errorf("failed to remove container %s: %w", name, err)
		}
	}

	return nil
}
