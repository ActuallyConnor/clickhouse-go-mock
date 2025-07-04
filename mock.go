package mock

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ClickHouse/clickhouse-go/v2/lib/proto"
	"github.com/pkg/errors"
	"time"
)

// ConnMock is a mock implementation of clickhouse driver.Conn interface
type ConnMock struct {
	Rows *driver.Rows
	Row  *driver.Row
}

func (c ConnMock) ServerVersion() (*driver.ServerVersion, error) {
	utc, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println("err: ", err.Error())
	}
	return &driver.ServerVersion{
		Name:        "name",
		DisplayName: "display-name",
		Revision:    0,
		Version: proto.Version{
			Major: 0,
			Minor: 0,
			Patch: 0,
		},
		Timezone: utc,
	}, nil
}

func (c ConnMock) Select(ctx context.Context, dest any, query string, args ...any) error {
	return nil
}

func (c ConnMock) Query(ctx context.Context, query string, args ...any) (driver.Rows, error) {
	if c.Rows != nil {
		return *c.Rows, nil
	}

	return nil, errors.New("no rows set")
}

func (c ConnMock) QueryRow(ctx context.Context, query string, args ...any) driver.Row {
	if c.Row != nil {
		return *c.Row
	}

	panic("no row set")
}

func (c ConnMock) PrepareBatch(ctx context.Context, query string, opts ...driver.PrepareBatchOption) (driver.Batch, error) {
	//TODO implement me
	panic("implement me")
}

func (c ConnMock) Exec(ctx context.Context, query string, args ...any) error {
	return nil
}

func (c ConnMock) AsyncInsert(ctx context.Context, query string, wait bool, args ...any) error {
	return nil
}

func (c ConnMock) Ping(ctx context.Context) error {
	return nil
}

func (c ConnMock) Stats() driver.Stats {
	//TODO implement me
	panic("implement me")
}

func (c ConnMock) Close() error {
	return nil
}

func (c ConnMock) Contributors() []string {
	return make([]string, 0)
}
