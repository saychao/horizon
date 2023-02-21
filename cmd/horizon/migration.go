package main

import (
	"database/sql"
	"github.com/SafeRE-IT/horizon/db2"
	"log"
	"os"
	"strconv"

	"github.com/SafeRE-IT/horizon/ingest2"

	"github.com/SafeRE-IT/horizon/db2/history/schema"
	"github.com/spf13/cobra"
)

type Migrator func(*sql.DB, db2.MigrateDir, int) (int, error)

func migrateDB(cmd *cobra.Command, args []string, dbConnectionURL string, migrator Migrator) {

	// Allow invokations with 1 or 2 args.  All other args counts are erroneous.
	if len(args) < 1 || len(args) > 2 {
		cmd.Usage()
		os.Exit(1)
	}

	dir := args[0]
	count := 0

	// If a second arg is present, parse it to an int and use it as the count
	// argument to the migration call.
	if len(args) == 2 {
		var err error
		count, err = strconv.Atoi(args[1])
		if err != nil {
			log.Println(err)
			cmd.Usage()
			os.Exit(1)
		}
	}
	migrate(dir, count, migrator, dbConnectionURL)
}

func migrate(direction string, count int, migrator Migrator, dbConnectionURL string) {
	dir := db2.MigrateDir(direction)

	db, err := sql.Open("postgres", dbConnectionURL)
	if err != nil {
		log.Fatal(err)
	}

	applied, err := migrator(db, dir, count)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Applied %d migration", applied)
	}
}

func tryToEmptyDB() {
	db, err := sql.Open("postgres", conf.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var dbIngestVersion = 0
	// we expect that migration up have been executed
	row := db.QueryRow("select ingest_version from support_params")

	err = row.Scan(&dbIngestVersion)
	if (err != nil) && (err != sql.ErrNoRows) {
		log.Fatal(err)
	}

	if dbIngestVersion != ingest2.CurrentIngestVersion {
		migrate("down", 0, schema.Migrate, conf.DatabaseURL)
		migrate("up", 0, schema.Migrate, conf.DatabaseURL)
		_, err = db.Exec("insert into support_params (ingest_version) values ($1)", ingest2.CurrentIngestVersion)
		if err != nil {
			log.Fatal(err)
		}
	}
}
