package snippet

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var spMap = make(map[string]*Snippet)

func init() {
	loadSnippet()
}

func GetSnippet(shortCut string) (*Snippet, bool) {
	if sp, isExist := spMap[shortCut]; isExist {
		return sp, true
	}
	return nil, false
}

func GetSnippetList() []*Snippet {
	spList := make([]*Snippet, 0)
	for _, sp := range spMap {
		spList = append(spList, sp)
	}
	return spList
}

func loadSnippet() {
	err := filepath.Walk("./snippet", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}
		snippet, err := readSnippetYaml(path)
		if err != nil {
			return err
		}
		if _, isExist := spMap[snippet.ShortCut]; isExist {
			return errors.New(fmt.Sprintf("snippet 关键字冲突：%s", snippet.ShortCut))
		}
		spMap[snippet.ShortCut] = snippet
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func readSnippetYaml(path string) (*Snippet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s Snippet
	err = yaml.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

type Variable struct {
	Name    string `yaml:"name"`
	Default string `yaml:"default"`
	Val     string
}

type Snippet struct {
	CMD          string     `yaml:"cmd"`
	ShortCut     string     `yaml:"short_cut"`
	VariableList []Variable `yaml:"variable"`
	Tag          []string   `yaml:"tag"`
}

func (s *Snippet) String() string {
	return fmt.Sprintf("\nCMD: %s\n ShortCut: %s", s.CMD, s.ShortCut)
}

func (s *Snippet) GetExecCMD(variable ...Variable) string {
	execCMD := s.CMD
	for _, vt := range variable {
		if vt.Val != "" {
			execCMD = strings.ReplaceAll(execCMD, vt.Name, vt.Val)
		} else {
			execCMD = strings.ReplaceAll(execCMD, vt.Name, vt.Default)
		}
	}
	return execCMD
}
