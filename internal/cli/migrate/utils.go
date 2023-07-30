package migrate

import (
	"os"
	"strconv"

	"bookstore/log"
)

func fetchEnv(name, fallback string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}

	return val
}

func parseCount(sCount string) uint64 {
	count, err := strconv.ParseUint(sCount, 10, strconv.IntSize)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read migrations count")
	}

	return count
}
