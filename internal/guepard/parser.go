package guepard

import (
	"strings"

	"github.com/MhemedAbderrahmen/lazyguepard/pkg/models"
)

// ParseDeployments parses the output from 'guepard list deployments'
func ParseDeployments(output string) []models.Deployment {
	deployments := []models.Deployment{}
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "ID") || strings.HasPrefix(line, "--") {
			continue
		}

		// For now, create a simple deployment from the raw line
		// We'll parse this more thoroughly later
		deployment := models.Deployment{
			Name: line,
		}
		deployments = append(deployments, deployment)
	}

	return deployments
}

// ParseBranches parses the output from 'guepard list branches'
func ParseBranches(output string) []models.Branch {
	branches := []models.Branch{}
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "ID") || strings.HasPrefix(line, "--") {
			continue
		}

		// Check if this is the current branch (marked with *)
		isCurrent := strings.HasPrefix(line, "*")
		if isCurrent {
			line = strings.TrimPrefix(line, "*")
			line = strings.TrimSpace(line)
		}

		branch := models.Branch{
			Name:      line,
			IsCurrent: isCurrent,
		}
		branches = append(branches, branch)
	}

	return branches
}
