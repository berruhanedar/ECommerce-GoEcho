package internal

import (
	"context"
	"errors"
	"tesodev-korpes/OrderService/config"
	"tesodev-korpes/OrderService/internal/types"
	"tesodev-korpes/pkg/customError"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) *Repository {
	return &Repository{collection: collection}
}

func (r *Repository) Create(ctx context.Context, order *types.Order) (string, error) {
	if order.Id == "" {
		order.Id = uuid.New().String()
	}
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return "", err
	}
	return order.Id, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*types.Order, error) {
	var order types.Order
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *Repository) GetAllOrders(ctx context.Context, findOptions *options.FindOptions) ([]types.Order, error) {
	var orders []types.Order

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *Repository) UpdateStatusByID(ctx context.Context, id string, status string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *Repository) Cancel(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     config.OrderStatus.Canceled,
			"updated_at": time.Now(),
		},
	}

	filter := bson.M{"_id": id}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("order not found")
	}

	return nil
}

func (r *Repository) SoftDeleteByID(ctx context.Context, id string) error {

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"is_delete":  true,
			"updated_at": time.Now(),
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *Repository) FindPriceWithMatchingDiscount(ctx context.Context, orderID string, role string) (*types.OrderPriceInfo, error) {
	var order types.Order
	filter := bson.M{"_id": orderID}
	err := r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, customError.NewNotFound(customError.OrderNotFound)
		}
		return nil, err
	}
	var matchingDiscount *types.Discount
	now := time.Now()
	for _, discount := range order.Discounts {
		if discount != nil && discount.Role == role && discount.StartDate.Before(now) && discount.EndDate.After(now) {
			matchingDiscount = discount
			break
		}
	}
	result := &types.OrderPriceInfo{
		TotalPrice: order.TotalPrice,
		Discount:   matchingDiscount,
	}
	return result, nil
}
