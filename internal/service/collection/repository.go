package collection

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	pgx "github.com/jackc/pgx/v5"
)

const (
	collectionTableName     = "collection"
	collectionGameTableName = "collection_game"
)

var (
	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type repository struct {
	db *pgx.Conn
}

func newRepository(db *pgx.Conn) *repository {
	return &repository{db: db}
}

func (r *repository) listCollections(ctx context.Context, userID int64) ([]Collection, error) {
	query, args, err := psql.
		Select("id", "name", "pinned", "created_at", "updated_at").
		From(collectionTableName).
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("pinned DESC, name ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var c Collection
		if err := pgxscan.ScanRow(&c, rows); err != nil {
			return nil, err
		}
		c.UserID = userID
		collections = append(collections, c)
	}

	for i := range collections {
		gameIDs, err := r.getCollectionGameIDs(ctx, collections[i].ID)
		if err != nil {
			return nil, err
		}
		collections[i].GameIDs = gameIDs
	}

	return collections, nil
}

func (r *repository) getCollectionGameIDs(ctx context.Context, collectionID int64) ([]int64, error) {
	query, args, err := psql.
		Select("game_id").
		From(collectionGameTableName).
		Where(squirrel.Eq{"collection_id": collectionID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gameIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		gameIDs = append(gameIDs, id)
	}

	return gameIDs, nil
}

func (r *repository) getCollection(ctx context.Context, collectionID, userID int64) (Collection, error) {
	query, args, err := psql.
		Select("id", "name", "pinned", "created_at", "updated_at").
		From(collectionTableName).
		Where(squirrel.Eq{"id": collectionID, "user_id": userID}).
		ToSql()
	if err != nil {
		return Collection{}, err
	}

	var c Collection
	err = pgxscan.Get(ctx, r.db, &c, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Collection{}, ErrCollectionNotFound
		}
		return Collection{}, err
	}

	c.UserID = userID
	c.GameIDs, err = r.getCollectionGameIDs(ctx, c.ID)
	if err != nil {
		return Collection{}, err
	}

	return c, nil
}

func (r *repository) createCollection(ctx context.Context, userID int64, req Collection) (Collection, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Collection{}, err
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			panic(err)
		}
	}()

	query, args, err := psql.
		Insert(collectionTableName).
		Columns("user_id", "name", "pinned").
		Values(userID, req.Name, req.Pinned).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return Collection{}, err
	}

	var c = Collection{
		UserID: userID,
		Name:   req.Name,
		Pinned: req.Pinned,
	}
	err = tx.QueryRow(ctx, query, args...).Scan(&c.ID)
	if err != nil {
		return Collection{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Collection{}, err
	}

	return c, nil
}

func (r *repository) updateCollection(ctx context.Context, collectionID, userID int64, req Collection) (Collection, error) {
	query, args, err := psql.
		Update(collectionTableName).
		Set("name", req.Name).
		Set("pinned", req.Pinned).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": collectionID, "user_id": userID}).
		Suffix("RETURNING updated_at").
		ToSql()
	if err != nil {
		return Collection{}, err
	}

	var updatedAt time.Time
	err = r.db.QueryRow(ctx, query, args...).Scan(&updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Collection{}, ErrCollectionNotFound
		}
		return Collection{}, err
	}

	// Получаем обновленную коллекцию
	return r.getCollection(ctx, collectionID, userID)
}

func (r *repository) deleteCollection(ctx context.Context, collectionID, userID int64) error {
	query, args, err := psql.
		Delete(collectionTableName).
		Where(squirrel.Eq{"id": collectionID, "user_id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCollectionNotFound
	}

	return nil
}

func (r *repository) addGameToCollection(ctx context.Context, collectionID, gameID, userID int64) error {
	return pgx.BeginTxFunc(ctx, r.db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		// Проверяем, что коллекция принадлежит пользователю
		query, args, err := psql.
			Select("1").
			From(collectionTableName).
			Where(squirrel.Eq{"id": collectionID, "user_id": userID}).
			ToSql()
		if err != nil {
			return err
		}

		var exists int
		err = r.db.QueryRow(ctx, query, args...).Scan(&exists)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrForbidden
			}
			return err
		}

		// Добавляем игру в коллекцию
		query, args, err = psql.
			Insert(collectionTableName).
			Columns("collection_id", "game_id").
			Values(collectionID, gameID).
			Suffix("ON CONFLICT DO NOTHING").
			ToSql()
		if err != nil {
			return err
		}

		_, err = r.db.Exec(ctx, query, args...)
		return err
	})
}

func (r *repository) removeGameFromCollection(ctx context.Context, collectionID, gameID, userID int64) error {
	return pgx.BeginTxFunc(ctx, r.db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		// Проверяем, что коллекция принадлежит пользователю
		query, args, err := psql.
			Select("1").
			From(collectionTableName).
			Where(squirrel.Eq{"id": collectionID, "user_id": userID}).
			ToSql()
		if err != nil {
			return err
		}

		var exists int
		err = tx.QueryRow(ctx, query, args...).Scan(&exists)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrForbidden
			}
			return err
		}

		// Удаляем игру из коллекции
		query, args, err = psql.
			Delete(collectionTableName).
			Where(squirrel.Eq{"collection_id": collectionID, "game_id": gameID}).
			ToSql()
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return err
		}
		return tx.Commit(ctx)
	})
}
