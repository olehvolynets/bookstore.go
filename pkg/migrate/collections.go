package migrate

import (
	"context"
	"log"
	"sort"
)

// Returns a slice of pending DB migrations, sorted from the oldest one to the most recent.
func (eng *Engine) PendingMigrations() (Migrations, error) {
	currentVer, err := eng.CurrentVersion()
	if err != nil {
		return nil, err
	}

	localMigrations, err := eng.ScanDir()
	if err != nil {
		return nil, err
	}

	pendingMigrations := Migrations{}
	for _, mig := range localMigrations {
		if mig.VersionUint() > currentVer.VersionUint() {
			pendingMigrations = append(pendingMigrations, mig)
		}
	}

	sort.Sort(pendingMigrations)

	return pendingMigrations, nil
}

func (eng *Engine) AppliedMigrations() (Migrations, error) {
	remoteMigrations, err := eng.remoteVersions()
	if err != nil {
		return nil, err
	}

	localMigrations, err := eng.ScanDir()
	if err != nil {
		return nil, err
	}

	remoteMigHash := make(map[string]*Migration, len(remoteMigrations))
	for _, m := range remoteMigrations {
		mp := new(Migration)
		*mp = m

		remoteMigHash[m.Version] = mp
	}

	for _, m := range localMigrations {
		remoteMig, ok := remoteMigHash[m.Version]
		if !ok {
			continue
		}

		remoteMig.Name = m.Name
		remoteMig.UpFile = m.UpFile
		remoteMig.DownFile = m.DownFile
	}

	appliedMigrations := Migrations{}
	for _, m := range remoteMigHash {
		appliedMigrations = append(appliedMigrations, *m)
	}

	sort.Sort(appliedMigrations)

	return appliedMigrations, nil
}

// List of migration versions stored in the database.
func (eng *Engine) remoteVersions() (Migrations, error) {
	ctx := context.TODO()
	rows, err := eng.db.Query(ctx, "SELECT * FROM schema_migrations;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appliedMigrations := Migrations{}
	for rows.Next() {
		m := Migration{}

		if err := rows.Scan(&m.Version, &m.AppliedAt); err != nil {
			log.Println("failed to scan a row")
		}

		appliedMigrations = append(appliedMigrations, m)
	}

	return appliedMigrations, nil
}
