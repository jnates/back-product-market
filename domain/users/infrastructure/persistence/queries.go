package persistence

const (
	InsertUser      = "INSERT INTO public.users (user_id,user_name,user_identifier,user_email,user_password,user_type_identifier) VALUES ($1,$2,$3,$4,$5,$6)"
	SelectUser      = "SELECT user_id, user_name, user_identifier, user_email, user_password, user_type_identifier  FROM public.users WHERE user_id = $1"
	SelectLoginUser = "SELECT user_id, user_name, user_identifier, user_email, user_password, user_type_identifier  FROM public.users WHERE user_name = $1"
	SelectUsers     = "SELECT user_id, user_name, user_identifier, user_email, user_password, user_type_identifier  FROM public.users"
)
