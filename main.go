package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/adinovcina/golang-setup/api/handlers"

	"github.com/adinovcina/golang-setup/tools/network/http"
	"github.com/redis/go-redis/v9"

	"github.com/adinovcina/golang-setup/config"
	"github.com/adinovcina/golang-setup/services"
	"github.com/adinovcina/golang-setup/tools/logger"
	"github.com/adinovcina/golang-setup/tools/mysql"
	r "github.com/adinovcina/golang-setup/tools/redis"
	"github.com/rs/zerolog"

	mysqlstore "github.com/adinovcina/golang-setup/store/mysql"
	redisstore "github.com/adinovcina/golang-setup/store/redis"
)

func main() {
	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() { <-c; cancel() }()

	// Instantiate a new type to represent our application.
	// This type lets us shared setup code with our end-to-end tests.
	m := NewMain()

	// Execute program.
	if err := m.Run(); err != nil {
		m.Close()
		logger.Error().Err(err).Msg("failed to run the app")
		os.Exit(1)
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Main represents the main application struct.
type Main struct {
	conf *config.Config

	HTTPServer *http.Server

	DB          *mysql.DB
	RedisClient *redis.Client
}

// NewMain creates a new instance of Main.
func NewMain() *Main {
	// Set local time to UTC
	time.Local = time.UTC

	// Load application configuration
	conf, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load config")
	}

	level, err := zerolog.ParseLevel(conf.Service.LogLevel)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse log level")
	}

	// Set global log level
	zerolog.SetGlobalLevel(level)

	// Connect to the MySQL database
	connectionString := &mysql.ConnectionData{
		Address:  conf.Database.Address,
		Name:     conf.Database.Name,
		Password: conf.Database.Password,
		Port:     conf.Database.Port,
		Username: conf.Database.Username,
	}

	db, err := mysql.Init(connectionString)
	log.Println(err)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init MySQL")
	}

	// Connect to the Redis
	redisConnectionString := r.ConnectionData{
		Address:  conf.Redis.Address,
		Database: conf.Redis.Database,
		Password: conf.Redis.Password,
	}

	redisClient, err := r.Init(redisConnectionString)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init Redis")
	}

	return &Main{
		conf:        conf,
		DB:          db,
		RedisClient: redisClient,
		HTTPServer:  http.NewServer(conf),
	}
}

// Run executes the main application logic.
func (main *Main) Run() error {
	// Initialize the database store
	mysqlStore := mysqlstore.New(main.DB.GetDB(), main.conf)

	// Initialize the in-memory store
	redisStore := redisstore.New(main.RedisClient)

	// Run database migrations if enabled
	if main.conf.Database.MigrationEnabled {
		err := main.DB.Migrate(main.conf.Database.MigrationFolder, main.conf.Service.Environment)
		if err != nil {
			return err
		}
	}

	// Initialize third-party services
	appServices := services.Init(main.conf)

	// Start the HTTP server and listen for incoming requests
	go func() {
		logger.Fatal().
			Err(handlers.
				Attach(main.HTTPServer, mysqlStore, redisStore, main.conf, appServices).
				Serve())
	}()

	return nil
}

// Close gracefully stops the program.
func (main *Main) Close() error {
	if main.HTTPServer != nil {
		if err := main.HTTPServer.Close(); err != nil {
			return err
		}
	}

	if main.DB != nil {
		if err := main.DB.Close(); err != nil {
			return err
		}
	}

	if main.RedisClient != nil {
		if err := main.RedisClient.Close(); err != nil {
			return err
		}
	}

	return nil
}
