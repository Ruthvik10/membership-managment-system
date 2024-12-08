package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Ruthvik10/membership-managment-system/internal/db/postgres"
	"github.com/Ruthvik10/membership-managment-system/internal/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type application struct {
	store  store
	db     struct{ dbURL string }
	logger logger
	server struct {
		addr string
	}
}

func (app *application) registerRoutes() *echo.Echo {
	var e = echo.New()
	v1 := e.Group("/api/v1")
	{
		app.registerHealthCheckRoutes(v1)
		app.registerMemberRoutes(v1)
		app.registerSportRoutes(v1)
	}

	return e
}

func (app *application) runServer(e *echo.Echo) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(app.server.addr); err != nil && err != http.ErrServerClosed {
			app.logger.WriteError("shutting down the server", err, nil)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		app.logger.WriteFatal("shutting down the server", err, nil)
	}
}

func (app *application) openDB() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), app.db.dbURL)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {

	app := &application{
		logger: log.NewZLogger(os.Stdout),
	}

	cfg, err := newConfig(".")
	if err != nil {
		app.logger.WriteFatal("Error loading config", err, nil)
	}

	app.db.dbURL = cfg.DBURL
	app.server.addr = cfg.APIAddr

	conn, err := app.openDB()
	if err != nil {
		app.logger.WriteFatal("Error connecting to the database:", err, nil)
	}

	app.logger.WriteInfo("Database connected", nil)

	memberStore := postgres.NewMemberStore(conn)
	sportStore := postgres.NewSportStore(conn)

	storeRegistry := struct {
		*postgres.MemberStore
		*postgres.SportStore
	}{
		memberStore,
		sportStore,
	}
	app.store = storeRegistry

	e := app.registerRoutes()
	app.runServer(e)
}
