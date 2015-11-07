// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9,!solaris

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

    "github.com/odeke-em/watnot/src"
)

func main() {
	flag.Parse()
	argv := flag.Args()
	argc := len(argv)
	if argc < 1 {
		fmt.Fprintf(os.Stderr, "expecting a file path at least\n")
		os.Exit(-1)
	}

    absPaths := []string{}

    for _, relPath := range argv {
	    absPath, err := filepath.Abs(relPath)
	    if err == nil {
            absPaths = append(absPaths, absPath)
        } else {
		    fmt.Fprintf(os.Stderr, "%q %v\n", absPath, err)
		    os.Exit(-2)
	    }

        
    }

	watnot.CatOnWriteChange(absPaths...)
}
