package scheduler

import (
	"errors"
	"fmt"

	"github.com/crillab/gophersat/bf"
	"github.com/mudler/openqa-scheduler/encoder"
)

type Scheduler struct {
	WorkerCollection *encoder.WorkerColl
	TestCollection   *encoder.TestColl
}

func NewScheduler(WorkerColl *encoder.WorkerColl, TestColl *encoder.TestColl) *Scheduler {

	return &Scheduler{WorkerCollection: WorkerColl, TestCollection: TestColl}
}

func (s *Scheduler) Schedule() (map[string]bool, error) {

	f := bf.True

	for _, t := range s.TestCollection.Tests {

		var vars []bf.Formula = make([]bf.Formula, 0)

		for _, w := range s.WorkerCollection.Workers {

			if w.Satisfies(t) { // encoding filter by class

				// If we accept this test, not going to accept others
				var doesnotaccept []bf.Formula = make([]bf.Formula, 0)
				for _, t2 := range s.TestCollection.Tests {
					if t2.Name != t.Name {
						doesnotaccept = append(doesnotaccept, bf.Not(bf.Var(w.Name+".accepts."+t2.Name)))
					}
				}

				// But if we accept it, task and worker goes together, and accepts is true
				final_formula := bf.And(bf.Var(w.Encode()), bf.Var(t.Encode()), bf.Var(w.Name+".accepts."+t.Name))

				for _, i := range doesnotaccept {
					// For each of it, bind it to not acceptance of other tasks
					final_formula = bf.And(final_formula, i)

				}
				vars = append(vars, final_formula) // bf.And(bf.Var(w.Encode()), bf.Var(t.Encode()), bf.Var(w.Name+".accepts."+t.Name), doesnotaccept...))
			}
		}

		f = bf.And(f, bf.Or(vars...))
	}

	fmt.Println(f)
	model := bf.Solve(f)
	if model == nil {
		return model, errors.New("Error: cannot assign tests to workers")
	}
	return model, nil
}
