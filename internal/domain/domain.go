package domain

import "os"

type SftpContent struct {
	Path string
	Info os.FileInfo
	Dir  bool
	File bool
}

type ProfileConfig struct {
	Operation      string `yaml:"operation"`
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	PrivateKey     string `yaml:"private_key"`
	LocalBasePath  string `yaml:"local_base_path"`
	RemoteBasePath string `yaml:"remote_base_path"`
}

type Config struct {
	Profiles map[string]ProfileConfig `yaml:"profiles"`
}
