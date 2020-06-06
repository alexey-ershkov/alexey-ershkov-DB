package repository

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/thread"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"time"
)

type Repository struct {
	db *pgx.ConnPool
}

func NewThreadRepository(db *pgx.ConnPool) thread.Repository {
	return &Repository{
		db: db,
	}
}

func (rep *Repository) InsertInto(tx *pgx.Tx, th *models.Thread) error {
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
	row := tx.QueryRow(
		"thread_insert_into",
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
	_, err := tx.Exec(
		"forum_users_insert_into",
		th.Forum,
		th.Author,
	)
	if err != nil {
		return err
	}
	return nil
}

func (rep *Repository) GetCreated(tx *pgx.Tx, th *models.Thread) error {
	row := tx.QueryRow(
		"thread_get_created",
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

func (rep *Repository) GetBySlug(tx *pgx.Tx, th *models.Thread) error {
	row := tx.QueryRow(
		"thread_get_by_slug",
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

func (rep *Repository) GetBySlugOrId(tx *pgx.Tx, th *models.Thread) error {
	slug := sql.NullString{}
	created := sql.NullTime{}
	votes := sql.NullInt64{}
	err := tx.QueryRow(
		"thread_get_by_slug_or_id",
		th.Slug).Scan(
		&th.Id,
		&th.Title,
		&th.Message,
		&created,
		&slug,
		&th.Author,
		&th.Forum,
		&votes,
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

func (rep *Repository) InsertIntoVotes(tx *pgx.Tx, v *models.Vote) error {
	err := tx.QueryRow(
		"votes_insert_into",
		v.Vote,
		v.Nickname,
		v.Thread,
	).Scan(&v.Thread)
	if err != nil {
		return err
	}
	return nil
}

func (rep *Repository) GetVotes(tx *pgx.Tx, th *models.Thread, v *models.Vote) error {
	votes := sql.NullInt64{}
	err := tx.QueryRow(
		"votes_get",
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

//TODO можно переписать на prepared statement
func (rep *Repository) Update(tx *pgx.Tx, th *models.Thread) error {
	slug := sql.NullString{}
	created := sql.NullTime{}
	votes := sql.NullInt64{}
	args := make([]string, 0)
	sqlStr := "UPDATE thread SET "
	if th.Message != "" {
		sqlStr += "message = $1 "
		args = append(args, th.Message)
	}
	if th.Title != "" {
		if len(args) == 1 {
			sqlStr += ","
		}
		sqlStr += " title = $%d "
		args = append(args, th.Title)
		sqlStr = fmt.Sprintf(sqlStr, len(args))
	}
	if len(args) == 0 {
		err := tx.QueryRow("thread_get_by_slug_or_id",
			th.Slug).Scan(
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
	} else {
		sqlStr += "WHERE id::citext = $%d or slug = $%d RETURNING id, title, message, created, slug, usr, forum"
		args = append(args, th.Slug)
		sqlStr = fmt.Sprintf(sqlStr, len(args), len(args))
		var err error
		if len(args) == 2 {
			err = tx.QueryRow(sqlStr, args[0], args[1]).Scan(
				&th.Id,
				&th.Title,
				&th.Message,
				&created,
				&slug,
				&th.Author,
				&th.Forum,
			)
		} else {
			err = tx.QueryRow(sqlStr, args[0], args[1], args[2]).Scan(
				&th.Id,
				&th.Title,
				&th.Message,
				&created,
				&slug,
				&th.Author,
				&th.Forum,
			)
		}
		if err != nil {
			return err
		}
		if created.Valid {
			th.Created = created.Time.Format(time.RFC3339Nano)
		}
		if slug.Valid {
			th.Slug = slug.String
		}
		err = tx.QueryRow(
			"votes_get",
			th.Id,
		).Scan(&votes)
		if err != nil {
			logrus.Error("SQL", err)
		}
		if votes.Valid {
			th.Votes = votes.Int64
		}
	}
	return nil
}

//TODO Можно периписать на prepared statements
func (rep *Repository) GetPosts(tx *pgx.Tx, th *models.Thread, desc, sort, limit, since string) ([]models.Post, error) {
	posts := make([]models.Post, 0)
	var sqlString string
	if sort == "tree" {

		sqlString = "SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread " +
			"FROM post p " +
			"WHERE p.thread = $1 "

		if since != "" {
			if desc == "true" {
				sqlString += "AND p.path::bigint[] < (SELECT path FROM post WHERE id = " + since + " )::bigint[] "
			} else {
				sqlString += "AND p.path::bigint[] > (SELECT path FROM post WHERE id = " + since + " )::bigint[] "
			}
		}

		sqlString += "ORDER BY p.path "

		if desc == "true" {
			sqlString += "DESC "
		}

		if limit != "" {
			sqlString += "LIMIT " + limit + " "
		}
	} else if sort == "parent_tree" {

		sqlString = "SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM " +
			"(" +
			"   SELECT * FROM post p2 " +
			"   WHERE p2.thread = $1 AND p2.parent = 0 "

		if since != "" {
			if desc == "true" {
				sqlString += "AND p2.path[1] < (SELECT path[1] FROM post WHERE id = " + since + " ) "
			} else {
				sqlString += "AND p2.path[1] > (SELECT path[1] FROM post WHERE id = " + since + " ) "
			}
		}
		sqlString += "ORDER BY p2.path "
		if desc == "true" {
			sqlString += "DESC "
		}
		if limit != "" {
			sqlString += "LIMIT " + limit + " "
		}
		sqlString += ") AS prt " +
			"JOIN post p ON prt.path[1] = p.path[1] " +
			"ORDER BY p.path[1] "
		if desc == "true" {
			sqlString += "DESC "
		}
		sqlString += ", p.path "

	} else {

		sqlString = "SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread " +
			"FROM post p " +
			"WHERE p.thread = $1 "

		if since != "" {
			if desc == "true" {
				sqlString += "AND p.id < " + since + " "
			} else {
				sqlString += "AND p.id > " + since + " "
			}
		}
		sqlString += "ORDER BY p.created "
		if desc == "true" {
			sqlString += "DESC "
		}
		sqlString += ", p.id "
		if desc == "true" {
			sqlString += "DESC "
		}
		if limit != "" {
			sqlString += "LIMIT " + limit + " "
		}
	}
	rows, err := tx.Query(sqlString, th.Id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		created := sql.NullTime{}
		p := models.Post{}
		err := rows.Scan(
			&p.Id,
			&p.Author,
			&created,
			&p.Forum,
			&p.IsEdited,
			&p.Message,
			&p.Parent,
			&p.Thread,
		)
		if err != nil {
			return nil, err
		}
		if created.Valid {
			p.Created = created.Time.Format(time.RFC3339Nano)
		}
		posts = append(posts, p)
	}
	rows.Close()
	return posts, nil
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
	_, err := rep.db.Prepare("thread_insert_into",
		"INSERT INTO thread (usr, created, forum, message, title, slug) VALUES ($1, $2, $3, $4, $5, $6)"+
			"ON CONFLICT DO NOTHING "+
			"RETURNING id",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_users_insert_into",
		"INSERT INTO forum_users (forum, nickname) "+
			"VALUES ($1,$2) "+
			"ON CONFLICT (forum,nickname) DO NOTHING ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_get_created",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug  "+
			"FROM thread t "+
			"JOIN forum f on t.forum = f.slug "+
			"WHERE t.usr = $1 AND t.forum = $2 AND t.message = $3 AND t.title = $4",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_get_by_slug",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug  "+
			"FROM thread t "+
			"JOIN forum f on t.forum = f.slug "+
			"WHERE t.slug = $1",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_get_by_slug_or_id",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug, SUM(v.vote)::integer "+
			"FROM thread t "+
			"JOIN forum f on t.forum = f.slug "+
			"LEFT JOIN vote v on t.id = v.thread "+
			"WHERE t.slug = $1 OR t.id::citext = $1 "+
			"GROUP BY f.slug, t.id",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("votes_insert_into",
		"INSERT INTO vote (vote, usr, thread) VALUES ($1 , $2, $3) "+
			"ON CONFLICT (usr,thread) "+
			"DO UPDATE SET vote = excluded.vote "+
			"RETURNING thread",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("votes_get",
		"SELECT SUM(v.vote) from vote v "+
			"WHERE v.thread = $1",
	)
	if err != nil {
		return err
	}

	return nil
}
