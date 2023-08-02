package migrate

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"

	"bookstore/internal/db"
	"bookstore/pkg/migrate"
)

var Command = &cobra.Command{
	Use:   "migrate COMMAND",
	Short: "DB migrations management",
}

var migrEngine *migrate.Engine

func init() {
	ctx := context.Background()

	config := db.NewConfig(ctx, "postgres")
	config.Attrs.Add("sslmode", fetchEnv("DB_SSLMODE", "disable"))

	conn, err := pgx.Connect(ctx, config.String())
	if err != nil {
		panic(err)
	}
	migrEngine = migrate.New(conn)

	Command.AddCommand(upCommand)
	Command.AddCommand(downCommand)
	Command.AddCommand(redoCommand)
	Command.AddCommand(createMigrationCommand)
	Command.AddCommand(statusCommand)
	Command.AddCommand(versionCommand)
}
