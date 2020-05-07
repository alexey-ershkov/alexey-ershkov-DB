package repository

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/thread"
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"time"
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
	slug := &sql.NullString{}
	if th.Slug != "" {
		slug.String = th.Slug
		slug.Valid = true
	}
	created := &sql.NullString{}
	if th.Created != "" {
		created.String = th.Created
		created.Valid = true
	}
	row := rep.db.QueryRow(
		"INSERT INTO thread (usr, created, forum, message, title, slug) VALUES ($1, $2, $3, $4, $5, $6)"+
			"RETURNING id",
		th.Author,
		created,
		th.Forum,
		th.Message,
		th.Title,
		slug,
	)
	var info int64
	if err := row.Scan(&info); err != nil {
		return err
	}
	return nil
}

func (rep *Repository) GetCreated(th *models.Thread) error {
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
	created := sql.NullTime{}
	slug := sql.NullString{}
	if err := row.Scan(
		&th.Id,
		&th.Title,
		&th.Message,
		&created,
		&slug,
		&th.Author,
		&th.Forum,
	); err != nil {
		return err
	}
	if created.Valid {
		th.Created = created.Time.Format(time.RFC3339Nano)
	}
	if slug.Valid {
		th.Slug = slug.String
	}
	return nil
}

func (rep *Repository) GetBySlug(th *models.Thread) error {
	row := rep.db.QueryRow(
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug  "+
			"FROM thread t "+
			"JOIN forum f on t.forum = f.slug "+
			"WHERE t.slug = $1",
		th.Slug,
	)
	created := sql.NullTime{}
	slug := sql.NullString{}
	if err := row.Scan(
		&th.Id,
		&th.Title,
		&th.Message,
		&created,
		&slug,
		&th.Author,
		&th.Forum,
	); err != nil {
		return err
	}
	if created.Valid {
		th.Created = created.Time.Format(time.RFC3339Nano)
	}
	if slug.Valid {
		th.Slug = slug.String
	}
	return nil
}

func (rep *Repository) GetBySlugOrId(th *models.Thread) error {
	slug := sql.NullString{}
	created := sql.NullTime{}
	votes := sql.NullInt64{}
	err := rep.db.QueryRow("SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug, SUM(v.vote)::integer "+
		"FROM thread t "+
		"JOIN forum f on t.forum = f.slug "+
		"LEFT JOIN vote v on t.id = v.thread "+
		"WHERE t.slug = $1 OR t.id::citext = $1 "+
		"GROUP BY f.slug, t.id", th.Slug).Scan(
		&th.Id, &th.Title, &th.Message, &created, &slug, &th.Author, &th.Forum, &votes,
	)
	if err != nil {
		return err
	}
	if created.Valid {
		th.Created = created.Time.Format(time.RFC3339Nano)
	}
	if slug.Valid {
		th.Slug = slug.String
	}
	if votes.Valid {
		th.Votes = votes.Int64
	}
	return nil
}

func (rep *Repository) InsertIntoVotes(v *models.Vote) error {
	err := rep.db.QueryRow(
		"INSERT INTO vote (vote, usr, thread) VALUES ($1 , $2, $3) "+
			"ON CONFLICT (usr,thread) "+
			"DO UPDATE SET vote = excluded.vote "+
			"RETURNING thread",
		v.Vote,
		v.Nickname,
		v.Thread,
	).Scan(&v.Thread)
	if err != nil {
		return err
	}
	return nil
}

func (rep *Repository) GetVotes(th *models.Thread, v *models.Vote) error {
	votes := sql.NullInt64{}
	err := rep.db.QueryRow(
		"SELECT SUM(v.vote) from vote v "+
			"WHERE v.thread = $1",
		v.Thread,
	).Scan(&votes)
	if err != nil {
		logrus.Error("SQL", err)
	}
	if votes.Valid {
		th.Votes = votes.Int64
	}
	return nil
}
