package repository

import (
	"bytes"
	"log"

	"github.com/jackc/pgx"
	"main/models"
)


func (store *DBStore) CreateForum(forum *models.Forum) (*models.Forum, error) {
	tx, err := store.DB.Begin()
	if err != nil {
	}

	forumExisting := models.Forum{}

	if err = tx.QueryRow("selectForumQuery", &forum.Slug).
		Scan(&forumExisting.Slug, &forumExisting.Title,
		&forumExisting.Posts, &forumExisting.Threads, &forumExisting.Moderator); err == nil {

		tx.Rollback()
		return &forumExisting, models.ForumAlreadyExists
	}

	if err := tx.QueryRow("INSERT INTO forum (slug, title, moderator) VALUES ($1, $2, (SELECT nickname FROM users WHERE nickname=$3)) RETURNING moderator::TEXT", &forum.Slug, &forum.Title, &forum.Moderator).
		Scan(&forum.Moderator); err != nil {
		tx.Rollback()
		return nil, models.UserNotFound
	}

	tx.Commit()
	return forum, nil
}

func (store *DBStore) GetForumDetails(slug interface{}) (*models.Forum, error) {
	forum := models.Forum{}

	err := store.DB.QueryRow("selectForumQuery", &slug).
		Scan(&forum.Slug, &forum.Title, &forum.Posts, &forum.Threads, &forum.Moderator)

	if err != nil {
		return nil, err
	}

	return &forum, err
}


func (store *DBStore) CreateThread(forumSlug interface{}, threadDetails *models.Thread) (*models.Thread, error) {
	tx, err := store.DB.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Rollback()

	user := models.User{}
	var userID, forumID int
	var realForumSlug string

	if err = tx.QueryRow("SELECT id, nickname::text, email::text, about, fullname FROM users WHERE nickname = $1", &threadDetails.User_nick).
		Scan(&userID, &user.Nickname, &user.Email, &user.About, &user.Fullname); err != nil {
		log.Println(err, *forumSlug.(*interface{}))
		tx.Rollback()

		return nil, models.UserNotFound
	}

	if err = tx.QueryRow("SELECT id, slug::text FROM forum WHERE slug = $1", &forumSlug).
		Scan(&forumID, &realForumSlug); err != nil {
		log.Println(err)
		tx.Rollback()

		return nil, models.ForumNotFound
	}

	if err = tx.QueryRow("INSERT INTO thread (slug, title, message, forum_id, forum_slug, user_id, user_nick, created) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING RETURNING id", threadDetails.Slug, &threadDetails.Title,
		&threadDetails.Message, forumID, &realForumSlug, userID, &user.Nickname, &threadDetails.Created).
		Scan(&threadDetails.Id); err != nil {

		existingThread := models.Thread{}

		if err = tx.QueryRow("getThreadBySlug", threadDetails.Slug).
			Scan(&existingThread.Id, &existingThread.Slug, &existingThread.Title,
			&existingThread.Message, &existingThread.Forum_slug, &existingThread.User_nick, &existingThread.Created,
			&existingThread.Votes_count); err == nil {

			tx.Rollback()
			return &existingThread, models.ThreadAlreadyExists
		}

		tx.Rollback()
		log.Fatalln(err)
	}

	if _, err = tx.Exec("insertIntoForumUsers", forumID, user.Nickname, user.Email, user.About, user.Fullname); err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	threadDetails.Forum_slug = realForumSlug
	threadDetails.User_nick = user.Nickname

	tx.Commit()
	return threadDetails, nil
}

func (store *DBStore) GetForumThreads(slug interface{}, limit []byte, since []byte, desc []byte) (*models.ThreadArr, error) {
	var err error
	var rows *pgx.Rows

	if since == nil {
		if bytes.Equal([]byte("true"), desc) {
			rows, err = store.DB.Query("gftLimitDesc", slug, limit)
		} else {
			rows, err = store.DB.Query("gftLimit", slug, limit)
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			rows, err = store.DB.Query("gftCreatedLimitDesc", slug, since, limit)
		} else {
			rows, err = store.DB.Query("gftCreatedLimit", slug, since, limit)
		}
	}

	if err != nil {
		log.Fatalln(err)
	}

	var threads models.ThreadArr

	for rows.Next() {
		thread := models.Thread{}

		if err = rows.Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message,
			&thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count); err != nil {
			log.Fatalln(err)
		}

		threads = append(threads, &thread)
	}
	rows.Close()

	if len(threads) == 0 {
		var forumID int
		if err = store.DB.QueryRow("getForumIDBySlug", &slug).Scan(&forumID); err != nil {
			return nil, models.ForumNotFound
		}
	}

	return &threads, nil
}

func (store *DBStore) GetForumUsers(slug interface{}, limit []byte, since []byte, desc []byte) (*models.UsersArr, error) {
	var err error

	var rows *pgx.Rows

	if since == nil {
		if bytes.Equal([]byte("true"), desc) {
			rows, err = store.DB.Query("gfuLimitDesc", slug, limit)
		} else {
			rows, err = store.DB.Query("gfuLimit", slug, limit)
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			rows, err = store.DB.Query("gfuSinceLimitDesc", slug, since, limit)
		} else {
			rows, err = store.DB.Query("gfuSinceLimit", slug, since, limit)
		}
	}

	if err != nil {
		rows.Close()
	}
	var users models.UsersArr

	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.Nickname, &user.Email, &user.About, &user.Fullname); err != nil {
			rows.Close()
			log.Fatalln(err)
		}
		users = append(users, &user)
	}

	rows.Close()

	if len(users) == 0 {
		var forumID int
		if err = store.DB.QueryRow("getForumIDBySlug", &slug).Scan(&forumID); err != nil {
			return nil, models.ForumNotFound
		}
	}

	return &users, nil
}
