package migrate

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"bookstore/log"
)

var statusCommand = &cobra.Command{
	Use:   "status",
	Args:  cobra.ExactArgs(0),
	Short: "Dump the migration status for the current DB",
	Run:   statusCb,
}

var statusHeader string

const (
	separator = "-----------------+-------+--------------------------------+---------------------+--------"
	rowFormat = "%s | %5s | %30s | %s | %s\n"
)

func init() {
	statusHeader = fmt.Sprintf(
		"%16s | %s | %30s | %19s | %7s\n%s",
		"Version",
		"State",
		"Name",
		"Applied At",
		"Files",
		separator,
	)
}

func statusCb(_ *cobra.Command, _ []string) {
	migs, err := migrEngine.Status()
	if err != nil {
		log.Fatal().Err(err).Msg("error while retrieving migration status")
	}

	sort.Sort(sort.Reverse(migs))

	fmt.Println(statusHeader)

	for _, m := range migs {
		state := "up"
		if m.AppliedAt.Equal(time.Time{}) {
			state = "down"
		}

		var filesPresence string
		if m.UpFile != "" && m.DownFile != "" {
			filesPresence = "up/down"
		} else if m.UpFile != "" {
			filesPresence = "up"
		} else if m.DownFile != "" {
			filesPresence = "down"
		} else {
			filesPresence = "NONE"
		}

		fmt.Printf(
			rowFormat,
			m.Version,
			state,
			m.Name,
			m.AppliedAt.Format(time.DateTime),
			filesPresence,
		)
	}
}
