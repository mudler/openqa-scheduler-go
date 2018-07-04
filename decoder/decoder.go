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

package decoder

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	common "github.com/mudler/openqa-scheduler-go/common"

	encoder "github.com/mudler/openqa-scheduler-go/encoder"
)

type Decoder struct{}
type Assignment struct {
	Worker *encoder.Worker
	Test   *encoder.Test
	State  string
	Value  bool
}

func (a *Assignment) Encode() string {
	return fmt.Sprintf(common.AssignFmt, a.Test.Encode(), a.Worker.Encode(), a.State)
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func NewAssignment(t *encoder.Test, w *encoder.Worker, state string, value bool) *Assignment {
	return &Assignment{Worker: w, Test: t, State: state, Value: value}
}
func DecodeTest(test string) (*encoder.Test, error) {
	test_att := strings.Split(test, common.TestSep)
	t := &encoder.Test{}

	t.Name = test_att[0]

	if len(test_att) > 2 {
		test_worker_class := test_att[1]
		test_parent := test_att[2]
		test_parallel := test_att[3]
		t.Parent = test_parent

		wc := strings.Split(test_worker_class, common.WorkerClassSep)
		p := strings.Split(test_parallel, common.TestParallelSep)
		t.WorkerClass = wc
		t.Parallel = p
	}

	return t, nil
}

func DecodeWorker(worker string) (*encoder.Worker, error) {
	worker_att := strings.Split(worker, common.WorkerSep)
	NameInstance := worker_att[0]
	WorkerClasses := worker_att[1]

	nameI := strings.Split(NameInstance, common.WorkerInstSep)
	wc := strings.Split(WorkerClasses, common.WorkerClassSep)
	instance, err := strconv.Atoi(nameI[1])
	if err != nil {
		return &encoder.Worker{}, err
	}

	return &encoder.Worker{Name: nameI[0], Instance: instance, WorkerClass: wc}, nil
}

func (d *Decoder) DecodeAssignment(assignment string) (*Assignment, error) {
	if !strings.Contains(assignment, common.AssignSep) {
		return &Assignment{}, errors.New("Decode error: malformed string")
	}
	data_row := strings.Split(assignment, common.AssignSep)
	if len(data_row) != 3 {
		return &Assignment{}, errors.New("Decode error: malformed string")
	}
	t, err := DecodeTest(data_row[0])
	if err != nil {
		return &Assignment{}, err
	}
	w, err := DecodeWorker(data_row[1])
	if err != nil {
		return &Assignment{}, err
	}

	// Assume true
	return NewAssignment(t, w, data_row[2], true), nil
}

func (d *Decoder) DecodeModel(model map[string]bool) []*Assignment {
	return DecodeModel(model)
}

func DecodeModel(model map[string]bool) []*Assignment {
	d := NewDecoder()
	ass := make([]*Assignment, 0)
	for k, v := range model {
		if a, err := d.DecodeAssignment(k); err == nil {
			a.Value = false
			if v {
				a.Value = true
			}
			if a.State == common.STATE_CURRENT {
				ass = append(ass, a)
			} // Else, there was a state transition between Initial state and current run
		}
	}
	return ass
}
