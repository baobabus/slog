// Copyright 2016 Aleksey Blinov. All rights reserved.

package slog

import (
	"log"
	"os"
)

type Facility interface {
	OpenLogs(level Priority) (map[Priority]*log.Logger, error)
}

type fFile struct {
	path string
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
				res[pri] = log.New(this.file, pri.Tag(), log.Ldate|log.Ltime)
			} else {
				res[pri] = log.New(this.file, pri.Tag(), log.Ldate|log.Ltime|log.Lshortfile)
			}
		} else {
			res[pri] = drain.Logger()
		}
	}
	return res, nil
}
