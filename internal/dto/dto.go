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
    ID uint64 `json:"id"`
	Name string `json:"name"`
}

type CreateProjectDTO struct {
    
    Name     string `json:"name"`
    TenantName string `json:"tenant_name"`
    
}
type DetailsProjectDTO struct {
    ID        uint64 `json:"id"`
    Name      string `json:"name"`
    TenantID  uint64 `json:"tenant_id"`
    TenantName string `json:"tenant_name"`
}

type CreateTableDTO struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Type        string `json:"type"`
}
type DetailsChainDTO struct {
    ID          uint64 `json:"id"`
    Name        string `json:"chain_name"`
    Type        string `json:"chain_type"`
    Priority    int    `json:"chain_priority"`
    Policy       string `json:"chain_policy"`
    ProjectID  uint64 `json:"project_id"`
    ProjectName   string `json:"project_name"`
    TableID     uint64 `json:"table_id"`
    TableName     string `json:"table_name"`
    Rules        []DetailsRuleDTO `json:"rules"`
}

type DetailsRuleDTO struct {
    Protocol    string `json:"protocol"`
    Port        int    `json:"port"`
    Action      string `json:"action"`
}
type DetailsTableDTO struct {
    ID          uint64 `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Type        string `json:"type"`
    Chains      []DetailsChainDTO `json:"chains"`
}

type CreateChainDTO struct {
    ID         uint64 `json:"id"`
    Name        string `json:"name"`
    Type        string `json:"type"`
    Priority    int    `json:"priority"`
    Policy       string `json:"policy"`
    ProjectID uint64 `json:"project_id"`
    ProjectName   string `json:"project_name"`
    TableID     uint64 `json:"table_id"`
	TableName     string `json:"table_name"`
}

type CreateServiceDTO struct {
    Name   string `json:"name"`
    Port   int    `json:"port"`
}

type CreateNetworkObjectDTO struct {
    Name    string `json:"name"`
    Address string `json:"address"`
}

type CreateRuleDTO struct {
    ChainName           string   `json:"chain_name"`
    Protocol            string   `json:"protocol"`
    Port                int      `json:"port"`
    Action              string   `json:"action"`
    ServiceRules        []CreateServiceDTO `json:"service_rules"` 
    NetworkObjectRules  []CreateNetworkObjectDTO  `json:"network_object_rules"` 
}