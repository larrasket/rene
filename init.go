package main

import (
	"database/sql"
	_ "embed"
	"errors"
	"net/http"
	"os"

	twauth "github.com/dghubble/oauth1/twitter"

	"github.com/dghubble/oauth1"
	_ "github.com/mattn/go-sqlite3"
)

type readFunc func(errMesg string) string

type properties struct {
	message  string
	required bool
	readFunc
}

type field struct {
	key, value string
	properties
}

type account struct {
	tok, sec, username string
	client             *http.Client
}

var Accounts []account
var modDes = `
Moderator can be used to add entries to the publishing database in production
(i.e. without editing) via sending a direct massage for the accounts. You can
leave this empty: `

var actDes = `
Do you want to set active hours? If the answery is not n, René will only work in
the timeframe that you will set. (You can customize this per account later or
change it entirely). Please answer either with n or the starting time frame, you
will be promoted to enter the end time after this.  Please use the 24h fromat
(example hh:mm): `

var consumerDes = `Enter twitter consumer key for your application: `
var consumerSecDes = `Enter twitter secret key for your application: `
var Fields = []field{
	{"moderator_account", "", properties{modDes, false, nil}},
	{"active_hours", "00:00:00:00", properties{actDes, false, timeFrameReader}},
	{"consumer_key", "", properties{consumerDes, true, nil}},
	{"consumer_secret", "", properties{consumerSecDes, true, nil}},
}

var AuthConfig *oauth1.Config

func init() {
	err := SetLogger()
	if err != nil {
		logger.Error(`
Couldn't create/open the log file, logs will not be saved due to the following
error:`, err)
	}

	err = SetDb()
	if err != nil {
		logger.Error("Couldn't initialize the database", err)
		os.Exit(1)
	}

	// read configuration
	confRead, err := Db.Prepare(`SELECT value FROM metadata WHERE key = ?`)
	if err != nil {
		logger.Error(`Couldn't connect to database`, err)
		os.Exit(1)
	}

	confWrite, err :=
		Db.Prepare(`INSERT INTO metadata (key, value) VALUES (?, ?)`)

	if err != nil {
		logger.Error(`Couldn't connect to database`, err)
		os.Exit(1)
	}

	defer confRead.Close()

	defer func(rows *sql.Stmt) {
		err := rows.Close()
		if err != nil {
			logger.Error(`Not able to close sql buffer`, err)
		}
	}(confRead)

	var notComplteInit bool
	for k := range Fields {
		err = confRead.QueryRow(Fields[k].key).Scan(&Fields[k].value)
		if errors.Is(err, sql.ErrNoRows) {
			notComplteInit = true
			err = Fields[k].FieldReader()
			if err != nil {
				logger.Info(
					`Couldn't read the field %s due to the following error: %s`,
					Fields[k].key, err)
				if Fields[k].required {
					os.Exit(1)
				}
			}
			_, err = confWrite.Exec(Fields[k].key, Fields[k].value)
			if err != nil {
				logger.Error(`Couldn't write to the database`, err)
				os.Exit(1)
			}
		} else if err != nil {
			logger.Error(`Couldn't connect to database`, err)
			os.Exit(1)
		}

	}
	rows, err := Db.Query(
		`SELECT username, account_token, account_secret FROM accounts`)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error("error while reading accounts", err)
		os.Exit(1)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error(`Not able to close sql buffer`, err)
		}
	}(rows)

	for rows.Next() {
		var acc account
		err = rows.Scan(&acc.username, &acc.tok, &acc.sec)
		if err != nil {
			logger.Error("error while reading accounts", err)
			os.Exit(1)
		}
		Accounts = append(Accounts, acc)
	}

	// setup auth config
	AuthConfig = &oauth1.Config{
		ConsumerKey:    Fields[2].value,
		ConsumerSecret: Fields[3].value,
		CallbackURL:    "oob",
		Endpoint:       twauth.AuthorizeEndpoint,
	}

	if notComplteInit {
		os.Exit(0)
	}
}
