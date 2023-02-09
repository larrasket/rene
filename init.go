package main

import (
	"database/sql"
	_ "embed"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type readFunc func(errMesg string) string

type properties struct {
	message  string
	required bool
	readFunc
}

type Field struct {
	key, value string
	properties
}

var modDes = `
Moderator can be used to add entries to the publishing database in production
(i.e. without editing) via sending a direct massage for the accounts. You can
leave this empty empty.`

var actDes = `
Do you want to set active hours? If the answery is n, René will only work in the
timeframe that you will set. (You can customize this per account later or change
it entirely). Please answer either with n or the starting time frame, you will
be promoted to enter the end time after this. Please use the 24h fromat like
this: hh:mm`

var MainConfig = []Field{{"moderator_account", "", properties{modDes,
	false, nil}}, {"active_hours", "00:00:00:00", properties{actDes, false,
	timeFrameReader}}}

func init() {
	SetLogger()
	err := SetDb()
	if err != nil {
		os.Exit(1)
	}

	confRead, err := Db.Prepare(`SELECT * FROM metadata WHERE key = ?`)
	if err != nil {
		log.Fatal(err)
	}

	confWrite, err :=
		Db.Prepare(`INSERT INTO metadata (key, value) VALUES (?, ?)`)

	if err != nil {
		log.Fatal(err)
	}

	defer confRead.Close()
	for _, f := range MainConfig {
		err = confRead.QueryRow(f.key).Scan(f.value)
		if err == sql.ErrNoRows {
			err = f.FieldReader()
			if err != nil {
				logger.Info(`
Couldn't read the field %s due to the following error: %s`, f.key, err)
				if f.required {
					log.Fatal("aborting..")
				}
			}
			_, err = confWrite.Exec(f.key, f.value)
			if err != nil {
				log.Fatal(`Couldn't write to the database`, err)
			}
		}
	}
}
