package repository

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/user"
	"database/sql"
	"fmt"
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

func (rep *Repository) InsertInto(user *models.User) error {
	var info string
	about := &sql.NullString{}
	if user.About != "" {
		about.String = user.About
		about.Valid = true
	}
	err := rep.db.QueryRow(
		"INSERT INTO usr (email, fullname, nickname, about) VALUES ($1, $2, $3, $4) RETURNING email",
		user.Email,
		user.Fullname,
		user.Nickname,
		about,
	).Scan(&info)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (rep *Repository) GetByNickname(user *models.User) error {
	row := rep.db.QueryRow(
		`SELECT u.email, u.fullname, u.nickname, u.about FROM usr u WHERE nickname = $1`,
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

func (rep *Repository) GetByNicknameOrEmail(user *models.User) ([]models.User, error) {
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

func (rep *Repository) Update(user *models.User) error {
	var info string
	sqlStr := "UPDATE usr SET " +
		"email = $1, " +
		"nickname = $2"
	args := make([]string, 0)
	args = append(args, user.Email, user.Nickname)
	if user.Fullname != "" {
		sqlStr += ", fullname = $%d "
		args = append(args, user.Fullname)
		sqlStr = fmt.Sprintf(sqlStr, len(args))
	}
	if user.About != "" {
		sqlStr += ", about = $%d "
		args = append(args, user.About)
		sqlStr = fmt.Sprintf(sqlStr, len(args))
	}
	sqlStr += "WHERE nickname = $2 RETURNING email"
	var err error
	switch len(args) {
	case 2:
		err = rep.db.QueryRow(sqlStr, args[0], args[1]).Scan(&info)
	case 3:
		err = rep.db.QueryRow(sqlStr, args[0], args[1], args[2]).Scan(&info)
	case 4:
		err = rep.db.QueryRow(sqlStr, args[0], args[1], args[2], args[3]).Scan(&info)
	}
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

func (rep *Repository) GetStatus(s *models.Status) error {
	rows, err := rep.db.Query(
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
