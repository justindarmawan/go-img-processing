package bootstrap

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
)

type Container struct {
	ctx context.Context
}

func Init() *Container {
	c := &Container{
		ctx: context.Background(),
	}
	c.initConfig()
	return c
}

func (c *Container) initConfig() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
	}
}
