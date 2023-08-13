package migrate

import (
	"context"
	"errors"
)

type applyMigrationCb func(version string) error

func (eng *Engine) applyMigration(version, query string, cb applyMigrationCb) error {
	ctx := context.Background()
	tx, err := eng.db.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr == nil {
			return err
		}

		return errors.Join(err, rbErr)
	}

	err = cb(version)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr == nil {
			return err
		}

		return errors.Join(err, rbErr)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

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
