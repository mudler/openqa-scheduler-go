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

func TestAdd(t *testing.T) {
	t1 := NewTest("t1")

	t1.AddWorkerClass("qemu32")
	t1.AddWorkerClass("qemu64")
	t1.AddParallel("t2")
	t1.SetParent("t2")

	if t1.Encode() != "t1#qemu32,qemu64#t2#t2" {
		t.Fatal("Encode mismatch", t1.Encode())
	}

	if tests[0].Name != "t1" {
		t.Fatal("Test not added to the test list", tests)
	}

	t2, err := DecodeTest(t1.Encode())
	if err != nil {
		t.Fatal("Error while decoding test", err)
	}

	if t2.Name != "t1" {
		t.Fatal("Test not Decoded correctly", t2)
	}
	if !t2.RequiresWorkerClass("qemu32") {
		t.Fatal("Test requires qemu32")
	}

}

func TestTestsColl(t *testing.T) {

	coll := NewTestColl()

	coll.NewTest("t1")

	if coll.Tests[0].Name != "t1" {
		t.Error("Test not added to collection")
	}

}
