package main

import (
	"Backend/config"
	"Backend/responses"
	"Backend/utils"
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/valkey-io/valkey-go"
	"golang.org/x/oauth2"
	"log/slog"
	"os"
	"sync"
)

const version = "1.0.0"

type application struct {
	wg     sync.WaitGroup
	logger *slog.Logger

	db    *sql.DB
	vkDB  valkey.Client
	oauth *oauth2.Config

	multiLLM *config.MultiLLM

	util      *utils.Utils
	responses *responses.ErrorResponses
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := config.NewDB()
	if err != nil {
		logger.Error(err.Error())
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			logger.Error(err.Error())
		}
	}(db)

	valkeyDB, err := config.NewValkeyDB()
	if err != nil {
		logger.Error(err.Error())
	}
	defer valkeyDB.Close()

	if err != nil {
		logger.Error(err.Error())
	}

	multiLLM, err := config.NewAI()
	if err != nil {
		logger.Error(err.Error())
	}

	util := utils.NewUtils(logger)
	app := &application{
		logger:    logger,
		db:        db,
		vkDB:      valkeyDB,
		oauth:     config.NewGoogleOAuth(),
		util:      util,
		responses: responses.NewErrorResponses(logger, util),
		multiLLM:  multiLLM,
	}

	if err := app.serve(); err != nil {
		logger.Error(err.Error())
	}
	os.Exit(1)
}
