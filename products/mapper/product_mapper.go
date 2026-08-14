// Package mapper converts between persistence entities and API models for the products domain.
package mapper

import "backend_crudgo/products/models"

// ToProductModel converts a persistence entity into the domain/API model.
func ToProductModel(db *models.ProductDB) *models.Product {
	return &models.Product{
		ProductID:          db.ProductID,
		ProductName:        db.ProductName,
		ProductAmount:      db.ProductAmount,
		ProductPrice:       db.ProductPrice,
		ProductUserCreated: db.ProductUserCreated,
		ProductDateCreated: db.ProductDateCreated,
		ProductUserModify:  db.ProductUserModify,
		ProductDateModify:  db.ProductDateModify,
	}
}
