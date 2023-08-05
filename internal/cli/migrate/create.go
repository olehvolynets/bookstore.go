package migrate

import (
	"github.com/spf13/cobra"

	"bookstore/internal/logging"
)

var createFailureMsg = "nutrition: failed to create a DB migration"

var createMigrationCommand = &cobra.Command{
	Use:   "create NAME",
	Args:  cobra.ExactArgs(1),
	Short: "Creates new migration file with the current timestamp",
	Run:   createMigrationCb,
}

func createMigrationCb(_ *cobra.Command, args []string) {
	log := logging.NewLoggerFromEnv()

	fileNames, err := migrEngine.CreateMigration(args[0])
	if err != nil {
		log.Fatal().Err(err).Msg(createFailureMsg)
	}

	log.Info().Str("filename", fileNames[0]).Msg("created")
	log.Info().Str("filename", fileNames[1]).Msg("created")
}
