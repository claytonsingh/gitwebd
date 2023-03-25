package main

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Context struct {
	*gin.Context
	repository *Repo
	namespace  string
	repo       string
}

func (this *Context) BuildUri(parts ...string) string {
	return BuildUri(this.Context, parts...)
}

func BuildUri(c *gin.Context, parts ...string) string {
	u := url.URL{
		Scheme: "http",
		Host:   c.Request.Host,
		Path:   strings.Join(parts, "/"),
	}
	if c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		u.Scheme = "https"
	}
	return u.String()
}

func FormatRef(c *gin.Context, m plumbing.Hash, path string) gin.H {
	hash := m.String()
	return gin.H{
		"sha": hash,
		"url": BuildUri(c, "repos", c.Param("namespace"), c.Param("repo"), path, hash),
	}
}

func FormatCommitRef(c *gin.Context, m plumbing.Hash) gin.H {
	return FormatRef(c, m, "git/commits")
}

func FormatCommitFull(c *gin.Context, m *object.Commit) gin.H {
	hash := m.Hash.String()
	p := make([]any, 0)
	m.Parents().ForEach(func(m *object.Commit) error {
		p = append(p, FormatCommitRef(c, m.Hash))
		return nil
	})
	return gin.H{
		"sha":       hash,
		"url":       BuildUri(c, "repos", c.Param("namespace"), c.Param("repo"), "git/commits", hash),
		"message":   m.Message,
		"author":    FormatSignature(c, m.Author),
		"committer": FormatSignature(c, m.Committer),
		"parents":   p,
		"tree":      FormatTreeRef(c, m.TreeHash),
	}
}

func FormatFile(c *gin.Context, f *object.File, m plumbing.Hash) gin.H {
	return gin.H{
		"url": BuildUri(c, "repos", c.Param("namespace"), c.Param("repo"), "git/blobs", f.Hash.String()),
		//"raw_url": BuildUri(c, "repos", c.Param("namespace"), c.Param("repo"), "git/blobs", f.Hash.String(), "raw"),
		"raw_url": BuildUri(c, "repos", c.Param("namespace"), c.Param("repo"), "raw", m.String(), f.Name),
		"sha":     f.Hash.String(),
		"mode":    f.Mode.String(),
		"path":    f.Name,
		"type":    "blob",
		"size":    f.Size,
	}
}

func FormatSignature(c *gin.Context, a object.Signature) gin.H {
	return gin.H{
		"date":  a.When.UTC(),
		"name":  a.Name,
		"email": a.Email,
	}
}

func FormatTreeRef(c *gin.Context, m plumbing.Hash) gin.H {
	return FormatRef(c, m, "git/trees")
}

func FormatBranchRef(c Context, branch *plumbing.Reference) gin.H {
	name := branch.Name().Short()
	return gin.H{
		"name":   name,
		"url":    c.BuildUri("repos", c.namespace, c.repo, "branches", name),
		"commit": FormatCommitRef(c.Context, branch.Hash()),
	}
}

func FormatBranchFull(c Context, branch *plumbing.Reference, commit *object.Commit) gin.H {
	name := branch.Name().Short()
	return gin.H{
		"name":   name,
		"url":    c.BuildUri("repos", c.namespace, c.repo, "branches", name),
		"commit": FormatCommitFull(c.Context, commit),
	}
}

func FormatTagRef(c Context, ref *plumbing.Reference, tag *object.Tag) gin.H {
	if tag == nil {
		name := ref.Name().Short()
		return gin.H{
			"sha":     nil,
			"tag":     name,
			"url":     c.BuildUri("repos", c.namespace, c.repo, "tags", name),
			"commit":  FormatCommitRef(c.Context, ref.Hash()),
			"message": "",
			"tagger":  nil,
		}
	} else {
		name := ref.Name().Short()
		return gin.H{
			"sha":     ref.Hash().String(),
			"tag":     name,
			"url":     c.BuildUri("repos", c.namespace, c.repo, "tags", name),
			"commit":  FormatCommitRef(c.Context, tag.Target),
			"message": tag.Message,
			"tagger":  FormatSignature(c.Context, tag.Tagger),
		}
	}
}

func FormatTagFull(c Context, tag *plumbing.Reference, commit *object.Commit) gin.H {
	name := tag.Name().Short()
	return gin.H{
		"tag":    name,
		"url":    c.BuildUri("repos", c.namespace, c.repo, "tags", name),
		"commit": FormatCommitFull(c.Context, commit),
	}
}

func FormatTree(c Context, tree *object.Tree, recursive bool) gin.H {
	r := make([]any, 0)
	if recursive {
		tree.Files().ForEach(func(f *object.File) error {
			r = append(r, FormatFile(c.Context, f, tree.Hash))
			return nil
		})
	} else {
		for _, entry := range tree.Entries {
			if entry.Mode.IsFile() {
				blob, _ := object.GetBlob(c.repository.Repo.Storer, entry.Hash)
				file := object.NewFile(entry.Name, entry.Mode, blob)
				r = append(r, FormatFile(c.Context, file, tree.Hash))
			} else if entry.Mode == filemode.Dir {
				r = append(r, gin.H{
					"path": entry.Name,
					"mode": entry.Mode,
					"type": "tree",
					"sha":  entry.Hash.String(),
					"url":  c.BuildUri("repos", c.namespace, c.repo, "git/trees", entry.Hash.String()),
				})
			} else if entry.Mode == filemode.Submodule {
				r = append(r, gin.H{
					"path": entry.Name,
					"mode": entry.Mode,
					"type": "submodule",
					"sha":  entry.Hash.String(),
					// "url":  c.BuildUri("repos", c.namespace, c.repo, "git/trees", entry.Hash.String()),
				})
			}
		}
	}
	return gin.H{
		"url":  c.BuildUri("repos", c.namespace, c.repo, "git/trees", tree.Hash.String()),
		"sha":  tree.Hash.String(),
		"tree": r,
	}
}
