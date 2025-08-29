package repository

import (
	"context"
	"fmt"

	"github.com/aliffatulmf/mkp-eticket-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TerminalRepository interface {
	List(ctx context.Context) ([]model.Terminal, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Terminal, error)
	Create(ctx context.Context, terminal *model.Terminal) error
	Update(ctx context.Context, terminal *model.Terminal) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type terminalRepository struct {
	db *pgxpool.Pool
}

func NewTerminalRepository(db *pgxpool.Pool) TerminalRepository {
	return &terminalRepository{db: db}
}

func (r *terminalRepository) List(ctx context.Context) ([]model.Terminal, error) {
	query := `SELECT id, code, name, address, is_active, created_at, updated_at FROM terminals ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query terminals: %w", err)
	}
	defer rows.Close()

	terms, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Terminal])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return terms, nil
}

func (r *terminalRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Terminal, error) {
	query := `SELECT id, code, name, address, is_active, created_at, updated_at FROM terminals WHERE id = $1`

	var terminal model.Terminal
	err := r.db.QueryRow(ctx, query, id).Scan(
		&terminal.ID, &terminal.Code, &terminal.Name, &terminal.Address,
		&terminal.IsActive, &terminal.CreatedAt, &terminal.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get terminal: %w", err)
	}

	return &terminal, nil
}

func (r *terminalRepository) Create(ctx context.Context, terminal *model.Terminal) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO terminals (id, code, name, address, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = tx.Exec(ctx, query,
		terminal.ID, terminal.Code, terminal.Name, terminal.Address,
		terminal.IsActive, terminal.CreatedAt, terminal.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create terminal: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *terminalRepository) Update(ctx context.Context, terminal *model.Terminal) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `UPDATE terminals SET name = $2, address = $3, is_active = $4, updated_at = $5 WHERE id = $1`

	_, err = tx.Exec(ctx, query,
		terminal.ID, terminal.Name, terminal.Address, terminal.IsActive, terminal.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update terminal: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *terminalRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `DELETE FROM terminals WHERE id = $1`

	_, err = tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete terminal: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
