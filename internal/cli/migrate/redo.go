package migrate

import (
	"github.com/spf13/cobra"

	"bookstore/log"
)

var redoCommand = &cobra.Command{
	Use:   "redo",
	Args:  cobra.ExactArgs(0),
	Short: "Re-run the latest migration",
	Run:   redoCb,
}

func redoCb(_ *cobra.Command, _ []string) {
	migs, err := migrEngine.Down(1)
	if err != nil {
		log.Error().Err(err).Msg("error while rolling back migrations")
	}

	for _, m := range migs {
		log.Info().Str("version", m.Version).Msg("successfully rolled back a migration")
	}

	migs, err = migrEngine.Up(1)
	if err != nil {
		log.Error().Err(err).Msg("error while applting migrations")
	}

	for _, m := range migs {
		log.Info().Str("version", m.Version).Msg("successfully applied a migration")
	}
}
