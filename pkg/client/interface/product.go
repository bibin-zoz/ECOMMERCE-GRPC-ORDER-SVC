package interfaces

type ProductClient interface {
	ProductStockMinus(productID, stock int) error
}
