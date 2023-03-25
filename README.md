Provides a read-only restful interface for accessing data from Git repositories (local to the server).

Modeled off the GitHub API for compatibility (see https://docs.github.com/en/rest).

# Configuration
| Environment Variable | Default value    | Description              |
|----------------------|------------------|--------------------------|
| `LISTEN`             | `127.0.0.1:8080` | IP and port to listen on |

# Notes
- All returned urls will be the same as the requested host header.
- If you are using https set `X-Forwarded-Proto: https` upstream on your reverse proxy.

# Indexes
```
GET /repos
[
    {
        "description": "",
        "name": "test",
        "url": "http://127.0.0.1:8080/repos/local/self"
    }
]
```

```
GET http://127.0.0.1:8080/repos/{namespace}/{repo_name}/
GET http://127.0.0.1:8080/repos/local/self
{
    "description": "",
    "indexes": {
        "branches": "http://127.0.0.1:8080/repos/local/self/branches",
        "tags": "http://127.0.0.1:8080/repos/local/self/tags"
    },
    "name": "test",
    "url": "http://127.0.0.1:8080/repos/test"
}
```

# Branches
Retrieves a list of branches
```
GET http://127.0.0.1:8080/repos/{namespace}/{repo_name}/branches
GET http://127.0.0.1:8080/repos/local/self/branches
[
    {
        "commit": {
            "sha": "426c438cbcbae58bba8ca8fe5fa102d1c9628890",
            "url": "http://127.0.0.1:8080/repos/local/self/git/commits/426c438cbcbae58bba8ca8fe5fa102d1c9628890"
        },
        "name": "master",
        "url": "http://127.0.0.1:8080/repos/local/self/branches/master"
    }
]
```

Retrieves a branch
```
GET http://127.0.0.1:8080/repos/local/{namespace}/{repo_name}/branches/{branch_name}
GET http://127.0.0.1:8080/repos/local/self/branches/master
{
    "commit": {
        "author": {
            "date": "2023-03-25T06:23:51Z",
            "email": "claytonsingh@gmail.com",
            "name": "Clayton Singh"
        },
        "committer": {
            "date": "2023-03-25T06:23:51Z",
            "email": "claytonsingh@gmail.com",
            "name": "Clayton Singh"
        },
        "message": "Initial API\n",
        "parents": [
            {
                "sha": "47278128e45e123f5dcf5e088259884ad81229ab",
                "url": "http://127.0.0.1:8080/repos/local/self/git/commits/47278128e45e123f5dcf5e088259884ad81229ab"
            }
        ],
        "sha": "426c438cbcbae58bba8ca8fe5fa102d1c9628890",
        "tree": {
            "sha": "056392b008ba5142234b9c900d2f768f5122790e",
            "url": "http://127.0.0.1:8080/repos/local/self/git/trees/056392b008ba5142234b9c900d2f768f5122790e"
        },
        "url": "http://127.0.0.1:8080/repos/local/self/git/commits/426c438cbcbae58bba8ca8fe5fa102d1c9628890"
    },
    "name": "master",
    "url": "http://127.0.0.1:8080/repos/local/self/branches/master"
}
```

# Tags
Retrieves a list of tags
```
GET http://127.0.0.1:8080/repos/local/{namespace}/{repo_name}/tags
GET http://127.0.0.1:8080/repos/local/self/tags
[
    {
        "commit": {
            "sha": "47278128e45e123f5dcf5e088259884ad81229ab",
            "url": "http://127.0.0.1:8080/repos/local/self/git/commits/47278128e45e123f5dcf5e088259884ad81229ab"
        },
        "message": "",
        "sha": null,
        "tag": "ABC",
        "tagger": null,
        "url": "http://127.0.0.1:8080/repos/local/self/tags/ABC"
    },
    {
        "commit": {
            "sha": "47278128e45e123f5dcf5e088259884ad81229ab",
            "url": "http://127.0.0.1:8080/repos/local/self/git/commits/47278128e45e123f5dcf5e088259884ad81229ab"
        },
        "message": "Tag message\n",
        "sha": "b9f62cc7e8a12e8ae08157984dbebb21d7875140",
        "tag": "DEF",
        "tagger": {
            "date": "2023-03-24T22:19:02Z",
            "email": "claytonsingh@gmail.com",
            "name": "Clayton Singh"
        },
        "url": "http://127.0.0.1:8080/repos/local/self/tags/DEF"
    }
]
```

Retrieve a tag
```
GET http://127.0.0.1:8080/repos/local/{namespace}/{repo_name}/tags/{tag_name}
GET http://127.0.0.1:8080/repos/local/self/tags/ABC
{
    "commit": {
        "sha": "47278128e45e123f5dcf5e088259884ad81229ab",
        "url": "http://127.0.0.1:8080/repos/local/self/git/commits/47278128e45e123f5dcf5e088259884ad81229ab"
    },
    "message": "",
    "sha": null,
    "tag": "ABC",
    "tagger": null,
    "url": "http://127.0.0.1:8080/repos/local/self/tags/ABC"
}
```

# Commits
Retrieve a commit
```
GET http://127.0.0.1:8080/repos/local/{namespace}/{repo_name}/commits/{commit_sha}
GET http://127.0.0.1:8080/repos/local/self/git/commits/426c438cbcbae58bba8ca8fe5fa102d1c9628890
{
    "author": {
        "date": "2023-03-25T06:23:51Z",
        "email": "claytonsingh@gmail.com",
        "name": "Clayton Singh"
    },
    "committer": {
        "date": "2023-03-25T06:23:51Z",
        "email": "claytonsingh@gmail.com",
        "name": "Clayton Singh"
    },
    "message": "Initial API\n",
    "parents": [
        {
            "sha": "47278128e45e123f5dcf5e088259884ad81229ab",
            "url": "http://127.0.0.1:8080/repos/local/self/git/commits/47278128e45e123f5dcf5e088259884ad81229ab"
        }
    ],
    "sha": "426c438cbcbae58bba8ca8fe5fa102d1c9628890",
    "tree": {
        "sha": "056392b008ba5142234b9c900d2f768f5122790e",
        "url": "http://127.0.0.1:8080/repos/local/self/git/trees/056392b008ba5142234b9c900d2f768f5122790e"
    },
    "url": "http://127.0.0.1:8080/repos/local/self/git/commits/426c438cbcbae58bba8ca8fe5fa102d1c9628890"
}
```

# Trees
Retrieves a tree

When in recursive mode directories are no longer returned and all files are shown.
```
optional: ?recursive=yes (default=no, non-recursive, accepts=0, false, no, 1, true, yes)
```

```
GET http://127.0.0.1:8080/repos/local/{namespace}/{repo_name}/git/trees/{tree_sha}
GET http://127.0.0.1:8080/repos/local/self/git/trees/056392b008ba5142234b9c900d2f768f5122790e
{
    "sha": "056392b008ba5142234b9c900d2f768f5122790e",
    "tree": [
        {
            "mode": "0100644",
            "path": "formatters.go",
            "raw_url": "http://127.0.0.1:8080/repos/local/self/raw/056392b008ba5142234b9c900d2f768f5122790e/formatters.go",
            "sha": "574f907b46e67c5c0453f04adaed0461a6556ce0",
            "size": 5134,
            "type": "blob",
            "url": "http://127.0.0.1:8080/repos/local/self/git/blobs/574f907b46e67c5c0453f04adaed0461a6556ce0"
        },
        {
            "mode": "0100644",
            "path": "go.mod",
            "raw_url": "http://127.0.0.1:8080/repos/local/self/raw/056392b008ba5142234b9c900d2f768f5122790e/go.mod",
            "sha": "b9414b519819f2512cbed8cf35acf6c5dfe88d83",
            "size": 2067,
            "type": "blob",
            "url": "http://127.0.0.1:8080/repos/local/self/git/blobs/b9414b519819f2512cbed8cf35acf6c5dfe88d83"
        },
        {
            "mode": "0100644",
            "path": "go.sum",
            "raw_url": "http://127.0.0.1:8080/repos/local/self/raw/056392b008ba5142234b9c900d2f768f5122790e/go.sum",
            "sha": "a1ea7e2aca4e462f6eca5cfddb90b9adff2371ff",
            "size": 18889,
            "type": "blob",
            "url": "http://127.0.0.1:8080/repos/local/self/git/blobs/a1ea7e2aca4e462f6eca5cfddb90b9adff2371ff"
        },
        {
            "mode": "0100644",
            "path": "main.go",
            "raw_url": "http://127.0.0.1:8080/repos/local/self/raw/056392b008ba5142234b9c900d2f768f5122790e/main.go",
            "sha": "747e005002e8b3e6e78c9c93f65bb6aa05abac57",
            "size": 6939,
            "type": "blob",
            "url": "http://127.0.0.1:8080/repos/local/self/git/blobs/747e005002e8b3e6e78c9c93f65bb6aa05abac57"
        }
    ],
    "url": "http://127.0.0.1:8080/repos/local/self/git/trees/056392b008ba5142234b9c900d2f768f5122790e"
}
```

# Blobs
Retrieves a blob with metadata with content base64 encoded.
```
GET http://127.0.0.1:8080/repos/local/{namespace}/{repo_name}/git/blobs/{blob_sha}
GET http://127.0.0.1:8080/repos/local/self/git/blobs/b9414b519819f2512cbed8cf35acf6c5dfe88d83
{
    "content": "...",
    "sha": "b9414b519819f2512cbed8cf35acf6c5dfe88d83",
    "size": 2067,
    "url": "http://127.0.0.1:8080/repos/local/self/git/blobs/b9414b519819f2512cbed8cf35acf6c5dfe88d83"
}
```

Retrieves a blobs content from a tree and file path.
```
GET http://127.0.0.1:8080/repos/local/{namespace}/{repo_name}/raw/{tree_sha}/{file_path}
GET http://127.0.0.1:8080/repos/local/self/raw/056392b008ba5142234b9c900d2f768f5122790e/go.mod
```
