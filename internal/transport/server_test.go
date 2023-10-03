package transport

import (
	"l0/internal/config"
	"l0/internal/service/cache"
	"l0/internal/service/order"
	"l0/internal/service/stan"
	"l0/internal/service/storage"
	"os"
	"testing"
)

var app *App

func TestMain(m *testing.M) {

	cfg := config.New()

	fakeStorage := storage.NewFake()

	orderSvc := order.New(
		fakeStorage,
		cache.New(),
	)

	stanSvc := stan.New(cfg.Stan, orderSvc)

	app = New(cfg, orderSvc, stanSvc)

	fakeStorage.Create("test-key", "test-val")
	app.orderSvc.Create("test-key-new", "test-val-new")

	code := m.Run()
	os.Exit(code)
}
