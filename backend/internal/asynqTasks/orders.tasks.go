package asynqTasks

import (
	"context"
	"encoding/json"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders"
	ordersTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type OrderAsynqTasks struct {
	orderRepo orders.OrderRepository
	logger    *zap.SugaredLogger
}

func NewOrderAsynqTasks(orderRepo orders.OrderRepository, logger *zap.SugaredLogger) *OrderAsynqTasks {
	return &OrderAsynqTasks{
		orderRepo: orderRepo,
		logger:    logger,
	}
}

func (h *OrderAsynqTasks) HandleOrderPaymentFailureTask(ctx context.Context, t *asynq.Task) error {
	var payload ordersTypes.WaitingOrderPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	err := h.orderRepo.HandlePaymentFailure(payload.OrderID)
	if err != nil {
		if errors.Is(err, ordersTypes.ErrOrderNotFound) {
			h.logger.Warnf("ℹ️ Order %d already deleted, skipping deferred deletion", payload.OrderID)
			return nil
		}

		if errors.Is(err, ordersTypes.ErrInappropriateOrderStatus) {
			h.logger.Warnf("ℹ️ Order %d is already paid, skipping deferred deletion", payload.OrderID)
			return nil
		}

		return err
	}

	return nil
}
