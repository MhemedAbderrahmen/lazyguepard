package guepard

  import (
      "bytes"
      "fmt"
      "os/exec"
      "strings"
  )

  // Client wraps the Guepard CLI
  type Client struct {
      guepardPath string
  }

  // NewClient creates a new Guepard CLI client
  func NewClient() (*Client, error) {
      // Check if guepard is in PATH
      path, err := exec.LookPath("guepard")
      if err != nil {
          return nil, fmt.Errorf("guepard CLI not found in PATH: %w", err)
      }

      return &Client{
          guepardPath: path,
      }, nil
  }

  // Execute runs a guepard command and returns the output
  func (c *Client) Execute(args ...string) (string, error) {
      cmd := exec.Command(c.guepardPath, args...)

      var stdout, stderr bytes.Buffer
      cmd.Stdout = &stdout
      cmd.Stderr = &stderr

      err := cmd.Run()
      if err != nil {
          return "", fmt.Errorf("command failed: %w\nstderr: %s", err, stderr.String())
      }

      return strings.TrimSpace(stdout.String()), nil
  }

  // ListDeployments lists all deployments
  func (c *Client) ListDeployments() (string, error) {
      return c.Execute("list", "deployments")
  }

  // ListBranches lists branches for a deployment
  func (c *Client) ListBranches(deploymentID string) (string, error) {
      return c.Execute("list", "branches", "-x", deploymentID)
  }

  // ListCommits lists commits for a deployment
  func (c *Client) ListCommits(deploymentID string) (string, error) {
      return c.Execute("list", "commits", "-x", deploymentID)
  }

