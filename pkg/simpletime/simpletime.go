package simpletime

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jackc/pgtype"
	"time"
)

type SimpleTime time.Time

func (st SimpleTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(st).Format("02.01.2006"))
	return []byte(stamp), nil
}

func (st *SimpleTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("02.01.2006", string(bytes.Trim(data, "\"")))
	if err != nil {
		return err
	}

	*st = SimpleTime(t)
	return nil
}

func (st *SimpleTime) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		return errors.New("NULL values can't be decoded. Scan into a &*SimpleTime to handle NULLs")
	}

	dec := &pgtype.Date{}
	if err := dec.DecodeBinary(ci, src); err != nil {
		return err
	}

	*st = SimpleTime(dec.Time)

	return nil
}

func (st SimpleTime) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {
	date := pgtype.Date{Time: time.Time(st), Status: pgtype.Present}
	return date.EncodeBinary(ci, buf)
}
