package main

import (
	"io/ioutil"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server   ServerConfig `yaml:server`
	Database DbConfig     `yaml:database`
	Logger   LoggerConfig `yaml:logger`
	Stat     StatConfig   `yaml:stat`
	Redis    RedisConf    `yaml:"redis"`
}

type ServerConfig struct {
	Endpoint         string        `yaml:"endpoint"`
	HttpReadTimeout  time.Duration `yaml:"http_read_timeout"`
	HttpWriteTimeout time.Duration `yaml:"http_write_timeout"`
	MaxHeadSize      uint32        `yaml:"max_head_size"`
}

type RedisConf struct {
	Servers               []string      `yaml:"servers"`
	ConnTimeout           time.Duration `yaml:"connTimeout"`
	ExpiredTime           time.Duration `yaml:"expiredTime"`
	EmptyCacheExpiredTime time.Duration `yaml:"emptyCacheExpiredTime"`
	Prefix                string        `yaml:"prefix"`
}

type DbConfig struct {
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	User           string        `yaml:"user"`
	Password       string        `yaml:"password"`
	DbName         string        `yaml:"db_name"`
	MaxOpenConnNum int           `yaml:"max_open_conn_num"`
	MaxIdleConnNum int           `yaml:"max_idle_conn_num"`
	KeepAliveTime  time.Duration `yaml:"keep_alive_time"`
	ConnTimeout    int           `yaml:"conn_timeout"`
	ReadTimeout    int           `yaml:"read_timeout"`
	WriteTimeout   int           `yaml:"write_timeout"`
}

type StatConfig struct {
	CycleSec      uint32 `yaml:"output_cycle"`
	StatFile      string `yaml:"stat_file"`
	RemoteMode    uint32 `yaml:"remote_mode"`
	QueueCapacity uint32 `yaml:"queue_capacity"`
	RoutineCount  uint32 `yaml:"routine_count"`
}

type LoggerConfig struct {
	Level            zap.AtomicLevel `yaml:"level"`
	OutputPaths      []string        `yaml:"outputPaths"`
	ErrorOutputPaths []string        `yaml:"errorOutputPaths"`
	Filename         string          `yaml:"filename"`
	MaxSize          int             `yaml:"max_size"`
	MaxBackups       int             `yaml:"max_backups"`
	MaxAge           int             `yaml:"max_age"`
}

func newConfig() *Config {
	return &Config{}
}

func (cfg *Config) load(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return err
	}

	return nil
}
