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

func (rep *Repository) GetThreads(tx *pgx.Tx, f *models.Forum, desc, limit, since string) ([]models.Thread, error) {
	ths := make([]models.Thread, 0)
	var rows *pgx.Rows
	var err error

	if since == "" {
		if desc == "true" {
			since = "infinity"
		} else {
			since = "-infinity"
		}
	}
	if desc == "true" {
		if limit != "" {
			rows, err = tx.Query("forum_get_threads_desc_with_limit", f.Slug, since, limit)
		} else {
			rows, err = tx.Query("forum_get_threads_desc", f.Slug, since)
		}
	} else {
		if limit != "" {
			rows, err = tx.Query("forum_get_threads_asc_with_limit", f.Slug, since, limit)
		} else {
			rows, err = tx.Query("forum_get_threads_asc", f.Slug, since)
		}
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		created := sql.NullTime{}
		slug := sql.NullString{}
		th := models.Thread{}
		votes := sql.NullInt64{}
		err = rows.Scan(&th.Id, &th.Title, &th.Message, &created, &slug, &th.Author, &th.Forum, &votes)
		if err != nil {
			return nil, err
		}
		if slug.Valid {
			th.Slug = slug.String
		}
		if votes.Valid {
			th.Votes = votes.Int64
		}
		if created.Valid {
			th.Created = created.Time.Format(time.RFC3339Nano)
		}
		ths = append(ths, th)
	}
	rows.Close()
	return ths, nil
}

func (rep *Repository) GetUsers(tx *pgx.Tx, f *models.Forum, desc, limit, since string) ([]models.User, error) {
	usr := make([]models.User, 0)
	var rows *pgx.Rows
	var err error

	switch true {
	case desc != "true" && since == "" && limit == "":
		rows, err = tx.Query("forum_get_users", f.Slug)
	case desc == "true" && since == "" && limit == "":
		rows, err = tx.Query("forum_get_users_desc", f.Slug)
	case desc != "true" && since != "" && limit == "":
		rows, err = tx.Query("forum_get_users_asc_with_since", f.Slug, since)
	case desc == "true" && since != "" && limit == "":
		rows, err = tx.Query("forum_get_users_desc_with_since", f.Slug, since)
	case desc != "true" && since == "" && limit != "":
		rows, err = tx.Query("forum_get_users_with_limit", f.Slug, limit)
	case desc == "true" && since == "" && limit != "":
		rows, err = tx.Query("forum_get_users_desc_with_limit", f.Slug, limit)
	case desc != "true" && since != "" && limit != "":
		rows, err = tx.Query("forum_get_users_asc_with_since_with_limit", f.Slug, since, limit)
	case desc == "true" && since != "" && limit != "":
		rows, err = tx.Query("forum_get_users_desc_with_since_with_limit", f.Slug, since, limit)
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
			"RETURNING title",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_by_slug",
		"SELECT f.posts, f.slug, f.threads,f.title, f.usr "+
			"FROM forum f "+
			"WHERE f.slug = $1 ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_threads_desc",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug, t.votes FROM thread t "+
			"JOIN forum f on t.forum = f.slug "+
			"WHERE t.forum = $1 AND t.created <=  $2::timestamptz "+
			"ORDER BY t.created DESC ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_threads_desc_with_limit",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug, t.votes FROM thread t "+
			"JOIN forum f on t.forum = f.slug "+
			"WHERE t.forum = $1 AND t.created <=  $2::timestamptz "+
			"ORDER BY t.created DESC "+
			"LIMIT $3",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_threads_asc",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug, t.votes FROM thread t "+
			"JOIN forum f on t.forum = f.slug "+
			"WHERE t.forum = $1 AND t.created >=  $2::timestamptz "+
			"ORDER BY t.created ",
	)
	if err != nil {
		return err
	}
	////TODO переписать на t.forum = $1 везде
	////TODO убрать JOIN
	_, err = rep.db.Prepare("forum_get_threads_asc_with_limit",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug, t.votes FROM thread t "+
			"JOIN forum f on t.forum = f.slug "+
			"WHERE t.forum = $1 AND t.created >=  $2::timestamptz "+
			"ORDER BY t.created "+
			"LIMIT $3 ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_users",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 "+
			"ORDER BY u.nickname ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_users_with_limit",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 "+
			"ORDER BY u.nickname "+
			"LIMIT $2 ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_users_desc",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 "+
			"ORDER BY u.nickname DESC ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_users_desc_with_limit",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 "+
			"ORDER BY u.nickname DESC "+
			"LIMIT $2 ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_users_desc_with_since_with_limit",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 AND u.nickname < $2 "+
			"ORDER BY u.nickname DESC "+
			"LIMIT $3 ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_users_desc_with_since",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 AND u.nickname < $2 "+
			"ORDER BY u.nickname DESC",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_users_asc_with_since_with_limit",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 AND u.nickname > $2 "+
			"ORDER BY u.nickname "+
			"LIMIT $3 ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_get_users_asc_with_since",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 AND u.nickname > $2 "+
			"ORDER BY u.nickname ",
	)
	if err != nil {
		return err
	}

	return nil
}
