package repos

import "cvbuilder/db"

type Repos struct {
	db *db.DB

	CV        *CV
	User      *User
	Job       *Job
	CVVariant *CVVariant
}

func Init(db *db.DB) (*Repos, error) {
	return &Repos{
		db:        db,
		CV:        InitCVRepo(db),
		User:      InitUserRepo(db),
		Job:       InitJobRepo(db),
		CVVariant: InitCVVariantRepo(db),
	}, nil
}
