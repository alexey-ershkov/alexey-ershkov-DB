package repository

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/user"
	"database/sql"
	"github.com/jackc/pgx"
)

type Repository struct {
	db *pgx.ConnPool
}

func NewUserRepo(db *pgx.ConnPool) user.Repository {
	return &Repository{
		db: db,
	}
}

func (rep *Repository) InsertInto(tx *pgx.Tx, user *models.User) error {
	var info string
	about := &sql.NullString{}
	if user.About != "" {
		about.String = user.About
		about.Valid = true
	}
	sqlUserInsertInto := "INSERT INTO usr (email, fullname, nickname, about) " +
		"VALUES ($1, $2, $3, $4) " +
		"ON CONFLICT DO NOTHING " +
		"RETURNING email"
	err := tx.QueryRow(
		sqlUserInsertInto,
		user.Email,
		user.Fullname,
		user.Nickname,
		about,
	).Scan(&info)
	if err != nil {
		return err
	}
	return nil
}

func (rep *Repository) GetByNickname(tx *pgx.Tx, user *models.User) error {
	sqlUserGetByNickname := "SELECT u.email, u.fullname, u.nickname, u.about " +
		"FROM usr u " +
		"WHERE nickname = $1 "
	row := tx.QueryRow(
		sqlUserGetByNickname,
		user.Nickname,
	)
	fullname := &sql.NullString{}
	about := &sql.NullString{}
	if err := row.Scan(&user.Email, fullname, &user.Nickname, about); err != nil {
		return err
	}
	if about.Valid {
		user.About = about.String
	}
	if fullname.Valid {
		user.Fullname = fullname.String
	}
	return nil
}

func (rep *Repository) GetByNicknameOrEmail(tx *pgx.Tx, user *models.User) ([]models.User, error) {
	sqlGetUserByNicknameOrEmail := "SELECT u.email, u.fullname, u.nickname, u.about " +
		"FROM usr u " +
		"WHERE nickname = $1 OR email = $2"
	rows, err := tx.Query(
		sqlGetUserByNicknameOrEmail,
		user.Nickname,
		user.Email,
	)
	if err != nil {
		return nil, err
	}
	users := make([]models.User, 0)
	for rows.Next() {
		fullname := &sql.NullString{}
		about := &sql.NullString{}
		if err := rows.Scan(&user.Email, fullname, &user.Nickname, about); err != nil {
			return nil, err
		}
		if about.Valid {
			user.About = about.String
		}
		if fullname.Valid {
			user.Fullname = fullname.String
		}
		users = append(users, *user)
	}
	rows.Close()
	return users, nil
}

func (rep *Repository) Update(tx *pgx.Tx, user *models.User) error {
	sqlUserUpdate := "UPDATE usr SET " +
		"email = $1, " +
		"nickname = $2, " +
		"fullname = $3, " +
		"about = $4 " +
		"WHERE nickname = $2 RETURNING email"
	_, err := tx.Exec(sqlUserUpdate,
		user.Email,
		user.Nickname,
		user.Fullname,
		user.About,
	)
	if err != nil {
		return err
	}
	return nil
}

func (rep *Repository) DeleteAll() error {
	_, err := rep.db.Exec(
		"DELETE FROM usr",
	)
	if err != nil {
		return err
	}
	return nil
}

func (rep *Repository) GetStatus(tx *pgx.Tx, s *models.Status) error {
	rows, err := tx.Query(
		"SELECT count(*) FROM forum " +
			"UNION ALL " +
			"SELECT count(*) " +
			"FROM post " +
			"UNION ALL " +
			"SELECT count(*) FROM thread " +
			"UNION ALL " +
			"SELECT count(*) FROM usr",
	)
	if err != nil {
		return err
	}
	i := 0
	for rows.Next() {
		var err error
		switch i {
		case 0:
			err = rows.Scan(&s.Forum)
		case 1:
			err = rows.Scan(&s.Post)
		case 2:
			err = rows.Scan(&s.Thread)
		case 3:
			err = rows.Scan(&s.User)
		}
		if err != nil {
			rows.Close()
			return err
		}
		i++
	}
	rows.Close()
	return nil
}

func (rep *Repository) CreateTx() (*pgx.Tx, error) {
	tx, err := rep.db.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (rep *Repository) CommitTx(tx *pgx.Tx) error {
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *Repository) Prepare() error {
	return nil
}
