package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"Homework-1/internal/cache"
	"Homework-1/internal/database"
	"Homework-1/internal/model/abstract"
	"Homework-1/internal/model/box"
	"Homework-1/internal/model/order"
	"Homework-1/pkg/constants"
	"Homework-1/pkg/tracing"
)

var (
	totalIssuedOrders = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "total_issued_orders",
			Help: "Total number of issued orders",
		},
	)
)

// OrderUseCase is
type OrderUseCase struct {
	repo  database.Datastore
	cache cache.Store
}

// NewOrderUseCase is
func NewOrderUseCase(repo database.Datastore, cache cache.Store) *OrderUseCase {
	return &OrderUseCase{repo: repo, cache: cache}
}

// CreateReceiveOrder is
func (o *OrderUseCase) CreateReceiveOrder(ctx context.Context, request order.Request) error {
	log.Println("[order][useCase][CreateReceiveOrder]")
	tracer := otel.Tracer("[order][useCase]")
	ctx, span := tracer.Start(ctx, "[CreateReceiveOrder]")
	defer span.End()

	cacheArgument := abstract.CacheArgument{
		ObjectType: "box",
		ObjectID:   request.BoxID,
	}

	var response box.AllResponse

	boxValue, redErr := o.cache.Get(ctx, cacheArgument)
	if redErr == nil && json.Unmarshal(boxValue, &response) == nil && !((response.IsCheck && response.Weight > request.Weight) || !response.IsCheck) {
		err := o.repo.OrderRepo().CreateReceiveOrder(ctx, request.ToStorage())
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		span.SetStatus(codes.Ok, "Successfully created a order")
		return nil
	}

	if err := o.repo.WithTransaction(ctx, func(db database.Datastore) error {
		boxData, err := db.BoxRepo().GetBox(ctx, request.BoxID)
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}
		response = boxData.ToServer()

		if !((response.IsCheck && response.Weight > request.Weight) || !response.IsCheck) {
			tracing.ErrorTracer(span, errors.New("invalid weight request"))
			return err
		}
		err = db.OrderRepo().CreateReceiveOrder(ctx, request.ToStorage())
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		return nil
	}); err != nil {
		tracing.ErrorTracer(span, err)
		return err
	}

	marshaledData, err := json.Marshal(response)
	if err != nil {
		tracing.ErrorTracer(span, err)
		log.Printf("[order][usecase][CreateReceiveOrder] json.Marshal: %v", err)
		return nil
	}

	if err = o.cache.Set(ctx, cacheArgument, marshaledData, constants.BoxTimeDuration); err != nil {
		tracing.ErrorTracer(span, err)
		log.Printf("[order][usecase][CreateReceiveOrder] o.cache.Set: %v", err)
	}

	span.SetStatus(codes.Ok, "Successfully created a order")
	return nil
}

// IssueOrders is
func (o *OrderUseCase) IssueOrders(ctx context.Context, request order.RequestOrderIDs) (float64, error) {
	log.Println("[order][useCase][IssueOrders]")
	tracer := otel.Tracer("[order][useCase]")
	ctx, span := tracer.Start(ctx, "[IssueOrders]")
	defer span.End()

	uniqueOrderIDs := make(map[int64]bool)

	for _, orderID := range request.OrderIDs {
		uniqueOrderIDs[orderID] = true
	}

	var totalCost float64
	if err := o.repo.WithTransaction(ctx, func(db database.Datastore) error {
		clientID, err := db.OrderRepo().GetClientID(ctx, request.OrderIDs[0])
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		for orderID := range uniqueOrderIDs {
			err = db.OrderRepo().UpdateIssueOrder(ctx, orderID, clientID)
			if err != nil {
				tracing.ErrorTracer(span, err)
				return err
			}

			boxDataWithOrderWeight, err := db.OrderRepo().BoxDataByOrderID(ctx, orderID)
			if err != nil {
				tracing.ErrorTracer(span, err)
				return err
			}

			if boxDataWithOrderWeight.IsCheck && boxDataWithOrderWeight.Weight > boxDataWithOrderWeight.OrderWeight {
				totalCost += boxDataWithOrderWeight.Cost
			}

		}

		return nil
	}); err != nil {
		tracing.ErrorTracer(span, err)
		return 0, err
	}

	totalIssuedOrders.Add(float64(len(uniqueOrderIDs)))

	span.SetStatus(codes.Ok, "Successfully issued Orders")
	return totalCost, nil
}

// ReturnedOrders is
func (o *OrderUseCase) ReturnedOrders(
	ctx context.Context,
	request abstract.Page,
) (abstract.PaginatedResponse[order.ReturnedResponse], error) {
	log.Println("[order][useCase][ReturnedOrders]")
	tracer := otel.Tracer("[order][useCase]")
	ctx, span := tracer.Start(ctx, "[ReturnedOrders]")
	defer span.End()

	var returnedOrdersData []order.ReturnedData
	var err error

	var returnedListResponse abstract.PaginatedResponse[order.ReturnedResponse]

	err = o.repo.WithTransaction(ctx, func(db database.Datastore) error {
		var count int64
		count, err = db.OrderRepo().CountReturnedOrders(ctx)
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		returnedListResponse.TotalItems = count

		returnedOrdersData, err = db.OrderRepo().ListReturnedOrders(ctx, request.ToStorage())
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		return nil
	})
	if err != nil {
		tracing.ErrorTracer(span, err)
		return abstract.PaginatedResponse[order.ReturnedResponse]{}, err
	}

	returnedOrdersList := lo.Map(
		returnedOrdersData,
		func(item order.ReturnedData, _ int) order.ReturnedResponse {
			return item.ToServer()
		},
	)

	returnedListResponse.Items = returnedOrdersList
	returnedListResponse.CurrentPage = request.CurrentPage
	returnedListResponse.ItemsPerPage = int64(len(returnedOrdersList))

	span.SetStatus(codes.Ok, "Successfully got list of returned orders")
	return returnedListResponse, nil
}

