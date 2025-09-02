package internal

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tesodev-korpes/pkg/customError"

	"tesodev-korpes/OrderService/internal/types"
	"tesodev-korpes/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// @title Order Service API
// @version 1.0
// @description API for managing order data
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
type Handler struct {
	service  *Service
	validate *validator.Validate
}

func NewHandler(e *echo.Echo, service *Service, clientMongo *mongo.Client) *Handler {
	validate := validator.New()

	handler := &Handler{
		service:  service,
		validate: validate,
	}

	g := e.Group("/order")
	g.POST("", handler.Create)
	g.GET("/:id", handler.GetByID)
	g.PATCH("/:id/ship", handler.ShipOrder)
	g.PATCH("/:id/deliver", handler.DeliverOrder)
	g.DELETE("/cancel/:id", handler.CancelOrder)
	g.GET("/list", handler.GetAllOrders)

	return handler
}

// Create godoc
// @Summary Create a new order
// @Description Create a new order with the given data
// @Tags orders
// @Accept json
// @Produce json
// @Param order body types.CreateOrderRequestModel true "Order to create"
// @Success 201 {object} types.OrderResponseModel "Returns created order details"
// @Failure 400 {object} customError.AppError "Invalid request body"
// @Failure 500 {object} customError.AppError "Internal server error"
// @Router /orders [post]
func (h *Handler) Create(c echo.Context) error {
	var req types.CreateOrderRequestModel
	if err := c.Bind(&req); err != nil {
		return customError.NewBadRequest(customError.InvalidOrderBody)
	}

	token := c.Request().Header.Get("Authorization") // ← kullanıcının JWT’si

	order := FromCreateOrderRequest(&req)

	createdID, err := h.service.Create(c.Request().Context(), order, token) // ← token eklendi
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return customError.NewNotFound(customError.CustomerNotFound)
		}
		return customError.NewInternal(customError.OrderServiceError, err)
	}

	createdOrder, err := h.service.GetByID(c.Request().Context(), createdID, token)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return customError.NewNotFound(customError.OrderNotFound)
		}
		return customError.NewInternal(customError.OrderServiceError, err)
	}
	fmt.Println(createdOrder)

	return c.JSON(http.StatusCreated, createdOrder)
}

func (h *Handler) GetByID(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")
	if !pkg.IsValidUUID(id) {
		return customError.NewBadRequest(customError.InvalidOrderID)
	}

	token := c.Request().Header.Get("Authorization")                              // ← token al
	orderWithCustomer, err := h.service.GetByID(c.Request().Context(), id, token) // ← token ver
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return customError.NewNotFound(customError.OrderNotFound)
		}
		customError.LogErrorWithCorrelation(err, correlationID)
		return customError.NewInternal(customError.OrderServiceError, err)
	}

	customError.LogInfoWithCorrelation("Order with customer fetched", correlationID)
	return c.JSON(http.StatusOK, orderWithCustomer)
}

// GetByID godoc
// @Summary Get order by ID
// @Description Retrieve an order with customer details by its unique ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Success 200 {object} types.OrderWithCustomerResponse "Order details with customer info"
// @Failure 400 {object} customError.AppError "Invalid ID format"
// @Failure 404 {object} customError.AppError "Order not found"
// @Failure 500 {object} customError.AppError "Internal server error"
// @Router /orders/{id} [get]

