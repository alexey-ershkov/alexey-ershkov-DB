package repository

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/user"
	"github.com/jackc/pgx"
)

type Repository struct {
	db *pgx.Conn
}

func (rep *Repository) InsertInto(user *models.User) error {
	var info string
	err := rep.db.QueryRow(
		"INSERT INTO usr (email, fullname, nickname, about) VALUES ($1, $2, $3, $4) RETURNING email",
		user.Email,
		user.Fullname,
		user.Nickname,
		user.About,
	).Scan(&info)
	if err != nil {
		return err
	}
	return nil
}

func (rep *Repository) GetUserByNickname(user *models.User) error {
	row := rep.db.QueryRow(
		`SELECT u.email, u.fullname, u.nickname, u.about FROM usr u WHERE nickname = $1`,
		user.Nickname,
	)
	if err := row.Scan(&user.Email, &user.Fullname, &user.Nickname, &user.About); err != nil {
		return err
	}
	return nil
}

func (rep *Repository) GetUsersByNicknameOrEmail(user *models.User) ([]models.User, error) {
	rows, err := rep.db.Query(
		"SELECT u.email, u.fullname, u.nickname, u.about FROM usr u WHERE nickname = $1 OR email = $2",
		user.Nickname,
		user.Email,
	)
	if err != nil {
		return nil, err
	}
	users := make([]models.User, 0)
	for rows.Next() {
		if err := rows.Scan(&user.Email, &user.Fullname, &user.Nickname, &user.About); err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	rows.Close()
	return users, nil
}

func NewUserRepo(db *pgx.Conn) user.Repository {
	return &Repository{
		db: db,
	}
}
