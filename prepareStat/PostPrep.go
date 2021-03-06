package prepareStat

import (
	"log"

	"github.com/jackc/pgx"
)

const getPostDetailsQuery = `SELECT id, user_nick::TEXT, message, created, forum_slug::TEXT, thread_id, is_edited, parent FROM post WHERE id=$1`

const getPostsFlatSinceLimit = `SELECT id,
	user_nick::TEXT,
	message,
	created,
	forum_slug::TEXT,
	thread_id,
	is_edited,
	parent
FROM post
WHERE thread_id=$1
	AND id > $3::TEXT::INTEGER
ORDER BY id
LIMIT $2::TEXT::BIGINT`

const getPostsFlatSinceLimitDesc = `SELECT id,
	user_nick::TEXT,
	message,
	created,
	forum_slug::TEXT,
	thread_id,
	is_edited,
	parent
FROM post
WHERE thread_id=$1
	AND id < $3::TEXT::INTEGER
ORDER BY id DESC
LIMIT $2::TEXT::BIGINT`

const getPostsFlatLimitDesc = `SELECT id,
	user_nick::TEXT,
	message,
	created,
	forum_slug::TEXT,
	thread_id,
	is_edited,
	parent
FROM post
WHERE thread_id=$1
ORDER BY id DESC
LIMIT $2::TEXT::BIGINT`

const getPostsFlatLimit = `SELECT id,
	user_nick::TEXT,
	message,
	created,
	forum_slug::TEXT,
	thread_id,
	is_edited,
	parent
FROM post
WHERE thread_id=$1
ORDER BY id
LIMIT $2::TEXT::BIGINT`

const getPostsTreeSinceLimit = `SELECT id,
	user_nick::TEXT,
	message,
	created,
	forum_slug::TEXT,
	thread_id,
	is_edited,
	parent
FROM post
WHERE thread_id = $1
	AND parents > (SELECT parents FROM post WHERE id = $3::TEXT::INTEGER)
ORDER BY parents
LIMIT $2::TEXT::BIGINT`

const getPostsTreeSinceLimitDesc = `SELECT id,
	user_nick::TEXT,
	message,
	created,
	forum_slug::TEXT,
	thread_id,
	is_edited,
	parent
FROM post
WHERE thread_id = $1
	AND parents < (SELECT parents FROM post WHERE id = $3::TEXT::INTEGER)
ORDER BY parents DESC
LIMIT $2::TEXT::BIGINT`

const getPostsTreeLimit = `SELECT id,
	user_nick::TEXT,
	message,
	created,
	forum_slug::TEXT,
	thread_id,
	is_edited,
	parent
FROM post
WHERE thread_id = $1
ORDER BY parents
LIMIT $2::TEXT::BIGINT`

const getPostsTreeLimitDesc = `SELECT id,
	user_nick::TEXT,
	message,
	created,
	forum_slug::TEXT,
	thread_id,
	is_edited,
	parent
FROM post
WHERE thread_id = $1
ORDER BY parents DESC
LIMIT $2::TEXT::BIGINT`

const getPostsParentTreeSinceLimit = `SELECT p.id,
	p.user_nick::TEXT,
	p.message,
	p.created,
	p.forum_slug::TEXT,
	p.thread_id,
	p.is_edited,
	p.parent
FROM post p
JOIN (
	SELECT id
	FROM post
	WHERE parent=0
		AND thread_id = $1
		AND main_parent > (SELECT main_parent
			FROM post
			WHERE id = $3::TEXT::INTEGER)
	ORDER BY id
	LIMIT $2::TEXT::INTEGER) s
ON p.main_parent=s.id
ORDER BY p.parents`


const getPostsParentTreeSinceLimitDesc = `SELECT p.id,
	p.user_nick::TEXT,
	p.message,
	p.created,
	p.forum_slug::TEXT,
	p.thread_id,
	p.is_edited,
	p.parent
FROM post p
JOIN (
	SELECT id
	FROM post
	WHERE parent=0
		AND thread_id = $1
		AND main_parent < (SELECT main_parent
			FROM post
			WHERE id = $3::TEXT::INTEGER)
	ORDER BY id DESC
	LIMIT $2::TEXT::INTEGER) s
ON p.main_parent=s.id
ORDER BY p.parents[1] DESC, p.parents[2:]`


const getPostsParentTreeLimitDesc = `SELECT p.id,
	p.user_nick::TEXT,
	p.message,
	p.created,
	p.forum_slug::TEXT,
	p.thread_id,
	p.is_edited,
	p.parent
FROM post p
JOIN (
	SELECT id
	FROM post
	WHERE parent=0 AND thread_id = $1
	ORDER BY id DESC
	LIMIT $2::TEXT::INTEGER) s
ON p.main_parent=s.id
ORDER BY p.parents[1] DESC, p.parents[2:]`


const getPostsParentTreeLimit = `SELECT p.id,
	p.user_nick::TEXT,
	p.message,
	p.created,
	p.forum_slug::TEXT,
	p.thread_id,
	p.is_edited,
	p.parent
FROM post p
JOIN (
	SELECT id
	FROM post
	WHERE parent=0 AND thread_id = $1
	ORDER BY id
	LIMIT $2::TEXT::INTEGER) s
ON p.main_parent=s.id
ORDER BY p.parents`

const selectParentAndParents = `SELECT thread_id,
	parents
FROM post
WHERE id = $1`

const insertIntoPost = `INSERT INTO post(id, user_nick, message, created, forum_slug, thread_id, parent, parents, main_parent)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING created`

func PreparePost(tx *pgx.ConnPool) {
	if _, err := tx.Prepare("getPostDetailsQuery", getPostDetailsQuery); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("insertIntoPost", insertIntoPost); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("selectParentAndParents", selectParentAndParents); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getPostsFlatLimit", getPostsFlatLimit); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getPostsFlatLimitDesc", getPostsFlatLimitDesc); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getPostsFlatSinceLimit", getPostsFlatSinceLimit); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getPostsFlatSinceLimitDesc", getPostsFlatSinceLimitDesc); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getPostsTreeSinceLimit", getPostsTreeSinceLimit); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getPostsTreeSinceLimitDesc", getPostsTreeSinceLimitDesc); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getPostsTreeLimit", getPostsTreeLimit); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getPostsTreeLimitDesc", getPostsTreeLimitDesc); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getPostsParentTreeSinceLimit", getPostsParentTreeSinceLimit); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getPostsParentTreeSinceLimitDesc", getPostsParentTreeSinceLimitDesc); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getPostsParentTreeLimitDesc", getPostsParentTreeLimitDesc); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getPostsParentTreeLimit", getPostsParentTreeLimit); err != nil {
		log.Fatalln(err)
	}
}
