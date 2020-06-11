package repository

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/post"
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"time"
)

type PostRepository struct {
	db *pgx.ConnPool
}

func NewPostRepository(db *pgx.ConnPool) post.Repository {
	return &PostRepository{
		db: db,
	}
}

func (rep *PostRepository) InsertInto(tx *pgx.Tx, p []*models.Post, th *models.Thread) error {
	created := sql.NullTime{}
	_, err := tx.Prepare("post_insert_into",
		"INSERT INTO post (usr, message,  parent, thread, forum, created) "+
			"VALUES ($1, $2, $3, $4, $5, current_timestamp) "+
			"RETURNING id, created",
	)
	if err != nil {
		logrus.Fatal(err)
	}
	for iter := range p {
		p[iter].Forum = th.Forum
		p[iter].Thread = th.Id
		var err error
		err = tx.QueryRow(
			"post_insert_into",
			p[iter].Author,
			p[iter].Message,
			p[iter].Parent,
			p[iter].Thread,
			p[iter].Forum,
		).Scan(&p[iter].Id, &created)
		if err != nil {
			return err
		}
		if created.Valid {
			p[iter].Created = created.Time.Format(time.RFC3339Nano)
		}
	}
	sqlForumPostUpdate := "UPDATE forum  SET posts = (posts + $1) " +
		"where slug = $2"
	if len(p) > 0 {
		_, err := tx.Exec(sqlForumPostUpdate, len(p), p[0].Forum)
		if err != nil {
			logrus.Error("Error while update post count: " + err.Error())
		}
	}
	return nil
}

func (rep *PostRepository) GetById(tx *pgx.Tx, p *models.Post) error {
	created := sql.NullTime{}
	err := tx.QueryRow(
		"post_get_by_id",
		p.Id,
	).Scan(&p.Author, &created, &p.Forum, &p.IsEdited, &p.Message, &p.Parent, &p.Thread, &p.Path)
	if err != nil {
		return err
	}
	if created.Valid {
		p.Created = created.Time.Format(time.RFC3339Nano)
	}
	return nil
}
func (rep *PostRepository) Update(tx *pgx.Tx, p *models.Post) error {
	created := sql.NullTime{}
	sqlString := "UPDATE post SET message = $1, isEdited = true " +
		"WHERE id = $2 " +
		"RETURNING usr, created, forum, isEdited, message, parent, thread"
	err := tx.QueryRow(
		sqlString,
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

func (rep *PostRepository) CreateTx() (*pgx.Tx, error) {
	tx, err := rep.db.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (rep *PostRepository) CommitTx(tx *pgx.Tx) error {
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *PostRepository) Prepare() error {
	_, err := rep.db.Prepare("post_get_by_id",
		"SELECT p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread, p.path "+
			"FROM post p "+
			"WHERE p.id = $1",
	)
	if err != nil {
		return err
	}

	return nil
}
