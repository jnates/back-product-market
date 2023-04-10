package persistence

const (
	InsertProduct  = "INSERT INTO public.products (product_name,product_amount,product_price,product_user_created,product_date_created,product_user_modify,product_date_modify) VALUES ($1,$2,$3,$4,'NOW()',1,'NOW()')"
	SelectProduct  = "SELECT product_id, product_name, product_amount,product_price, product_user_created, product_date_created, product_user_modify, product_date_modify FROM public.products WHERE product_id = $1"
	SelectProducts = "SELECT product_id, product_name, product_amount,product_price, product_user_created, product_date_created, product_user_modify, product_date_modify FROM public.products"
	UpdateProduct  = "UPDATE public.products SET product_name = $1, product_amount = $2, product_user_modify = $3, product_date_modify = NOW() WHERE product_id = $4"
	DeleteProduct  = "DELETE FROM public.products WHERE product_id = $1"
)
