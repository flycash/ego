package elog

import (
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
)

type Option func(c *Container)

type Container struct {
	config *Config
	name   string
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		panic(err)
	}
	c.name = key
	return c
}

func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	if eapp.IsDevelopmentMode() {
		c.config.Debug = true           // 调试模式，终端输出
		c.config.EnableAsync = false    // 调试模式，同步输出
		c.config.EnableAddCaller = true // 调试模式，增加行号输出
	}

	if c.config.encoderConfig == nil {
		c.config.encoderConfig = DefaultZapConfig()
	}

	if c.config.Debug {
		c.config.encoderConfig = DefaultDebugConfig()
	}

	if eapp.EnableLoggerAddApp() {
		c.config.fields = append(c.config.fields, FieldApp(eapp.Name()))
	}

	logger := newLogger(c.name, c.config)
	// 如果名字不为空，加载动态配置
	if c.name != "" {
		// c.name 为配置name
		logger.AutoLevel(c.name + ".level")
	}

	return logger
}
