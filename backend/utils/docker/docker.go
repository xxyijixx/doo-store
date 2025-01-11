package docker

import (
	"context"
	"doo-store/backend/utils/cmd"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Client struct {
	cli *client.Client
}

func NewClient() (Client, error) {
	// query Docker sock path
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithHost("unix:///var/run/docker.sock"), client.WithAPIVersionNegotiation())
	if err != nil {
		return Client{}, err
	}

	return Client{
		cli: cli,
	}, nil
}

func (c Client) Close() {
	_ = c.cli.Close()
}

func (c Client) GetClient() *client.Client {
	return c.cli
}

func NewDockerClient() (*client.Client, error) {
	// query Docker sock path
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithHost("unix:///var/run/docker.sock"), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (c Client) ListContainersByName(names []string) ([]types.Container, error) {
	var (
		options  container.ListOptions
		namesMap = make(map[string]bool)
		res      []types.Container
	)
	options.All = true
	if len(names) > 0 {
		var array []filters.KeyValuePair
		for _, n := range names {
			namesMap["/"+n] = true
			array = append(array, filters.Arg("name", n))
		}
		options.Filters = filters.NewArgs(array...)
	}
	containers, err := c.cli.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}
	for _, con := range containers {
		if _, ok := namesMap[con.Names[0]]; ok {
			res = append(res, con)
		}
	}
	return res, nil
}
func (c Client) ListAllContainers() ([]types.Container, error) {
	var (
		options container.ListOptions
	)
	options.All = true
	containers, err := c.cli.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (c Client) CreateNetwork(name string) error {
	_, err := c.cli.NetworkCreate(context.Background(), name, network.CreateOptions{
		Driver: "bridge",
	})
	return err
}

func (c Client) DeleteImage(imageID string) error {
	if _, err := c.cli.ImageRemove(context.Background(), imageID, image.RemoveOptions{Force: true}); err != nil {
		return err
	}
	return nil
}

func (c Client) InspectContainer(containerID string) (types.ContainerJSON, error) {
	return c.cli.ContainerInspect(context.Background(), containerID)
}

func (c Client) PullImage(imageName string, force bool) (string, error) {
	if !force {
		exist, err := c.CheckImageExist(imageName)
		if err != nil {
			return "", err
		}
		if exist {
			return "", nil
		}
	}
	reader, err := c.cli.ImagePull(context.Background(), imageName, image.PullOptions{})
	if err != nil {
		return "", err
	}
	stdout, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

func (c Client) GetImageIDByName(imageName string) (string, error) {
	filter := filters.NewArgs()
	filter.Add("reference", imageName)
	list, err := c.cli.ImageList(context.Background(), image.ListOptions{
		Filters: filter,
	})
	if err != nil {
		return "", err
	}
	if len(list) > 0 {
		return list[0].ID, nil
	}
	return "", nil
}

func (c Client) GetImageFirstExposedPortByName(imageName string) (int, error) {
	imageInspect, _, err := c.cli.ImageInspectWithRaw(context.Background(), imageName)
	if err != nil {
		return 0, err
	}
	portNum := 0
	for port := range imageInspect.Config.ExposedPorts {
		portStr := strings.Split(string(port), "/")[0]
		portNum, err = strconv.Atoi(portStr)
		if err != nil {
			return 0, err
		}
		break
	}
	return portNum, nil
}

func (c Client) CheckImageExist(imageName string) (bool, error) {
	filter := filters.NewArgs()
	filter.Add("reference", imageName)
	list, err := c.cli.ImageList(context.Background(), image.ListOptions{
		Filters: filter,
	})
	if err != nil {
		return false, err
	}
	return len(list) > 0, nil
}

func (c Client) NetworkExist(name string) bool {
	var options network.ListOptions
	options.Filters = filters.NewArgs(filters.Arg("name", name))
	networks, err := c.cli.NetworkList(context.Background(), options)
	if err != nil {
		return false
	}
	return len(networks) > 0
}

func (c Client) CopyFileToContainer(containerId, srcFile, dstFile string) error {
	cli, err := NewDockerClient()
	if err != nil {
		return nil
	}
	defer cli.Close()
	_, err = os.Open(srcFile)
	if err != nil {
		return err
	}

	_, err = cmd.Execf("docker cp %s %s:%s", srcFile, containerId, dstFile)

	return err
}

func (c Client) RemoveFileFormContainer(containerId, filePath string) error {
	// 创建一个执行命令的配置
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"rm", "-f", filePath},
	}

	// 创建执行命令
	execIDResp, err := c.cli.ContainerExecCreate(context.Background(), containerId, execConfig)
	if err != nil {
		return fmt.Errorf("error creating exec: %v", err)
	}

	// 执行命令
	execAttachResp, err := c.cli.ContainerExecAttach(context.Background(), execIDResp.ID, container.ExecStartOptions{})
	if err != nil {
		return fmt.Errorf("error attaching to exec: %v", err)
	}
	defer execAttachResp.Close()

	// 读取命令输出
	outputDone := make(chan error)
	go func() {
		_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, execAttachResp.Reader)
		outputDone <- err
	}()

	// 等待命令执行完成
	err = <-outputDone
	if err != nil && err != io.EOF {
		return fmt.Errorf("error during command execution: %v", err)
	}

	return nil
}

