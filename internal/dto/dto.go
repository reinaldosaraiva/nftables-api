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
type CreateTableDTO struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Type        string `json:"type"`
}
type CreateChainDTO struct {
    Name        string `json:"name"`
    Type        string `json:"type"`
    Priority    int    `json:"priority"`
    Policy       string `json:"policy"`
    ProjectID   uint64 `json:"project_id"`
	TableID     uint64 `json:"table_id"`
}