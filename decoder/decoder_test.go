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
	"testing"

	"github.com/mudler/openqa-scheduler-go/encoder"
)

func TestAssignment(t *testing.T) {
	// Simple test to ensure we are not panic'ing for nothing
	tests := encoder.NewTestColl()
	workers := encoder.NewWorkerColl()

	w1 := workers.NewWorker("mudler")
	workers.NewWorker("mudler_away")
	t1 := tests.NewTest("lunch")
	tests.NewTest("hiking")

	w1.AddWorkerClass("developer")
	t1.AddWorkerClass("developer")

	assignment := NewAssignment(t1, w1, "current", true)

	if assignment.Encode() != "lunch#developer##@mudler:0#developer@current" {
		t.Error("Different Encode", assignment.Encode())
	}

	a, err := NewDecoder().DecodeAssignment("lunch#developer##@mudler:0#developer@current")
	if err != nil {
		t.Error(err)
	}
	if a.Value != true {
		t.Error("Default value for assignment should be true")
	}
	if a.Test.Name != t1.Name {
		t.Error("Encode/Decode differs", a.Test.Name, t1.Name)
	}
	if a.Encode() != assignment.Encode() {
		t.Error("Encode/Decode differs", a.Encode(), assignment.Encode())
	}

	a, err = NewDecoder().DecodeAssignment("fail")
	if err == nil {
		t.Error("Decode didn't failed")
	}

	_, err = NewDecoder().DecodeAssignment("fail@")
	if err == nil {
		t.Error("Decode didn't failed")
	}
}

func TestWorkerDecode(t *testing.T) {
	workers := encoder.NewWorkerColl()

	w := workers.NewWorker("w1")

	w.Instance = 20
	w.AddWorkerClass("qemu32")
	w.AddWorkerClass("qemu64")

	w2, err := DecodeWorker(w.Encode())
	if err != nil {
		t.Fatal("Error while decoding worker", err)
	}
	if w2.Name != "w1" {
		t.Fatal("Test not added to the test list", w2)
	}

	if w2.Name != "w1" {
		t.Fatal("Test not added to the test list", w2)
	}
	if !w2.ProvidesWorkerClass("qemu32") {
		t.Fatal("Worker provides qemu32")
	}
}

func TestTestDecode(t *testing.T) {
	tests := encoder.NewTestColl()

	w := tests.NewTest("mudler")

	w.AddWorkerClass("qemu32")
	w.AddWorkerClass("qemu64")

	w2, err := DecodeTest(w.Encode())
	if err != nil {
		t.Fatal("Error while decoding worker", err)
	}
	if w2.Name != "mudler" {
		t.Fatal("Test not added to the test list", w2)
	}

	if !w2.RequiresWorkerClass("qemu32") {
		t.Fatal("Worker provides qemu32")
	}
}

func TestDecodeModel(t *testing.T) {
	tests := encoder.NewTestColl()
	workers := encoder.NewWorkerColl()
	d := NewDecoder()
	w1 := workers.NewWorker("mudler")
	workers.NewWorker("mudler_away")
	t1 := tests.NewTest("lunch")
	tests.NewTest("hiking")

	w1.AddWorkerClass("developer")
	t1.AddWorkerClass("developer")

	assignment := NewAssignment(t1, w1, "current", false)
	model := map[string]bool{assignment.Encode(): false}
	ass := DecodeModel(model)

	if ass[0].Encode() != assignment.Encode() {
		t.Error("Failed decoding model ", model)
	}

	assignment = NewAssignment(t1, w1, "current", true)
	model = map[string]bool{assignment.Encode(): true}
	ass = DecodeModel(model)

	if ass[0].Encode() != assignment.Encode() {
		t.Error("Failed decoding model ", model)
	}
	if !ass[0].Value {
		t.Error("Failed decoding model ", model)
	}

	assignment = NewAssignment(t1, w1, "old", true)
	model = map[string]bool{assignment.Encode(): true}
	ass = d.DecodeModel(model) // It should just parse the "current" ones
	if len(ass) > 0 {
		t.Error("DecodeModel gave result also if aren't current")
	}
}
