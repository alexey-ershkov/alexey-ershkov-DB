package repository

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/thread"
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"strconv"
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

func (rep *Repository) InsertIntoForumUsers (tx *pgx.Tx, forum, nickname string) {
	var buffer string

	err := tx.QueryRow("get_forum_user", forum, nickname).Scan(&buffer)
	if err != nil {
		_, err = rep.db.Exec(
			"forum_users_insert_into",
			forum,
			nickname,
		)
		if err != nil {
			logrus.Error("Insert into forum users " + err.Error())
		}
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
	if err := row.Scan(&th.Id); err != nil {
		return err
	}

	rep.InsertIntoForumUsers(tx, th.Forum, th.Author)

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
		&th.Votes,
	); err != nil {
		return err
	}
	if created.Valid {
		th.Created = created.Time.Format(time.RFC3339Nano)
	}
	if slug.Valid {
		th.Slug = slug.String
	} else {
		th.Slug = ""
	}
	return nil
}

func (rep *Repository) GetById(tx *pgx.Tx, th *models.Thread) error {

	row := tx.QueryRow(
		"thread_get_by_id",
		th.Id,
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
		&th.Votes,
	); err != nil {
		return err
	}
	if created.Valid {
		th.Created = created.Time.Format(time.RFC3339Nano)
	}
	if slug.Valid {
		th.Slug = slug.String
	} else {
		th.Slug = ""
	}

	return nil
}

func (rep *Repository) GetBySlugOrId(tx *pgx.Tx, th *models.Thread) error {

	Id, err := strconv.ParseInt(th.Slug, 10, 64)
	if err == nil {
		th.Id = Id
		th.Slug = ""
	}

	if th.Slug != "" {
		err = rep.GetBySlug(tx, th)
	} else {
		err = rep.GetById(tx, th)
	}

	if err != nil {
		return err
	}

	return nil
}

func (rep *Repository) InsertIntoVotes(tx *pgx.Tx, th *models.Thread, v *models.Vote) error {
	var vote int32
	vote = 0

	err := tx.QueryRow(
		"votes_get_info",
		v.Nickname,
		v.Thread,
	).Scan(&vote)


	if vote == 0 {
		err = tx.QueryRow(
			"votes_insert_into",
			v.Vote,
			v.Nickname,
			v.Thread,
		).Scan(&v.Thread)
		th.Votes += int64(v.Vote)
	} else {
		if vote != v.Vote {
			err = tx.QueryRow(
				"votes_update",
				v.Vote,
				v.Nickname,
				v.Thread,
			).Scan(&v.Thread)
			th.Votes += 2*int64(v.Vote)
		}
	}

	if err != nil {
		return err
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
		err = rep.GetBySlugOrId(tx, th)
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
			&votes,
		)
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
			&votes,
		)
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
			&votes,
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


func (rep *Repository) Prepare() error {
	_, err := rep.db.Prepare("thread_insert_into",
		"INSERT INTO thread (usr, created, forum, message, title, slug) VALUES ($1, $2, $3, $4, $5, $6)"+
			"ON CONFLICT DO NOTHING "+
			"RETURNING id",
	)
	if err != nil {
		return err
	}

	//TODO проверка на существование записи в таблице forum users
	_, err = rep.db.Prepare("get_forum_user",
		"SELECT nickname FROM forum_users " +
			"WHERE forum = $1 AND nickname = $2 ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("forum_users_insert_into",
		"INSERT INTO forum_users (forum, nickname) "+
			"VALUES ($1,$2) ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_get_by_slug",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, t.forum, t.votes  "+
			"FROM thread t "+
			"WHERE t.slug = $1",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_get_by_id",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, t.forum, t.votes "+
			"FROM thread t "+
			"WHERE t.id = $1 ",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("votes_insert_into",
		"INSERT INTO vote (vote, usr, thread) VALUES ($1 , $2, $3) "+
			"RETURNING thread",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("votes_update",
		"UPDATE vote SET vote = $1 WHERE usr = $2 and thread = $3 "+
			"RETURNING thread",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("votes_get_info",
		"SELECT vote FROM vote " +
			"WHERE usr = $1 and thread = $2 ",
	)
	if err != nil {
		return err
	}


	_, err = rep.db.Prepare("thread_update_all",
		"UPDATE thread SET "+
			"message = $1, "+
			"title = $2 "+
			"WHERE id::citext = $3 or slug = $3 "+
			"RETURNING id, title, message, created, slug, usr, forum, votes",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_update_message",
		"UPDATE thread SET "+
			"message = $1 "+
			"WHERE id::citext = $2 or slug = $2 "+
			"RETURNING id, title, message, created, slug, usr, forum, votes",
	)
	if err != nil {
		return err
	}

	_, err = rep.db.Prepare("thread_update_title",
		"UPDATE thread SET "+
			"title = $1 "+
			"WHERE id::citext = $2 or slug = $2 "+
			"RETURNING id, title, message, created, slug, usr, forum, votes",
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
