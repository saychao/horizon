package schema

import (
	"database/sql"
	"errors"
	"github.com/saychao/horizon/db2"

	migrate "github.com/rubenv/sql-migrate"
)

//go:generate go-bindata -nometadata -ignore .+\.go$ -pkg schema -o bindata.go ./...

// Migrations represents all of the schema migration for horizon
var Migrations migrate.MigrationSource = &migrate.AssetMigrationSource{
	Asset:    Asset,
	AssetDir: AssetDir,
	Dir:      "migrations",
}

// Migrate performs schema migration.  Migrations can occur in one of three
// ways:
//
// - up: migrations are performed from the currently installed version upwards.
// If count is 0, all unapplied migrations will be run.
//
// - down: migrations are performed from the current version downard. If count
// is 0, all applied migrations will be run in a downard direction.
//
// - redo: migrations are first ran downard `count` times, and then are rand
// upward back to the current version at the start of the process. If count is
// 0, a count of 1 will be assumed.
func Migrate(db *sql.DB, dir db2.MigrateDir, count int) (int, error) {
	switch dir {
	case db2.MigrateUp:
		return migrate.ExecMax(db, "postgres", Migrations, migrate.Up, count)
	case db2.MigrateDown:
		return migrate.ExecMax(db, "postgres", Migrations, migrate.Down, count)
	case db2.MigrateRedo:

		if count == 0 {
			count = 1
		}

		down, err := migrate.ExecMax(db, "postgres", Migrations, migrate.Down, count)
		if err != nil {
			return down, err
		}

		return migrate.ExecMax(db, "postgres", Migrations, migrate.Up, down)
	default:
		return 0, errors.New("Invalid migration direction")
	}
}
