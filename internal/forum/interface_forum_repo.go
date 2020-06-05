package forum

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"github.com/jackc/pgx"
)

type Repository interface {
	InsertInto(tx *pgx.Tx, forum *models.Forum) error
	GetBySlug(tx *pgx.Tx, forum *models.Forum) error
	GetThreads(tx *pgx.Tx, forum *models.Forum, desc, limit, since string) ([]models.Thread, error)
	GetUsers(tx *pgx.Tx, forum *models.Forum, desc, limit, since string) ([]models.User, error)
	CreateTx() (*pgx.Tx, error)
	CommitTx(tx *pgx.Tx) error
}
