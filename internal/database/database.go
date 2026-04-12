package database

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbname = "./app.db"

type DB struct {
	DB      *gorm.DB
	Session Session
	Ctx     context.Context
}

func StartSession(sessionType string) (*DB, error) {
	ctx := context.Background()
	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	runMigrations(db)
	// create session
	session := Session{SessionType: sessionType}
	err = gorm.G[Session](db).Create(ctx, &session)
	if err != nil {
		return nil, err
	}
	wrapper := &DB{db, session, ctx}
	jsonData, err := json.Marshal(session)
	fmt.Println(string(jsonData))
	return wrapper, nil
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&Session{},
		&Language{},
		&LanguageData{},
		&Runtime{},
		&RuntimeData{},
		&GithubStar{},
		&DockerImage{},
		&ImageSize{},
		&ContainerRun{},
	)
}

func EndSession(db *DB, sessionErr error) error {
	var code int64 = 0
	if sessionErr != nil {
		code = 1
		errMsg := sessionErr.Error()
		db.Session.Error = &errMsg
	}
	db.Session.ExitCode = &code
	tx := db.DB.Save(&db.Session)
	return tx.Error
}
