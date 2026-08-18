package plugin

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"

	protocol "github.com/weirdGuy/kastor/protocol/v1"
)

//go:embed scaffold/*
var scaffoldFS embed.FS

func (Handler) Scaffold(_ context.Context, request *protocol.ScaffoldRequest) (*protocol.ScaffoldResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("anthropic: scaffold request is missing")
	}
	var names []string
	err := fs.WalkDir(scaffoldFS, "scaffold", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			names = append(names, strings.TrimPrefix(name, "scaffold/"))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("anthropic: read scaffold: %w", err)
	}
	sort.Strings(names)
	response := &protocol.ScaffoldResponse{}
	for _, name := range names {
		data, err := scaffoldFS.ReadFile("scaffold/" + name)
		if err != nil {
			return nil, fmt.Errorf("anthropic: read scaffold %s: %w", name, err)
		}
		if name == "kastor.hcl" && request.VersionConstraint != "" {
			data = []byte(strings.Replace(string(data), `version = "~> 0.1"`, `version = `+strconv.Quote(request.VersionConstraint), 1))
		}
		response.Files = append(response.Files, protocol.File{Path: name, Data: data})
	}
	return response, nil
}
