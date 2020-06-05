package post

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"github.com/jackc/pgx"
)

type Repository interface {
	InsertInto(post []*models.Post) error
	GetById(post *models.Post) error
	Update(post *models.Post) error
	CreateTx() (*pgx.Tx, error)
	CommitTx(tx *pgx.Tx) error
}
