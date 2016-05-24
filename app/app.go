package Dingo

import (
	"fmt"
	"os"

	"github.com/dingoblog/dingo/app/handler"
	"github.com/dingoblog/dingo/app/model"
	"github.com/dingoblog/dingo/app/modules/posts"
	"github.com/facebookgo/inject"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// TODO: remove this global var once app depencency injection structure is done
var db *gorm.DB

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func Init(dbPath, privKey, pubKey string) {
	model.InitializeKey(privKey, pubKey)
	if err := model.Initialize(dbPath, fileExists(dbPath)); err != nil {
		err = fmt.Errorf("failed to intialize db: %v", err)
		panic(err)
	}
	fmt.Printf("Database is used at %s\n", dbPath)

	// We'll keep another instance of the db connection because we need to refactor the model package before removing it
	// This should me moved somewhere inside Run instead of using global vars. TODO: Improve config system
	dbi, err := gorm.Open("sqlite3", dbPath)

	if err != nil {
		panic(err)
	}

	db = dbi
}

func Run(portNumber string) {

	var g inject.Graph
	var h handler.Module
	var posts posts.Module

	err := g.Provide(
		&inject.Object{Value: db, Complete: true},
		&inject.Object{Value: &posts},
	)

	if err != nil {
		err = fmt.Errorf("failed to provide di: %v", err)
		panic(err)
	}

	// Share graph object with handler module
	h.Populate(g)

	app := h.Initialize()
	fmt.Printf("Application Started on port %s\n", portNumber)
	app.Run(":" + portNumber)
}
