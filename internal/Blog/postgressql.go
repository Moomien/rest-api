package Blog

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSql struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(p *pgxpool.Pool) *PostgresSql {
	return &PostgresSql{
		pool: p,
	}
}

func (pg *PostgresSql) Save(data Blog) error {
	ctx := context.Background()
	//начинаем транзакцию
	tx, err := pg.pool.Begin(ctx)
	if err != nil {
		return err
	}
	//откат изменений назад если что то пойдет не так
	defer tx.Rollback(ctx)

	if data.ID == 0 {
		//логика INSERT
		var blogID int
		blog := `INSERT INTO blog(title, content, category) VALUES($1, $2, $3) RETURNING id;`
		err = tx.QueryRow(ctx, blog, data.Title, data.Content, data.Category).Scan(&blogID)
		if err != nil {
			return err
		}

		for _, tagName := range data.Tags {
			var tagID int
			tags := `INSERT INTO tags(name) VALUES($1)
				 ON CONFLICT(name) DO UPDATE SET name = EXCLUDED.name
				 RETURNING id`
			err = tx.QueryRow(ctx, tags, tagName).Scan(&tagID)
			if err != nil {
				return err
			}
			blog_tags := `INSERT INTO blog_tags(blog_id, tag_id) VALUES($1, $2)
					  ON CONFLICT DO NOTHING`
			_, err = tx.Exec(ctx, blog_tags, blogID, tagID)
			if err != nil {
				return err
			}
		}
	} else {
		// ЛОГИКА UPDATE
		// 1. Обновляем поля title, content по data.ID
		// 2. Удаляем старые теги: DELETE FROM blog_tags WHERE blog_id = data.ID
		// 3. Записываем новые теги из слайса data.Tags
		updated_time := time.Now()
		blog := `UPDATE blog SET title = $1, content = $2, category = $3, updated_at = $4 WHERE id = $5`
		_, err = tx.Exec(ctx, blog, data.Title, data.Content, data.Category, updated_time, data.ID)
		if err != nil {
			return err
		}
		blog_tags := `DELETE FROM blog_tags WHERE blog_id = $1`
		_, err = tx.Exec(ctx, blog_tags, data.ID)
		for _, tagName := range data.Tags {
			var tagID int
			tags := `INSERT INTO tags(name) VALUES($1)
				 ON CONFLICT(name) DO UPDATE SET name = EXCLUDED.name
				 RETURNING id`
			err = tx.QueryRow(ctx, tags, tagName).Scan(&tagID)
			if err != nil {
				return err
			}
			blog_tags := `INSERT INTO blog_tags(blog_id, tag_id) VALUES($1, $2)
					  ON CONFLICT DO NOTHING`
			_, err = tx.Exec(ctx, blog_tags, data.ID, tagID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}
