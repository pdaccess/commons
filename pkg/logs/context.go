package logs

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func LogWithContext(context string) func(values ...interface{}) {
	return func(values ...interface{}) {
		event := log.Info().
			Str("context", context)

		for i, value := range values {
			event.Interface(fmt.Sprintf("%d", i), value)
		}

		event.Send()
	}
}
