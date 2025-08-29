package repository

import (
	"context"
	"fmt"

	"github.com/aliffatulmf/mkp-eticket-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.Admin, error)
	FindByID(ctx context.Context, id int) (*model.Admin, error)
	Create(ctx context.Context, admin *model.Admin) error
	Delete(ctx context.Context, id int) error
}

type adminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) FindByUsername(ctx context.Context, username string) (*model.Admin, error) {
	query := `SELECT id, username, password FROM admins WHERE username = $1`

	var admin model.Admin
	err := r.db.QueryRow(ctx, query, username).Scan(&admin.ID, &admin.Username, &admin.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin by username: %w", err)
	}

	return &admin, nil
}

func (r *adminRepository) FindByID(ctx context.Context, id int) (*model.Admin, error) {
	query := `SELECT id, username, password FROM admins WHERE id = $1`

	var admin model.Admin
	err := r.db.QueryRow(ctx, query, id).Scan(&admin.ID, &admin.Username, &admin.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin by id: %w", err)
	}

	return &admin, nil
}

func (r *adminRepository) Create(ctx context.Context, admin *model.Admin) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO admins (username, password) VALUES ($1, $2) RETURNING id`

	err = tx.QueryRow(ctx, query, admin.Username, admin.Password).Scan(&admin.ID)
	if err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *adminRepository) Delete(ctx context.Context, id int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `DELETE FROM admins WHERE id = $1`

	result, err := tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete admin: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("admin with id %d not found", id)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
