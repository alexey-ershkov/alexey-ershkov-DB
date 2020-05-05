package repository

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/thread"
	"github.com/jackc/pgx"
)

type Repository struct {
	db *pgx.Conn
}

func NewThreadRepository(db *pgx.Conn) thread.Repository {
	return &Repository{
		db: db,
	}
}

func (rep *Repository) InsertInto(th *models.Thread) error {
	row := rep.db.QueryRow(
		"INSERT INTO thread (usr, created, forum, message, title, slug) VALUES ($1, $2, $3, $4, $5, $6)"+
			"RETURNING id",
		th.Author,
		th.Created,
		th.Forum,
		th.Message,
		th.Title,
		th.Slug,
	)
	var info int64
	if err := row.Scan(&info); err != nil {
		return err
	}
	return nil
}

func (rep *Repository) Get(th *models.Thread) error {
	row := rep.db.QueryRow(
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug  "+
			"FROM thread t "+
			"JOIN forum f on t.forum = f.slug "+
			"WHERE t.usr = $1 AND t.forum = $2 AND t.message = $3 AND t.title = $4",
		th.Author,
		th.Forum,
		th.Message,
		th.Title,
	)
	if err := row.Scan(
		&th.Id,
		&th.Title,
		&th.Message,
		&th.Created,
		&th.Slug,
		&th.Author,
		&th.Forum,
	); err != nil {
		return err
	}
	return nil
}
