package migrate

import (
	"strconv"

	"bookstore/internal/logging"
)

func parseCount(sCount string) uint64 {
	log := logging.NewLoggerFromEnv()

	count, err := strconv.ParseUint(sCount, 10, strconv.IntSize)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read migrations count")
	}

	return count
}