// ShipOrder godoc
// @Summary Ship an order
// @Description Mark an order as shipped by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 500 {object} customError.AppError "Internal server error"
// @Router /orders/{id}/ship [put]
func (h *Handler) ShipOrder(c echo.Context) error {
	id := c.Param("id")
	if !pkg.IsValidUUID(id) {
		return customError.NewBadRequest(customError.InvalidOrderID)
	}
	err := h.service.ShipOrder(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return customError.NewNotFound(customError.OrderNotFound)
		}

		var appErr *customError.AppError
		if errors.As(err, &appErr) {
			if appErr.Code == customError.ErrorDefinitions[customError.OrderStatusConflict].TypeCode {
				return err
			}

		}

		return customError.NewInternal(customError.OrderServiceError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order shipped successfully"})
}

// DeliverOrder godoc
// @Summary Deliver an order
// @Description Changes the order status to DELIVERED if current status is SHIPPED
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string "Order delivered successfully"
// @Failure 400 {object} customError.AppError "Invalid ID format"
// @Failure 404 {object} customError.AppError "Order not found"
// @Failure 409 {object} customError.AppError "Invalid order state for delivery"
// @Failure 500 {object} customError.AppError "Internal server error"
// @Security ApiKeyAuth
// @Router /order/{id}/deliver [put]
func (h *Handler) DeliverOrder(c echo.Context) error {
	id := c.Param("id")
	if !pkg.IsValidUUID(id) {
		return customError.NewBadRequest(customError.InvalidOrderID)
	}

	err := h.service.DeliverOrder(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return customError.NewNotFound(customError.OrderNotFound)
		}

		var appErr *customError.AppError
		if errors.As(err, &appErr) {
			if appErr.Code == customError.ErrorDefinitions[customError.OrderStatusConflict].TypeCode {
				return err
			}

		}

		return customError.NewInternal(customError.OrderServiceError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order delivered successfully"})
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel an order by its ID.
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} customError.AppError "Invalid ID format"
// @Failure 404 {object} customError.AppError "Order not found"
// @Failure 500 {object} customError.AppError "Internal server error"
// @Router /orders/{id}/cancel [put]
func (h *Handler) CancelOrder(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")
	if !pkg.IsValidUUID(id) {
		return customError.NewBadRequest(customError.InvalidOrderID)
	}

	err := h.service.CancelOrder(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return customError.NewNotFound(customError.OrderNotFound)
		}

		var appErr *customError.AppError
		if errors.As(err, &appErr) {
			if appErr.Code == customError.ErrorDefinitions[customError.OrderStatusConflict].TypeCode {
				return err
			}

		}

		customError.LogErrorWithCorrelation(err, correlationID)
		return customError.NewInternal(customError.OrderServiceError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order cancelled successfully. "})
}

// DeleteOrder godoc
// @Summary Soft delete an order by ID
// @Description Marks the order as deleted without removing it permanently.
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string "Deletion success message"
// @Failure 404 {object} customError.AppError "Order not found"
// @Failure 409 {object} customError.AppError "Conflict error"
// @Failure 500 {object} customError.AppError "Internal server error"
// @Router /orders/{id} [delete]
func (h *Handler) DeleteOrder(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")
	if !pkg.IsValidUUID(id) {
		return customError.NewBadRequest(customError.InvalidOrderID)
	}

	err := h.service.DeleteOrder(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return customError.NewNotFound(customError.OrderNotFound)
		}

		var appErr *customError.AppError
		if errors.As(err, &appErr) {
			if appErr.Code == customError.ErrorDefinitions[customError.OrderStatusConflict].TypeCode {
				return err
			}

		}

		customError.LogErrorWithCorrelation(err, correlationID)
		return customError.NewInternal(customError.OrderServiceError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order deleted (soft delete) successfully."})
}

// GetAllOrders godoc
// @Summary List all orders with pagination
// @Description Retrieve a paginated list of all orders.
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{} "Returns list of orders with pagination"
// @Failure 404 {object} customError.AppError "No orders found"
// @Failure 500 {object} customError.AppError "Internal server error"
// @Router /orders [get]
func (h *Handler) GetAllOrders(c echo.Context) error {
	params := types.Pagination{
		Page:  1,
		Limit: 10,
	}

	if p := c.QueryParam("page"); p != "" {
		if pageInt, err := strconv.Atoi(p); err == nil && pageInt > 0 {
			params.Page = pageInt
		}
	}

	if l := c.QueryParam("limit"); l != "" {
		if limitInt, err := strconv.Atoi(l); err == nil && limitInt > 0 {
			params.Limit = limitInt
		}
	}

	orders, err := h.service.GetAllOrders(c.Request().Context(), params)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return customError.NewNotFound(customError.OrderNotFound)
		}
		return customError.NewInternal(customError.OrderServiceError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": orders})
}

func (h *Handler) GetPremiumOrderPrice(c echo.Context) error {
	orderID := c.Param("id")
	if orderID == "" {
		return customError.NewBadRequest(customError.EmptyOrderID)
	}

	if !pkg.IsValidUUID(orderID) {
		return customError.NewBadRequest(customError.InvalidOrderID)
	}

	result, err := h.service.CalculatePremiumFinalPrice(c.Request().Context(), orderID)
	if err != nil {

		var appErr *customError.AppError
		if errors.As(err, &appErr) {
			if appErr.Code == customError.ErrorDefinitions[customError.OrderNotFound].TypeCode {
				return err
			}

			return customError.NewInternal(customError.OrderServiceError, err)
		}
	}
	return c.JSON(http.StatusOK, result)
}

func (h *Handler) GetNonPremiumOrderPrice(c echo.Context) error {
	orderID := c.Param("id")

	if orderID == "" {
		return customError.NewBadRequest(customError.EmptyOrderID)
	}
	if !pkg.IsValidUUID(orderID) {
		return customError.NewBadRequest(customError.InvalidOrderID)
	}

	result, err := h.service.CalculateNonPremiumFinalPrice(c.Request().Context(), orderID)
	if err != nil {

		var appErr *customError.AppError
		if errors.As(err, &appErr) {
			if appErr.Code == customError.ErrorDefinitions[customError.OrderNotFound].TypeCode {
				return err
			}

			return customError.NewInternal(customError.OrderServiceError, err)
		}
	}

	return c.JSON(http.StatusOK, result)
}
