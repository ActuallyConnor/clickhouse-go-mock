package mock

import (
	"context"
	"testing"
)

func TestConnMock_ServerVersion(t *testing.T) {
	conn := ConnMock{}
	version, err := conn.ServerVersion()
	if err != nil {
		t.Errorf("ServerVersion() error = %v, want nil", err)
	}
	if version.Name != "name" {
		t.Errorf("ServerVersion().Name = %v, want %v", version.Name, "name")
	}
	if version.DisplayName != "display-name" {
		t.Errorf("ServerVersion().DisplayName = %v, want %v", version.DisplayName, "display-name")
	}
	if version.Timezone.String() != "America/New_York" {
		t.Errorf("ServerVersion().Timezone = %v, want America/New_York", version.Timezone)
	}
}

func TestConnMock_Query(t *testing.T) {
	tests := []struct {
		name    string
		rows    *Rows
		wantErr bool
	}{
		{
			name:    "with nil rows",
			rows:    nil,
			wantErr: true,
		},
		{
			name: "with rows",
			rows: &Rows{
				Data: [][]interface{}{
					{"test"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := ConnMock{
				Rows: tt.rows,
			}
			rows, err := conn.Query(context.Background(), "SELECT 1")
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return // Skip Close() if we got an error
			}
			defer rows.Close()
			if !tt.wantErr && rows == nil {
				t.Error("Query() returned nil rows when error wasn't expected")
			}
		})
	}
}

func TestConnMock_QueryRow(t *testing.T) {
	t.Run("with row", func(t *testing.T) {
		row := &Row{
			Data: []interface{}{"test"},
		}
		conn := ConnMock{
			Row: row,
		}
		result := conn.QueryRow(context.Background(), "SELECT 1")
		if result == nil {
			t.Error("QueryRow() returned nil")
		}
	})

	t.Run("with nil row", func(t *testing.T) {
		conn := ConnMock{}
		defer func() {
			if r := recover(); r == nil {
				t.Error("QueryRow() with nil row should panic")
			}
		}()
		conn.QueryRow(context.Background(), "SELECT 1")
	})
}

func TestRow_Scan(t *testing.T) {
	tests := []struct {
		name    string
		row     Row
		dest    []interface{}
		wantErr bool
	}{
		{
			name: "successful scan",
			row: Row{
				Data: []interface{}{"test", 42},
			},
			dest:    []interface{}{new(string), new(int)},
			wantErr: false,
		},
		{
			name: "mismatched types",
			row: Row{
				Data: []interface{}{"test"},
			},
			dest:    []interface{}{new(int)},
			wantErr: true,
		},
		{
			name: "wrong number of destinations",
			row: Row{
				Data: []interface{}{"test"},
			},
			dest:    []interface{}{new(string), new(string)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.row.Scan(tt.dest...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type testStruct struct {
	Field1 string
	Field2 int
}

func TestRow_ScanStruct(t *testing.T) {
	tests := []struct {
		name    string
		row     Row
		dest    interface{}
		wantErr bool
	}{
		{
			name: "successful struct scan",
			row: Row{
				Data: []interface{}{"test", 42},
			},
			dest:    &testStruct{},
			wantErr: false,
		},
		{
			name: "nil destination",
			row: Row{
				Data: []interface{}{"test", 42},
			},
			dest:    nil,
			wantErr: true,
		},
		{
			name: "non-pointer destination",
			row: Row{
				Data: []interface{}{"test", 42},
			},
			dest:    testStruct{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.row.ScanStruct(tt.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnMock_Contributors(t *testing.T) {
	conn := ConnMock{}
	contributors := conn.Contributors()
	if len(contributors) != 0 {
		t.Errorf("Contributors() returned %v elements, want 0", len(contributors))
	}
}

func TestConnMock_Close(t *testing.T) {
	conn := ConnMock{}
	if err := conn.Close(); err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}

func TestConnMock_Ping(t *testing.T) {
	conn := ConnMock{}
	if err := conn.Ping(context.Background()); err != nil {
		t.Errorf("Ping() error = %v, want nil", err)
	}
}

func TestConnMock_Exec(t *testing.T) {
	conn := ConnMock{}
	if err := conn.Exec(context.Background(), "SELECT 1"); err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
}

func TestConnMock_AsyncInsert(t *testing.T) {
	conn := ConnMock{}
	if err := conn.AsyncInsert(context.Background(), "INSERT INTO table VALUES (1)", true); err != nil {
		t.Errorf("AsyncInsert() error = %v, want nil", err)
	}
}
