package utils

import (
	"io"

	"github.com/gromey/proto-rest/logger"
)

func Closer(c io.Closer) {
	if err := c.Close(); err != nil {
		if logger.InLevel(logger.LevelError) {
			logger.Error(err)
		}
	}
}
