package post

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"github.com/jackc/pgx"
)

type Repository interface {
	InsertInto(tx *pgx.Tx, post []*models.Post) error
	GetById(tx *pgx.Tx, post *models.Post) error
	Update(tx *pgx.Tx, post *models.Post) error
	CreateTx() (*pgx.Tx, error)
	CommitTx(tx *pgx.Tx) error
}
