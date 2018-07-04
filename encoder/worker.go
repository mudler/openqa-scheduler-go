// Copyright © 2018 SUSE LLC
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see <http://www.gnu.org/licenses/>.

package encoder

import (
	"fmt"
	"strconv"
	"strings"
)

var workers []*Worker

type WorkerColl struct {
	Workers []*Worker
}

func NewWorkerColl() *WorkerColl {
	return &WorkerColl{}
}

func (coll *WorkerColl) NewWorker(name string) *Worker {
	t := &Worker{Name: name}
	coll.AddWorker(t)
	return t
}

func (coll *WorkerColl) AddWorker(t *Worker) {
	coll.Workers = append(coll.Workers, t)
}

const workerEncodeFormat = "%s:%d#%s"

type Worker struct {
	Name        string
	Instance    int
	WorkerClass []string
}

func NewWorker(name string) *Worker {
	w := &Worker{Name: name}
	AddWorker(w)
	return w
}

func AddWorker(w *Worker) {
	workers = append(workers, w)
}

func (w *Worker) ProvidesWorkerClass(s string) bool {
	for _, w1 := range w.WorkerClass {
		if w1 == s {
			return true
		}
	}

	return false
}

func (w *Worker) Satisfies(t *Test) bool {
	for _, w1 := range w.WorkerClass {
		if t.RequiresWorkerClass(w1) {
			return true
		}
	}

	return false
}

func DecodeWorker(worker string) (*Worker, error) {
	worker_att := strings.Split(worker, "#")
	NameInstance := worker_att[0]
	WorkerClasses := worker_att[1]

	nameI := strings.Split(NameInstance, ":")
	wc := strings.Split(WorkerClasses, ",")
	instance, err := strconv.Atoi(nameI[1])
	if err != nil {
		return &Worker{}, err
	}

	return &Worker{Name: nameI[0], Instance: instance, WorkerClass: wc}, nil
}

func (w *Worker) AddWorkerClass(wc string) {
	w.WorkerClass = append(w.WorkerClass, wc)
}

func (w *Worker) Encode() string {
	w_class := strings.Join(w.WorkerClass, ",")

	return fmt.Sprintf(workerEncodeFormat, w.Name, w.Instance, w_class)
}
