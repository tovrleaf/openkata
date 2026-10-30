// Command site-version prints the release the checked-out commit contains,
// such as v0.14.0, for the deploy scripts to build into the binaries.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tovrleaf/openkata/internal/siteversion"
)

func main() {
	dir := flag.String("dir", ".", "repository directory")
	require := flag.Bool("require", false, "fail if the commit contains no release tag")
	flag.Parse()

	if err := run(*dir, *require, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "site-version:", err)
		os.Exit(1)
	}
}

// run prints the version. With require, a commit that contains no release
// tag is an error, so a deploy cannot ship "Release dev", for instance
// from a clone that was made without its tags.
func run(dir string, require bool, out io.Writer) error {
	version := siteversion.FromGit(dir)
	if require && version == siteversion.Dev {
		return errors.New("no release tag (vMAJOR.MINOR.PATCH) in the history of this commit; fetch the tags")
	}
	_, err := fmt.Fprintln(out, version)
	return err
}
