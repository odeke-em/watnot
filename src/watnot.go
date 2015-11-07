// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9,!solaris

package watnot

import (
	"bufio"
	"fmt"
	"log"
	"os"

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
		return fmt.Errorf("cat cannot open %s", p)
	}

	scanner := bufio.NewScanner(f)

	// Clear the screen first
	fmt.Printf("\033[2J\033[;H")

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil
}

func CatOnWriteChange(paths ...string) {
	// Initial cat, at least once
    for _, p := range paths {
	    cat(p)
    }
   
    changesChan := NewWatcher(paths...)
    for p := range changesChan {
        cat(p)
    }
}

func NewWatcher(paths ...string) (changesChan chan string) {
    changesChan = make(chan string)

    doneCount := uint64(0)
    doneChan  := make(chan bool)

    for _, p := range paths {
	    watcher, err := fsnotify.NewWatcher()

	    if err != nil {
		    log.Printf("watcher for %q: err %v\n", p, err)
            continue     
        }

	    err = watcher.Add(p)
	    if err != nil {
		    log.Printf("%q error: %v\n", p, err)
            watcher.Close()
            continue
        }

        doneCount += 1

	    go func(w *fsnotify.Watcher, done chan bool) {
            defer w.Close()

		    for {
			    select {
			    case event := <-w.Events:
				    if (event.Op & fsnotify.Write) == fsnotify.Write {
					    changesChan <- event.Name
				    }

			    case err := <-w.Errors:
                    if err != nil {
				        log.Println("error:", err)
                    }
			    default:
				    continue
			    }
		    }

            done <- true
	    }(watcher, doneChan)
    }

    go func(done chan bool, n uint64) {
        fmt.Println("n", n)
        for i := uint64(0); i < n; i++ { 
            <- done
        }

        close(changesChan)
    }(doneChan, doneCount)

    return
}
