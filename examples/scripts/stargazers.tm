
// Returns a list of stargazers for a Github repository
func get_stargazers(owner='golang', repo='go') {
    url := 'https://api.github.com/repos/{owner}/{repo}/stargazers'
    return fetch(url, {
        headers: {
            'Accept': 'application/vnd.github+json',
            'X-GitHub-Api-Version': '2022-11-28',
        },
    }).json()
}

// Print the login and url for stargazers of github.com/golang/go
get_stargazers().each(func(item) {
    print('login: {item["login"]} url: {item["url"]}')
})
