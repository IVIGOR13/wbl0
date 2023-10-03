package storage

import (
	"context"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldCreateOrder(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close(context.Background())

	storage := &Storage{
		Db: mock,
	}

	mock.ExpectExec("insert into orders").
		WithArgs("test-key", "test-val").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err = storage.Create("test-key", "test-val"); err != nil {
		t.Errorf("%s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetOrder(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close(context.Background())

	storage := &Storage{
		Db: mock,
	}

	rows := mock.NewRows([]string{"data"}).
		AddRow("test-val")

	mock.ExpectQuery("select data from orders").WithArgs("test-key").WillReturnRows(rows)

	data, err := storage.Get("test-key")
	if err != nil {
		t.Errorf("%s", err)
	}

	assert.Equal(t, data, "test-val")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetOrderList(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close(context.Background())

	storage := &Storage{
		Db: mock,
	}

	rows := mock.NewRows([]string{"order_uid", "data"}).
		AddRow("test-key", "test-val").
		AddRow("1", "2").
		AddRow("3", "4")

	mock.ExpectQuery("select order_uid, data from orders").WillReturnRows(rows)

	orderPairs := storage.GetAll()

	assert.Equal(t, orderPairs[0].OrderUID, "test-key")
	assert.Equal(t, orderPairs[0].Data, "test-val")

	assert.Equal(t, orderPairs[1].OrderUID, "1")
	assert.Equal(t, orderPairs[1].Data, "2")

	assert.Equal(t, orderPairs[2].OrderUID, "3")
	assert.Equal(t, orderPairs[2].Data, "4")

	assert.Equal(t, len(orderPairs), 3)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
