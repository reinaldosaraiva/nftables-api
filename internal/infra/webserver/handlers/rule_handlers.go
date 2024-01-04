package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/reinaldosaraiva/nftables-api/internal/dto"
	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/reinaldosaraiva/nftables-api/internal/infra/database"
	"gorm.io/gorm"
)

type RuleHandler struct {
    RuleDB         database.RuleInterface
    ChainDB        database.ChainInterface
    ServiceDB      database.ServiceInterface
    NetworkObjectDB database.NetworkObjectInterface
}

func NewRuleHandler(ruleDB database.RuleInterface, chainDB database.ChainInterface, serviceDB database.ServiceInterface, networkObjectDB database.NetworkObjectInterface) *RuleHandler {
    return &RuleHandler{
        RuleDB:         ruleDB,
        ChainDB:        chainDB,
        ServiceDB:      serviceDB,
        NetworkObjectDB: networkObjectDB,
    }
}

func (h *RuleHandler) createOrFindChain(ChainName string) (*dto.DetailsChainDTO, error) {

    chain,err := h.ChainDB.FindByName(ChainName)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }


    return chain, nil
}

func (h *RuleHandler) createOrFindService(dto dto.CreateServiceDTO) (*entity.Service, error) {

    service,err := h.ServiceDB.FindByName(dto.Name)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }

    if errors.Is(err, gorm.ErrRecordNotFound) {
        newService := entity.Service{
            Name: dto.Name,
            Port: dto.Port,
        }
        err = h.ServiceDB.Create(&newService)
        if err != nil {
            return nil, err
        }
        return &newService, nil
    }

    return service, nil
}

func (h *RuleHandler) createOrFindNetworkObject(dto dto.CreateNetworkObjectDTO) (*entity.NetworkObject, error) {
    
    net_obj,err := h.NetworkObjectDB.FindByName(dto.Name)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }


    if errors.Is(err, gorm.ErrRecordNotFound) {
        newNetworkObject := entity.NetworkObject{
            Name:    dto.Name,
            Address: dto.Address,
        }
        err = h.NetworkObjectDB.Create(&newNetworkObject)
        if err != nil {
            return nil, err
        }
        return &newNetworkObject, nil
    }

    return net_obj, nil
}


