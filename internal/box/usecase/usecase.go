package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"Homework-1/internal/cache"
	"Homework-1/internal/database"
	"Homework-1/internal/model/abstract"
	"Homework-1/internal/model/box"
	"Homework-1/pkg/constants"
	"Homework-1/pkg/errlst"
	"Homework-1/pkg/tracing"
)

// BoxUseCase is
type BoxUseCase struct {
	repo  database.Datastore
	cache cache.Store
}

// NewBoxUseCase is
func NewBoxUseCase(repo database.Datastore, cache cache.Store) *BoxUseCase {
	return &BoxUseCase{repo: repo, cache: cache}
}

// CreateBox is
func (b *BoxUseCase) CreateBox(ctx context.Context, request box.Request) (int64, error) {
	log.Println("[box][useCase][CreateBox]")
	tracer := otel.Tracer("[box_v1][useCase]")
	ctx, span := tracer.Start(ctx, "[CreateBox]")
	defer span.End()

	id, err := b.repo.BoxRepo().CreateBox(ctx, request.ToStorage())
	if err != nil {
		tracing.ErrorTracer(span, err)
		return -1, err
	}

	span.SetStatus(codes.Ok, "Box created successfully")
	return id, nil
}

// DeleteBoxByID is
func (b *BoxUseCase) DeleteBoxByID(ctx context.Context, boxID int64) error {
	log.Println("[box][useCase][DeleteBoxByID]")
	tracer := otel.Tracer("[box_v1][useCase]")
	ctx, span := tracer.Start(ctx, "[DeleteBoxByID]")
	defer span.End()

	err := b.repo.BoxRepo().DeleteBoxByID(ctx, boxID)
	if err != nil {
		tracing.ErrorTracer(span, err)
		return err
	}

	span.SetStatus(codes.Ok, "Successfully deleted box by ID")
	return nil
}

// ListBoxes is
func (b *BoxUseCase) ListBoxes(ctx context.Context, boxPage abstract.Page) (abstract.PaginatedResponse[box.AllResponse], error) {
	log.Println("[box][useCase][ListBoxes]")
	tracer := otel.Tracer("[box_v1][useCase]")
	ctx, span := tracer.Start(ctx, "[ListBoxes]")
	defer span.End()

	var boxAllData []box.AllData
	var err error
	var boxListResponse abstract.PaginatedResponse[box.AllResponse]

	err = b.repo.WithTransaction(ctx, func(db database.Datastore) error {
		var count int64
		count, err = db.BoxRepo().CountBoxes(ctx)
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		boxListResponse.TotalItems = count

		boxAllData, err = db.BoxRepo().ListBoxes(ctx, boxPage.ToStorage())
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		return nil
	})
	if err != nil {
		tracing.ErrorTracer(span, err)
		return abstract.PaginatedResponse[box.AllResponse]{}, err
	}

	boxList := lo.Map(
		boxAllData,
		func(item box.AllData, _ int) box.AllResponse {
			return item.ToServer()
		},
	)

	boxListResponse.Items = boxList
	boxListResponse.CurrentPage = boxPage.CurrentPage
	boxListResponse.ItemsPerPage = int64(len(boxList))

	span.SetStatus(codes.Ok, "Successfully got list of box")
	return boxListResponse, nil
}

// GetBox is
func (b *BoxUseCase) GetBox(ctx context.Context, boxID int64) (box.AllResponse, error) {
	log.Println("[box][useCase][GetBox]")
	tracer := otel.Tracer("[box_v1][useCase]")
	ctx, span := tracer.Start(ctx, "[GetBox]")
	defer span.End()

	cacheArgument := abstract.CacheArgument{
		ObjectType: "box",
		ObjectID:   boxID,
	}

	var response box.AllResponse

	cachedValue, err := b.cache.Get(ctx, cacheArgument)
	if err == nil && json.Unmarshal(cachedValue, &response) == nil {
		span.SetStatus(codes.Ok, "Successfully got Box")
		return response, err
	}

	boxData, err := b.repo.BoxRepo().GetBox(ctx, boxID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tracing.ErrorTracer(span, errlst.ErrBoxNotFound)
			return box.AllResponse{}, errlst.ErrBoxNotFound
		}

		tracing.ErrorTracer(span, err)
		return box.AllResponse{}, err
	}

	response = boxData.ToServer()

	marshaledData, err := json.Marshal(response)
	if err != nil {
		log.Printf("[box][usecase][GetBox] json.Marshal: %v", err)
		tracing.ErrorTracer(span, err)
		return boxData.ToServer(), nil
	}

	err = b.cache.Set(ctx, cacheArgument, marshaledData, constants.BoxTimeDuration)
	if err != nil {
		log.Printf("[box][usecase][GetBox] b.cache.Set: %v", err)
		tracing.ErrorTracer(span, err)
	}

	span.SetStatus(codes.Ok, "Successfully got Box")
	return response, nil
}
