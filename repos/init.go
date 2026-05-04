package repos

import "cvbuilder/db"

type Repos struct {
	db *db.DB
}

func Init(db *db.DB) (*Repos, error) {
	return &Repos{
		db: db,
	}, nil
}
