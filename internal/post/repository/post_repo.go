package repository

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/post"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"time"
)

type PostRepository struct {
	db *pgx.Conn
}

func NewPostRepository(db *pgx.Conn) post.Repository {
	return &PostRepository{
		db: db,
	}
}

func (rep *PostRepository) InsertInto(p []*models.Post) error {
	tx, err := rep.db.Begin()
	if err != nil {
		logrus.Error("SQL: cannot start TX")
		return tools.SqlError
	}
	created := sql.NullTime{}
	for _, val := range p {
		err := tx.QueryRow(
			"INSERT INTO post (usr, message,  parent, thread, forum, created) "+
				"VALUES ($1, $2, $3, $4, $5, current_timestamp) "+
				"RETURNING id, created",
			val.Author,
			val.Message,
			val.Parent,
			val.Thread,
			val.Forum,
		).Scan(&val.Id, &created)
		if err != nil {
			return err
		}
		if created.Valid {
			val.Created = created.Time.Format(time.RFC3339Nano)
		}
	}
	if err := tx.Commit(); err != nil {
		logrus.Error("SQL cannot commit TX")
		return tools.SqlError
	}
	return nil
}

func (rep *PostRepository) GetById(p *models.Post) error {
	created := sql.NullTime{}
	err := rep.db.QueryRow(
		"SELECT p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.id = $1",
		p.Id,
	).Scan(&p.Author, &created, &p.Forum, &p.IsEdited, &p.Message, &p.Parent, &p.Thread)
	if err != nil {
		return err
	}
	if created.Valid {
		p.Created = created.Time.Format(time.RFC3339Nano)
	}
	return nil
}
func (rep *PostRepository) Update(p *models.Post) error {
	created := sql.NullTime{}
	err := rep.db.QueryRow(
		"UPDATE post SET message = $1, isEdited = true "+
			"WHERE id = $2 "+
			"RETURNING usr, created, forum, isEdited, message, parent, thread",
		p.Message,
		p.Id,
	).Scan(&p.Author, &created, &p.Forum, &p.IsEdited, &p.Message, &p.Parent, &p.Thread)
	if err != nil {
		return err
	}
	if created.Valid {
		p.Created = created.Time.Format(time.RFC3339Nano)
	}
	return nil
}
