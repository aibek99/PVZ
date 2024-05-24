package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"Homework-1/internal/connection"
	"Homework-1/internal/model/abstract"
	"Homework-1/internal/model/box"
	"Homework-1/pkg/errlst"
)

// BoxRepository is
type BoxRepository struct {
	psqlDB connection.DB
}

// NewBoxPGRepository is
func NewBoxPGRepository(psqlDB connection.DB) *BoxRepository {
	return &BoxRepository{
		psqlDB: psqlDB,
	}
}

// BoxData is
func (b *BoxRepository) BoxData(ctx context.Context, id int64) (box.AllResponse, error) {
	log.Println("[box_v1][repository][BoxData]")
	var boxData box.AllData
	err := b.psqlDB.Get(
		ctx,
		&boxData,
		"SELECT id, name, cost, is_check, weight, created_at, updated_at FROM box WHERE deleted_at IS NULL AND id = $1",
		id,
	)

	if err != nil {
		return box.AllResponse{}, err
	}

	return boxData.ToServer(), nil
}

// CreateBox is
func (b *BoxRepository) CreateBox(ctx context.Context, box box.Data) (int64, error) {
	log.Println("[box_v1][repository][CreateBox]")

	var id int64
	row := b.psqlDB.QueryRow(
		ctx,
		"INSERT INTO box(name, cost, is_check, weight) VALUES ($1,$2,$3,$4) RETURNING ID;",
		box.Name,
		box.Cost,
		box.IsCheck,
		box.Weight,
	)

	err := row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("p.psqlDB.QueryRow: %w", err)
	}

	return id, nil
}

// DeleteBoxByID is
func (b *BoxRepository) DeleteBoxByID(ctx context.Context, id int64) error {
	log.Println("[box_v1][repository][DeleteBoxByID]")

	result, err := b.psqlDB.Execute(
		ctx,
		"UPDATE box SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL",
		time.Now(),
		id,
	)
	if err != nil {
		return fmt.Errorf("p.psqlDB.ExecContext: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if rows == 0 {
		return errlst.ErrBoxNotFound
	}

	return nil
}

// CountBoxes is
func (b *BoxRepository) CountBoxes(ctx context.Context) (int64, error) {
	log.Println("[box_v1][repository][CountBoxes]")
	var totalCount int64

	err := b.psqlDB.Get(
		ctx,
		&totalCount,
		"SELECT COUNT(*) FROM box WHERE deleted_at IS NULL",
	)
	if err != nil {
		return 0, fmt.Errorf("b.psqlDB.Get: %w", err)
	}

	return totalCount, nil
}

// ListBoxes is
func (b *BoxRepository) ListBoxes(ctx context.Context, boxPagination abstract.PageData) ([]box.AllData, error) {
	log.Println("[box_v1][repository][ListBoxes]")
	offset := (boxPagination.CurrentPage - 1) * boxPagination.ItemsPerPage
	var boxAllData []box.AllData

	err := b.psqlDB.Select(
		ctx,
		&boxAllData,
		"SELECT id, name, cost, is_check, weight, created_at, updated_at FROM box WHERE deleted_at IS NULL "+
			"ORDER BY created_at DESC OFFSET $1 LIMIT $2",
		offset,
		boxPagination.ItemsPerPage,
	)
	if err != nil {
		return []box.AllData{}, fmt.Errorf("p.psqlDB.Select: %w", err)
	}

	return boxAllData, nil
}

// GetBox is
func (b *BoxRepository) GetBox(ctx context.Context, boxID int64) (box.AllData, error) {
	log.Println("[box_v1][repository][GetBox]")

	var boxAllData box.AllData

	err := b.psqlDB.Get(
		ctx,
		&boxAllData,
		"Select id, name, cost, is_check, weight, created_at, updated_at FROM box Where id=$1 AND deleted_at IS NULL",
		boxID,
	)
	if err != nil {
		return box.AllData{}, fmt.Errorf("p.psqlDB.Get: %w", err)
	}

	return boxAllData, nil
}
