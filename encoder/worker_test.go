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

import "testing"

func TestWorkerAdd(t *testing.T) {
	w := NewWorker("w1")

	w.Instance = 20
	w.AddWorkerClass("qemu32")
	w.AddWorkerClass("qemu64")

	if w.Encode() != "w1:20#qemu32,qemu64" {
		t.Fatal("Encode mismatch", w.Encode())
	}

	if workers[0].Name != "w1" {
		t.Fatal("Test not added to the test list", workers)
	}
}

func TestWorkerColl(t *testing.T) {

	coll := NewWorkerColl()

	w1 := coll.NewWorker("w1")

	w1.Instance = 2
	w1.AddWorkerClass("wc")
	t1 := NewTest("name")
	t1.AddWorkerClass("wc")

	if coll.Workers[0].Name != "w1" {
		t.Error("Worker not added to collection")
	}

	if !w1.Satisfies(t1) {
		t.Error("Worker doesn't satisfies the test")
	}

}
