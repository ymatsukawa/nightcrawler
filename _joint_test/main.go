package main

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ymatsukawa/nightcrawler"
	// d "github.com/ymatsukawa/nightcrawler/detector"
)

func main() {
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	// or
	// baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})

	logger := slog.New(nightcrawler.NewSlogHandler(baseHandler, nil))
	// or
	// suppress := []int{d.SelectMany}
	// logger := slog.New(nightcrawler.NewSlogHandler(baseHandler, suppress))
	gormLogger := newSlogGormLogger(logger)

	db, err := connectDB(gormLogger)
	if err != nil {
		logger.Error("db connection failed", slog.Any("err", err))
		os.Exit(1)
	}

	if err := migrateAndSeed(db); err != nil {
		logger.Error("migrate/seed failed", slog.Any("err", err))
		os.Exit(1)
	}
	logger.Info("migrate & seed done")

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())

	e.GET("/select-many", selectManyHandler(db))
	e.GET("/index-noworth", indexNoWorthHandler(db))
	e.GET("/subquery", subqueryHandler(db))
	e.GET("/heavy-calc", heavyCalcHandler(db))
	e.GET("/using-function", usingFunctionHandler(db))
	e.GET("/nplusone", nPlusOneHandler(db))
	e.GET("/no-slow-query", noSloQueryHandler(db))
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, echo.Map{
			"endpoints": []string{
				"/select-many",
				"/index-noworth",
				"/subquery",
				"/heavy-calc",
				"/using-function",
				"/nplusone",
				"/no-slow-query",
			},
		})
	})

	addr := ":" + env("PORT", "8080")
	if err := e.Start(addr); err != nil {
		logger.Error("server stopped", slog.Any("err", err))
		os.Exit(1)
	}
}
