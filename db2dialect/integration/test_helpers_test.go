package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"strings"
	"testing"

	_ "github.com/ibmdb/go_ibm_db"
	"github.com/ibmdb/go_ibm_db/db2dialect"
)

type testConfig struct {
	Host     string `json:"HOSTNAME"`
	Port     string `json:"PORT"`
	Database string `json:"DATABASE"`
	UID      string `json:"UID"`
	PWD      string `json:"PWD"`
}

func Createconnection() *sql.DB {
	connectionString := os.Getenv("DB2_CONNSTR")
	if connectionString == "" {
		config := testConfig{}
		if file, err := os.Open("../../testdata/config.json"); err == nil {
			defer file.Close()
			_ = json.NewDecoder(file).Decode(&config)
		}

		lookup := func(name, fallback string) string {
			if value, ok := os.LookupEnv(name); ok {
				return value
			}
			return fallback
		}
		connectionString = "PROTOCOL=tcpip;HOSTNAME=" + lookup("DB2_HOSTNAME", config.Host) +
			";PORT=" + lookup("DB2_PORT", config.Port) +
			";DATABASE=" + lookup("DB2_DATABASE", config.Database) +
			";UID=" + lookup("DB2_USER", config.UID) +
			";PWD=" + lookup("DB2_PASSWD", config.PWD)
	}

	db, _ := sql.Open("go_ibm_db", connectionString)
	return db
}

func GetDialect(db *sql.DB) *db2dialect.Dialect {
	var dialect *db2dialect.Dialect
	switch strings.ToUpper(strings.TrimSpace(os.Getenv("DB2_TARGET_PLATFORM"))) {
	case "ZOS", "Z/OS", "MAINFRAME":
		dialect = db2dialect.NewZOS()
	case "IBMI", "AS400", "ISERIES":
		dialect = db2dialect.NewIBMi()
	case "LUW":
		dialect = db2dialect.NewLUW()
	default:
		dialect = db2dialect.New()
	}
	dialect.Init(db)
	return dialect
}

func OpenTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db := Createconnection()
	if db == nil {
		t.Fatal("failed to create SQL connection")
	}
	if err := db.Ping(); err != nil {
		db.Close()
		t.Skipf("DB2 integration database unavailable: %v", err)
	}
	return db
}
