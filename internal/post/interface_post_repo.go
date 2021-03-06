package post

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"github.com/jackc/pgx"
)

type Repository interface {
	InsertInto(tx *pgx.Tx, post []*models.Post, thread *models.Thread) error
	GetById(tx *pgx.Tx, post *models.Post) error
	Update(tx *pgx.Tx, post *models.Post) error
	CreateTx() (*pgx.Tx, error)
	Prepare() error
}
