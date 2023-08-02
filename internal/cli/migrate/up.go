package migrate

import (
	"github.com/spf13/cobra"

	"bookstore/internal/log"
)

var upCommand = &cobra.Command{
	Use:   "up [COUNT]",
	Args:  cobra.RangeArgs(0, 1),
	Short: "Apply pending migrations",
	Run:   upCb,
}

func upCb(_ *cobra.Command, args []string) {
	var count uint64 = 0
	if len(args) > 0 {
		count = parseCount(args[0])
	}

	migs, err := migrEngine.Up(count)
	if err != nil {
		log.Error().Err(err).Msg("error while applting migrations")
	}

	for _, m := range migs {
		log.Info().Str("version", m.Version).Msg("successfully applied a migration")
	}
}
