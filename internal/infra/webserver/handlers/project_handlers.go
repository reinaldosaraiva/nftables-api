// project_handlers.go
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
type ProjectHandler struct {
	ProjectDB database.ProjectInterface
	TenantDB  database.TenantInterface
}

func NewProjectHandler(projectDB database.ProjectInterface, tenantDB database.TenantInterface) *ProjectHandler {
	return &ProjectHandler{
		ProjectDB: projectDB,
		TenantDB:  tenantDB,
	}
}

// Create Project godoc
// @Summary Create a new Project
// @Description Create Projects
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param request body dto.CreateProjectDTO true "project request"
// @Success 201
// @Failure 400,500 {object} Error
// @Router /projects [post]
// @Security ApiKeyAuth
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
    var projectDTO dto.CreateProjectDTO
    err := json.NewDecoder(r.Body).Decode(&projectDTO)
	fmt.Println(projectDTO)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    tenant, err := h.TenantDB.FindByName(projectDTO.TenantName)
    if err != nil {
        // Se o Tenant não for encontrado, retorna 404
        if errors.Is(err, gorm.ErrRecordNotFound) {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        // Outros erros do banco de dados
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    project := &entity.Project{
        Name:     projectDTO.Name,
        TenantID: tenant.ID,
    }

    err = h.ProjectDB.Create(project)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}
// Get a project by ID godoc
// @Summary Get a project by ID
// @Description Get a project by ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path int true "project ID"
// @Success 200 {object} dto.CreateProjectDTO
// @Failure 400,404
// @Router /projects/{id} [get]
// @Security ApiKeyAuth
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	project, err := h.ProjectDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

// GetProjectsWithFilters godoc
// @Summary Get projects with optional filters
// @Description Get projects filtered by ID or name
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id query int false "Project ID"
// @Param name query string false "Project Name"
// @Success 200 {array} dto.CreateProjectDTO 
// @Failure 400 "Invalid parameter format"
// @Failure 404 "Project not found"
// @Router /projects/filter [get]
// @Security ApiKeyAuth
func (h *ProjectHandler) GetProjectsWithFilters(w http.ResponseWriter, r *http.Request) {
    queryValues := r.URL.Query()
    name := queryValues.Get("name")
    idStr := queryValues.Get("id")

    var project *dto.DetailsProjectDTO
    var err error

    if idStr != "" {
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        project, err = h.ProjectDB.FindByID(uint64(id))
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        
    } else if name != "" {
        project, err = h.ProjectDB.FindByName(name)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        
    } else {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(project)
}

// Update a project by ID godoc
// @Summary Update a project by ID
// @Description Update a project by ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path int true "project id"
// @Param request body dto.CreateProjectDTO true "project request"
// @Success 200
// @Failure 400,404,500 {object} Error
// @Router /projects/{id} [put]
// @Security ApiKeyAuth
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var project entity.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	project.ID = uint(id)
	err = h.ProjectDB.Update(&project)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete Project by ID godoc
// @Summary Delete Project by ID
// @Description Delete Project by ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path int true "project id"
// @Success 200
// @Failure 400,404,500 {object} Error
// @Router /projects/{id} [delete]
// @Security ApiKeyAuth
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)

	err = h.ProjectDB.Delete(id)
	fmt.Println(err)
	if err != nil {
			
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
		w.WriteHeader(http.StatusInternalServerError)
		}
        return
	}

	w.WriteHeader(http.StatusOK)
}

// Get all Projects godoc
// @Summary Get all Projects
// @Description Get all Projects
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Success 200 {object} []dto.CreateProjectDTO
// @Failure 400,500 {object} Error
// @Router /projects [get]
// @Security ApiKeyAuth
func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
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

	projects, err := h.ProjectDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projects)
}
