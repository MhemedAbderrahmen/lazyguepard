package state

  import (
      "github.com/MhemedAbderrahmen/lazyguepard/pkg/models"
  )

  // AppState holds the current application state
  type AppState struct {
      Deployments       []models.Deployment
      SelectedDeployment *models.Deployment
      Branches          []models.Branch
      Commits           []models.Commit
  }

  // NewAppState creates a new application state
  func NewAppState() *AppState {
      return &AppState{
          Deployments: []models.Deployment{},
          Branches:    []models.Branch{},
          Commits:     []models.Commit{},
      }
  }
