package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cameronharro/hs-workflow-tester/internal/actiondefinition"
)

func getActionDefs(root string) ([]actiondefinition.ActionDefinition, error) {
	filePaths, err := getActionDefFiles(".")
	if err != nil {
		return nil, err
	}

	actionDefinitions := make([]actiondefinition.ActionDefinition, len(filePaths))
	for i, path := range filePaths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		def, err := actiondefinition.Parse(data)
		if err != nil {
			return nil, err
		}

		actionDefinitions[i] = def
	}

	return actionDefinitions, nil
}

func getActionDefFiles(root string) ([]string, error) {
	var workflowActionsDirPath string
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.Name() == "workflow-actions" {
			workflowActionsDirPath = path
			return fs.SkipAll
		}
		if d.Name() == ".git" || d.Name() == "node_modules" {
			return fs.SkipDir
		}
		if strings.Contains(path, "hsproject.json") {
			dir, _ := filepath.Split(path)
			workflowActionsDirPath = dir
		}
		if workflowActionsDirPath != "" && !strings.Contains(path, workflowActionsDirPath) {
			return fs.SkipDir
		}
		return nil
	})

	matches, err := filepath.Glob(workflowActionsDirPath + "/*-hsmeta.json")
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, errors.New("No action definition files found")
	}
	return matches, nil
}
