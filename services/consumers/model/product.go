package model

import (
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Name     string
	Category string
	Brand    string
	Prices   []float32
}

type ProductNotification struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Brand    string  `json:"brand"`
	Price    float32 `json:"price"`
}

type ProductQuantity struct {
	Product  primitive.ObjectID
	Quantity int
}

func (p *Product) GetCurrentPrice() float32 {
	return p.Prices[len(p.Prices)-1]
}

func (p *Product) Similarity(otherProduct *Product) float32 {
	if p.ID == otherProduct.ID {
		return 1
	}

	baseSimilarity := float32(0)

	if p.Category == otherProduct.Category {
		baseSimilarity += 0.25
	}

	if p.Brand == otherProduct.Brand {
		baseSimilarity += 0.25
	}

	lev := metrics.NewLevenshtein()
	lev.CaseSensitive = false
	nameSimilarity := float32(strutil.Similarity(p.Name, otherProduct.Name, lev)) / 2

	return baseSimilarity + nameSimilarity
}

func (p *Product) SimilarityMultiple(otherProducts []Product) float32 {
	fit := float32(0)

	for _, product := range otherProducts {
		similarity := p.Similarity(&product)
		fit += similarity
	}

	return fit
}
