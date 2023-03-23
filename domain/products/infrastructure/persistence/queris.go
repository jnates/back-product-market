package persistence

const (
	InsertProduct  = "INSERT INTO public.products (product_id,product_name,product_amount,product_user_created,product_date_created,product_user_modify,product_date_modify) VALUES ($1,$2,$3,$4,'NOW()',$5,'NOW()')"
	SelectProduct  = "SELECT product_id, product_name, product_amount, product_user_created, product_date_created, product_user_modify, product_date_modify FROM public.products WHERE product_id = $1"
	SelectProducts = "SELECT product_id, product_name, product_amount, product_user_created, product_date_created, product_user_modify, product_date_modify FROM public.products"
)
