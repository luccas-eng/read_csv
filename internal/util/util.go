package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

// NullBool type represent sql.NullBool type from database/sql package
type NullBool sql.NullBool

// NullFloat64 type represent sql.NullFloat64 type from database/sql package
type NullFloat64 sql.NullFloat64

// NullInt64 type represent sql.NullInt64 type from database/sql package
type NullInt64 sql.NullInt64

// NullString type represent sql.NullString type from database/sql package
type NullString sql.NullString

//NullTime type represent sql.NullTime type from database/sql package
type NullTime mysql.NullTime

// MarshalJSON implements the MarshalJSON interface for NullBool
func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return json.Marshal(nil)
}

// Scan implements the Scanner interface for NullBool
func (n *NullBool) Scan(value interface{}) error {
	var i sql.NullBool
	if err := i.Scan(value); err != nil {
		return fmt.Errorf("util.(n *NullBool)Scan()=>i.Scan(): %w", err)
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*n = NullBool{i.Bool, false}
	} else {
		*n = NullBool{i.Bool, true}
	}
	return nil
}

// MarshalJSON implements the MarshalJSON interface for NullFloat64
func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Float64)
	}
	return json.Marshal(nil)
}

func (n NullFloat64) ToString(prec int) string {
	if n.Valid {
		return strconv.FormatFloat(n.Float64, 'f', prec, 64)
	}
	return "0.00"
}

// Scan implements the Scanner interface for NullFloat64
func (n *NullFloat64) Scan(value interface{}) error {
	var i sql.NullFloat64
	if err := i.Scan(value); err != nil {
		return fmt.Errorf("util.(n *NullFloat64)Scan()=>i.Scan(): %w", err)
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*n = NullFloat64{i.Float64, false}
	} else {
		*n = NullFloat64{i.Float64, true}
	}
	return nil
}

// MarshalJSON implements the MarshalJSON interface for NullInt64
func (n NullInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return json.Marshal(nil)
}

// Scan implements the Scanner interface for NullInt64
func (n *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return fmt.Errorf("util.(n *NullInt64)Scan()=>i.Scan(): %w", err)
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*n = NullInt64{i.Int64, false}
	} else {
		*n = NullInt64{i.Int64, true}
	}
	return nil
}

func (n NullInt64) ToString() string {
	if n.Valid {
		return strconv.FormatInt(n.Int64, 10)
	}
	return "0"
}

// MarshalJSON implements the MarshalJSON interface for NullString
func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

// Scan implements the Scanner interface for NullString
func (n *NullString) Scan(value interface{}) error {
	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return fmt.Errorf("util.(n *NullString)Scan()=>i.Scan(): %w", err)
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*n = NullString{i.String, false}
	} else {
		*n = NullString{i.String, true}
	}
	return nil
}

// Scan implements the Scanner interface for NullTime
func (n *NullTime) Scan(value interface{}) error {
	var i mysql.NullTime
	if err := i.Scan(value); err != nil {
		return fmt.Errorf("util.(n *NullTime)Scan()=>i.Scan(): %w", err)
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*n = NullTime{i.Time, false}
	} else {
		*n = NullTime{i.Time, true}
	}
	return nil
}

// MarshalJSON implements the MarshalJSON interface for NullTime
func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return json.Marshal(nil)
}

func (n *NullString) NewString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func (n *NullBool) NewBool(b bool) sql.NullBool {
	if b == false {
		return sql.NullBool{}
	}
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}

func (n *NullInt64) NewInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func (n *NullFloat64) NewFloat64(f float64) sql.NullFloat64 {
	if f == 0.0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}
}

func (n *NullTime) NewTime(t time.Time) mysql.NullTime {
	if t.IsZero() {
		return mysql.NullTime{}
	} else {
		return mysql.NullTime{
			Time:  t,
			Valid: true,
		}
	}
}
