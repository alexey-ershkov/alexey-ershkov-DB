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

func (rep *Repository) Update(tx *pgx.Tx, th *models.Thread) error {
	slug := sql.NullString{}
	created := sql.NullTime{}
	votes := sql.NullInt64{}
	var err error
	switch true {
	case th.Message == "" && th.Title == "":
		err = tx.QueryRow("thread_get_by_slug_or_id",
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
	case th.Message != "" && th.Title == "":
		err = tx.QueryRow("thread_update_message",
			th.Message,
			th.Slug).Scan(
			&th.Id,
			&th.Title,
			&th.Message,
			&created,
			&slug,
			&th.Author,
			&th.Forum,
		)
		e := tx.QueryRow(
			"votes_get",
			th.Id,
		).Scan(&votes)
		if e != nil {
			logrus.Error("SQL", err)
		}
	case th.Message == "" && th.Title != "":
		err = tx.QueryRow("thread_update_title",
			th.Title,
			th.Slug).Scan(
			&th.Id,
			&th.Title,
			&th.Message,
			&created,
			&slug,
			&th.Author,
			&th.Forum,
		)
		e := tx.QueryRow(
			"votes_get",
			th.Id,
		).Scan(&votes)
		if e != nil {
			logrus.Error("SQL", err)
		}
	case th.Message != "" && th.Title != "":
		err = tx.QueryRow("thread_update_all",
			th.Message,
			th.Title,
			th.Slug).Scan(
			&th.Id,
			&th.Title,
			&th.Message,
			&created,
			&slug,
			&th.Author,
			&th.Forum,
		)
		e := tx.QueryRow(
			"votes_get",
			th.Id,
		).Scan(&votes)
		if e != nil {
			logrus.Error("SQL", err)
		}
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
	if votes.Valid {
		th.Votes = votes.Int64
	}

	return nil
}

func (rep *Repository) GetPosts(tx *pgx.Tx, th *models.Thread, desc, sort, limit, since string) ([]models.Post, error) {
	posts := make([]models.Post, 0)
	var err error
	var rows *pgx.Rows
	if sort == "tree" {
		switch true {
		case desc != "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_posts_tree_asc", th.Id)
		case desc == "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_posts_tree_desc", th.Id)
		case desc != "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_posts_tree_asc_with_since", th.Id, since)
		case desc == "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_posts_tree_desc_with_since", th.Id, since)
		case desc != "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_posts_tree_asc_with_limit", th.Id, limit)
		case desc == "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_posts_tree_desc_with_limit", th.Id, limit)
		case desc != "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_posts_tree_asc_with_since_with_limit", th.Id, since, limit)
		case desc == "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_posts_tree_desc_with_since_with_limit", th.Id, since, limit)
		}
	} else if sort == "parent_tree" {
		switch true {
		case desc != "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_posts_parent_asc", th.Id)
		case desc == "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_posts_parent_desc", th.Id)
		case desc != "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_posts_parent_asc_with_since", th.Id, since)
		case desc == "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_posts_parent_desc_with_since", th.Id, since)
		case desc != "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_posts_parent_asc_with_limit", th.Id, limit)
		case desc == "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_posts_parent_desc_with_limit", th.Id, limit)
		case desc != "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_posts_parent_asc_with_since_with_limit", th.Id, since, limit)
		case desc == "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_posts_parent_desc_with_since_with_limit", th.Id, since, limit)
		}
	} else {
		switch true {
		case desc != "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_post_flat_asc", th.Id)
		case desc == "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_post_flat_desc", th.Id)
		case desc != "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_post_flat_asc_with_since", th.Id, since)
		case desc == "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_post_flat_desc_with_since", th.Id, since)
		case desc != "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_post_flat_asc_with_limit", th.Id, limit)
		case desc == "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_post_flat_desc_with_limit", th.Id, limit)
		case desc != "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_post_flat_asc_with_since_with_limit", th.Id, since, limit)
		case desc == "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_post_flat_desc_with_since_with_limit", th.Id, since, limit)
		}
	}

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

	_, err = rep.db.Prepare("thread_update_all",
		"UPDATE thread SET "+
			"message = $1, "+
			"title = $2 "+
			"WHERE id::citext = $3 or slug = $3 "+
			"RETURNING id, title, message, created, slug, usr, forum",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_update_message",
		"UPDATE thread SET "+
			"message = $1 "+
			"WHERE id::citext = $2 or slug = $2 "+
			"RETURNING id, title, message, created, slug, usr, forum",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_update_title",
		"UPDATE thread SET "+
			"title = $1 "+
			"WHERE id::citext = $2 or slug = $2 "+
			"RETURNING id, title, message, created, slug, usr, forum",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_tree_asc",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.path ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_tree_desc",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.path DESC ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_tree_asc_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.path "+
			"LIMIT $2 ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_tree_desc_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.path DESC "+
			"LIMIT $2 ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_tree_asc_with_since",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.path::bigint[] > (SELECT path FROM post WHERE id = $2 )::bigint[] "+
			"ORDER BY p.path ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_tree_desc_with_since",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.path::bigint[] < (SELECT path FROM post WHERE id = $2 )::bigint[] "+
			"ORDER BY p.path DESC ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_tree_asc_with_since_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.path::bigint[] > (SELECT path FROM post WHERE id = $2 )::bigint[] "+
			"ORDER BY p.path "+
			"LIMIT $3",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_tree_desc_with_since_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.path::bigint[] < (SELECT path FROM post WHERE id = $2 )::bigint[] "+
			"ORDER BY p.path DESC "+
			"LIMIT $3",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_parent_asc",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	ORDER BY p2.path "+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_parent_desc",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	ORDER BY p2.path DESC "+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] DESC , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_parent_asc_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	ORDER BY p2.path "+
			"	LIMIT $2"+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_parent_desc_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	ORDER BY p2.path DESC "+
			"	LIMIT $2"+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] DESC , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_parent_asc_with_since",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	AND p2.path[1] > (SELECT path[1] FROM post WHERE id = $2 ) "+
			"	ORDER BY p2.path "+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_parent_desc_with_since",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	AND p2.path[1] < (SELECT path[1] FROM post WHERE id = $2 ) "+
			"	ORDER BY p2.path DESC "+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] DESC , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_parent_asc_with_since_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	AND p2.path[1] > (SELECT path[1] FROM post WHERE id = $2 ) "+
			"	ORDER BY p2.path "+
			"	LIMIT $3"+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_posts_parent_desc_with_since_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"("+
			"   SELECT * FROM post p2 "+
			"   WHERE p2.thread = $1 AND p2.parent = 0 "+
			"	AND p2.path[1] < (SELECT path[1] FROM post WHERE id = $2 ) "+
			"	ORDER BY p2.path DESC "+
			"	LIMIT $3"+
			") "+
			"AS prt "+
			"JOIN post p ON prt.path[1] = p.path[1] "+
			"ORDER BY p.path[1] DESC , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_post_flat_asc",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.created, p.id",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_post_flat_desc",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.created DESC , p.id DESC ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_post_flat_asc_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.created, p.id "+
			"LIMIT $2 ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_post_flat_desc_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 "+
			"ORDER BY p.created DESC , p.id DESC "+
			"LIMIT $2",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_post_flat_asc_with_since",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.id > $2 "+
			"ORDER BY p.created, p.id",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_post_flat_desc_with_since",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.id < $2 "+
			"ORDER BY p.created DESC , p.id DESC ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_post_flat_asc_with_since_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.id > $2 "+
			"ORDER BY p.created, p.id "+
			"LIMIT $3 ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_post_flat_desc_with_since_with_limit",
		"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"FROM post p "+
			"WHERE p.thread = $1 AND p.id < $2 "+
			"ORDER BY p.created DESC , p.id DESC "+
			"LIMIT $3",
	)
	if err != nil {
		return err
	}

	return nil
}