func (c Client) MoveFileInContainer(containerId, srcPath, dstPath string) error {
	// 创建一个执行命令的配置
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"mv", srcPath, dstPath},
	}

	// 创建执行命令
	execIDResp, err := c.cli.ContainerExecCreate(context.Background(), containerId, execConfig)
	if err != nil {
		return fmt.Errorf("error creating exec: %v", err)
	}

	// 执行命令
	execAttachResp, err := c.cli.ContainerExecAttach(context.Background(), execIDResp.ID, container.ExecStartOptions{})
	if err != nil {
		return fmt.Errorf("error attaching to exec: %v", err)
	}
	defer execAttachResp.Close()

	// 读取命令输出
	outputDone := make(chan error)
	go func() {
		_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, execAttachResp.Reader)
		outputDone <- err
	}()

	// 等待命令执行完成
	err = <-outputDone
	if err != nil && err != io.EOF {
		return fmt.Errorf("error during command execution: %v", err)
	}

	return nil
}

func (c Client) FileExistsInContainer(containerId, filePath string) (bool, error) {
	// 创建一个执行命令的配置
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"test", "-f", filePath},
	}

	// 创建执行命令
	execIDResp, err := c.cli.ContainerExecCreate(context.Background(), containerId, execConfig)
	if err != nil {
		return false, fmt.Errorf("error creating exec: %v", err)
	}

	// 执行命令
	execAttachResp, err := c.cli.ContainerExecAttach(context.Background(), execIDResp.ID, container.ExecStartOptions{})
	if err != nil {
		return false, fmt.Errorf("error attaching to exec: %v", err)
	}
	defer execAttachResp.Close()

	// 读取命令输出
	outputDone := make(chan error)
	go func() {
		_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, execAttachResp.Reader)
		outputDone <- err
	}()

	// 等待命令执行完成
	err = <-outputDone
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("error during command execution: %v", err)
	}

	// 获取命令执行状态
	inspectResp, err := c.cli.ContainerExecInspect(context.Background(), execIDResp.ID)
	if err != nil {
		return false, fmt.Errorf("error inspecting exec: %v", err)
	}

	// ExitCode为0表示文件存在
	return inspectResp.ExitCode == 0, nil
}

func (c Client) MoveFileWithCheck(containerId, srcPath, dstPath string) error {
	// 检查源文件是否存在
	exists, err := c.FileExistsInContainer(containerId, srcPath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("source file does not exist: %s", srcPath)
	}

	// 检查目标文件是否存在
	exists, err = c.FileExistsInContainer(containerId, dstPath)
	if err != nil {
		return err
	}
	if exists {
		// 如果目标文件存在，生成一个新的文件名
		ext := filepath.Ext(dstPath)
		base := strings.TrimSuffix(dstPath, ext)
		timestamp := time.Now().Format("20060102150405")
		dstPath = fmt.Sprintf("%s_%s%s", base, timestamp, ext)
	}

	// 移动文件
	return c.MoveFileInContainer(containerId, srcPath, dstPath)
}
