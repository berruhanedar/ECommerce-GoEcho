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
			"status":     config.OrderStatus.Canceled, //types.OrderCanceled
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

// func (r *Repository) FindPriceWithMatchingDiscount(ctx context.Context, orderID string, role string) (*types.OrderPriceInfo, error) {

// 	now := time.Now()

// 	pipeline := mongo.Pipeline{
// 		{{Key: "$match", Value: bson.D{
// 			{Key: "_id", Value: orderID},
// 		}}},
// 		{{Key: "$project", Value: bson.D{
// 			{Key: "total_price", Value: 1},
// 			{Key: "discount", Value: bson.D{
// 				{Key: "$filter", Value: bson.D{
// 					{Key: "input", Value: "$discount"},
// 					{Key: "as", Value: "d"},
// 					{Key: "cond", Value: bson.D{
// 						{Key: "$and", Value: bson.A{
// 							bson.D{{Key: "$eq", Value: bson.A{"$$d.role", role}}},
// 							bson.D{{Key: "$lte", Value: bson.A{"$$d.start_date", now}}},
// 							bson.D{{Key: "$gte", Value: bson.A{"$$d.end_date", now}}},
// 						}},
// 					}},
// 				}},
// 			}},
// 		}}},
// 	}

// 	cursor, err := r.collection.Aggregate(ctx, pipeline)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	if !cursor.Next(ctx) {
// 		return nil, customError.NewNotFound(customError.OrderNotFound)
// 	}

// 	var tempResult types.AggregationResult
// 	if err := cursor.Decode(&tempResult); err != nil {
// 		return nil, err
// 	}

// 	finalResult := &types.OrderPriceInfo{
// 		TotalPrice: tempResult.TotalPrice,
// 	}

// 	if len(tempResult.Discount) > 0 {
// 		finalResult.Discount = &tempResult.Discount[0]
// 	}

// 	return finalResult, nil
// }

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

// func (r *Repository) FindPriceWithMatchingDiscount(ctx context.Context, orderID string) (*types.OrderPriceInfo, error) {
// 	var orderData types.OrderPriceInfo
// 	projection := options.FindOne().SetProjection(bson.M{
// 		"total_price": 1,
// 		"discount":    1,
// 	})
// 	err := r.collection.FindOne(ctx, bson.M{"_id": orderID}, projection).Decode(&orderData)
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return nil, customError.NewNotFound(customError.OrderNotFound)
// 		}
// 		return nil, err
// 	}
// 	return &orderData, nil
// }
