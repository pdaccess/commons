package logs

import (
	"github.com/rs/zerolog/log"
	testcontainers "github.com/testcontainers/testcontainers-go"
)

type ContainerLogger struct {
	context string
}

func NewContainerLogger(context string) *ContainerLogger {
	return &ContainerLogger{
		context: context,
	}
}

func (c *ContainerLogger) Accept(l testcontainers.Log) {
	log.Info().
		Str("ctx", c.context).
		Msg(string(l.Content))
}
