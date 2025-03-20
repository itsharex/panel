package data

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tnb-labs/panel/internal/app"
	"github.com/tnb-labs/panel/internal/biz"
	"github.com/tnb-labs/panel/pkg/shell"
)

type containerComposeRepo struct{}

func NewContainerComposeRepo() biz.ContainerComposeRepo {
	return &containerComposeRepo{}
}

// List 列出所有编排文件名
func (r *containerComposeRepo) List() ([]string, error) {
	dir := filepath.Join(app.Root, "server", "compose")
	var names []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		names = append(names, filepath.Base(path))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return names, nil
}

// Get 获取编排文件和环境变量内容
func (r *containerComposeRepo) Get(name string) (string, string, error) {
	content, _ := os.ReadFile(filepath.Join(app.Root, "server", "compose", name, "docker-compose.yml"))
	env, _ := os.ReadFile(filepath.Join(app.Root, "server", "compose", name, ".env"))
	return string(content), string(env), nil // 有意忽略错误，这样可以允许新建文件
}

// Create 创建编排文件
func (r *containerComposeRepo) Create(name, compose, env string) error {
	dir := filepath.Join(app.Root, "server", "compose", name)
	if err := os.MkdirAll(dir, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "docker-compose.yml"), []byte(compose), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(env), 0644); err != nil {
		return err
	}

	return nil
}

// Up 启动编排
func (r *containerComposeRepo) Up(name string, force bool) error {
	file := filepath.Join(app.Root, "server", "compose", name, "docker-compose.yml")
	cmd := "docker compose -f %s up -d"
	if force {
		cmd += " --pull always" // 强制拉取镜像
	}
	_, err := shell.Execf(cmd, file)
	return err
}

// Down 停止编排
func (r *containerComposeRepo) Down(name string) error {
	file := filepath.Join(app.Root, "server", "compose", name, "docker-compose.yml")
	_, err := shell.Execf("docker compose -f %s down", file)
	return err
}

// Remove 删除编排
func (r *containerComposeRepo) Remove(name string) error {
	dir := filepath.Join(app.Root, "server", "compose", name)
	return os.RemoveAll(dir)
}
