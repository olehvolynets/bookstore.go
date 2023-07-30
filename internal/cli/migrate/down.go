package migrate

import (
	"github.com/spf13/cobra"

	"bookstore/log"
)

var downCommand = &cobra.Command{
	Use:   "down [COUNT]",
	Args:  cobra.RangeArgs(0, 1),
	Short: "Roll back migrations",
	Run:   downCb,
}

func downCb(_ *cobra.Command, args []string) {
	var count uint64 = 0
	if len(args) > 0 {
		count = parseCount(args[0])
	}

	migs, err := migrEngine.Down(count)
	if err != nil {
		log.Error().Err(err).Msg("error while rolling back migrations")
	}

	for _, m := range migs {
		log.Info().Str("version", m.Version).Msg("successfully rolled back a migration")
	}
}
