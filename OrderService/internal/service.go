package internal

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"tesodev-korpes/OrderService/config"
	"tesodev-korpes/OrderService/internal/types"
	"tesodev-korpes/pkg/client"
	"tesodev-korpes/pkg/customError"
)

type Service struct {
	repo   *Repository
	client *client.Client
}

func NewService(repo *Repository, client *client.Client) *Service {
	return &Service{
		repo:   repo,
		client: client,
	}
}

func (s *Service) Create(ctx context.Context, order *types.Order, token string) (string, error) {
	if order.CustomerId == "" {
		return "", errors.New("customerId not found")
	}

	customer, err := s.fetchCustomerByID(order.CustomerId, token)
	if err != nil || customer == nil {
		return "", err
	}

	id, err := s.repo.Create(ctx, order)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Service) GetByID(ctx context.Context, id string, token string) (*types.OrderWithCustomerResponse, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	customer, err := s.fetchCustomerByID(order.CustomerId, token)
	if err != nil {
		return nil, err
	}

	return ToOrderWithCustomerResponse(ToOrderResponse(order), customer), nil
}

func (s *Service) GetAllOrders(ctx context.Context, pagination types.Pagination) ([]*types.OrderResponseModel, error) {
	skip := (pagination.Page - 1) * pagination.Limit

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pagination.Limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	orders, err := s.repo.GetAllOrders(ctx, findOptions)
	if err != nil {
		return nil, err
	}

	var response []*types.OrderResponseModel
	for _, order := range orders {
		resp := ToOrderResponse(&order)
		if resp != nil {
			response = append(response, resp)
		}
	}

	return response, nil
}

func (s *Service) ShipOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if order.Status != config.OrderStatus.Ordered {
		return customError.NewConflict(customError.OrderStatusConflict, order.Status, config.OrderStatus.Ordered)
	}
	return s.repo.UpdateStatusByID(ctx, id, config.OrderStatus.Shipped)
}

func (s *Service) DeliverOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if order.Status != config.OrderStatus.Shipped {
		return customError.NewConflict(customError.OrderStatusConflict, order.Status, config.OrderStatus.Shipped)
	}
	return s.repo.UpdateStatusByID(ctx, id, config.OrderStatus.Delivered)
}

func (s *Service) CancelOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	switch order.Status {
	case config.OrderStatus.Ordered, config.OrderStatus.Delivered, config.OrderStatus.Canceled:
		return customError.NewConflict(customError.OrderStatusConflict, order.Status, config.OrderStatus.Canceled)
	}
	return s.repo.UpdateStatusByID(ctx, id, config.OrderStatus.Canceled)
}

func (s *Service) DeleteOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != config.OrderStatus.Delivered {
		return customError.NewConflict(customError.OrderStatusConflict, order.Status, config.OrderStatus.Delivered)
	}
	if order.Status != config.OrderStatus.Canceled {
		return customError.NewConflict(customError.OrderStatusConflict, order.Status, config.OrderStatus.Canceled)
	}

	return s.repo.SoftDeleteByID(ctx, id)
}

func (s *Service) CalculatePremiumFinalPrice(ctx context.Context, orderID string) (*types.FinalPriceResult, error) {
	repoResult, err := s.repo.FindPriceWithMatchingDiscount(ctx, orderID, "premium")
	if err != nil {
		return nil, err
	}

	result := s.calculatePriceFromRepoResult(repoResult)

	return result, nil
}

func (s *Service) CalculateNonPremiumFinalPrice(ctx context.Context, orderID string) (*types.FinalPriceResult, error) {
	repoResult, err := s.repo.FindPriceWithMatchingDiscount(ctx, orderID, "non-premium")
	if err != nil {
		return nil, err
	}

	fmt.Println(repoResult.Discount)
	result := s.calculatePriceFromRepoResult(repoResult)

	return result, nil
}

func (s *Service) fetchCustomerByID(customerID, token string) (*types.CustomerResponseModel, error) {
	if customerID == "" {
		return nil, errors.New("customerID empty")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	if token != "" {
		headers["Authorization"] = token
	}

	var customer types.CustomerResponseModel
	if err := s.client.Get("/customer/"+customerID, headers, &customer); err != nil {
		return nil, err
	}
	return &customer, nil
}

func calculateTotalPrice(items []types.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.UnitPrice
	}
	return total
}

func (s *Service) calculatePriceFromRepoResult(repoResult *types.OrderPriceInfo) *types.FinalPriceResult {
	totalPrice := repoResult.TotalPrice
	var discountAmount float64 = 0.0
	var discountType string

	if repoResult.Discount != nil {
		d := repoResult.Discount
		discountType = d.Type

		switch d.Type {
		case "percentage":
			discountAmount = totalPrice * (d.Value / 100)
		case "fixed-amount":
			discountAmount = d.Value
		default:
			discountAmount = 0.0
			discountType = "unknown"
		}
	}

	finalPrice := totalPrice - discountAmount
	if finalPrice < 0 {
		finalPrice = 0
	}

	return &types.FinalPriceResult{
		OriginalPrice:   totalPrice,
		DiscountApplied: discountAmount,
		FinalPrice:      finalPrice,
		DiscountType:    discountType,
	}
}
