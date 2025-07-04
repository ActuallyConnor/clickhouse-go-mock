package mock

import (
	"fmt"
	"reflect"
)

type Row struct {
	Data []interface{}
}

func (r Row) Err() error {
	return nil
}

func (r Row) Scan(dest ...any) error {
	if len(dest) != len(r.Data) {
		return fmt.Errorf("expected %d destination arguments, got %d", len(r.Data), len(dest))
	}

	for i, d := range dest {
		v := reflect.ValueOf(d)
		if v.Kind() != reflect.Ptr {
			return fmt.Errorf("destination argument at index %d must be a pointer", i)
		}

		dv := reflect.ValueOf(r.Data[i])
		ev := v.Elem()

		if !dv.Type().AssignableTo(ev.Type()) {
			return fmt.Errorf("cannot assign value of type %v to destination of type %v at index %d", dv.Type(), ev.Type(), i)
		}

		ev.Set(dv)
	}

	return nil
}

func (r Row) ScanStruct(dest any) error {
	if dest == nil {
		return fmt.Errorf("destination cannot be nil")
	}

	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer to struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a pointer to struct")
	}

	t := v.Type()
	if t.NumField() != len(r.Data) {
		return fmt.Errorf("expected %d fields in struct, got %d", len(r.Data), t.NumField())
	}

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			return fmt.Errorf("field %s cannot be set", t.Field(i).Name)
		}

		dv := reflect.ValueOf(r.Data[i])
		if !dv.Type().AssignableTo(field.Type()) {
			return fmt.Errorf("cannot assign value of type %v to field %s of type %v", dv.Type(), t.Field(i).Name, field.Type())
		}

		field.Set(dv)
	}

	return nil
}
