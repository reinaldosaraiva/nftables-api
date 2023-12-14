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

func setupDatabaseForProject(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to open database: %v", err)
    }

    err = db.AutoMigrate(&entity.Project{}, &entity.Tenant{}, &entity.Chain{})
    if err != nil {
        t.Fatalf("Failed to migrate database: %v", err)
    }

    return db
}

func TestCreateProject(t *testing.T) {
    db := setupDatabaseForProject(t)
    projectDB := NewProjectDB(db)
    project := &entity.Project{Name: "Project 1", TenantID: 1} 

    err := projectDB.Create(project)
    assert.NoError(t, err)
    assert.NotZero(t, project.ID)
}

func TestFindProjectByID(t *testing.T) {
    db := setupDatabaseForProject(t)
    projectDB := NewProjectDB(db)

    project := &entity.Project{Name: "Project 1", TenantID: 1}
    err := projectDB.Create(project)
    assert.NoError(t, err)

    foundProject, err := projectDB.FindByID(uint64(project.ID))
    assert.NoError(t, err)
    assert.NotNil(t, foundProject)
    assert.Equal(t, "Project 1", foundProject.Name)
}

func TestUpdateProject(t *testing.T) {
    db := setupDatabaseForProject(t)
    projectDB := NewProjectDB(db)

    project := &entity.Project{Name: "Project 1", TenantID: 1}
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
    db := setupDatabaseForProject(t)
    projectDB := NewProjectDB(db)

    for i := 0; i < 10; i++ {
        project := &entity.Project{Name: fmt.Sprintf("Project %d", i), TenantID: 1}
        err := projectDB.Create(project)
        assert.NoError(t, err)
    }

    projects, err := projectDB.FindAll(1, 5, "asc")
    assert.NoError(t, err)
    assert.Len(t, projects, 5)
}

func TestDeleteProject(t *testing.T) {
    db := setupDatabaseForProject(t)
    projectDB := NewProjectDB(db)

    project := &entity.Project{Name: "Project to Delete", TenantID: 1}
    err := projectDB.Create(project)
    assert.NoError(t, err)

    err = projectDB.Delete(uint64(project.ID))
    assert.NoError(t, err)

    _, err = projectDB.FindByID(uint64(project.ID))
    assert.Error(t, err)
}
