package repository

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/forum"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"database/sql"
	"github.com/jackc/pgx"
	"time"
)

type Repository struct {
	db *pgx.ConnPool
}

func NewForumRepository(db *pgx.ConnPool) forum.Repository {
	return &Repository{db: db}
}

func (rep *Repository) InsertInto(tx *pgx.Tx, f *models.Forum) error {
	row := tx.QueryRow(
		"forum_insert_into",
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

func (rep *Repository) GetBySlug(tx *pgx.Tx, f *models.Forum) error {
	row := tx.QueryRow(
		"forum_get_by_slug",
		f.Slug,
	)
	if err := row.Scan(&f.Posts, &f.Slug, &f.Threads, &f.Title, &f.User); err != nil {
		return err
	}
	return nil
}

//TODO можно переписать на prepared statement
func (rep *Repository) GetThreads(tx *pgx.Tx, f *models.Forum, desc, limit, since string) ([]models.Thread, error) {
	ths := make([]models.Thread, 0)
	var sqlStr string
	if desc == "true" {
		sqlStr = "SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug FROM thread t " +
			"JOIN forum f on t.forum = f.slug " +
			"WHERE f.slug = $1 AND t.created <=  $2::timestamp AT TIME ZONE '0'" +
			"ORDER BY t.created DESC "
	} else {
		sqlStr = "SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug FROM thread t " +
			"JOIN forum f on t.forum = f.slug " +
			"WHERE f.slug = $1 AND t.created >=  $2::timestamp AT TIME ZONE '0'" +
			"ORDER BY t.created "
	}
	if limit != "" {
		sqlStr += "LIMIT " + limit
	}
	if since == "" {
		if desc == "true" {
			since = "infinity"
		} else {
			since = "-infinity"
		}
	}
	rows, err := tx.Query(sqlStr, f.Slug, since)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		created := sql.NullTime{}
		slug := sql.NullString{}
		th := models.Thread{}
		err = rows.Scan(&th.Id, &th.Title, &th.Message, &created, &slug, &th.Author, &th.Forum)
		if err != nil {
			return nil, err
		}
		if slug.Valid {
			th.Slug = slug.String
		}
		if created.Valid {
			th.Created = created.Time.Format(time.RFC3339Nano)
		}
		ths = append(ths, th)
	}
	rows.Close()
	return ths, nil
}

//TODO можно переписать на prepared statement
func (rep *Repository) GetUsers(tx *pgx.Tx, f *models.Forum, desc, limit, since string) ([]models.User, error) {
	usr := make([]models.User, 0)

	sqlQuery := "SELECT u.email, u.fullname, u.nickname, u.about " +
		"FROM forum_users " +
		"JOIN usr u on forum_users.nickname = u.nickname " +
		"WHERE forum = $1 "

	if since != "" {
		if desc == "true" {
			sqlQuery += "AND u.nickname < $2 "
		} else {
			sqlQuery += "AND u.nickname > $2 "
		}
	}

	sqlQuery += "ORDER BY u.nickname "

	if desc == "true" {
		sqlQuery += "DESC "
	}

	if limit != "" {
		sqlQuery += "LIMIT " + limit
	}

	var rows *pgx.Rows
	var err error
	if since != "" {
		rows, err = tx.Query(sqlQuery, f.Slug, since)
	} else {
		rows, err = tx.Query(sqlQuery, f.Slug)
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.Email, &u.Fullname, &u.Nickname, &u.About)
		if err != nil {
			rows.Close()
			return nil, err
		}
		usr = append(usr, u)
	}
	rows.Close()
	return usr, nil
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
	_, err := rep.db.Prepare("forum_insert_into",
		"INSERT INTO forum (slug, title, usr) "+
			"VALUES ($1, $2, $3) "+
			"ON CONFLICT DO NOTHING "+
			"RETURNING title",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_by_slug",
		"SELECT count(p), f.slug, (SELECT count(*) FROM forum f2 JOIN thread t2 on f2.slug = t2.forum WHERE f2.slug = $1), f.title, u.nickname FROM forum f "+
			"LEFT JOIN thread t on f.slug = t.forum "+
			"LEFT JOIN post p on t.id = p.thread "+
			"JOIN usr u on f.usr = u.nickname "+
			"WHERE f.slug = $1 "+
			"GROUP BY f.slug, u.nickname")
	if err != nil {
		return err
	}

	return nil
}
