package model

import (
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID
	Description string
	Category    string
	Prices      []float32
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
		baseSimilarity = 0.5
	}

	lev := metrics.NewLevenshtein()
	lev.CaseSensitive = false
	descriptionSimilarity := float32(strutil.Similarity(p.Description, otherProduct.Description, lev)) / 2

	return baseSimilarity + descriptionSimilarity
}

func (p *Product) SimilarityMultiple(otherProducts []*Product) float32 {
	fit := float32(0)

	for _, product := range otherProducts {
		similarity := p.Similarity(product)
		fit += similarity
	}

	return fit
}
