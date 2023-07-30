package migrate

import "context"

type applyMigrationCb func(version string) error

func (eng *Engine) applyMigration(version, query string, cb applyMigrationCb) error {
	ctx := context.Background()
	tx, err := eng.db.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = cb(version)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)

	return nil
}

func (eng *Engine) onApplyMigration(version string) error {
	ctx := context.Background()
	_, err := eng.db.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version)

	return err
}

func (eng *Engine) onRollBackMigration(version string) error {
	ctx := context.Background()
	_, err := eng.db.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", version)

	return err
}
