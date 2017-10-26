package docker

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/normegil/connectionutils"
	"github.com/normegil/interval"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const DOCKER_ADDRESS string = "127.0.0.1"
const maxWaitTime = 5 * time.Second
const stepWaitTime = 10 * time.Millisecond

type Options struct {
	Name   string
	Image  string
	Ports  []PortBinding
	Logger Logger
}

type PortBinding struct {
	Protocol         string
	Internal         int
	ExternalInterval string
}

type ContainerInfo struct {
	Address net.IP
	Ports   map[PortBinding]int
}

func New(options Options) (*ContainerInfo, func() error, error) {
	var l Logger = &DefaultLogger{}
	if nil != options.Logger {
		l = options.Logger
	}

	l.Print("New docker client from environment")
	client, err := docker.NewEnvClient()
	if nil != err {
		return nil, nil, errors.Wrap(err, "MongoDB: Could not create docker client")
	}

	if err = pullImage(client, options.Image); err != nil {
		return nil, nil, errors.Wrap(err, "Downloading image: "+options.Image)
	}

	ip := net.ParseIP(DOCKER_ADDRESS)
	if err := checkOptions(options); err != nil {
		return nil, nil, errors.New("Docker instance cannot be used without a external port")
	}

	containerName := options.Name + "-" + uuid.NewV4().String()
	dockerPorts, err := selectPorts(ip, options.Ports)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Selecting ports")
	}
	portBindings := toDockerPortBindings(ip, dockerPorts)
	l.Printf("Port Bindings: %+v", portBindings)

	l.Printf("Creating container: %+v", containerName)
	ctx := context.Background()
	containerInfo, err := client.ContainerCreate(ctx, &container.Config{
		Image:        options.Image,
		ExposedPorts: toExposedPorts(options.Ports),
	}, &container.HostConfig{
		PortBindings: portBindings,
	}, nil, containerName)
	if nil != err {
		return nil, nil, errors.Wrap(err, "Could not create container ("+containerName+")")
	}
	for _, warning := range containerInfo.Warnings {
		l.Print(warning)
	}

	l.Print("Starting container: " + containerName)
	if err := client.ContainerStart(ctx, containerInfo.ID, types.ContainerStartOptions{}); nil != err {
		return nil, nil, errors.Wrap(err, "Could not start container ("+containerName+")")
	}

	l.Print("Waiting for container: " + containerName)
	reachablePorts := dockerPorts[options.Ports[0]]
	if err := waitContainer(client, containerInfo.ID, DOCKER_ADDRESS+":"+strconv.Itoa(reachablePorts), maxWaitTime); nil != err {
		return nil, nil, errors.Wrap(err, "Container not started withing time limit")
	}
	l.Print("Container started: " + containerName)

	return &ContainerInfo{
			Address: ip,
			Ports:   dockerPorts,
		}, func() error {
			l.Print("Removing container: " + containerName)
			ctx := context.Background()
			if err := client.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); nil != err {
				return errors.Wrap(err, "MongoDB: Could not remove "+containerName)
			}
			return nil
		}, nil
}

func pullImage(client *docker.Client, id string) error {
	images, err := client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return errors.Wrap(err, "Listing images")
	}
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == id {
				return nil
			}
		}
	}

	closer, err := client.ImagePull(context.Background(), id, types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "Pulling image: "+id)
	}
	defer closer.Close()
	return nil
}

func checkOptions(options Options) error {
	if nil == options.Ports || 0 == len(options.Ports) {
		return errors.New("At least one port should be open for external communication")
	}
	return nil
}

func toExposedPorts(ports []PortBinding) nat.PortSet {
	exposed := make(map[nat.Port]struct{})
	for _, binding := range ports {
		exposed[nat.Port(strconv.Itoa(binding.Internal)+"/"+binding.Protocol)] = struct{}{}
	}
	return nat.PortSet(exposed)
}

func selectPorts(address net.IP, possiblePorts []PortBinding) (map[PortBinding]int, error) {
	used := make([]int, 0)
	toReturn := make(map[PortBinding]int)
	for _, binding := range possiblePorts {
		interval, err := interval.ParseIntervalInteger(binding.ExternalInterval)
		if err != nil {
			return nil, errors.Wrapf(err, "Parsing %s", binding.ExternalInterval)
		}
		selected := connectionutils.SelectPortExcluding(address, *interval, used)
		toReturn[binding] = selected.Port
	}
	return toReturn, nil
}

func toDockerPortBindings(address net.IP, ports map[PortBinding]int) map[nat.Port][]nat.PortBinding {
	toReturn := make(map[nat.Port][]nat.PortBinding)
	for binding, selectedPort := range ports {
		toReturn[nat.Port(strconv.Itoa(binding.Internal)+"/"+binding.Protocol)] = []nat.PortBinding{
			{
				//HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(selectedPort), // + "/" + binding.Protocol,
			},
		}
	}
	return toReturn
}

func waitContainer(client *docker.Client, containerID string, hostport string, maxWait time.Duration) error {
	if err := waitStarted(client, containerID, maxWait); nil != err {
		return err
	}
	if err := waitReachable(hostport, maxWait); nil != err {
		return err
	}
	return nil
}

func waitReachable(hostport string, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		c, err := net.Dial("tcp", hostport)
		if nil == err {
			return c.Close()
		}
		time.Sleep(stepWaitTime)
	}
	return fmt.Errorf("Could not reach %s {WaitingTime: %+v}", hostport, maxWait)
}

func waitStarted(client *docker.Client, containerID string, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		ctx := context.Background()
		c, err := client.ContainerInspect(ctx, containerID)
		if err != nil {
			break
		}
		if c.State.Running {
			return nil
		}
		time.Sleep(stepWaitTime)
	}
	return fmt.Errorf("Container not started: %s {WaitingTime: %+v}", containerID, maxWait)
}
