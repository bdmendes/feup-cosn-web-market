package model

type ProductQuantity struct {
	ProductID interface{} `json:"product_id"`
	Quantity  int         `json:"quantity"`
}
