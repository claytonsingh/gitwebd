package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Repos []struct {
		Namespace   string `yaml:"namespace"`
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Path        string `yaml:"path"`
	}
}

type Repo struct {
	Namespace   string `yaml:"namespace"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
	Url         string
	Repo        *git.Repository
}

func LoadConfiguration() (map[string]*Repo, error) {
	var this Configuration
	config_path := os.Getenv("CONFIG_PATH")
	if config_path == "" {
		config_path = "config.yaml"
	}

	yamlFile, err := ioutil.ReadFile(config_path)
	if err != nil {
		log.Printf("Error loading config file")
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &this)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*Repo, 0)
	for _, repo := range this.Repos {
		slug := fmt.Sprintf("%s/%s", repo.Namespace, repo.Name)
		result[slug] = &Repo{
			Namespace:   repo.Namespace,
			Name:        repo.Name,
			Description: repo.Description,
			Path:        repo.Path,
			Url:         slug,
		}
	}

	return result, nil
}

func LoadRepos(this map[string]*Repo) error {
	repos := make(map[string]*git.Repository, 0)
	for _, repo := range this {
		if path, err := filepath.Abs(repo.Path); err == nil {
			if r, ok := repos[path]; ok {
				repo.Repo = r
			} else {
				if r, err := git.PlainOpen(path); err != nil {
					log.Printf("Unable to load repo %s/%s at %s\n", repo.Namespace, repo.Name, repo.Path)
				} else {
					repo.Repo = r
					repos[path] = r
				}
			}
		} else {
			log.Printf("Unable to load repo %s/%s at %s\n", repo.Namespace, repo.Name, repo.Path)
		}
	}
	return nil
}
