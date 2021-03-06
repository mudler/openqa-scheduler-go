package main

import (
	"fmt"

	encoder "github.com/mudler/openqa-scheduler-go/encoder"
	scheduler "github.com/mudler/openqa-scheduler-go/scheduler"
)

func main() {
	tests := encoder.NewTestColl()
	workers := encoder.NewWorkerColl()

	w1 := workers.NewWorker("worker1")
	w2 := workers.NewWorker("worker2")
	w3 := workers.NewWorker("worker3")
	w4 := workers.NewWorker("worker4")

	t1 := tests.NewTest("t1")
	t2 := tests.NewTest("t2")
	t3 := tests.NewTest("t3")
	t4 := tests.NewTest("t4")

	w1.AddWorkerClass("qemu64")
	w1.AddWorkerClass("qemu32")
	w1.Instance = 1

	w2.AddWorkerClass("qemu64")
	w2.AddWorkerClass("qemu32")
	w2.Instance = 1

	w3.AddWorkerClass("qemu32")
	w3.Instance = 1

	w4.AddWorkerClass("qemu64")
	w4.Instance = 1

	t1.AddWorkerClass("qemu32")
	t2.AddWorkerClass("qemu64")
	t3.AddWorkerClass("qemu32")
	t4.AddWorkerClass("qemu64")

	s := scheduler.NewScheduler(workers, tests)
	model, err := s.ScheduleDecode()
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	for _, a := range model {
		if a.Value {
			fmt.Println("Test:", a.Test.Name, "Assigned to worker:", a.Worker.Name)
		}
	}

}