// UpdateAcceptOrder is
func (o *OrderUseCase) UpdateAcceptOrder(ctx context.Context, request order.RequestWithClientID) error {
	log.Println("[order][useCase][ReturnedOrders]")
	tracer := otel.Tracer("[order][useCase]")
	ctx, span := tracer.Start(ctx, "[UpdateAcceptOrder]")
	defer span.End()

	if err := o.repo.WithTransaction(ctx, func(db database.Datastore) error {
		err := db.OrderRepo().UpdateAcceptOrder(ctx, request.ToStorage())
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		return nil
	}); err != nil {
		tracing.ErrorTracer(span, err)
		return err
	}
	span.SetStatus(codes.Ok, "Successfully accept Order")
	return nil
}

// DeleteReturnedOrder is
func (o *OrderUseCase) DeleteReturnedOrder(ctx context.Context, orderID int64) error {
	log.Println("[order][useCase][DeleteReturnedOrder]")
	tracer := otel.Tracer("[order][useCase]")
	ctx, span := tracer.Start(ctx, "[DeleteReturnedOrder]")
	defer span.End()

	if err := o.repo.WithTransaction(ctx, func(db database.Datastore) error {
		err := db.OrderRepo().DeleteReturnOrder(ctx, orderID)
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		return nil
	}); err != nil {
		tracing.ErrorTracer(span, err)
		return err
	}
	span.SetStatus(codes.Ok, "Successfully deleted Order")
	return nil
}

// OrderList is
func (o *OrderUseCase) OrderList(ctx context.Context, request abstract.Page) (abstract.PaginatedResponse[order.AllResponse], error) {
	log.Println("[order][useCase][OrderList]")
	tracer := otel.Tracer("[order][useCase]")
	ctx, span := tracer.Start(ctx, "[OrderList]")
	defer span.End()

	var listOrderData []order.AllResponseData
	var err error
	var allOrdersList abstract.PaginatedResponse[order.AllResponse]

	err = o.repo.WithTransaction(ctx, func(db database.Datastore) error {
		var count int64
		count, err = db.OrderRepo().CountOrders(ctx)
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		allOrdersList.TotalItems = count

		listOrderData, err = db.OrderRepo().ListOrders(ctx, request.ToStorage())
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		return nil
	})
	if err != nil {
		tracing.ErrorTracer(span, err)
		return abstract.PaginatedResponse[order.AllResponse]{}, err
	}

	orderList := lo.Map(
		listOrderData,
		func(item order.AllResponseData, _ int) order.AllResponse {
			return item.ToServer()
		},
	)

	allOrdersList.Items = orderList
	allOrdersList.CurrentPage = request.CurrentPage
	allOrdersList.ItemsPerPage = int64(len(orderList))

	span.SetStatus(codes.Ok, "Successfully got order list")
	return allOrdersList, nil
}

// UniqueClientsList is
func (o *OrderUseCase) UniqueClientsList(
	ctx context.Context,
	request abstract.Page,
) (abstract.PaginatedResponse[order.ListUniqueClients], error) {
	log.Println("[order][useCase][UniqueClientsList]")
	tracer := otel.Tracer("[order][useCase]")
	ctx, span := tracer.Start(ctx, "[UniqueClientsList]")
	defer span.End()

	var uniqueClientListData []order.ListUniqueClientsData
	var err error

	var uniqueClientListResponse abstract.PaginatedResponse[order.ListUniqueClients]

	err = o.repo.WithTransaction(ctx, func(db database.Datastore) error {
		var count int64
		count, err = db.OrderRepo().CountUniqueClients(ctx)
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		uniqueClientListResponse.TotalItems = count

		uniqueClientListData, err = db.OrderRepo().ListUniqueClients(ctx, request.ToStorage())
		if err != nil {
			tracing.ErrorTracer(span, err)
			return err
		}

		return nil
	})
	if err != nil {
		tracing.ErrorTracer(span, err)
		return abstract.PaginatedResponse[order.ListUniqueClients]{}, err
	}

	uniqueClientList := lo.Map(
		uniqueClientListData,
		func(item order.ListUniqueClientsData, _ int) order.ListUniqueClients {
			return item.ToServer()
		},
	)

	uniqueClientListResponse.Items = uniqueClientList
	uniqueClientListResponse.CurrentPage = request.CurrentPage
	uniqueClientListResponse.ItemsPerPage = int64(len(uniqueClientList))

	span.SetStatus(codes.Ok, "Successfully got unique client list")

	return uniqueClientListResponse, nil
}
