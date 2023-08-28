package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/wwqdrh/gokit/ostool"
)

func ImageExist(name string) error {
	images, err := ImageList(types.ImageListOptions{All: true})
	if err != nil {
		return err
	}

	for _, item := range images {
		for _, tag := range item.RepoTags {
			if tag == name {
				return nil
			}
		}
	}

	return ErrImageNotExist
}

func ImageList(opts types.ImageListOptions) ([]types.ImageSummary, error) {
	// cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	return cli.ImageList(context.Background(), opts)
}

func ImagePull(name string, opts types.ImagePullOptions) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	out, err := cli.ImagePull(context.Background(), name, opts)
	if err != nil {
		return "", err
	}

	buffer := bytes.Buffer{}
	buf := make([]byte, 100)
	for {
		n, err := out.Read(buf)
		if n > 0 {
			fmt.Println(string(buf[:n]))
			buffer.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
	}

	return buffer.String(), nil
}

func ImageDelete(name string, opts types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return cli.ImageRemove(context.Background(), name, opts)
}

func ImageBuild(dockerfile string, imageName string) error {
	workdir := path.Dir(dockerfile)
	return ostool.RunCmdStd(fmt.Sprintf("cd %s; docker build -t %s .", workdir, imageName))
}

type ImageFileEntry struct {
	Source string
	Target string
}

func ImageBuildFromStr(content string, imageName string, entrys []ImageFileEntry) error {
	dir, err := os.MkdirTemp("", "edge_board_build_")
	if err != nil {
		return errors.Wrapf(err, "创建临时文件夹失败")
	}
	defer func() { os.RemoveAll(dir) }()

	if err := os.WriteFile(path.Join(dir, "Dockerfile"), []byte(content), os.ModePerm); err != nil {
		return errors.Wrapf(err, "创建Dockerfile失败")
	}

	for _, item := range entrys {
		_, err := copyFile(path.Join(dir, item.Target), item.Source)
		if err != nil {
			return errors.Wrapf(err, "复制文件失败")
		}
	}

	return ImageBuild(path.Join(dir, "Dockerfile"), imageName)
}

func copyFile(dstFileName string, srcFileName string) (written int64, err error) {
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		return -1, errors.Wrapf(err, "open file err")
	}
	defer srcFile.Close()

	if err := os.MkdirAll(path.Dir(dstFileName), os.ModePerm); err != nil {
		return -1, err
	}
	dstFile, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return -1, errors.Wrapf(err, "open file err")
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}
