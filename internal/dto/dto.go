package dto


type CreateUserDTO struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginDTO struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type GetJWTOutput struct {
	AccessToken string `json:"access_token"`
}	

type CreateTenantDTO struct {
	Name string `json:"name"`
}

type CreateProjectDTO struct {
    Name     string `json:"name"`
    TenantID uint64 `json:"tenant_id"`
}
