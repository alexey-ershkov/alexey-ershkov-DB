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

func (rep *Repository) InsertInto(tx *pgx.Tx, th *models.Thread) error {
	sqlThreadInsertInto := "INSERT INTO thread (usr, created, forum, message, title, slug) VALUES ($1, $2, $3, $4, $5, $6)" +
		"ON CONFLICT DO NOTHING " +
		"RETURNING id"
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
		sqlThreadInsertInto,
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

func (rep *Repository) GetCreated(tx *pgx.Tx, th *models.Thread) error {
	sqlGetThreadCreated := "SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug  " +
		"FROM thread t " +
		"JOIN forum f on t.forum = f.slug " +
		"WHERE t.usr = $1 AND t.forum = $2 AND t.message = $3 AND t.title = $4"
	row := tx.QueryRow(
		sqlGetThreadCreated,
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

func (rep *Repository) GetBySlugOrId(tx *pgx.Tx, th *models.Thread) error {
	Id, err := strconv.ParseInt(th.Slug, 10, 64)
	if err == nil {
		th.Slug = ""
		th.Id = Id
	}
	slug := sql.NullString{}
	created := sql.NullTime{}
	votes := sql.NullInt64{}
	if th.Slug != "" {
		sqlGetThreadBySlug := "SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug, t.votes " +
			"FROM thread t " +
			"JOIN forum f on t.forum = f.slug " +
			"WHERE t.slug  = $1 "
		err = tx.QueryRow(
			sqlGetThreadBySlug,
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
	} else {
		sqlGetThreadById := "SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug, t.votes " +
			"FROM thread t " +
			"JOIN forum f on t.forum = f.slug " +
			"WHERE t.id = $1 "
		err = tx.QueryRow(
			sqlGetThreadById,
			th.Id).Scan(
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
	} else {
		th.Slug = ""
	}
	if votes.Valid {
		th.Votes = votes.Int64
	}
	return nil
}

func (rep *Repository) InsertIntoVotes(tx *pgx.Tx, v *models.Vote) error {
	sqlInsertVote := "INSERT INTO vote (vote, usr, thread) VALUES ($1 , $2, $3) " +
		"ON CONFLICT (usr,thread) " +
		"DO UPDATE SET vote = excluded.vote " +
		"RETURNING thread"
	err := tx.QueryRow(
		sqlInsertVote,
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
	sqlGetVotes := "SELECT votes from thread " +
		"WHERE id = $1"
	votes := sql.NullInt64{}
	err := tx.QueryRow(
		sqlGetVotes,
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
		sqlGetThreadBySlugOrId := "SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug, t.votes " +
			"FROM thread t " +
			"JOIN forum f on t.forum = f.slug " +
			"WHERE t.id::citext = $1 OR t.slug  = $1 "
		err = tx.QueryRow(sqlGetThreadBySlugOrId,
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
		sqlThreadUpdateMessage := "UPDATE thread SET " +
			"message = $1 " +
			"WHERE id::citext = $2 or slug = $2 " +
			"RETURNING id, title, message, created, slug, usr, forum, votes"
		err = tx.QueryRow(sqlThreadUpdateMessage,
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
		sqlThreadUpdateTitle := "UPDATE thread SET " +
			"title = $1 " +
			"WHERE id::citext = $2 or slug = $2 " +
			"RETURNING id, title, message, created, slug, usr, forum, votes"
		err = tx.QueryRow(sqlThreadUpdateTitle,
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
		sqlThreadUpdateAll := "UPDATE thread SET " +
			"message = $1, " +
			"title = $2 " +
			"WHERE id::citext = $3 or slug = $3 " +
			"RETURNING id, title, message, created, slug, usr, forum, votes"
		err = tx.QueryRow(sqlThreadUpdateAll,
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
	//var err error
	//_, err = rep.db.Prepare("thread_posts_tree_asc",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 "+
	//		"ORDER BY p.path ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_tree_desc",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 "+
	//		"ORDER BY p.path DESC ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_tree_asc_with_limit",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 "+
	//		"ORDER BY p.path "+
	//		"LIMIT $2 ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_tree_desc_with_limit",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 "+
	//		"ORDER BY p.path DESC "+
	//		"LIMIT $2 ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_tree_asc_with_since",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 AND p.path::bigint[] > (SELECT path FROM post WHERE id = $2 )::bigint[] "+
	//		"ORDER BY p.path ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_tree_desc_with_since",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 AND p.path::bigint[] < (SELECT path FROM post WHERE id = $2 )::bigint[] "+
	//		"ORDER BY p.path DESC ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_tree_asc_with_since_with_limit",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 AND p.path::bigint[] > (SELECT path FROM post WHERE id = $2 )::bigint[] "+
	//		"ORDER BY p.path "+
	//		"LIMIT $3",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_tree_desc_with_since_with_limit",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 AND p.path::bigint[] < (SELECT path FROM post WHERE id = $2 )::bigint[] "+
	//		"ORDER BY p.path DESC "+
	//		"LIMIT $3",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_parent_asc",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
	//		"("+
	//		"   SELECT * FROM post p2 "+
	//		"   WHERE p2.thread = $1 AND p2.parent = 0 "+
	//		"	ORDER BY p2.path "+
	//		") "+
	//		"AS prt "+
	//		"JOIN post p ON prt.path[1] = p.path[1] "+
	//		"ORDER BY p.path[1] , p.path ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_parent_desc",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
	//		"("+
	//		"   SELECT * FROM post p2 "+
	//		"   WHERE p2.thread = $1 AND p2.parent = 0 "+
	//		"	ORDER BY p2.path DESC "+
	//		") "+
	//		"AS prt "+
	//		"JOIN post p ON prt.path[1] = p.path[1] "+
	//		"ORDER BY p.path[1] DESC , p.path ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_parent_asc_with_limit",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
	//		"("+
	//		"   SELECT * FROM post p2 "+
	//		"   WHERE p2.thread = $1 AND p2.parent = 0 "+
	//		"	ORDER BY p2.path "+
	//		"	LIMIT $2"+
	//		") "+
	//		"AS prt "+
	//		"JOIN post p ON prt.path[1] = p.path[1] "+
	//		"ORDER BY p.path[1] , p.path ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_parent_desc_with_limit",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
	//		"("+
	//		"   SELECT * FROM post p2 "+
	//		"   WHERE p2.thread = $1 AND p2.parent = 0 "+
	//		"	ORDER BY p2.path DESC "+
	//		"	LIMIT $2"+
	//		") "+
	//		"AS prt "+
	//		"JOIN post p ON prt.path[1] = p.path[1] "+
	//		"ORDER BY p.path[1] DESC , p.path ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_parent_asc_with_since",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
	//		"("+
	//		"   SELECT * FROM post p2 "+
	//		"   WHERE p2.thread = $1 AND p2.parent = 0 "+
	//		"	AND p2.path[1] > (SELECT path[1] FROM post WHERE id = $2 ) "+
	//		"	ORDER BY p2.path "+
	//		") "+
	//		"AS prt "+
	//		"JOIN post p ON prt.path[1] = p.path[1] "+
	//		"ORDER BY p.path[1] , p.path ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_parent_desc_with_since",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
	//		"("+
	//		"   SELECT * FROM post p2 "+
	//		"   WHERE p2.thread = $1 AND p2.parent = 0 "+
	//		"	AND p2.path[1] < (SELECT path[1] FROM post WHERE id = $2 ) "+
	//		"	ORDER BY p2.path DESC "+
	//		") "+
	//		"AS prt "+
	//		"JOIN post p ON prt.path[1] = p.path[1] "+
	//		"ORDER BY p.path[1] DESC , p.path ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_parent_asc_with_since_with_limit",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
	//		"("+
	//		"   SELECT * FROM post p2 "+
	//		"   WHERE p2.thread = $1 AND p2.parent = 0 "+
	//		"	AND p2.path[1] > (SELECT path[1] FROM post WHERE id = $2 ) "+
	//		"	ORDER BY p2.path "+
	//		"	LIMIT $3"+
	//		") "+
	//		"AS prt "+
	//		"JOIN post p ON prt.path[1] = p.path[1] "+
	//		"ORDER BY p.path[1] , p.path ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_posts_parent_desc_with_since_with_limit",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
	//		"("+
	//		"   SELECT * FROM post p2 "+
	//		"   WHERE p2.thread = $1 AND p2.parent = 0 "+
	//		"	AND p2.path[1] < (SELECT path[1] FROM post WHERE id = $2 ) "+
	//		"	ORDER BY p2.path DESC "+
	//		"	LIMIT $3"+
	//		") "+
	//		"AS prt "+
	//		"JOIN post p ON prt.path[1] = p.path[1] "+
	//		"ORDER BY p.path[1] DESC , p.path ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_post_flat_asc",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 "+
	//		"ORDER BY p.created, p.id",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_post_flat_desc",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 "+
	//		"ORDER BY p.created DESC , p.id DESC ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_post_flat_asc_with_limit",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 "+
	//		"ORDER BY p.created, p.id "+
	//		"LIMIT $2 ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_post_flat_desc_with_limit",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 "+
	//		"ORDER BY p.created DESC , p.id DESC "+
	//		"LIMIT $2",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_post_flat_asc_with_since",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 AND p.id > $2 "+
	//		"ORDER BY p.created, p.id",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_post_flat_desc_with_since",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 AND p.id < $2 "+
	//		"ORDER BY p.created DESC , p.id DESC ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_post_flat_asc_with_since_with_limit",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 AND p.id > $2 "+
	//		"ORDER BY p.created, p.id "+
	//		"LIMIT $3 ",
	//)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = rep.db.Prepare("thread_post_flat_desc_with_since_with_limit",
	//	"SELECT p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
	//		"FROM post p "+
	//		"WHERE p.thread = $1 AND p.id < $2 "+
	//		"ORDER BY p.created DESC , p.id DESC "+
	//		"LIMIT $3",
	//)
	//if err != nil {
	//	return err
	//}

	return nil
}
