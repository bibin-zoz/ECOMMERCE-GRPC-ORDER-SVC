package interfaces

import "order/service/pkg/util/models"

type CartClient interface {
	GetAllItemsFromCart(userID int) ([]models.Cart, error)
	DoesCartExist(userID int) (bool, error)
	TotalAmountInCart(userID int) (float64, error)
	UpdateCartAfterOrder(userID, productID int, quantity float64) error


}