// CreateRule godoc
// @Summary Create a new Rule
// @Description Create Rule with Chain, Service, and Network Object relations
// @Tags Rules
// @Accept json
// @Produce json
// @Param request body dto.CreateRuleDTO true "Rule request"
// @Success 201
// @Failure 400,500 {object} Error
// @Router /rules [post]
// @Security ApiKeyAuth
func (h *RuleHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
    var ruleDTO dto.CreateRuleDTO
    fmt.Println(ruleDTO)
    if err := json.NewDecoder(r.Body).Decode(&ruleDTO); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    rule, err := h.buildRuleFromDTO(ruleDTO)
    fmt.Println(rule)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := h.RuleDB.Create(rule); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func (h *RuleHandler) buildRuleFromDTO(dto dto.CreateRuleDTO) (*entity.Rule, error) {
    fmt.Println("CHAIN NAME: "+dto.ChainName)
    chain, err := h.createOrFindChain(dto.ChainName)
    if err != nil {
        return nil, err
    }
    if chain == nil {
        return nil, errors.New("chain not found") 
    }

    rule := &entity.Rule{
        ChainID: chain.ID,
        Protocol: dto.Protocol,
        Port: dto.Port,
        Action: dto.Action,
    }

    for _, serviceDTO := range dto.ServiceRules {
        service, err := h.createOrFindService(serviceDTO)
        if err != nil {
            return nil, err
        }
        rule.ServiceRules = append(rule.ServiceRules, *service)
    }

    for _, networkObjectDTO := range dto.NetworkObjectRules {
        networkObject, err := h.createOrFindNetworkObject(networkObjectDTO)
        if err != nil {
            return nil, err
        }
        rule.NetworkObjectRules = append(rule.NetworkObjectRules, *networkObject)
    }

    return rule, nil
}

// GetRule godoc
// @Summary Get a rule by ID
// @Description Get a single rule by its ID
// @Tags Rules
// @Accept json
// @Produce json
// @Param id path int true "Rule ID"
// @Success 200 {object} dto.CreateRuleDTO 
// @Failure 400,404 {object} Error
// @Router /rules/{id} [get]
// @Security ApiKeyAuth
func (h *RuleHandler) GetRule(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    rule, err := h.RuleDB.FindByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(rule)
}

// GetRulesWithFilters godoc
// @Summary Get rules with optional filters
// @Description Get rules filtered by ID or name
// @Tags Rules
// @Accept  json
// @Produce  json
// @Param id query int false "Rule ID"
// @Param name query string false "Rule Name"
// @Success 200 {array} []dto.CreateRuleDTO 
// @Failure 400 "Invalid parameter format"
// @Failure 404 "Rule not found"
// @Router /rules/filter [get]
// @Security ApiKeyAuth
func (h *RuleHandler) GetRulesWithFilters(w http.ResponseWriter, r *http.Request) {
    queryValues := r.URL.Query()
    name := queryValues.Get("name")
    idStr := queryValues.Get("id")

    var rules []entity.Rule

    if idStr != "" {
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
            http.Error(w, "Invalid ID format", http.StatusBadRequest)
            return
        }
        rule, err := h.RuleDB.FindByID(id)
        if err != nil {
            http.Error(w, "Rule not found", http.StatusNotFound)
            return
        }
        rules = append(rules, *rule)

    } else if name != "" {
        rule, err := h.RuleDB.FindByName(name)
        if err != nil {
            http.Error(w, "Rule not found", http.StatusNotFound)
            return
        }
        rules = append(rules, *rule)

    } else {
        http.Error(w, "Invalid query parameters", http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(rules)
}


// UpdateRule godoc
// @Summary Update a rule
// @Description Update an existing rule
// @Tags Rules
// @Accept json
// @Produce json
// @Param id path int true "Rule ID"
// @Param rule body dto.CreateRuleDTO true "Update Rule"
// @Success 200 
// @Failure 400,404,500 {object} Error
// @Router /rules/{id} [put]
// @Security ApiKeyAuth
func (h *RuleHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var ruleDTO dto.CreateRuleDTO
    err = json.NewDecoder(r.Body).Decode(&ruleDTO)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    rule, err := h.RuleDB.FindByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    rule.Protocol = ruleDTO.Protocol
    rule.Port = ruleDTO.Port
    rule.Action = ruleDTO.Action

    err = h.RuleDB.Update(rule)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(rule)
}

// DeleteRule godoc
// @Summary Delete a rule
// @Description Delete a rule by its ID
// @Tags Rules
// @Accept json
// @Produce json
// @Param id path int true "Rule ID"
// @Success 200 {object} string "OK"
// @Failure 400,404,500 {object} Error
// @Router /rules/{id} [delete]
// @Security ApiKeyAuth
func (h *RuleHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err = h.RuleDB.Delete(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Rule deleted successfully")
}

// GetRules godoc
// @Summary Get all rules
// @Description Get a list of all rules
// @Tags Rules
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param sort query string false "Sort order"
// @Success 200 {array} []dto.CreateRuleDTO 
// @Failure 400,500 {object} Error
// @Router /rules [get]
// @Security ApiKeyAuth
func (h *RuleHandler) GetRules(w http.ResponseWriter, r *http.Request) {
    pageStr := r.URL.Query().Get("page")
    limitStr := r.URL.Query().Get("limit")
    sort := r.URL.Query().Get("sort")

    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 {
        limit = 10
    }

    rules, err := h.RuleDB.FindAll(page, limit, sort)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(rules)
}

