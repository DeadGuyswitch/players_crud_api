package players

import (
	"context"
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

//func (r *PostgreSQLPGXRepository) List(ctx context.Context) ([]Player, error) {
//	rows, err := r.db.Query("SELECT * from players")
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var playersList []Player
//}
