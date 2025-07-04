package mock

import (
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"io"
	"reflect"
	"sync"
)

var (
	rowStates = sync.Map{}
	stateID   int64
	stateMu   sync.Mutex
)

type Rows struct {
	Data       [][]interface{}
	currentRow int32
	id         int64
}

func NewRows(data [][]interface{}) *Rows {
	stateMu.Lock()
	stateID++
	id := stateID
	stateMu.Unlock()

	rows := &Rows{
		Data: data,
		id:   id,
	}
	rowStates.Store(id, int32(0))
	return rows
}

func (r Rows) Next() bool {
	val, _ := rowStates.Load(r.id)
	current := val.(int32)
	if current >= int32(len(r.Data)) {
		return false
	}
	rowStates.Store(r.id, current+1)
	return true
}

func (r Rows) Scan(dest ...any) error {
	val, _ := rowStates.Load(r.id)
	current := val.(int32) - 1 // subtract 1 because Next() already incremented
	if current < 0 || current >= int32(len(r.Data)) {
		return io.EOF
	}

	row := r.Data[current]
	for i, val := range row {
		if i >= len(dest) {
			break
		}
		reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(val))
	}
	return nil
}

func (r Rows) ScanStruct(dest any) error {
	if r.currentRow >= int32(len(r.Data)) {
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
	rowStates.Delete(r.id)
	return nil
}

func (r Rows) Err() error {
	return nil
}
