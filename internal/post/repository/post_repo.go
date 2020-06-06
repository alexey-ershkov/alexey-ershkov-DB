package repository

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/post"
	"database/sql"
	"github.com/jackc/pgx"
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

func (rep *PostRepository) InsertInto(tx *pgx.Tx, p []*models.Post) error {
	created := sql.NullTime{}
	for _, val := range p {
		var err error
		if val.Parent != 0 {
			err = tx.QueryRow(
				"INSERT INTO post (usr, message,  parent, thread, forum, path, created) "+
					"VALUES ($1, $2, $3, $4, $5, $6::BIGINT[], current_timestamp) "+
					"RETURNING id, created",
				val.Author,
				val.Message,
				val.Parent,
				val.Thread,
				val.Forum,
				val.Path,
			).Scan(&val.Id, &created)
		} else {
			err = tx.QueryRow(
				"INSERT INTO post (usr, message,  parent, thread, forum, created) "+
					"VALUES ($1, $2, $3, $4, $5, current_timestamp) "+
					"RETURNING id, created",
				val.Author,
				val.Message,
				val.Parent,
				val.Thread,
				val.Forum,
			).Scan(&val.Id, &created)
		}
		if err != nil {
			return err
		}
		if created.Valid {
			val.Created = created.Time.Format(time.RFC3339Nano)
		}
		_, err = tx.Exec(
			"INSERT INTO forum_users (forum, nickname) "+
				"VALUES ($1,$2) "+
				"ON CONFLICT (forum, nickname) DO NOTHING",
			val.Forum,
			val.Author,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rep *PostRepository) GetById(tx *pgx.Tx, p *models.Post) error {
	created := sql.NullTime{}
	err := tx.QueryRow(
		"SELECT p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread, p.path "+
			"FROM post p "+
			"WHERE p.id = $1",
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
	err := tx.QueryRow(
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
