package migrate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

var (
	directions  = [2]string{"up", "down"}
	formatRegex = regexp.MustCompile(`(\d{16})_(.+)\.(\w+)\.sql`)
	errCreation = errors.New("failed to create a DB migration")
)

type ErrPartialSuccess struct {
	Message   string
	Remaining Migrations
}

func (e ErrPartialSuccess) Error() string {
	return e.Message
}

type Engine struct {
	db  *pgx.Conn
	dir string
}

// Creates a migration file.
// Returns a list of created file names
func (eng *Engine) CreateMigration(name string) ([]string, error) {
	version := strconv.FormatInt(time.Now().UnixMicro(), 10)
	versionGlob := filepath.Join(eng.dir, version+"_*.sql")
	matches, err := filepath.Glob(versionGlob)
	if err != nil {
		return nil, errors.Join(errCreation, err)
	}

	if len(matches) > 0 {
		return nil, errors.Join(
			errCreation,
			fmt.Errorf("duplicate migration version: %s", version),
			err,
		)
	}

	if err = os.MkdirAll(eng.dir, os.ModePerm); err != nil {
		return nil, errors.Join(errCreation, err)
	}

	createdFiles := []string{}
	for _, direction := range directions {
		baseName := fmt.Sprintf("%s_%s.%s.sql", version, name, direction)
		fileName := filepath.Join(eng.dir, baseName)

		f, err := os.Create(fileName)
		if err != nil {
			if len(createdFiles) != 0 {
				for _, f := range createdFiles {
					os.Remove(f)
				}
			}

			return nil, errors.Join(errCreation, err)
		}
		f.Close()

		createdFiles = append(createdFiles, fileName)
	}

	return createdFiles, nil
}

func (eng *Engine) CurrentVersion() (Migration, error) {
	ctx := context.Background()
	row := eng.db.QueryRow(ctx, "SELECT * FROM schema_migrations ORDER BY version DESC LIMIT 1;")

	m := Migration{}

	err := row.Scan(&m.Version, &m.AppliedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Migration{Version: "0", AppliedAt: time.Time{}}, nil
		}

		return m, err
	}

	return m, nil
}

// Returns slice of migrations combining all local migrations as well as the migrations, that have been
// previously applied, but which files are missing.
func (eng *Engine) Status() (Migrations, error) {
	localMigrations, err := eng.ScanDir()
	if err != nil {
		return nil, err
	}

	migrationHash := make(map[string]*Migration, len(localMigrations))
	for _, mig := range localMigrations {
		m := new(Migration)
		*m = mig
		migrationHash[mig.Version] = m
	}

	remoteMigrations, err := eng.remoteVersions()
	if err != nil {
		return nil, err
	}

	for _, remoteMig := range remoteMigrations {
		pMig, ok := migrationHash[remoteMig.Version]

		if ok {
			pMig.AppliedAt = remoteMig.AppliedAt
		} else {
			m := new(Migration)
			*m = remoteMig

			migrationHash[remoteMig.Version] = m
		}
	}

	migrations := Migrations{}
	for _, mig := range migrationHash {
		migrations = append(migrations, *mig)
	}

	sort.Sort(migrations)

	return migrations, nil
}

func (eng *Engine) Up(count uint64) (appliedMigrations Migrations, outErr error) {
	pendingMigrations, err := eng.PendingMigrations()
	if err != nil {
		return nil, err
	}

	if count > 0 {
		pendingMigrations = pendingMigrations[:count]
	}

	appliedMigrations = Migrations{}
	for _, mig := range pendingMigrations {
		query, err := eng.readMigrationFile(mig.UpFile)
		if err != nil {
			outErr = err
			break
		}

		err = eng.applyMigration(mig.Version, string(query), eng.onApplyMigration)
		if err != nil {
			outErr = err
			break
		}

		appliedMigrations = append(appliedMigrations, mig)
	}

	if len(appliedMigrations) != 0 && len(appliedMigrations) != len(pendingMigrations) {
		partialSuccessError := ErrPartialSuccess{Message: "some DB migrations failed to apply"}
		partialSuccessError.Remaining = pendingMigrations[len(appliedMigrations):]
		outErr = errors.Join(partialSuccessError, outErr)
	}

	return appliedMigrations, outErr
}

func (eng *Engine) Down(count uint64) (rolledBack Migrations, outErr error) {
	appliedVersions, err := eng.AppliedMigrations()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		count = uint64(len(appliedVersions))
	}

	sort.Sort(sort.Reverse(appliedVersions))

	rolledBack = Migrations{}
	for _, mig := range appliedVersions {
		if count == 0 {
			break
		}

		query, err := eng.readMigrationFile(mig.DownFile)
		if err != nil {
			outErr = err
			break
		}

		err = eng.applyMigration(mig.Version, string(query), eng.onRollBackMigration)
		if err != nil {
			outErr = err
			break
		}

		count--
		rolledBack = append(rolledBack, mig)
	}

	if count != 0 {
		partialSuccessError := ErrPartialSuccess{Message: "some DB migrations failed to roll back"}
		partialSuccessError.Remaining = appliedVersions[len(rolledBack):]
		outErr = errors.Join(partialSuccessError, outErr)
	}

	return rolledBack, outErr
}

// Intended for internal usage, but who knows...
// Returns a slice of Migration's for each migration version present in the corresponding directory.
func (eng *Engine) ScanDir() (Migrations, error) {
	files, err := os.ReadDir(eng.dir)
	if err != nil {
		return nil, err
	}

	migrationHash := make(map[string]*Migration)
	for _, file := range files {
		scan := formatRegex.FindStringSubmatch(file.Name())

		version := scan[1]

		m, ok := migrationHash[version]
		if !ok {
			migrationHash[version] = &Migration{Version: version}
			m = migrationHash[version]
		}

		m.Name = scan[2]

		switch scan[3] {
		case upDirection:
			m.UpFile = file.Name()
		case downDirection:
			m.DownFile = file.Name()
		default:
			return nil, fmt.Errorf(
				"invalid migration file format - direction is invalid, %s",
				file.Name(),
			)
		}
	}

	migrations := Migrations{}
	for _, m := range migrationHash {
		migrations = append(migrations, *m)
	}

	return migrations, nil
}

func (eng *Engine) readMigrationFile(filename string) ([]byte, error) {
	filepath := filepath.Join(eng.dir, filename)
	query, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return query, nil
}
