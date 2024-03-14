package usecase

import (
	"errors"

	clienterface "order/service/pkg/client/interface"
	"order/service/pkg/domain"
	interfaces "order/service/pkg/repository/interface"
	services "order/service/pkg/usecase/interface"
	"order/service/pkg/util/models"

	"github.com/jinzhu/copier"
)

type orderUseCase struct {
	orderRepository   interfaces.OrderRepository
	cartRepository    clienterface.CartClient
	productRepository clienterface.ProductClient
}

func NewOrderUseCase(repository interfaces.OrderRepository, cartRepo clienterface.CartClient, productRepo clienterface.ProductClient) services.OrderUseCase {
	return &orderUseCase{
		orderRepository:   repository,
		cartRepository:    cartRepo,
		productRepository: productRepo,
	}
}
func (or *orderUseCase) OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (domain.OrderSuccessResponse, error) {
	var orderBody models.OrderIncoming
	err := copier.Copy(&orderBody, &orderFromCart)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	orderBody.UserID = userID
	cartExist, err := or.cartRepository.DoesCartExist(userID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if !cartExist {
		return domain.OrderSuccessResponse{}, errors.New("cart empty can't order")
	}

	addressExist, err := or.orderRepository.AddressExist(orderBody)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if !addressExist {
		return domain.OrderSuccessResponse{}, errors.New("address does not exist")
	}

	paymentExist, err := or.orderRepository.PaymentExist(orderBody)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if !paymentExist {
		return domain.OrderSuccessResponse{}, errors.New("payment does not exist")
	}

	cartItems, err := or.cartRepository.GetAllItemsFromCart(orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	total, err := or.cartRepository.TotalAmountInCart(orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	FinalPrice := total
	order_id, err := or.orderRepository.OrderItems(orderBody, FinalPrice)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if err := or.orderRepository.AddOrderProducts(order_id, cartItems); err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	orderSuccessResponse, err := or.orderRepository.GetBriefOrderDetails(order_id)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	var orderItemDetails domain.OrderItem
	for _, c := range cartItems {
		orderItemDetails.ProductID = c.ProductID
		orderItemDetails.Quantity = c.Quantity
		err := or.cartRepository.UpdateCartAfterOrder(userID, int(orderItemDetails.ProductID), orderItemDetails.Quantity)
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}
		err = or.productRepository.ProductStockMinus(int(orderItemDetails.ProductID), int(orderItemDetails.Quantity))
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}
	}
	return orderSuccessResponse, nil
}
func (or *orderUseCase) GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {

	fullOrderDetails, err := or.orderRepository.GetOrderDetails(userId, page, count)
	if err != nil {
		return []models.FullOrderDetails{}, err
	}
	return fullOrderDetails, nil

}
