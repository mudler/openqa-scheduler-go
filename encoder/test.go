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
	"strings"

	"github.com/mudler/openqa-scheduler-go/common"
)

var tests []*Test

type TestColl struct {
	Tests []*Test
}

func NewTestColl() *TestColl {
	return &TestColl{}
}
func (coll *TestColl) NewTest(name string) *Test {
	t := &Test{Name: name}
	coll.AddTest(t)
	return t
}

func (coll *TestColl) AddTest(t *Test) {
	coll.Tests = append(coll.Tests, t)
}

type Test struct {
	WorkerClass []string
	Name        string
	Parent      string
	Parallel    []string
}

func NewTest(name string) *Test {
	t := &Test{Name: name}
	AddTest(t)
	return t
}

func AddTest(t *Test) {
	tests = append(tests, t)
}

func (t *Test) AddWorkerClass(wc string) {
	t.WorkerClass = append(t.WorkerClass, wc)
}

func (t *Test) RequiresWorkerClass(s string) bool {
	for _, w := range t.WorkerClass {
		if w == s {
			return true
		}
	}

	return false
}

func (t *Test) SetParent(p string) {
	t.Parent = p
}

func (t *Test) AddParallel(p string) {
	t.Parallel = append(t.Parallel, p)
}

func (t *Test) Encode() string {
	w_class := strings.Join(t.WorkerClass, common.WorkerClassSep)
	p := strings.Join(t.Parallel, common.TestParallelSep)
	return fmt.Sprintf(common.TestEncodeFormat, t.Name, w_class, t.Parent, p)
}

// Actions: test1 is assigned at worker1

//
// test1 is scheduled in worker 2 if test1 worker class is contained in w2(worker_class)
// ->
// s^t_a == w2 if wc^t == ( wc'^w or wc''^w )
