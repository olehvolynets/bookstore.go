package migrate

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/jackc/pgx/v5"
)

//go:embed init_template.sql
var initTemplate []byte

type engineBuilder struct {
	engine *Engine
}

func NewBuilder() *engineBuilder {
	return &engineBuilder{engine: new(Engine)}
}

func (eb *engineBuilder) SetDb(db *pgx.Conn) *engineBuilder {
	eb.engine.db = db

	return eb
}

func (eb *engineBuilder) SetDir(dir string) *engineBuilder {
	path, err := filepath.Abs(dir)
	if err != nil {
		fmt.Println(err)

		panic(fmt.Sprintf("DB migration director is invalid (%s)", path))
	}

	eb.engine.dir = path

	return eb
}

func (eb *engineBuilder) Done() *Engine {
	err := eb.validate()
	if err != nil {
		panic(err)
	}
	err = eb.initEngine()
	if err != nil {
		panic(err)
	}

	return eb.engine
}

func (eb *engineBuilder) validate() error {
	if eb.engine.db == nil {
		return errors.New("DB unset for the migration engine")
	}

	if eb.engine.dir == "" {
		eb.engine.dir = "database/migrations"
	}

	return nil
}

func (eb *engineBuilder) initEngine() error {
	ctx := context.TODO()
	err := eb.engine.db.Ping(ctx)
	if err != nil {
		return errors.Join(errors.New("DB ping failed"), err)
	}

	_, err = eb.engine.db.Exec(ctx, string(initTemplate))
	if err != nil {
		return errors.Join(errors.New("failed to create control table (%v)"), err)
	}

	return nil
}
