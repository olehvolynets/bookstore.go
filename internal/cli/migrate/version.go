package migrate

import (
	"github.com/spf13/cobra"

	"bookstore/log"
)

var versionFailureMsg = "failed to retireve the DB migration version"

var versionCommand = &cobra.Command{
	Use:   "version",
	Args:  cobra.ExactArgs(0),
	Short: "Print the current version of the database",
	Run:   versionCb,
}

func versionCb(_ *cobra.Command, _ []string) {
	mig, err := migrEngine.CurrentVersion()
	if err != nil {
		log.Fatal().Err(err).Msg(versionFailureMsg)
	}

	if mig.Version == "0" {
		log.Info().Str("version", "NONE").Msg("no DB migrations have been applied so far")
	} else {
		log.Info().Str("version", mig.Version).Msg("current DB migration version")
	}
}
