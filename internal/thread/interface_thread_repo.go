package thread

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"github.com/jackc/pgx"
)

type Repository interface {
	InsertInto(tx *pgx.Tx, thread *models.Thread) error
	GetCreated(tx *pgx.Tx, thread *models.Thread) error
	GetBySlugOrId(tx *pgx.Tx, thread *models.Thread) error
	InsertIntoVotes(tx *pgx.Tx, vote *models.Vote) error
	GetVotes(tx *pgx.Tx, thread *models.Thread, vote *models.Vote) error
	Update(tx *pgx.Tx, thread *models.Thread) error
	GetPosts(tx *pgx.Tx, thread *models.Thread, desc, sort, limit, since string) ([]models.Post, error)
	CreateTx() (*pgx.Tx, error)
	CommitTx(tx *pgx.Tx) error
	Prepare() error
}
