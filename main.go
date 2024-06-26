package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Panic1[T any](a T, e error) T {
	if e != nil {
		panic(e)
	}
	return a
}

func main() {
	var repos map[string]*Repo
	{
		var err error
		repos, err = LoadConfiguration()
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		LoadRepos(repos)
	}

	WithRepo := func(handler func(Context)) func(*gin.Context) {
		return func(c *gin.Context) {
			namespace := c.Param("namespace")
			reponame := c.Param("repo")
			repository, ok := repos[fmt.Sprintf("%s/%s", namespace, reponame)]
			if ok && repository != nil {
				handler(Context{
					Context:    c,
					repository: repository,
					namespace:  namespace,
					repo:       reponame,
				})
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		}
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/repos")
	})

	r.GET("/repos", func(c *gin.Context) {
		r := make([]any, 0)
		for _, repo := range repos {
			r = append(r, gin.H{
				"name":        repo.Name,
				"namespace":   repo.Namespace,
				"description": repo.Description,
				"url":         BuildUri(c, "repos", repo.Url),
			})
		}
		c.IndentedJSON(http.StatusOK, r)
	})

	// Repo index
	r.GET("/repos/:namespace/:repo/", WithRepo(func(c Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"name":        c.repository.Name,
			"description": c.repository.Description,
			"url":         c.BuildUri("repos", c.repository.Name),
			"indexes": gin.H{
				"branches": c.BuildUri("repos", c.namespace, c.repo, "branches"),
				"tags":     c.BuildUri("repos", c.namespace, c.repo, "tags"),
				// "commits":  c.BuildUri("repos", c.namespace, c.repo, "commits"),
			},
		})
	}))

	// branch index
	r.GET("/repos/:namespace/:repo/branches", WithRepo(func(c Context) {
		r := make([]any, 0)
		if branches, err := c.repository.Repo.Branches(); err == nil {
			branches.ForEach(func(t *plumbing.Reference) error {
				r = append(r, FormatBranchRef(c, t))
				return nil
			})
			c.IndentedJSON(http.StatusOK, r)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}))

	r.GET("/repos/:namespace/:repo/branches/:branch", WithRepo(func(c Context) {
		if branch, err := c.repository.Repo.Storer.Reference(plumbing.NewBranchReferenceName(c.Param("branch"))); err == nil {
			if commit, err := c.repository.Repo.CommitObject(branch.Hash()); err == nil {
				c.IndentedJSON(http.StatusOK, FormatBranchFull(c, branch, commit))
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}))

	r.GET("/repos/:namespace/:repo/branches/:branch/raw/*path", WithRepo(func(c Context) {
		if branch, err := c.repository.Repo.Storer.Reference(plumbing.NewBranchReferenceName(c.Param("branch"))); err == nil {
			if commit, err := c.repository.Repo.CommitObject(branch.Hash()); err == nil {
				c.Redirect(http.StatusTemporaryRedirect, c.FormatRawUri(commit.TreeHash, c.Param("path")))
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}))

	// tag index
	r.GET("/repos/:namespace/:repo/tags", WithRepo(func(c Context) {
		r := make([]any, 0)
		if tags, err := c.repository.Repo.Tags(); err == nil {
			tags.ForEach(func(ref *plumbing.Reference) error {
				tag, err := c.repository.Repo.TagObject(ref.Hash())
				switch err {
				case nil:
					if tag.TargetType == plumbing.CommitObject {
						r = append(r, FormatTagRef(c, ref, tag))
					}
				case plumbing.ErrObjectNotFound:
					if _, err := c.repository.Repo.CommitObject(ref.Hash()); err == nil {
						r = append(r, FormatTagRef(c, ref, nil))
					}
				}
				return nil
			})
			c.IndentedJSON(http.StatusOK, r)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}))

	r.GET("/repos/:namespace/:repo/tags/:tag", WithRepo(func(c Context) {
		if ref, err := c.repository.Repo.Tag(c.Param("tag")); err == nil {
			tag, err := c.repository.Repo.TagObject(ref.Hash())
			switch err {
			case nil:
				if tag.TargetType == plumbing.CommitObject {
					c.IndentedJSON(http.StatusOK, FormatTagRef(c, ref, tag))
					return
				}
			case plumbing.ErrObjectNotFound:
				if _, err := c.repository.Repo.CommitObject(ref.Hash()); err == nil {
					c.IndentedJSON(http.StatusOK, FormatTagRef(c, ref, nil))
					return
				}
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}))

	r.GET("/repos/:namespace/:repo/tags/:tag/raw/*path", WithRepo(func(c Context) {
		if ref, err := c.repository.Repo.Tag(c.Param("tag")); err == nil {
			tag, err := c.repository.Repo.TagObject(ref.Hash())
			switch err {
			case nil:
				if tag.TargetType == plumbing.CommitObject {
					if commit, err := c.repository.Repo.CommitObject(tag.Target); err == nil {
						c.Redirect(http.StatusTemporaryRedirect, c.FormatRawUri(commit.TreeHash, c.Param("path")))
						return
					}
				}
			case plumbing.ErrObjectNotFound:
				if commit, err := c.repository.Repo.CommitObject(ref.Hash()); err == nil {
					c.Redirect(http.StatusTemporaryRedirect, c.FormatRawUri(commit.TreeHash, c.Param("path")))
					return
				}
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}))

	// commit index
	r.GET("/repos/:namespace/:repo/commits", WithRepo(func(c Context) {
		r := make([]any, 0)
		if commit, err := c.repository.Repo.CommitObjects(); err == nil {
			commit.ForEach(func(m *object.Commit) error {
				cmt := FormatCommitRef(c.Context, m.Hash)
				cmt["message"] = m.Message
				cmt["author"] = FormatSignature(c.Context, m.Author)
				cmt["committer"] = FormatSignature(c.Context, m.Committer)
				r = append(r, cmt)
				return nil
			})
			c.IndentedJSON(http.StatusOK, r)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}))

	r.GET("/repos/:namespace/:repo/commits/:sha", WithRepo(func(c Context) {
		sha := plumbing.NewHash(c.Param("sha"))
		if commit, err := c.repository.Repo.CommitObject(sha); err == nil {
			c.IndentedJSON(http.StatusOK, FormatCommitFull(c.Context, commit))
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}))

	r.GET("/repos/:namespace/:repo/commits/:sha/raw/*path", WithRepo(func(c Context) {
		sha := plumbing.NewHash(c.Param("sha"))
		if commit, err := c.repository.Repo.CommitObject(sha); err == nil {
			c.Redirect(http.StatusTemporaryRedirect, c.FormatRawUri(commit.TreeHash, c.Param("path")))
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}))

	// trees
	r.GET("/repos/:namespace/:repo/trees/:sha", WithRepo(func(c Context) {
		var recursive = false
		switch strings.ToLower(c.Query("recursive")) {
		case "1", "true", "yes":
			recursive = true
		}
		sha := plumbing.NewHash(c.Param("sha"))
		if tree, err := c.repository.Repo.TreeObject(sha); err == nil {
			// tree.Entries
			c.IndentedJSON(http.StatusOK, FormatTree(c, tree, recursive))
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}))

	// blobs
	r.GET("/repos/:namespace/:repo/blobs/:sha", WithRepo(func(c Context) {
		sha := plumbing.NewHash(c.Param("sha"))
		if blob, err := c.repository.Repo.BlobObject(sha); err == nil {
			if r, err := blob.Reader(); err == nil {
				defer r.Close()

				buff := new(bytes.Buffer)
				buff.ReadFrom(r)

				c.IndentedJSON(http.StatusOK, gin.H{
					"url":     c.BuildUri("repos", c.namespace, c.repo, "blobs", sha.String()),
					"sha":     sha.String(),
					"size":    blob.Size,
					"content": base64.RawStdEncoding.EncodeToString(buff.Bytes()),
				})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}))

	// raw file download
	r.GET("/repos/:namespace/:repo/raw/:sha/*path", WithRepo(func(c Context) {

		read := func(tree *object.Tree) error {
			if tree, err := tree.File(c.Param("path")[1:]); err == nil {
				if reader, err := tree.Blob.Reader(); err == nil {
					defer reader.Close()

					c.DataFromReader(http.StatusOK, tree.Blob.Size, "application/octet-stream", reader, nil)
					return nil
				} else {
					return err
				}
			} else {
				return err
			}
		}

		if tree, err := c.repository.Repo.TreeObject(plumbing.NewHash(c.Param("sha"))); err == nil {
			if read(tree) == nil {
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}))

	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = "127.0.0.1:8080"
	}
	log.Fatal(r.Run(listen)) // listen on 127.0.0.1:8080
}
