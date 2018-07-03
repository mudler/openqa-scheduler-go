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

	// AddWorker(w)

	if workers[0].Name != "w1" {
		t.Fatal("Test not added to the test list", workers)
	}

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
