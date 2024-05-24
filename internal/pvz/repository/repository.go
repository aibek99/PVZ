package repository

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"Homework-1/internal/connection"
	"Homework-1/internal/model/abstract"
	"Homework-1/internal/model/pvz"
	"Homework-1/pkg/errlst"
)

// PVZRepository is
type PVZRepository struct {
	psqlDB connection.DB
}

// NewPVZPGRepository is
func NewPVZPGRepository(psqlDB connection.DB) *PVZRepository {
	return &PVZRepository{
		psqlDB: psqlDB,
	}
}

// CreatePVZ is
func (p *PVZRepository) CreatePVZ(ctx context.Context, pvzData pvz.Data) error {
	log.Println("[pvz][repository][CreatePVZ]")

	_, err := p.psqlDB.Execute(ctx,
		"INSERT INTO pvz(name, address, contact) VALUES ($1,$2,$3);",
		pvzData.Name,
		pvzData.Address,
		pvzData.Contact,
	)
	if err != nil {
		return fmt.Errorf("p.psqlDB.ExecContext: %w", err)
	}

	return nil
}

// CheckPVZ is
func (p *PVZRepository) CheckPVZ(ctx context.Context, pvzData pvz.Data) error {
	log.Println("[pvz][repository][GetPVZ]")
	var exists bool

	err := p.psqlDB.QueryRow(ctx,
		"SELECT EXISTS (SELECT 1 FROM pvz WHERE name = $1 AND address = $2 AND contact = $3 AND deleted_at IS NULL)",
		pvzData.Name,
		pvzData.Address,
		pvzData.Contact,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("p.psqlDB.QueryRowContext: %w", err)
	}

	if exists {
		return errlst.ErrPVZAlreadyExists
	}

	return nil
}

// GetPVZ is
func (p *PVZRepository) GetPVZ(ctx context.Context, pvzID int64) (pvz.AllData, error) {
	log.Println("[pvz][repository][GetPVZ]")

	var pvzAllData pvz.AllData

	err := p.psqlDB.Get(
		ctx,
		&pvzAllData,
		"Select id, name, address, contact, created_at, updated_at FROM pvz Where  id=$1 AND deleted_at IS NULL",
		pvzID,
	)
	if err != nil {
		return pvz.AllData{}, fmt.Errorf("p.psqlDB.Get: %w", err)
	}

	return pvzAllData, nil
}

// ListPVZ is
func (p *PVZRepository) ListPVZ(
	ctx context.Context,
	pvzPaginationData abstract.PageData,
) ([]pvz.AllData, error) {
	log.Println("[pvz][repository][ListPVZ]")
	offset := (pvzPaginationData.CurrentPage - 1) * pvzPaginationData.ItemsPerPage
	var pvzAllData []pvz.AllData

	err := p.psqlDB.Select(
		ctx,
		&pvzAllData,
		"SELECT id, name, address, contact, created_at, updated_at FROM pvz WHERE deleted_at IS NULL "+
			"ORDER BY created_at DESC OFFSET $1 LIMIT $2",
		offset,
		pvzPaginationData.ItemsPerPage,
	)
	if err != nil {
		return []pvz.AllData{}, fmt.Errorf("p.psqlDB.Select: %w", err)
	}

	return pvzAllData, nil
}

// CountOfPVZ is
func (p *PVZRepository) CountOfPVZ(ctx context.Context) (int64, error) {
	log.Println("[pvz][repository][CountOfPVZ]")
	var totalCount int64

	err := p.psqlDB.Get(
		ctx,
		&totalCount,
		"SELECT COUNT(*) FROM pvz WHERE deleted_at IS NULL",
	)
	if err != nil {
		return 0, fmt.Errorf("b.psqlDB.Get: %w", err)
	}

	return totalCount, nil
}

// UpdatePVZ is
func (p *PVZRepository) UpdatePVZ(ctx context.Context, updatePVZData pvz.UpdateData) error {
	log.Println("[pvz][repository][UpdatePVZ]")

	query := "UPDATE pvz SET"
	setValues := make([]interface{}, 0)

	num := 1

	if updatePVZData.Name != nil {
		query += " name = $" + strconv.Itoa(num) + ","

		setValues = append(setValues, *updatePVZData.Name)
		num++
	}

	if updatePVZData.Address != nil {
		query += " address = $" + strconv.Itoa(num) + ","

		setValues = append(setValues, *updatePVZData.Address)
		num++
	}

	if updatePVZData.Contact != nil {
		query += " contact = $" + strconv.Itoa(num) + ","

		setValues = append(setValues, *updatePVZData.Contact)
		num++
	}
	query += " updated_at = NOW(),"
	query = strings.TrimSuffix(query, ",")
	query += " WHERE id = $" + strconv.Itoa(num)

	setValues = append(setValues, updatePVZData.ID)

	result, err := p.psqlDB.Execute(ctx, query, setValues...)
	if err != nil {
		return fmt.Errorf("p.psqlDB.ExecContext: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if rows == 0 {
		return errlst.ErrPVZNotFound
	}

	return nil
}

// DeletePVZByID is
func (p *PVZRepository) DeletePVZByID(ctx context.Context, pvzID int64) error {
	log.Println("[pvz][repository][DeletePVZ]")

	result, err := p.psqlDB.Execute(
		ctx,
		"UPDATE pvz SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL",
		time.Now(),
		pvzID,
	)
	if err != nil {
		return fmt.Errorf("p.psqlDB.ExecContext: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if rows == 0 {
		return errlst.ErrPVZNotFound
	}

	return nil
}
