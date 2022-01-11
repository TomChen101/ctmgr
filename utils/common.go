package utils

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

/*获取当前机器运行的所有容器names
*@return1 err error
*@return2 names []string
 */
func GetAllContainerInsNames() (error, []string) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err, nil
	}
	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err, nil
	}
	names := make([]string, 0)
	for _, container := range containers {
		if len(container.Names) == 0 {
			continue
		}
		names = append(names, container.Names[0][1:])
	}
	return nil, names
}
