package go_ibm_db

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/ibmdb/go_ibm_db/api"
)

func TestBindValueEmptyArray(t *testing.T) {
	tests := []struct {
		name    string
		value   driver.Value
		sqlType api.SQLSMALLINT
	}{
		{"int64", []int64{}, api.SQL_BIGINT},
		{"nil int64", []int64(nil), api.SQL_BIGINT},
		{"string", []string{}, api.SQL_WVARCHAR},
		{"nil string", []string(nil), api.SQL_WVARCHAR},
		{"bool", []bool{}, api.SQL_BIGINT},
		{"nil bool", []bool(nil), api.SQL_BIGINT},
		{"float64", []float64{}, api.SQL_DOUBLE},
		{"nil float64", []float64(nil), api.SQL_DOUBLE},
		{"date", []time.Time{}, api.SQL_TYPE_DATE},
		{"nil date", []time.Time(nil), api.SQL_TYPE_DATE},
		{"time", []time.Time{}, api.SQL_TYPE_TIME},
		{"nil time", []time.Time(nil), api.SQL_TYPE_TIME},
		{"timestamp", []time.Time{}, api.SQL_TYPE_TIMESTAMP},
		{"nil timestamp", []time.Time(nil), api.SQL_TYPE_TIMESTAMP},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Errorf("BindValue panicked: %v", recovered)
				}
			}()
			parameter := Parameter{SQLType: tt.sqlType}
			err := parameter.BindValue(api.SQLHSTMT(api.SQL_NULL_HSTMT), 0, tt.value)
			if err == nil || err.Error() != "empty array parameter" {
				t.Errorf("expected empty array parameter error, got %v", err)
			}
		})
	}
}
