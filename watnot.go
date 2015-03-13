// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9,!solaris

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-fsnotify/fsnotify"
)

func cat(p string) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	var fileInfo os.FileInfo
	fileInfo, err = f.Stat()
	if err != nil {
		return err
	}

	isdir := (fileInfo.Mode() & os.ModeDir) != 0
	if isdir {
		return fmt.Errorf("cannot cat %s", p)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil
}

func NewWatcher(p string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if (event.Op & fsnotify.Write) == fsnotify.Write {
					cat(event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(p)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func main() {
	flag.Parse()
	argv := flag.Args()
	argc := len(argv)
	if argc < 1 {
		fmt.Fprintf(os.Stderr, "expecting a file path at least")
		os.Exit(-1)
	}

	absPath, err := filepath.Abs(argv[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n")
		os.Exit(-2)
	}
	NewWatcher(absPath)
}
