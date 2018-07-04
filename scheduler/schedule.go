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

package scheduler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/crillab/gophersat/bf"
	"github.com/mudler/openqa-scheduler-go/encoder"
)

const assignSep = "@"
const assignFmt = "%s" + assignSep + "%s"
const stateSep = "!"
const stateFmt = "%s" + stateSep + "%s"

const STATE_RUNNING = "running"

type Scheduler struct {
	WorkerCollection *encoder.WorkerColl
	TestCollection   *encoder.TestColl
}

func NewScheduler(WorkerColl *encoder.WorkerColl, TestColl *encoder.TestColl) *Scheduler {
	return &Scheduler{WorkerCollection: WorkerColl, TestCollection: TestColl}
}

func (s *Scheduler) Assign(w *encoder.Worker, t *encoder.Test) string {
	return fmt.Sprintf(assignFmt, t.Encode(), w.Encode())
}

func (s *Scheduler) DecodeAssignment(assignment string) (*encoder.Worker, *encoder.Test, error) {
	if !strings.Contains(assignment, assignSep) {
		return &encoder.Worker{}, &encoder.Test{}, errors.New("Decode error: malformed string")
	}
	data_row := strings.Split(assignment, assignSep)
	if len(data_row) != 2 {
		return &encoder.Worker{}, &encoder.Test{}, errors.New("Decode error: malformed string")
	}
	t, err := encoder.DecodeTest(data_row[0])
	if err != nil {
		return &encoder.Worker{}, &encoder.Test{}, err
	}
	w, err := encoder.DecodeWorker(data_row[1])
	if err != nil {
		return &encoder.Worker{}, &encoder.Test{}, err
	}
	return w, t, nil
}

func (s *Scheduler) TaskState(t *encoder.Test, state string) string {
	return fmt.Sprintf(stateFmt, t.Encode(), state)
}

func (s *Scheduler) BuildFormula() bf.Formula {
	f := bf.True
	// TODO: This is very raw and all have at least to go to binary encoding and avoid wasting cycles
	// Optimization needed

	for _, t := range s.TestCollection.Tests {

		var vars []bf.Formula = make([]bf.Formula, 0)

		for _, w := range s.WorkerCollection.Workers {

			if w.Satisfies(t) { // encoding filter by class - remove unnecessary load from solver with simple check

				// If we accept this test, not going to accept others
				var doesnotaccept []bf.Formula = make([]bf.Formula, 0)
				for _, t2 := range s.TestCollection.Tests {
					if t2.Name != t.Name {
						doesnotaccept = append(doesnotaccept, bf.Not(bf.Var(s.Assign(w, t2))))
					}
				}

				// But if we accept it, task and worker goes together, and we set it to accepted
				final_formula := bf.And(bf.Var(w.Encode()), bf.Var(t.Encode()), bf.Var(s.TaskState(t, STATE_RUNNING)), bf.Var(s.Assign(w, t)))

				for _, i := range doesnotaccept {
					// For each of it, bind it to not acceptance of other tasks
					final_formula = bf.And(final_formula, i)
				}
				vars = append(vars, final_formula) // bf.And(bf.Var(w.Encode()), bf.Var(t.Encode()), bf.Var(w.Name+".accepts."+t.Name), doesnotaccept...))
			}
		}

		f = bf.And(f, bf.Or(vars...))
	}
	return f
}

func (s *Scheduler) Solve(f bf.Formula) (map[string]bool, bf.Formula, error) {
	model := bf.Solve(f)
	if model == nil {
		return model, f, errors.New("Error: cannot assign tests to workers")
	}
	return model, f, nil
}

func (s *Scheduler) Schedule() (map[string]bool, bf.Formula, error) {
	return s.Solve(s.BuildFormula())
}
