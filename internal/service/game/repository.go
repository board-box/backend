package game

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	pgx "github.com/jackc/pgx/v5"
)

const gameTableName = "game"

var (
	psql            = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	ErrGameNotFound = errors.New("game not found")
	ErrEmptyIDs     = errors.New("empty IDs")
)

type repository struct {
	db *pgx.Conn
}

func newRepository(db *pgx.Conn) *repository {
	return &repository{db: db}
}

func (r *repository) listGames(ctx context.Context) ([]Game, error) {
	query, args, err := psql.
		Select("id", "title", "description", "genre", "age", "person", "avg_time", "difficulty", "image", "rules").
		From(gameTableName).
		OrderBy("title ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var game Game
		err = pgxscan.ScanRow(&game, rows)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

func (r *repository) getGameById(ctx context.Context, id int64) (Game, error) {
	query, args, err := psql.
		Select("id", "title", "description", "genre", "age", "person", "avg_time", "difficulty", "image", "rules").
		From(gameTableName).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return Game{}, err
	}

	var game Game
	err = pgxscan.Get(ctx, r.db, &game, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Game{}, ErrGameNotFound
		}
		return Game{}, err
	}

	return game, nil
}

func (r *repository) getGamesByIds(ctx context.Context, ids []int64) ([]Game, error) {
	if len(ids) == 0 {
		return nil, ErrEmptyIDs
	}

	query, args, err := psql.
		Select("id", "title", "description", "genre", "age", "person", "avg_time", "difficulty", "image", "rules").
		From(gameTableName).
		Where(squirrel.Eq{"id": ids}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var game Game
		err = pgxscan.ScanRow(&game, rows)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return nil, ErrGameNotFound
	}

	return games, nil
}

func (r *repository) createGame(ctx context.Context, game Game) (int64, error) {
	query, args, err := psql.
		Insert(gameTableName).
		Columns("title", "description", "genre", "age", "person", "avg_time", "difficulty", "image", "rules").
		Values(game.Title, game.Description, game.Genre, game.Age, game.Person, game.AvgTime, game.Difficulty, game.Image, game.Rules).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) updateGame(ctx context.Context, game Game) error {
	query, args, err := psql.
		Update(gameTableName).
		Set("title", game.Title).
		Set("description", game.Description).
		Set("genre", game.Genre).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": game.ID}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return ErrGameNotFound
	}

	return nil
}

func (r *repository) deleteGame(ctx context.Context, id int64) error {
	query, args, err := psql.
		Delete(gameTableName).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return ErrGameNotFound
	}

	return nil
}
