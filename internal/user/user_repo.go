package user

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"github.com/jackc/pgx"
)

type Repository interface {
	InsertInto(tx *pgx.Tx, user *models.User) error
	GetByNickname(tx *pgx.Tx, user *models.User) error
	GetByNicknameOrEmail(tx *pgx.Tx, user *models.User) ([]models.User, error)
	Update(tx *pgx.Tx, user *models.User) error
	DeleteAll() error
	GetStatus(tx *pgx.Tx, status *models.Status) error
	CreateTx() (*pgx.Tx, error)
	CommitTx(tx *pgx.Tx) error
}
