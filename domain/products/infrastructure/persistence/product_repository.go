package persistence

import (
	"backend_crudgo/domain/products/domain/model"
	repoDomain "backend_crudgo/domain/products/domain/repository"
	"backend_crudgo/infrastructure/database"
	response "backend_crudgo/types"
	"errors"

	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
)

type sqlProductRepo struct {
	Conn *database.DataDB
}

// NewProductRepository Should initialize the dependencies for this service.
func NewProductRepository(Conn *database.DataDB) repoDomain.ProductRepository {
	return &sqlProductRepo{
		Conn: Conn,
	}
}

// CreateProduct takes in a context and a product as input and returns a CreateResponse and error.
// It prepares an SQL statement to insert a product into the database and then executes the query.
// If the query executes successfully, a success response is returned with a message "Product created".
func (sr *sqlProductRepo) CreateProduct(ctx context.Context, product *model.Product) (*response.CreateResponse, error) {
	var idResult string

	stmt, err := sr.Conn.DB.PrepareContext(ctx, InsertProduct)
	if err != nil {
		return &response.CreateResponse{}, err
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			log.Error().Msgf("Could not close testament : [error] %s", err.Error())
		}
	}()

	row := stmt.QueryRowContext(ctx, &product.ProductID, &product.ProductName, &product.ProductAmount,&product.ProductPrice , &product.ProductUserCreated, &product.ProductUserModify)

	if err = row.Scan(&idResult); err != sql.ErrNoRows {
		return &response.CreateResponse{}, err
	}

	return &response.CreateResponse{
		Message: "Product created",
	}, nil
}

// GetProduct takes in a context and an id as input and returns a GenericResponse and error.
// It prepares an SQL statement to select a product from the database by its id and then executes the query.
// If the query executes successfully, a success response is returned with the selected product.
func (sr *sqlProductRepo) GetProduct(ctx context.Context, id string) (*response.GenericResponse, error) {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, SelectProduct)
	if err != nil {
		return &response.GenericResponse{}, err
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			log.Error().Msgf("Could not close testament : [error] %s", err.Error())
		}
	}()

	row := stmt.QueryRowContext(ctx, id)
	product := &model.Product{}

	if err = row.Scan(&product.ProductID, &product.ProductName, &product.ProductAmount, &product.ProductPrice,
		&product.ProductUserCreated, &product.ProductDateCreated, &product.ProductUserModify, &product.ProductDateModify); err != nil {
		if err == sql.ErrNoRows {
			return &response.GenericResponse{Error: "Product not found"}, errors.New(" Product not found ")
		}
		return &response.GenericResponse{Error: err.Error()}, err
	}

	return &response.GenericResponse{
		Message: "Get product success",
		Product: product,
	}, nil
}

// GetProducts takes in a context as input and returns a GenericResponse and error.
// It prepares an SQL statement to select all products from the database and then executes the query.
// If the query executes successfully, a success response is returned with a list of selected products.
func (sr *sqlProductRepo) GetProducts(ctx context.Context) (*response.GenericResponse, error) {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, SelectProducts)
	if err != nil {
		return &response.GenericResponse{}, nil
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			log.Error().Msgf("Could not close testament : [error] %s", err.Error())
		}
	}()
	row, err := sr.Conn.DB.QueryContext(ctx, SelectProducts)

	var products []*model.Product
	for row.Next() {
		var product = &model.Product{}
		err = row.Scan(&product.ProductID, &product.ProductName, &product.ProductAmount, &product.ProductPrice,
			&product.ProductUserCreated, &product.ProductDateCreated, &product.ProductUserModify, &product.ProductDateModify)

		products = append(products, product)
	}
	if err != nil {
		return &response.GenericResponse{Error: err.Error()}, err
	}

	return &response.GenericResponse{
		Message: "Get product success",
		Product: products,
	}, nil
}

// UpdateProduct takes in a context, an id and a product as input and returns a GenericResponse and error.
// It prepares an SQL statement to update a product in the database and then executes the query.
// If the query executes successfully, a success response is returned with a message "Product updated".
func (sr *sqlProductRepo) UpdateProduct(ctx context.Context, id string, product *model.Product) (*response.GenericResponse, error) {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, UpdateProduct)
	if err != nil {
		return &response.GenericResponse{}, err
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			log.Error().Msgf("Could not close testament : [error] %s", err.Error())
		}
	}()

	if _, err = stmt.ExecContext(ctx, &product.ProductName, &product.ProductAmount, &product.ProductPrice, &product.ProductUserModify, id); err != nil {
		return &response.GenericResponse{Error: err.Error()}, err
	}

	return &response.GenericResponse{
		Message: "Product updated",
	}, nil
}

// DeleteProduct takes in a context and an id as input and returns a GenericResponse and error.
// It prepares an SQL statement to delete a product from the database by its id and then executes the query.
// If the query executes successfully, a success response is returned with a message "Product deleted".
func (sr *sqlProductRepo) DeleteProduct(ctx context.Context, id string) (*response.GenericResponse, error) {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, DeleteProduct)
	if err != nil {
		return &response.GenericResponse{}, err
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			log.Error().Msgf("Could not close testament : [error] %s", err.Error())
		}
	}()

	if _, err = stmt.ExecContext(ctx, id); err != nil {
		return &response.GenericResponse{Error: err.Error()}, err
	}

	return &response.GenericResponse{
		Message: "Product deleted",
	}, nil
}
