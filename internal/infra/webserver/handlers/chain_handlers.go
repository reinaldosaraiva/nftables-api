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

type ChainHandler struct {
    ChainDB  database.ChainInterface
    ProjectDB database.ProjectInterface
    TableDB   database.TableInterface
}

func NewChainHandler(chainDB database.ChainInterface, projectDB database.ProjectInterface, tableDB database.TableInterface) *ChainHandler {
    return &ChainHandler{
        ChainDB:  chainDB,
        ProjectDB: projectDB,
        TableDB:   tableDB,
    }
}

// CreateChain godoc
// @Summary Create a new Chain
// @Description Create Chains
// @Tags Chains
// @Accept  json
// @Produce  json
// @Param request body dto.CreateChainDTO true "chain request"
// @Success 201
// @Failure 400,500 {object} Error
// @Router /chains [post]
// @Security ApiKeyAuth
func (h *ChainHandler) CreateChain(w http.ResponseWriter, r *http.Request) {
    var chainDTO dto.CreateChainDTO
    err := json.NewDecoder(r.Body).Decode(&chainDTO)
    fmt.Println(chainDTO)
    fmt.Println(err)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    project, err := h.ProjectDB.FindByName(chainDTO.ProjectName)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    table, err := h.TableDB.FindByName(chainDTO.TableName)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    chain := &entity.Chain{
        Name:     chainDTO.Name,
        Type:     chainDTO.Type,
        Policy:   chainDTO.Policy,
        Priority: chainDTO.Priority,
        ProjectID: project.ID,
        TableID:  table.ID,
    }

    err = h.ChainDB.Create(chain)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}
// GetChain godoc
// @Summary Get a chain by ID
// @Description Get a chain by ID
// @Tags Chains
// @Accept  json
// @Produce  json
// @Param id path int true "Chain ID"
// @Success 200 {object} dto.CreateChainDTO 
// @Failure 400,404
// @Router /chains/{id} [get]
// @Security ApiKeyAuth
func (h *ChainHandler) GetChain(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    chain, err := h.ChainDB.FindByID(id)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(chain)
}

// GetChainsWithFilters godoc
// @Summary Get chains with optional filters
// @Description Get chains filtered by ID or name
// @Tags Chains
// @Accept  json
// @Produce  json
// @Param id query int false "Chain ID"
// @Param name query string false "Chain Name"
// @Success 200 {array} dto.CreateChainDTO 
// @Failure 400 "Invalid parameter format"
// @Failure 404 "Chain not found"
// @Router /chains/filter [get]
// @Security ApiKeyAuth
func (h *ChainHandler) GetChainsWithFilters(w http.ResponseWriter, r *http.Request) {
    queryValues := r.URL.Query()
    name := queryValues.Get("name")
    idStr := queryValues.Get("id")

    var chain *dto.DetailsChainDTO
    var err error

    if idStr != "" {
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        chain, err = h.ChainDB.FindByID(uint64(id))
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        
    } else if name != "" {
        chain, err = h.ChainDB.FindByName(name)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        
    } else {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(chain)
}

// UpdateChain godoc
// @Summary Update a chain by ID
// @Description Update a chain by ID
// @Tags Chains
// @Accept  json
// @Produce  json
// @Param id path int true "Chain ID"
// @Param request body dto.CreateChainDTO true "chain request"
// @Success 200
// @Failure 400,404,500 {object} Error
// @Router /chains/{id} [put]
// @Security ApiKeyAuth
func (h *ChainHandler) UpdateChain(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var chain entity.Chain
    err := json.NewDecoder(r.Body).Decode(&chain)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    chain.ID = uint(id)
    err = h.ChainDB.Update(&chain)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// DeleteChain godoc
// @Summary Delete a chain by ID
// @Description Delete a chain by ID
// @Tags Chains
// @Accept  json
// @Produce  json
// @Param id path int true "Chain ID"
// @Success 200
// @Failure 400,404,500 {object} Error
// @Router /chains/{id} [delete]
// @Security ApiKeyAuth
func (h *ChainHandler) DeleteChain(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    err = h.ChainDB.Delete(id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// GetChains godoc
// @Summary Get all chains
// @Description Get all chains
// @Tags Chains
// @Accept  json
// @Produce  json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Success 200 {object} []dto.CreateChainDTO 
// @Failure 400,500 {object} Error
// @Router /chains [get]
// @Security ApiKeyAuth
func (h *ChainHandler) GetChains(w http.ResponseWriter, r *http.Request) {
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

    chains, err := h.ChainDB.FindAll(page, limit, sort)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(chains)
}