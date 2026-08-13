// Package constants centralizes literal values used across the products domain
// (route params, table and column names) so they are defined once and referenced everywhere.
package constants

const (
	// ID is the route parameter name for the product identifier.
	ID = "id"

	// ProductsTable is the public.products table name.
	ProductsTable = "products"

	// ColumnProductID is the products table primary key column.
	ColumnProductID = "product_id"
	// ColumnProductName is the product name column.
	ColumnProductName = "product_name"
	// ColumnProductAmount is the product amount column.
	ColumnProductAmount = "product_amount"
	// ColumnProductPrice is the product price column.
	ColumnProductPrice = "product_price"
	// ColumnProductUserCreated is the column holding the creator's user ID.
	ColumnProductUserCreated = "product_user_created"
	// ColumnProductDateCreated is the creation timestamp column.
	ColumnProductDateCreated = "product_date_created"
	// ColumnProductUserModify is the column holding the last modifier's user ID.
	ColumnProductUserModify = "product_user_modify"
	// ColumnProductDateModify is the last modification timestamp column.
	ColumnProductDateModify = "product_date_modify"
)

// Columns lists every products column returned by SELECT queries, in scan order.
var Columns = []string{
	ColumnProductID,
	ColumnProductName,
	ColumnProductAmount,
	ColumnProductPrice,
	ColumnProductUserCreated,
	ColumnProductDateCreated,
	ColumnProductUserModify,
	ColumnProductDateModify,
}
