// project_db_test.go
package database

import (
	"fmt"
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDatabaseAndTenantForProject(t *testing.T) (*gorm.DB, *entity.Tenant) {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to open database: %v", err)
    }
    err = db.AutoMigrate(&entity.Chain{}, &entity.Rule{},&entity.Project{}, &entity.Table{},&entity.Tenant{},&entity.Service{},&entity.NetworkObject{},)
    if err != nil {
        t.Fatalf("Failed to migrate database: %v", err)
    }
    if err != nil {
        t.Fatalf("Failed to migrate database: %v", err)
    }
    tenant := &entity.Tenant{Name: "Tenant for Project Testing"}
    err = db.Create(tenant).Error
    if err != nil {
        t.Fatalf("Failed to create tenant: %v", err)
    }

    return db, tenant
}


func TestCreateProject(t *testing.T) {
    db,tenant := setupDatabaseAndTenantForProject(t)
    tenantDB := NewTenantDB(db)
    projectDB := NewProjectDB(db)
    tenant = &entity.Tenant{Name: "Tenant 1"}
    err := tenantDB.Create(tenant)
    assert.NoError(t, err)
    assert.NotZero(t, tenant.ID)
    project := &entity.Project{Name: "Project 1", TenantID: uint64(tenant.ID)}
    err = projectDB.Create(project)
    assert.NoError(t, err)
    assert.NotZero(t, project.ID)
}

func TestFindProjectByID(t *testing.T) {
    db, tenant := setupDatabaseAndTenantForProject(t)
    projectDB := NewProjectDB(db)

    project := &entity.Project{Name: "Project 1", TenantID: uint64(tenant.ID)}
    err := projectDB.Create(project)
    assert.NoError(t, err)

    foundProject, err := projectDB.FindByID(uint64(project.ID))
    assert.NoError(t, err)
    assert.NotNil(t, foundProject)
    assert.Equal(t, "Project 1", foundProject.Name)
}
func TestFindProjectByName(t *testing.T) {
    db, tenant := setupDatabaseAndTenantForProject(t)
    projectDB := NewProjectDB(db)

    project := &entity.Project{Name: "Project 1", TenantID: uint64(tenant.ID)}
    err := projectDB.Create(project)
    assert.NoError(t, err)

    foundProjects, err := projectDB.FindByName("Project 1")
    assert.NoError(t, err)
    assert.NotEmpty(t, foundProjects)
    // assert.Equal(t, "Project 1", foundProjects)
}

func TestUpdateProject(t *testing.T) {
    db, tenant := setupDatabaseAndTenantForProject(t)
    projectDB := NewProjectDB(db)

    project := &entity.Project{Name: "Project 1", TenantID: uint64(tenant.ID)}
    err := projectDB.Create(project)
    assert.NoError(t, err)

    project.Name = "Updated Project"
    err = projectDB.Update(project)
    assert.NoError(t, err)

    updatedProject, err := projectDB.FindByID(uint64(project.ID))
    assert.NoError(t, err)
    assert.Equal(t, "Updated Project", updatedProject.Name)
}

func TestFindAllProjects(t *testing.T) {
    db, tenant := setupDatabaseAndTenantForProject(t)
    projectDB := NewProjectDB(db)

    for i := 0; i < 10; i++ {
        project := &entity.Project{Name: fmt.Sprintf("Project %d", i), TenantID: uint64(tenant.ID)}
        err := projectDB.Create(project)
        assert.NoError(t, err)
    }

    projects, err := projectDB.FindAll(1, 5, "asc")
    assert.NoError(t, err)
    assert.Len(t, projects, 5)
}

func TestDeleteProject(t *testing.T) {
    db, tenant := setupDatabaseAndTenantForProject(t)
    projectDB := NewProjectDB(db)

    project := &entity.Project{Name: "Project to Delete", TenantID: uint64(tenant.ID)}
    err := projectDB.Create(project)
    assert.NoError(t, err)

    err = projectDB.Delete(uint64(project.ID))
    assert.NoError(t, err)

    _, err = projectDB.FindByID(uint64(project.ID))
    assert.Error(t, err)
}

func TestDeleteProjectCascade(t *testing.T) {
    db, tenant := setupDatabaseAndTenantForProject(t)
    projectDB := NewProjectDB(db)
    tableDB := NewTableDB(db)
    chainDB := NewChainDB(db)

    project := &entity.Project{Name: "Project with Chains", TenantID: uint64(tenant.ID)}
    err := projectDB.Create(project)
    assert.NoError(t, err)

    table := &entity.Table{Name: "mangle", Description: "Mangle table", Type: "mangle"}
    err = tableDB.Create(table)
    assert.NoError(t, err)

    for i := 0; i < 3; i++ {
        chain := &entity.Chain{Name: fmt.Sprintf("Chain %d", i), ProjectID: uint64(project.ID), TableID: uint64(table.ID)}
        err = chainDB.Create(chain)
        assert.NoError(t, err)
    }

    err = projectDB.Delete(uint64(project.ID))
    assert.NoError(t, err)

    _, err = projectDB.FindByID(uint64(project.ID))
    assert.Error(t, err)

    // chains, err := chainDB.FindAll(1, 10, "asc")
    // assert.NoError(t, err)
    // for _, c := range chains {
    //     assert.NotEqual(t, c.ProjectID, project.ID)
    // }
}

