package main

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"ecommerce-analysis/internal/config"
	"ecommerce-analysis/internal/repository"
	"ecommerce-analysis/internal/service"
	"ecommerce-analysis/internal/utils"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logger := utils.NewLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	// Connect to database
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		logger.Error("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database: %v", err)
		os.Exit(1)
	}

	// Set up repository and service
	repo := repository.NewRepository(db)
	analyzer := service.NewAnalyzer(repo, logger, cfg.Quantile)

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Shutting down...")
		db.Close()
		os.Exit(0)
	}()

	// Run analysis
	logger.Info("Starting customer revenue analysis...")
	if err := analyzer.AnalyzeCustomerRevenue(); err != nil {
		logger.Error("Analysis failed: %v", err)
		os.Exit(1)
	}

	logger.Info("Analysis completed successfully")
}
