package main

import (
	"path/filepath"

	"github.com/restic/restic/internal/archiver"
	"github.com/restic/restic/internal/debug"
	"github.com/restic/restic/internal/errors"
	"github.com/restic/restic/internal/fs"
	"github.com/restic/restic/internal/repository"
)

// rejectResticCache returns a RejectByNameFunc that rejects the restic cache
// directory (if set).
func rejectResticCache(repo *repository.Repository) (archiver.RejectByNameFunc, error) {
	if repo.Cache() == nil {
		return func(string) bool {
			return false
		}, nil
	}
	cacheBase := repo.Cache().BaseDir()

	if cacheBase == "" {
		return nil, errors.New("cacheBase is empty string")
	}

	return func(item string) bool {
		if fs.HasPathPrefix(cacheBase, item) {
			debug.Log("rejecting restic cache directory %v", item)
			return true
		}

		return false
	}, nil
}

// rejectResticRepo returns a RejectByNameFunc that rejects the active local
// restic repository directory (if backed by a local filesystem backend).
func rejectResticRepo(repo *repository.Repository) (archiver.RejectByNameFunc, error) {
	repoPath := repo.LocalPath()
	if repoPath == "" {
		return func(string) bool {
			return false
		}, nil
	}

	repoBase, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	}
	resolvedBase, err := filepath.EvalSymlinks(repoBase)
	if err != nil {
		resolvedBase = ""
	}

	return func(item string) bool {
		if fs.HasPathPrefix(repoBase, item) || (resolvedBase != "" && fs.HasPathPrefix(resolvedBase, item)) {
			debug.Log("rejecting restic repository directory %v", item)
			return true
		}

		return false
	}, nil
}
