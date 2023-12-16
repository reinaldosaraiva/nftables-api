package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/reinaldosaraiva/nftables-api/internal/dto"
	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/reinaldosaraiva/nftables-api/internal/infra/database"
)

type ChainHandler struct {
    ChainDB database.ChainInterface
}

func NewChainHandler(db database.ChainInterface) *ChainHandler {
    return &ChainHandler{ChainDB: db}
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
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    chain, err := entity.NewChain(chainDTO.Name, chainDTO.Description, chainDTO.Type, chainDTO.State, chainDTO.ProjectID)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
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
