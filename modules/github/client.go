package github

import (
	"context"
	"fmt"
	"encoding/json"

	"github.com/risor-io/risor/object"
	"github.com/google/go-github/v73/github"
)

const CLIENT object.Type = "github.client"

// Client wraps the GitHub client and provides Risor integration
type Client struct {
	base
	value *github.Client
}

func (c *Client) Type() object.Type {
	return CLIENT
}

func (c *Client) Inspect() string {
	return "github.client()"
}

func (c *Client) Interface() interface{} {
	return c.value
}

func (c *Client) Equals(other object.Object) object.Object {
	if c == other {
		return object.True
	}
	return object.False
}

func (c *Client) GetAttr(name string) (object.Object, bool) {
	switch name {
	// Repository operations
	case "get_repo":
		return object.NewBuiltin("get_repo", c.GetRepo), true
	case "list_repos":
		return object.NewBuiltin("list_repos", c.ListRepos), true
	case "list_repo_contents":
		return object.NewBuiltin("list_repo_contents", c.ListRepoContents), true
	case "get_repo_content":
		return object.NewBuiltin("get_repo_content", c.GetRepoContent), true
	
	// Pull request operations
	case "list_pull_requests":
		return object.NewBuiltin("list_pull_requests", c.ListPullRequests), true
	case "get_pull_request":
		return object.NewBuiltin("get_pull_request", c.GetPullRequest), true
	case "create_pull_request":
		return object.NewBuiltin("create_pull_request", c.CreatePullRequest), true
	case "list_pull_request_files":
		return object.NewBuiltin("list_pull_request_files", c.ListPullRequestFiles), true
	case "list_pull_request_commits":
		return object.NewBuiltin("list_pull_request_commits", c.ListPullRequestCommits), true
	
	// Commit operations
	case "list_commits":
		return object.NewBuiltin("list_commits", c.ListCommits), true
	case "get_commit":
		return object.NewBuiltin("get_commit", c.GetCommit), true
	
	// GitHub Actions operations
	case "list_workflow_runs":
		return object.NewBuiltin("list_workflow_runs", c.ListWorkflowRuns), true
	case "get_workflow_run":
		return object.NewBuiltin("get_workflow_run", c.GetWorkflowRun), true
	case "list_workflows":
		return object.NewBuiltin("list_workflows", c.ListWorkflows), true
	case "get_workflow":
		return object.NewBuiltin("get_workflow", c.GetWorkflow), true
	
	// User/Organization operations
	case "get_user":
		return object.NewBuiltin("get_user", c.GetUser), true
	case "list_user_repos":
		return object.NewBuiltin("list_user_repos", c.ListUserRepos), true
	case "list_org_repos":
		return object.NewBuiltin("list_org_repos", c.ListOrgRepos), true
	}
	return nil, false
}

func (c *Client) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("type error: cannot set %q on %s object", name, CLIENT)
}

func (c *Client) IsTruthy() bool {
	return true
}

func (c *Client) Cost() int {
	return 0
}

// Helper function to convert GitHub API objects to Risor objects
func asMap(value interface{}) object.Object {
	data, err := json.Marshal(value)
	if err != nil {
		return object.NewError(err)
	}
	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return object.NewError(err)
	}
	return object.FromGoType(dataMap)
}

// Repository operations

