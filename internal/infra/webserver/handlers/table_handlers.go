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

type TableHandler struct {
    TableDB database.TableInterface
}

func NewTableHandler(db database.TableInterface) *TableHandler {
    return &TableHandler{TableDB: db}
}

// CreateTable godoc
// @Summary Create a new Table
// @Description Create Tables
// @Tags Tables
// @Accept  json
// @Produce  json
// @Param request body dto.CreateTableDTO true "table request"
// @Success 201
// @Failure 400,500 {object} Error
// @Router /tables [post]
// @Security ApiKeyAuth
func (h *TableHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
    var tableDTO dto.CreateTableDTO
    err := json.NewDecoder(r.Body).Decode(&tableDTO)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    table, err := entity.NewTable(tableDTO.Name, tableDTO.Description, tableDTO.Type, tableDTO.State)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    err = h.TableDB.Create(table)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
}

// GetTable godoc
// @Summary Get a table by ID
// @Description Get a table by ID
// @Tags Tables
// @Accept  json
// @Produce  json
// @Param id path int true "Table ID"
// @Success 200 {object} dto.CreateTableDTO
// @Failure 400,404
// @Router /tables/{id} [get]
// @Security ApiKeyAuth
func (h *TableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    table, err := h.TableDB.FindByID(id)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(table)
}

// UpdateTable godoc
// @Summary Update a table by ID
// @Description Update a table by ID
// @Tags Tables
// @Accept  json
// @Produce  json
// @Param id path int true "Table ID"
// @Param request body dto.CreateTableDTO true "table request"
// @Success 200
// @Failure 400,404,500 {object} Error
// @Router /tables/{id} [put]
// @Security ApiKeyAuth
func (h *TableHandler) UpdateTable(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var table entity.Table
    err := json.NewDecoder(r.Body).Decode(&table)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    table.ID = uint(id)
    err = h.TableDB.Update(&table)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// DeleteTable godoc
// @Summary Delete a table by ID
// @Description Delete a table by ID
// @Tags Tables
// @Accept  json
// @Produce  json
// @Param id path int true "Table ID"
// @Success 200
// @Failure 400,404,500 {object} Error
// @Router /tables/{id} [delete]
// @Security ApiKeyAuth
func (h *TableHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    err = h.TableDB.Delete(id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// GetTables godoc
// @Summary Get all tables
// @Description Get all tables
// @Tags Tables
// @Accept  json
// @Produce  json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Success 200 {object} []dto.CreateTableDTO
// @Failure 400,500 {object} Error
// @Router /tables [get]
// @Security ApiKeyAuth
func (h *TableHandler) GetTables(w http.ResponseWriter, r *http.Request) {
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

    tables, err := h.TableDB.FindAll(page, limit, sort)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(tables)
}
