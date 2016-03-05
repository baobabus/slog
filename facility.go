// Copyright 2016 Aleksey Blinov. All rights reserved.

package slog

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Facility interface {
	OpenLogs(level Priority) (map[Priority]*log.Logger, error)
	Reopen() error
}

type fFile struct {
	path string
	mux  sync.RWMutex
	file *os.File
}

func NewStdFacility(file *os.File) (Facility, error) {
	return &fFile{file: file}, nil
}

func NewFileFacility(path string) (Facility, error) {
	return &fFile{path: path}, nil
}

func (this *fFile) OpenLogs(level Priority) (map[Priority]*log.Logger, error) {
	var err error
	if this.file == nil && len(this.path) > 0 {
		this.file, err = os.OpenFile(this.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			return nil, err
		}
	}
	res := make(map[Priority]*log.Logger, prioritiesCount)
	for pri := PriorityError; pri <= PriorityTrace; pri++ {
		if pri <= level {
			if pri < PriorityTrace {
				res[pri] = log.New(this, pri.Tag(), log.Ldate|log.Ltime)
			} else {
				res[pri] = log.New(this, pri.Tag(), log.Ldate|log.Ltime|log.Lshortfile)
			}
		} else {
			res[pri] = drain.Logger()
		}
	}
	return res, nil
}

func (this *fFile) Reopen() error {
	this.mux.Lock()
	defer this.mux.Unlock()
	if len(this.path) > 0 {
		if this.file == nil {
			return fmt.Errorf("unable to reopen file: %s", this.path)
		}
		if f, err := os.OpenFile(this.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640); err != nil {
			return err
		} else {
			this.file.Sync()
			this.file.Close()
			this.file = f
		}
	}
	return nil
}

func (this *fFile) Write(p []byte) (n int, err error) {
	this.mux.RLock()
	defer this.mux.RUnlock()
	if this.file != nil {
		return this.file.Write(p)
	} else {
		return 0, fmt.Errorf("not open: %s", this.path)
	}
}
