# github

The GitHub module provides a wrapper around the GitHub API, allowing you to
interact with repositories, pull requests, commits, workflows, and user information.

## Creating a Client

### github.client(token)

Creates a new GitHub client. If no token is provided, an unauthenticated client
is created with limited API access. For most operations, you'll need to provide
a personal access token.

```risor
// Unauthenticated client (limited access)
client := github.client()

// Authenticated client
client := github.client("your-github-token")
```

## Repository Operations

### client.get_repo(owner, repo)

Gets information about a repository.

```risor
client := github.client("your-token")
repo := client.get_repo("octocat", "Hello-World")
print(repo.name)
print(repo.description)
```

### client.list_repos(username, options?)

Lists repositories for a user.

```risor
client := github.client("your-token")
repos := client.list_repos("octocat", {
    per_page: 10,
    page: 1
})
for _, repo := range repos {
    print(repo.name)
}
```

### client.list_user_repos(username, options?)

Lists repositories for a user with additional filtering options.

```risor
client := github.client("your-token")
repos := client.list_user_repos("octocat", {
    type: "public",
    sort: "updated",
    direction: "desc"
})
```

### client.list_org_repos(org, options?)

Lists repositories for an organization.

```risor
client := github.client("your-token")
repos := client.list_org_repos("github", {
    type: "public",
    sort: "updated"
})
```

### client.list_repo_contents(owner, repo, path, options?)

Lists contents of a repository directory.

```risor
client := github.client("your-token")
contents := client.list_repo_contents("octocat", "Hello-World", "", {
    ref: "main"
})
for _, item := range contents {
    print(item.name, item.type)
}
```

### client.get_repo_content(owner, repo, path, options?)

Gets the content of a file in a repository.

```risor
client := github.client("your-token")
content := client.get_repo_content("octocat", "Hello-World", "README.md", {
    ref: "main"
})
print(content.name)
print(content.content)
```

## Pull Request Operations

### client.list_pull_requests(owner, repo, options?)

Lists pull requests for a repository.

```risor
client := github.client("your-token")
prs := client.list_pull_requests("octocat", "Hello-World", {
    state: "open",
    sort: "updated",
    direction: "desc"
})
for _, pr := range prs {
    print(pr.number, pr.title)
}
```

### client.get_pull_request(owner, repo, number)

Gets a specific pull request.

```risor
client := github.client("your-token")
pr := client.get_pull_request("octocat", "Hello-World", 1)
print(pr.title)
print(pr.state)
```

### client.create_pull_request(owner, repo, options)

Creates a new pull request.

```risor
client := github.client("your-token")
pr := client.create_pull_request("octocat", "Hello-World", {
    title: "New feature",
    body: "This PR adds a new feature",
    head: "feature-branch",
    base: "main"
})
print(pr.number)
```

### client.list_pull_request_files(owner, repo, number)

Lists files in a pull request.

```risor
client := github.client("your-token")
files := client.list_pull_request_files("octocat", "Hello-World", 1)
for _, file := range files {
    print(file.filename)
    print(file.status)
}
```

### client.list_pull_request_commits(owner, repo, number)

Lists commits in a pull request.

```risor
client := github.client("your-token")
commits := client.list_pull_request_commits("octocat", "Hello-World", 1)
for _, commit := range commits {
    print(commit.sha)
    print(commit.commit.message)
}
```

## Commit Operations

### client.list_commits(owner, repo, options?)

Lists commits for a repository.

```risor
client := github.client("your-token")
commits := client.list_commits("octocat", "Hello-World", {
    sha: "main",
    path: "README.md",
    author: "octocat"
})
for _, commit := range commits {
    print(commit.sha)
    print(commit.commit.message)
}
```

### client.get_commit(owner, repo, sha)

Gets a specific commit.

```risor
client := github.client("your-token")
commit := client.get_commit("octocat", "Hello-World", "abc123")
print(commit.commit.message)
print(commit.stats.additions)
print(commit.stats.deletions)
```

## GitHub Actions Operations

### client.list_workflow_runs(owner, repo, options?)

Lists workflow runs for a repository.

```risor
client := github.client("your-token")
runs := client.list_workflow_runs("octocat", "Hello-World", {
    actor: "octocat",
    branch: "main",
    event: "push",
    status: "success"
})
for _, run := range runs {
    print(run.id)
    print(run.status)
    print(run.conclusion)
}
```

### client.get_workflow_run(owner, repo, runId)

Gets a specific workflow run.

```risor
client := github.client("your-token")
run := client.get_workflow_run("octocat", "Hello-World", 123456)
print(run.status)
print(run.conclusion)
```

### client.list_workflows(owner, repo)

Lists workflows for a repository.

```risor
client := github.client("your-token")
workflows := client.list_workflows("octocat", "Hello-World")
for _, workflow := range workflows {
    print(workflow.id)
    print(workflow.name)
    print(workflow.path)
}
```

### client.get_workflow(owner, repo, workflowId)

Gets a specific workflow.

```risor
client := github.client("your-token")
workflow := client.get_workflow("octocat", "Hello-World", "ci.yml")
print(workflow.name)
print(workflow.path)
```

## User Operations

### client.get_user(username)

Gets information about a user.

```risor
client := github.client("your-token")
user := client.get_user("octocat")
print(user.name)
print(user.company)
print(user.location)
```

## Authentication

For most operations, you'll need to provide a GitHub personal access token. You can create one at https://github.com/settings/tokens

The token should have appropriate permissions for the operations you want to perform:
- `repo` - Full control of private repositories
- `public_repo` - Access to public repositories  
- `user` - Access to user profile information
- `read:org` - Read organization membership

## Error Handling

All functions return errors as Risor objects when API calls fail. Common error scenarios include:
- Rate limiting
- Authentication failures
- Resource not found
- Insufficient permissions

```risor
client := github.client("your-token")
repo := client.get_repo("owner", "nonexistent-repo")
if error(repo) {
    print("Error:", repo)
}
```

## Examples

### List Recent Commits

```risor
client := github.client("your-token")
commits := client.list_commits("octocat", "Hello-World")
for _, commit := range commits[:5] {
    print(commit.commit.author.date, commit.commit.message)
}
```

### Check PR Status

```risor
client := github.client("your-token")
pr := client.get_pull_request("octocat", "Hello-World", 1)
if pr.state == "open" {
    print("PR is open")
    files := client.list_pull_request_files("octocat", "Hello-World", 1)
    print("Files changed:", len(files))
}
```

### Monitor Workflow Runs

```risor
client := github.client("your-token")
runs := client.list_workflow_runs("octocat", "Hello-World")
failed_runs := []
for _, run := range runs {
    if run.conclusion == "failure" {
        failed_runs.append(run)
    }
}
print("Failed runs:", len(failed_runs))
```