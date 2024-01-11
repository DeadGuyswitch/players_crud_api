package players

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreSQLPGXRepository struct {
	db *pgxpool.Pool
}

func NewPostgreSQLPGXRepository(db *pgxpool.Pool) *PostgreSQLPGXRepository {
	return &PostgreSQLPGXRepository{
		db: db,
	}
}

func (r *PostgreSQLPGXRepository) Migrate(ctx context.Context) error {
	query := `	
CREATE TABLE IF NOT EXISTS players(
	 	id SERIAL PRIMARY KEY,
	 	kitNumber INTEGER NOT NULL,
	 	firstName TEXT NOT NULL,
	 	lastName TEXT NOT NULL,
	 	nationality TEXT NOT NULL,
	 	pos TEXT NOT NULL,
	)
`
	_, err := r.db.Exec(ctx, query)
	return err
}

func (r *PostgreSQLPGXRepository) Create(ctx context.Context, player Player) (*Player, error) {
	var id int64
	err := r.db.QueryRow(ctx,
		"INSERT INTO players(kitNumber, firstName, lastName, nationality, pos) values ($1, $2, $3, $4, $5) RETURNING id",
		player.KitNumber, player.FirstName, player.LastName, player.Nationality, player.Position).Scan(&id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	player.ID = id

	return &player, nil
}

func (r *PostgreSQLPGXRepository) List(ctx context.Context) ([]Player, error) {
	rows, err := r.db.Query(ctx, "SELECT * from players")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playersList []Player
	for rows.Next() {
		var player Player
		if err := rows.Scan(&player.ID, &player.FirstName, &player.LastName, &player.Nationality, &player.Position); err != nil {
			return nil, err
		}

		playersList = append(playersList, player)
	}
	return playersList, nil
}

func (r *PostgreSQLPGXRepository) Get(ctx context.Context, firstName string, lastName string) (*Player, error) {
	row := r.db.QueryRow(ctx, "SELECT * FROM players WHERE firstName = $1 AND lastName = $2", firstName, lastName)

	var player Player
	if err := row.Scan(&player.ID, &player.FirstName, &player.LastName, &player.Nationality, &player.Position); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExist
		}
		return nil, err
	}
	return &player, nil
}

func (r *PostgreSQLPGXRepository) Update(ctx context.Context, id int64, updated Player) (*Player, error) {
	res, err := r.db.Exec(ctx, "UPDATE players SET kitNumber = $1, firstName = $2, lastName = $3, nationality = $4, pos = $5",
		updated.KitNumber, updated.FirstName, updated.LastName, updated.Nationality, updated.Position)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}
	return &updated, nil
}

func (r *PostgreSQLPGXRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.Exec(ctx, "DELETE FROM players WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}
	return err
}
