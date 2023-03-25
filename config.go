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

type YamlConfiguration struct {
	Repos []YamlRepo
}

type YamlRepo struct {
	Namespace   string `yaml:"namespace"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
}

type Repo struct {
	Namespace   string
	Name        string
	Description string
	Path        string
	Url         string
	Repo        *git.Repository
}

func LoadConfiguration() (map[string]*Repo, error) {
	var this YamlConfiguration
	config_path := os.Getenv("CONFIG_PATH")

	yamlFile, err := ioutil.ReadFile(func() string {
		if config_path == "" {
			return "config.yaml"
		}
		return config_path
	}())

	if err == nil {
		// If we loaded the yaml file then unmarshal it
		err = yaml.Unmarshal(yamlFile, &this)
		if err != nil {
			return nil, err
		}
	} else {
		if config_path == "" {
			log.Printf("Error loading config file, trying working directory")
			// If no config was defined and loaded then default to loading the repo at the CWD
			cwd, _ := os.Getwd()
			this = YamlConfiguration{
				Repos: []YamlRepo{
					{
						Namespace:   "local",
						Name:        "self",
						Description: cwd,
						Path:        ".",
					},
				},
			}
		} else {
			log.Printf("Error loading config file")
			return nil, err
		}
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