// GetRepo gets information about a repository
func (c *Client) GetRepo(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("get_repo", 2, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	repository, _, apiErr := c.value.Repositories.Get(ctx, owner, repo)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	return asMap(repository)
}

// ListRepos lists repositories for a user or organization
func (c *Client) ListRepos(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 || len(args) > 2 {
		return object.NewArgsRangeError("list_repos", 1, 2, len(args))
	}
	
	username, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	var listOpts github.ListOptions
	if len(args) == 2 {
		opts, err := object.AsMap(args[1])
		if err != nil {
			return err
		}
		
		if perPage := opts.Get("per_page"); perPage != object.Nil {
			perPageInt, err := object.AsInt(perPage)
			if err != nil {
				return err
			}
			listOpts.PerPage = int(perPageInt)
		}
		if page := opts.Get("page"); page != object.Nil {
			pageInt, err := object.AsInt(page)
			if err != nil {
				return err
			}
			listOpts.Page = int(pageInt)
		}
	}
	
	repos, _, apiErr := c.value.Repositories.List(ctx, username, &github.RepositoryListOptions{
		ListOptions: listOpts,
	})
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	result := make([]object.Object, len(repos))
	for i, repo := range repos {
		result[i] = asMap(repo)
	}
	
	return object.NewList(result)
}

// ListRepoContents lists contents of a repository directory
func (c *Client) ListRepoContents(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 3 || len(args) > 4 {
		return object.NewArgsRangeError("list_repo_contents", 3, 4, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	path, err := object.AsString(args[2])
	if err != nil {
		return err
	}
	
	var opts *github.RepositoryContentGetOptions
	if len(args) == 4 {
		options, err := object.AsMap(args[3])
		if err != nil {
			return err
		}
		
		opts = &github.RepositoryContentGetOptions{}
		if ref := options.Get("ref"); ref != object.Nil {
			refStr, err := object.AsString(ref)
			if err != nil {
				return err
			}
			opts.Ref = refStr
		}
	}
	
	_, contents, _, apiErr := c.value.Repositories.GetContents(ctx, owner, repo, path, opts)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	result := make([]object.Object, len(contents))
	for i, content := range contents {
		result[i] = asMap(content)
	}
	
	return object.NewList(result)
}

// GetRepoContent gets the content of a file in a repository
func (c *Client) GetRepoContent(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 3 || len(args) > 4 {
		return object.NewArgsRangeError("get_repo_content", 3, 4, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	path, err := object.AsString(args[2])
	if err != nil {
		return err
	}
	
	var opts *github.RepositoryContentGetOptions
	if len(args) == 4 {
		options, err := object.AsMap(args[3])
		if err != nil {
			return err
		}
		
		opts = &github.RepositoryContentGetOptions{}
		if ref := options.Get("ref"); ref != object.Nil {
			refStr, err := object.AsString(ref)
			if err != nil {
				return err
			}
			opts.Ref = refStr
		}
	}
	
	content, _, _, apiErr := c.value.Repositories.GetContents(ctx, owner, repo, path, opts)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	return asMap(content)
}

// Pull request operations

// ListPullRequests lists pull requests for a repository
func (c *Client) ListPullRequests(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 || len(args) > 3 {
		return object.NewArgsRangeError("list_pull_requests", 2, 3, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	var opts *github.PullRequestListOptions
	if len(args) == 3 {
		options, err := object.AsMap(args[2])
		if err != nil {
			return err
		}
		
		opts = &github.PullRequestListOptions{}
		if state := options.Get("state"); state != object.Nil {
			stateStr, err := object.AsString(state)
			if err != nil {
				return err
			}
			opts.State = stateStr
		}
		if sort := options.Get("sort"); sort != object.Nil {
			sortStr, err := object.AsString(sort)
			if err != nil {
				return err
			}
			opts.Sort = sortStr
		}
		if direction := options.Get("direction"); direction != object.Nil {
			directionStr, err := object.AsString(direction)
			if err != nil {
				return err
			}
			opts.Direction = directionStr
		}
	}
	
	prs, _, apiErr := c.value.PullRequests.List(ctx, owner, repo, opts)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	result := make([]object.Object, len(prs))
	for i, pr := range prs {
		result[i] = asMap(pr)
	}
	
	return object.NewList(result)
}

// GetPullRequest gets a specific pull request
func (c *Client) GetPullRequest(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 3 {
		return object.NewArgsError("get_pull_request", 3, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	prNumber, err := object.AsInt(args[2])
	if err != nil {
		return err
	}
	
	pr, _, apiErr := c.value.PullRequests.Get(ctx, owner, repo, int(prNumber))
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	return asMap(pr)
}

// CreatePullRequest creates a new pull request
func (c *Client) CreatePullRequest(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 3 {
		return object.NewArgsError("create_pull_request", 3, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	options, err := object.AsMap(args[2])
	if err != nil {
		return err
	}
	
	newPR := &github.NewPullRequest{}
	
	if title := options.Get("title"); title != object.Nil {
		titleStr, err := object.AsString(title)
		if err != nil {
			return err
		}
		newPR.Title = &titleStr
	}
	
	if body := options.Get("body"); body != object.Nil {
		bodyStr, err := object.AsString(body)
		if err != nil {
			return err
		}
		newPR.Body = &bodyStr
	}
	
	if head := options.Get("head"); head != object.Nil {
		headStr, err := object.AsString(head)
		if err != nil {
			return err
		}
		newPR.Head = &headStr
	}
	
	if base := options.Get("base"); base != object.Nil {
		baseStr, err := object.AsString(base)
		if err != nil {
			return err
		}
		newPR.Base = &baseStr
	}
	
	pr, _, apiErr := c.value.PullRequests.Create(ctx, owner, repo, newPR)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	return asMap(pr)
}

// ListPullRequestFiles lists files in a pull request
func (c *Client) ListPullRequestFiles(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 3 {
		return object.NewArgsError("list_pull_request_files", 3, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	prNumber, err := object.AsInt(args[2])
	if err != nil {
		return err
	}
	
	files, _, apiErr := c.value.PullRequests.ListFiles(ctx, owner, repo, int(prNumber), nil)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	result := make([]object.Object, len(files))
	for i, file := range files {
		result[i] = asMap(file)
	}
	
	return object.NewList(result)
}

// ListPullRequestCommits lists commits in a pull request
func (c *Client) ListPullRequestCommits(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 3 {
		return object.NewArgsError("list_pull_request_commits", 3, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	prNumber, err := object.AsInt(args[2])
	if err != nil {
		return err
	}
	
	commits, _, apiErr := c.value.PullRequests.ListCommits(ctx, owner, repo, int(prNumber), nil)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	result := make([]object.Object, len(commits))
	for i, commit := range commits {
		result[i] = asMap(commit)
	}
	
	return object.NewList(result)
}

// Commit operations

// ListCommits lists commits for a repository
func (c *Client) ListCommits(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 || len(args) > 3 {
		return object.NewArgsRangeError("list_commits", 2, 3, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	var opts *github.CommitsListOptions
	if len(args) == 3 {
		options, err := object.AsMap(args[2])
		if err != nil {
			return err
		}
		
		opts = &github.CommitsListOptions{}
		if sha := options.Get("sha"); sha != object.Nil {
			shaStr, err := object.AsString(sha)
			if err != nil {
				return err
			}
			opts.SHA = shaStr
		}
		if path := options.Get("path"); path != object.Nil {
			pathStr, err := object.AsString(path)
			if err != nil {
				return err
			}
			opts.Path = pathStr
		}
		if author := options.Get("author"); author != object.Nil {
			authorStr, err := object.AsString(author)
			if err != nil {
				return err
			}
			opts.Author = authorStr
		}
	}
	
	commits, _, apiErr := c.value.Repositories.ListCommits(ctx, owner, repo, opts)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	result := make([]object.Object, len(commits))
	for i, commit := range commits {
		result[i] = asMap(commit)
	}
	
	return object.NewList(result)
}

// GetCommit gets a specific commit
func (c *Client) GetCommit(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 3 {
		return object.NewArgsError("get_commit", 3, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	sha, err := object.AsString(args[2])
	if err != nil {
		return err
	}
	
	commit, _, apiErr := c.value.Repositories.GetCommit(ctx, owner, repo, sha, nil)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	return asMap(commit)
}

// GitHub Actions operations

// ListWorkflowRuns lists workflow runs for a repository
func (c *Client) ListWorkflowRuns(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 || len(args) > 3 {
		return object.NewArgsRangeError("list_workflow_runs", 2, 3, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	var opts *github.ListWorkflowRunsOptions
	if len(args) == 3 {
		options, err := object.AsMap(args[2])
		if err != nil {
			return err
		}
		
		opts = &github.ListWorkflowRunsOptions{}
		if actor := options.Get("actor"); actor != object.Nil {
			actorStr, err := object.AsString(actor)
			if err != nil {
				return err
			}
			opts.Actor = actorStr
		}
		if branch := options.Get("branch"); branch != object.Nil {
			branchStr, err := object.AsString(branch)
			if err != nil {
				return err
			}
			opts.Branch = branchStr
		}
		if event := options.Get("event"); event != object.Nil {
			eventStr, err := object.AsString(event)
			if err != nil {
				return err
			}
			opts.Event = eventStr
		}
		if status := options.Get("status"); status != object.Nil {
			statusStr, err := object.AsString(status)
			if err != nil {
				return err
			}
			opts.Status = statusStr
		}
	}
	
	runs, _, apiErr := c.value.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opts)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	result := make([]object.Object, len(runs.WorkflowRuns))
	for i, run := range runs.WorkflowRuns {
		result[i] = asMap(run)
	}
	
	return object.NewList(result)
}

// GetWorkflowRun gets a specific workflow run
func (c *Client) GetWorkflowRun(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 3 {
		return object.NewArgsError("get_workflow_run", 3, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	runID, err := object.AsInt(args[2])
	if err != nil {
		return err
	}
	
	run, _, apiErr := c.value.Actions.GetWorkflowRunByID(ctx, owner, repo, runID)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	return asMap(run)
}

// ListWorkflows lists workflows for a repository
func (c *Client) ListWorkflows(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("list_workflows", 2, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	workflows, _, apiErr := c.value.Actions.ListWorkflows(ctx, owner, repo, nil)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	result := make([]object.Object, len(workflows.Workflows))
	for i, workflow := range workflows.Workflows {
		result[i] = asMap(workflow)
	}
	
	return object.NewList(result)
}

// GetWorkflow gets a specific workflow
func (c *Client) GetWorkflow(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 3 {
		return object.NewArgsError("get_workflow", 3, len(args))
	}
	
	owner, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	repo, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	
	workflowID, err := object.AsInt(args[2])
	if err != nil {
		return err
	}
	
	workflow, _, apiErr := c.value.Actions.GetWorkflowByID(ctx, owner, repo, workflowID)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	return asMap(workflow)
}

// User/Organization operations

// GetUser gets information about a user
func (c *Client) GetUser(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("get_user", 1, len(args))
	}
	
	username, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	user, _, apiErr := c.value.Users.Get(ctx, username)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	return asMap(user)
}

// ListUserRepos lists repositories for a user
func (c *Client) ListUserRepos(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 || len(args) > 2 {
		return object.NewArgsRangeError("list_user_repos", 1, 2, len(args))
	}
	
	username, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	var opts *github.RepositoryListOptions
	if len(args) == 2 {
		options, err := object.AsMap(args[1])
		if err != nil {
			return err
		}
		
		opts = &github.RepositoryListOptions{}
		if repoType := options.Get("type"); repoType != object.Nil {
			typeStr, err := object.AsString(repoType)
			if err != nil {
				return err
			}
			opts.Type = typeStr
		}
		if sort := options.Get("sort"); sort != object.Nil {
			sortStr, err := object.AsString(sort)
			if err != nil {
				return err
			}
			opts.Sort = sortStr
		}
		if direction := options.Get("direction"); direction != object.Nil {
			directionStr, err := object.AsString(direction)
			if err != nil {
				return err
			}
			opts.Direction = directionStr
		}
	}
	
	repos, _, apiErr := c.value.Repositories.List(ctx, username, opts)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	result := make([]object.Object, len(repos))
	for i, repo := range repos {
		result[i] = asMap(repo)
	}
	
	return object.NewList(result)
}

// ListOrgRepos lists repositories for an organization
func (c *Client) ListOrgRepos(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 || len(args) > 2 {
		return object.NewArgsRangeError("list_org_repos", 1, 2, len(args))
	}
	
	org, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	
	var opts *github.RepositoryListByOrgOptions
	if len(args) == 2 {
		options, err := object.AsMap(args[1])
		if err != nil {
			return err
		}
		
		opts = &github.RepositoryListByOrgOptions{}
		if repoType := options.Get("type"); repoType != object.Nil {
			typeStr, err := object.AsString(repoType)
			if err != nil {
				return err
			}
			opts.Type = typeStr
		}
		if sort := options.Get("sort"); sort != object.Nil {
			sortStr, err := object.AsString(sort)
			if err != nil {
				return err
			}
			opts.Sort = sortStr
		}
		if direction := options.Get("direction"); direction != object.Nil {
			directionStr, err := object.AsString(direction)
			if err != nil {
				return err
			}
			opts.Direction = directionStr
		}
	}
	
	repos, _, apiErr := c.value.Repositories.ListByOrg(ctx, org, opts)
	if apiErr != nil {
		return object.NewError(apiErr)
	}
	
	result := make([]object.Object, len(repos))
	for i, repo := range repos {
		result[i] = asMap(repo)
	}
	
	return object.NewList(result)
}

// New creates a new GitHub client wrapper
func New(client *github.Client) *Client {
	return &Client{value: client}
}