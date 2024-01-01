// tenant_handler.go
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

type TenantHandler struct {
	TenantDB database.TenantInterface
}

func NewTenantHandler(db database.TenantInterface) *TenantHandler {
	return &TenantHandler{TenantDB: db}
}
// Create Tenant godoc
// @Summary Create a new Tenant
// @Description Create Tenants
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Param request body dto.CreateTenantDTO true "tenant request"
// @Success 201
// @Failure 400,500 {object} Error
// @Router /tenants [post]
// @Security ApiKeyAuth
func (h *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var tenant dto.CreateTenantDTO
	err := json.NewDecoder(r.Body).Decode(&tenant)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	t, err:= entity.NewTenant(tenant.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.TenantDB.Create(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Get a tenant by ID godoc
// @Summary Get a tenant by ID
// @Description Get a tenant by ID
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Param ID path int true "tenant ID"
// @Success 200 {object} dto.CreateTenantDTO 
// @Failure 400,404
// @Router /tenants/{id} [get]
// @Security ApiKeyAuth
func (h *TenantHandler) GetTenantByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tenant, err := h.TenantDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenant)
}

// Get a tenant by name godoc
// @Summary Get a tenant by name
// @Description Get a tenant by name
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Param name path string true "tenant name"
// @Success 200 {object} dto.CreateTenantDTO 
// @Failure 400,404
// @Router /tenants/name/{name} [get]
// @Security ApiKeyAuth
func (h *TenantHandler) GetTenantByName(w http.ResponseWriter, r *http.Request) {
    name := chi.URLParam(r, "name")
	tenant, err := h.TenantDB.FindByName(name)
    if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenant)
}

// Get tenants with filters godoc
// @Summary Get tenants with optional filters
// @Description Get tenants filtered by ID or name
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Param id query int false "Tenant ID"
// @Param name query string false "Tenant Name"
// @Success 200 {array} dto.CreateTenantDTO 
// @Failure 400 "Invalid parameter format"
// @Failure 404 "Tenant not found"
// @Router /tenants/filter [get]
// @Security ApiKeyAuth
func (h *TenantHandler) GetTenantsWithFilters(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
    name := queryValues.Get("name")
    idStr := queryValues.Get("id")

	var tenant *dto.CreateTenantDTO
    var err error

    if idStr != "" {
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        tenant, err = h.TenantDB.FindByID(uint64(id))
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        
    } else if name != "" {
        tenant, err = h.TenantDB.FindByName(name)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        
    } else {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(tenant)
}

// Update a tenant by ID godoc
// @Summary Update a tenant by ID
// @Description Update a tenant by ID
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Param id path int true "tenant id"
// @Param request body dto.CreateTenantDTO true "tenant request"
// @Success 200
// @Failure 400,404,500 {object} Error
// @Router /tenants/{id} [put]
// @Security ApiKeyAuth
func (h *TenantHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var tenant entity.Tenant
	err := json.NewDecoder(r.Body).Decode(&tenant)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tenant.ID = uint(id)
	err = h.TenantDB.Update(&tenant)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete Tenant by ID godoc
// @Summary Delete Tenant a tenant by ID
// @Description Delete Tenant a tenant by ID
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Param id path int true "tenant id"
// @Success 200
// @Failure 400,404,500 {object} Error
// @Router /tenants/{id} [delete]
// @Security ApiKeyAuth
func (h *TenantHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.TenantDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Get all Tenants godoc
// @Summary Get all Tenants
// @Description Get all Tenants
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Success 200 {object} []dto.CreateTenantDTO 
// @Failure 400,500 {object} Error
// @Router /tenants [get]
// @Security ApiKeyAuth
func (h *TenantHandler) GetTenants(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}

	tenants, err := h.TenantDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenants)
}
