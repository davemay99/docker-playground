package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/kr/pretty"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	if len(containers) == 0 {
		fmt.Println("No running Docker containers")
		return
	}

	fmt.Println("Running container count", len(containers))

	containerJSON := Docker_ContainersToJSON()
	fmt.Println("Containers formatted as JSON")
	fmt.Printf("%s", containerJSON)

	for _, container := range containers {
		fmt.Println(container.ID)
		logs := Docker_LogsFromContainerID(container.ID)
		fmt.Println(logs)
	}
}

func Docker_LogsFromContainerID(id string) string {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true} // Timestamps, Details, Follow
	out, err := cli.ContainerLogs(ctx, id, options)
	if err != nil {
		panic(err)
	}

	pretty.Print(out)
	// written, err := io.Copy(os.Stdout, out)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%d bytes copied to os.Stdout\n", written)

	// var dstout, dsterr io.Writer
	// written, err := stdcopy.StdCopy(dstout, dsterr, out)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%d bytes demultiplexed to dstout and dsterr\n", written)

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	outString := buf.String()

	return outString
}

func Docker_ContainersToJSON() string {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(containers)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", bytes)
}
