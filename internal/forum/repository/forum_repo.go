package repository

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/forum"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"github.com/jackc/pgx"
)

type Repository struct {
	db *pgx.Conn
}

func NewForumRepository(db *pgx.Conn) forum.Repository {
	return &Repository{db: db}
}

func (rep *Repository) InsertInto(f *models.Forum) error {
	row := rep.db.QueryRow(
		"INSERT INTO forum (slug, title, usr) VALUES ($1, $2, $3) RETURNING title",
		f.Slug,
		f.Title,
		f.User,
	)
	var info string
	err := row.Scan(&info)
	if err != nil {
		return err
	}
	return nil
}

func (rep *Repository) GetBySlug(f *models.Forum) error {
	row := rep.db.QueryRow(
		"SELECT count(p), f.slug, (SELECT count(*) FROM forum f2 JOIN thread t2 on f2.slug = t2.forum WHERE f2.slug = $1), f.title, u.nickname FROM forum f "+
			"LEFT JOIN thread t on f.slug = t.forum "+
			"LEFT JOIN post p on t.id = p.thread "+
			"JOIN usr u on f.usr = u.nickname "+
			"WHERE f.slug = $1 "+
			"GROUP BY f.slug, u.nickname",
		f.Slug,
	)
	if err := row.Scan(&f.Posts, &f.Slug, &f.Threads, &f.Title, &f.User); err != nil {
		return err
	}
	return nil
}
