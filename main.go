package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"golang.org/x/mod/semver"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: newest-version <repo>")
		os.Exit(1)
	}
	config := &config.RemoteConfig{
		URLs: []string{os.Args[1]},
	}
	s := memory.NewStorage()
	remote := git.NewRemote(s, config)
	rfs, err := remote.List(&git.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list remote repo: %v\n", err)
		os.Exit(1)
	}

	var newest *plumbing.Reference
	for _, rf := range rfs {
		if !rf.Name().IsTag() {
			continue
		}
		if newest == nil {
			newest = rf
			continue
		}
		if semver.Compare(rf.Name().Short(), newest.Name().Short()) > 0 {
			newest = rf
		}
	}
	if newest == nil {
		fmt.Fprintln(os.Stderr, "no valid tags found")
		os.Exit(1)
	}
	fmt.Println(semver.Canonical(newest.Name().Short()))
}
