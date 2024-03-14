package interfaces

import (
	"order/service/pkg/domain"
	"order/service/pkg/util/models"
)

type OrderUseCase interface {
	OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (domain.OrderSuccessResponse, error)
	GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error)
}
