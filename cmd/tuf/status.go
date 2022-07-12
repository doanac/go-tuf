package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/flynn/go-docopt"
	"github.com/theupdateframework/go-tuf"
)

func init() {
	register("status", cmdStatus, `
usage: tuf status --expires <date> <role>

Check if the role's metadata will be expired on the given date. 

Options:
  --expires=<date>   Must be in one of the formats:
                     * RFC3339  - 2006-01-02T15:04:05Z07:00
                     * RFC822   - 02 Jan 06 15:04 MST
                     * UnixDate - Mon Jan _2 15:04:05 MST 2006
`)
}

func cmdStatus(args *docopt.Args, repo *tuf.Repo) error {
	role := args.String["<role>"]
	expiresStr := args.String["--expires"]
	if expiresStr == "" {
		return errors.New("--expires arg not set")
	}

	formats := []string{
		time.RFC3339,
		time.RFC822,
		time.UnixDate,
	}
	for _, fmt := range formats {
		expires, err := time.Parse(fmt, expiresStr)
		if err == nil {
			return repo.RoleStatus(role, expires)
		}
	}
	return fmt.Errorf("failed to parse --expires arg")
}
