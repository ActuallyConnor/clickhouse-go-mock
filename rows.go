package mock

import (
	"fmt"
	"io"
	"reflect"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Rows struct {
	Data       [][]interface{}
	currentRow int
}

func (r Rows) Next() bool {
	return true
}

func (r Rows) Scan(dest ...any) error {
	if r.currentRow >= len(r.Data) {
		return io.EOF
	}
	row := r.Data[r.currentRow]
	for i, val := range row {
		if i >= len(dest) {
			break
		}
		reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(val))
	}
	return nil
}

func (r Rows) ScanStruct(dest any) error {
	if r.currentRow >= len(r.Data) {
		return io.EOF
	}
	row := r.Data[r.currentRow]
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("must pass a pointer, not a value, to ScanStruct destination")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("must pass a pointer to struct, not %s", v.Kind())
	}
	for i := 0; i < v.NumField() && i < len(row); i++ {
		field := v.Field(i)
		if field.CanSet() {
			field.Set(reflect.ValueOf(row[i]))
		}
	}
	r.currentRow++
	return nil
}

func (r Rows) ColumnTypes() []driver.ColumnType {
	return []driver.ColumnType{}
}

func (r Rows) Totals(dest ...any) error {
	return nil
}

func (r Rows) Columns() []string {
	return []string{}
}

func (r Rows) Close() error {
	return nil
}

func (r Rows) Err() error {
	return nil
}
