package tuf

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/theupdateframework/go-tuf/data"
	"github.com/theupdateframework/go-tuf/pkg/keys"
	"github.com/theupdateframework/go-tuf/util"
)

func GitLocalStore(dir string, p util.PassphraseFunc) LocalStore {
	fsStore := fileSystemStore{
		dir:            dir,
		passphraseFunc: p,
		signerForKeyID: make(map[string]keys.Signer),
		keyIDsForRole:  make(map[string][]string),
	}
	return &gitStore{fsStore}
}

type gitStore struct {
	fsStore fileSystemStore
}

func (g *gitStore) GetMeta() (map[string]json.RawMessage, error) {
	return g.fsStore.GetMeta()
}

func (g *gitStore) SetMeta(name string, meta json.RawMessage) error {
	return g.fsStore.SetMeta(name, meta)
}

func (g *gitStore) FileIsStaged(name string) bool {
	return g.fsStore.FileIsStaged(name)
}

func (g *gitStore) WalkStagedTargets(paths []string, targetsFn TargetsWalkFunc) error {
	return g.fsStore.WalkStagedTargets(paths, targetsFn)
}

func (g *gitStore) Commit(consistentSnapshot bool, versions map[string]int64, hashes map[string]data.Hashes) error {
	if err := g.fsStore.Commit(consistentSnapshot, versions, hashes); err != nil {
		return err
	}
	cmd := exec.Command("git", "add", "repository/*")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to run `git add`: %w", err)
	}

	cmd = exec.Command("git", "status", "--porcelain")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to check for changes to git repository: %w", err)
	}

	if len(out) > 0 {
		// There is a change to commit
		cmd = exec.Command("git", "commit", "-m", "updated by go-tuf")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}
	return nil
}

func (g *gitStore) GetSigners(role string) ([]keys.Signer, error) {
	return g.fsStore.GetSigners(role)
}

func (g *gitStore) SaveSigner(role string, signer keys.Signer) error {
	return g.fsStore.SaveSigner(role, signer)
}

func (g *gitStore) SignersForKeyIDs(keyIDs []string) []keys.Signer {
	return g.fsStore.SignersForKeyIDs(keyIDs)
}

func (g *gitStore) Clean() error {
	return g.fsStore.Clean()
}
