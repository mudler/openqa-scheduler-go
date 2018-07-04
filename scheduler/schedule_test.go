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
	"fmt"
	"testing"

	"github.com/mudler/openqa-scheduler-go/encoder"
)

func TestSchedule(t *testing.T) {
	// Simple test to ensure we are not panic'ing for nothing
	tests := encoder.NewTestColl()
	workers := encoder.NewWorkerColl()

	w1 := workers.NewWorker("mudler")
	workers.NewWorker("mudler_away")
	t1 := tests.NewTest("lunch")
	tests.NewTest("hiking")

	w1.AddWorkerClass("developer")
	t1.AddWorkerClass("developer")

	s := NewScheduler(workers, tests)
	model, f, err := s.Schedule()
	fmt.Println(f)
	if err != nil {
		t.Error(err)
	}
	for k, v := range model {
		if w, test, err := s.DecodeAssignment(k); err == nil {
			if w.Name != "mudler" {
				t.Error("Failed solve, results: ", w, test, k, v)
			}
		}

	}
}
