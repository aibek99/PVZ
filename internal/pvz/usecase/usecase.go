package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"Homework-1/internal/cache"
	"Homework-1/internal/database"
	"Homework-1/internal/model/abstract"
	"Homework-1/internal/model/pvz"
	"Homework-1/pkg/errlst"
	"Homework-1/pkg/tracing"
)

const pvzTimeDuration = 24 * time.Hour

// PVZUseCase is
type PVZUseCase struct {
	repo  database.Datastore
	cache cache.Store
}

// NewPVZUseCase is
func NewPVZUseCase(repo database.Datastore, cache cache.Store) *PVZUseCase {
	return &PVZUseCase{repo: repo, cache: cache}
}

// CreatePVZ is
func (p *PVZUseCase) CreatePVZ(ctx context.Context, request pvz.Request) error {
	log.Println("[pvz][useCase][CreatePVZ]")
	tracer := otel.Tracer("[pvz_v1][useCase]")
	ctx, span := tracer.Start(ctx, "[CreatePVZ]")
	defer span.End()

	if err := p.repo.WithTransaction(ctx, func(db database.Datastore) error {
		err := db.PvzRepo().CheckPVZ(ctx, request.ToStorage())
		if err != nil {
			tracing.ErrorTracer(span, errlst.ErrPVZAlreadyExists)
			return errlst.ErrPVZAlreadyExists
		}

		return db.PvzRepo().CreatePVZ(ctx, request.ToStorage())
	}); err != nil {
		tracing.ErrorTracer(span, err)
		return err
	}

	span.SetStatus(codes.Ok, "PVZ created successfully")
	return nil
}

// GetPVZ is
func (p *PVZUseCase) GetPVZ(ctx context.Context, pvzID int64) (pvz.AllResponse, error) {
	log.Println("[pvz][useCase][GetPVZ]")
	tracer := otel.Tracer("[pvz_v1][useCase]")
	ctx, span := tracer.Start(ctx, "[GetPVZ]")
	defer span.End()

	cacheArgument := abstract.CacheArgument{
		ObjectType: "pvz",
		ObjectID:   pvzID,
	}

	var response pvz.AllResponse

	cachedValue, err := p.cache.Get(ctx, cacheArgument)
	if err == nil && json.Unmarshal(cachedValue, &response) == nil {
		span.SetStatus(codes.Ok, "Successfully got Box")
		return response, nil
	}

	var pvzData pvz.AllData

	pvzData, err = p.repo.PvzRepo().GetPVZ(ctx, pvzID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tracing.ErrorTracer(span, errlst.ErrPVZNotFound)
			return pvz.AllResponse{}, errlst.ErrPVZNotFound
		}
		tracing.ErrorTracer(span, err)
		return pvz.AllResponse{}, fmt.Errorf("p.repo.PvzRepo.GetPVZ: %w", err)
	}

	response = pvzData.ToPVZServer()

	marshaledData, err := json.Marshal(response)
	if err != nil {
		log.Printf("[pvz][usecase][GetPVZ] json.Marshal: %v", err)
		tracing.ErrorTracer(span, err)
		return pvzData.ToPVZServer(), err
	}

	err = p.cache.Set(ctx, cacheArgument, marshaledData, pvzTimeDuration)
	if err != nil {
		log.Printf("[pvz][usecase][GetPVZ] b.cache.Set: %v", err)
		tracing.ErrorTracer(span, err)
	}

	span.SetStatus(codes.Ok, "Successfully got Box")
	return response, nil
}

// DeletePVZByID is
func (p *PVZUseCase) DeletePVZByID(ctx context.Context, pvzID int64) error {
	log.Println("[pvz][useCase][DeletePVZ]")
	tracer := otel.Tracer("[pvz_v1][useCase]")
	ctx, span := tracer.Start(ctx, "[DeletePVZByID]")
	defer span.End()

	if err := p.repo.WithTransaction(ctx, func(db database.Datastore) error {
		return db.PvzRepo().DeletePVZByID(ctx, pvzID)
	}); err != nil {
		tracing.ErrorTracer(span, err)
		return err
	}

	span.SetStatus(codes.Ok, "Successfully deleted PVZ by ID")
	return nil
}

// UpdatePVZ is
func (p *PVZUseCase) UpdatePVZ(ctx context.Context, updatePVZRequest pvz.UpdateRequest) error {
	log.Println("[pvz][useCase][UpdatePVZ]")
	tracer := otel.Tracer("[pvz_v1][useCase]")
	ctx, span := tracer.Start(ctx, "[UpdatePVZ]")
	defer span.End()

	if err := p.repo.WithTransaction(ctx, func(db database.Datastore) error {
		return db.PvzRepo().UpdatePVZ(ctx, updatePVZRequest.ToStorage())
	}); err != nil {
		tracing.ErrorTracer(span, err)
		return err
	}

	span.SetStatus(codes.Ok, "Successfully updated PVZ")
	return nil
}

// ListPVZ is
func (p *PVZUseCase) ListPVZ(ctx context.Context, pvzPagination abstract.Page) (abstract.PaginatedResponse[pvz.AllResponse], error) {
	log.Println("[pvz][useCase][ListPVZ]")
	tracer := otel.Tracer("[pvz_v1][useCase]")
	ctx, span := tracer.Start(ctx, "[ListPVZ]")
	defer span.End()

	var pvzAllData []pvz.AllData
	var err error
	var pvzListResponse abstract.PaginatedResponse[pvz.AllResponse]

	if err = p.repo.WithTransaction(ctx, func(db database.Datastore) error {
		var count int64
		count, err = db.PvzRepo().CountOfPVZ(ctx)
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		pvzListResponse.TotalItems = count
		pvzAllData, err = db.PvzRepo().ListPVZ(ctx, pvzPagination.ToStorage())
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		return nil
	}); err != nil {
		tracing.ErrorTracer(span, err)
		return abstract.PaginatedResponse[pvz.AllResponse]{}, err
	}

	pvzList := lo.Map(
		pvzAllData,
		func(item pvz.AllData, _ int) pvz.AllResponse {
			return item.ToPVZServer()
		},
	)

	pvzListResponse.Items = pvzList
	pvzListResponse.CurrentPage = pvzPagination.CurrentPage
	pvzListResponse.ItemsPerPage = int64(len(pvzList))

	span.SetStatus(codes.Ok, "Successfully got list of PVZ")
	return pvzListResponse, nil
}
