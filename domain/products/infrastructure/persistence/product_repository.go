package persistence

import (
	"backend_crudgo/domain/products/domain/model"
	repoDomain "backend_crudgo/domain/products/domain/repository"
	"backend_crudgo/infrastructure/database"
	response "backend_crudgo/types"

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

func (sr *sqlProductRepo) CreateProductHandler(ctx context.Context, product *model.Product) (*response.CreateResponse, error) {
	var idResult string

	stmt, err := sr.Conn.DB.PrepareContext(ctx, InsertProduct)
	if err != nil {
		return &response.CreateResponse{}, err
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error().Msgf("Could not close testament : [error] %s", err.Error())
		}
	}()

	row := stmt.QueryRowContext(ctx, &product.ProductID, &product.ProductName, &product.ProductAmount, &product.ProductUserCreated, &product.ProductUserModify)
	err = row.Scan(&idResult)
	if err != sql.ErrNoRows {
		return &response.CreateResponse{}, err
	}

	GenericResponse := response.CreateResponse{
		Message: "Product created",
	}

	return &GenericResponse, nil
}

func (sr *sqlProductRepo) GetProductHandler(ctx context.Context, id string) (*response.GenericResponse, error) {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, SelectProduct)
	if err != nil {
		return &response.GenericResponse{}, err
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error().Msgf("Could not close testament : [error] %s", err.Error())
		}
	}()

	row := stmt.QueryRowContext(ctx, id)
	product := &model.Product{}

	err = row.Scan(&product.ProductID, &product.ProductName, &product.ProductAmount, &product.ProductUserCreated,
		&product.ProductDateCreated, &product.ProductUserModify, &product.ProductDateModify)
	if err != nil {
		return &response.GenericResponse{Error: err.Error()}, err
	}

	GenericResponse := &response.GenericResponse{
		Message: "Get product success",
		Product: product,
	}

	return GenericResponse, nil
}

func (sr *sqlProductRepo) GetProductsHandler(ctx context.Context) (*response.GenericResponse, error) {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, SelectProducts)
	if err != nil {
		return &response.GenericResponse{}, nil
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error().Msgf("Could not close testament : [error] %s", err.Error())
		}
	}()
	row, err := sr.Conn.DB.QueryContext(ctx, SelectProducts)

	var products []*model.Product
	for row.Next() {
		var product = &model.Product{}
		err = row.Scan(&product.ProductID, &product.ProductName, &product.ProductAmount, &product.ProductUserCreated,
			&product.ProductDateCreated, &product.ProductUserModify, &product.ProductDateModify)

		products = append(products, product)
	}
	if err != nil {
		return &response.GenericResponse{Error: err.Error()}, err
	}
	GenericResponse := &response.GenericResponse{
		Message: "Get product success",
		Product: products,
	}

	return GenericResponse, nil
}
