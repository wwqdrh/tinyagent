package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type VolumeSpec struct {
	Name    string   `json:"name"`
	Driver  string   `json:"driver"`
	Options []string `json:"options"`
	Labels  []string `json:"labels"`
}

func VolumeList(filter volume.ListOptions) (volume.ListResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return volume.ListResponse{}, err
	}
	defer cli.Close()

	return cli.VolumeList(context.Background(), filter)
}

func VolumeAdd(spec VolumeSpec) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	var array []filters.KeyValuePair
	array = append(array, filters.Arg("name", spec.Name))
	vos, _ := cli.VolumeList(context.TODO(), volume.ListOptions{Filters: filters.NewArgs(array...)})
	if len(vos.Volumes) != 0 {
		for _, v := range vos.Volumes {
			if v.Name == spec.Name {
				return ErrRecordExist
			}
		}
	}
	options := volume.CreateOptions{
		Name:       spec.Name,
		Driver:     spec.Driver,
		DriverOpts: stringsToMap(spec.Options),
		Labels:     stringsToMap(spec.Labels),
	}
	if _, err := cli.VolumeCreate(context.TODO(), options); err != nil {
		return err
	}
	return nil
}

func VolumeDelete(name string, force bool) error {
	return withCli(func(cli *client.Client) error {
		return cli.VolumeRemove(context.Background(), name, force)
	})
}

func VolumeDeleteByIDs(ids ...string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	for _, id := range ids {
		if err := cli.VolumeRemove(context.TODO(), id, true); err != nil {
			if strings.Contains(err.Error(), "volume is in use") {
				return errors.Wrapf(err, "id: %s", id)
			}
			return err
		}
	}
	return nil
}
