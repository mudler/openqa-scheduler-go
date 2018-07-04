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

	"github.com/mudler/openqa-scheduler-go/common"

	"github.com/crillab/gophersat/bf"
	"github.com/mudler/openqa-scheduler-go/decoder"
	"github.com/mudler/openqa-scheduler-go/encoder"
)

type Scheduler struct {
	WorkerCollection *encoder.WorkerColl
	TestCollection   *encoder.TestColl

	InitialState []*decoder.Assignment
}

func NewScheduler(WorkerColl *encoder.WorkerColl, TestColl *encoder.TestColl) *Scheduler {
	return &Scheduler{WorkerCollection: WorkerColl, TestCollection: TestColl}
}

func (s *Scheduler) Assign(w *encoder.Worker, t *encoder.Test) string {
	return decoder.NewAssignment(t, w, common.STATE_CURRENT, true).Encode()
}

func (s *Scheduler) AssignState(state string, w *encoder.Worker, t *encoder.Test) string {
	return decoder.NewAssignment(t, w, state, true).Encode()
}

func (s *Scheduler) TaskState(t *encoder.Test, state string) string {
	return fmt.Sprintf(common.StateFmt, t.Encode(), state)
}

func (s *Scheduler) BuildFormula() bf.Formula {
	f := bf.True

	// TODO: Consider split solving in multiple worker/tests chunks and re-run itself
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
				final_formula := bf.And(
					bf.Var(w.Encode()),
					bf.Var(t.Encode()),
					bf.Var(s.TaskState(t, common.STATE_RUNNING)),
					bf.Var(s.Assign(w, t)), bf.Not(bf.Var(s.AssignState(common.STATE_OLD, w, t))), // Assign if not already in initial state
				)
				//option 2: filter from encoding the tasks already running in InitialState
				for _, i := range doesnotaccept {
					// For each of it, bind it to not acceptance of other tasks
					final_formula = bf.And(final_formula, i)
				}
				vars = append(vars, final_formula) // bf.And(bf.Var(w.Encode()), bf.Var(t.Encode()), bf.Var(w.Name+".accepts."+t.Name), doesnotaccept...))
			}
		}

		f = bf.And(f, bf.Or(vars...))
	}

	var vars []bf.Formula = make([]bf.Formula, 0)
	// Apply initial state
	if s.InitialState != nil {
		for _, a := range s.InitialState {
			a.State = common.STATE_OLD
			if a.Value {
				vars = append(vars, bf.Var(a.Encode()))
			} else {
				vars = append(vars, bf.Not(bf.Var(a.Encode())))
			}
		}
		vars = append(vars, f)
		f = bf.And(vars...)
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

func (s *Scheduler) ScheduleDecode() ([]*decoder.Assignment, error) {
	model, _, err := s.Schedule()
	if err != nil {
		return []*decoder.Assignment{}, err
	}
	ass := decoder.DecodeModel(model)
	return ass, nil
}
