package main

import (
	"code.google.com/p/goprotobuf/proto"
	"flag"
	"fmt"
	"github.com/nenadl/gomesos"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	taskLimit := 5
	taskId := 0
	exit := make(chan bool)
	localExecutor, _ := executorPath()

	master := flag.String("master", "localhost:5050", "Location of leading Mesos master")
	executorUri := flag.String("executor-uri", localExecutor, "URI of executor executable")
	flag.Parse()

	executor := &gomesos.ExecutorInfo{
		ExecutorId: &gomesos.ExecutorID{Value: proto.String("default")},
		Command: &gomesos.CommandInfo{
			Value: proto.String("./example-executor"),
			Uris: []*gomesos.CommandInfo_URI{
				&gomesos.CommandInfo_URI{Value: executorUri},
			},
		},
		Name:   proto.String("Test Executor (Go)"),
		Source: proto.String("go_test"),
	}

	driver := gomesos.SchedulerDriver{
		Master: *master,
		Framework: gomesos.FrameworkInfo{
			Name: proto.String("GoFramework"),
			User: proto.String(""),
		},

		Scheduler: &gomesos.Scheduler{
			ResourceOffers: func(driver *gomesos.SchedulerDriver, offers []gomesos.Offer) {
				for _, offer := range offers {
					taskId++
					fmt.Printf("Launching task: %d\n", taskId)

					tasks := []gomesos.TaskInfo{
						gomesos.TaskInfo{
							Name: proto.String("go-task"),
							TaskId: &gomesos.TaskID{
								Value: proto.String("go-task-" + strconv.Itoa(taskId)),
							},
							SlaveId:  offer.SlaveId,
							Executor: executor,
							Resources: []*gomesos.Resource{
								gomesos.ScalarResource("cpus", 1),
								gomesos.ScalarResource("mem", 512),
							},
						},
					}

					driver.LaunchTasks(offer.Id, tasks)
				}
			},

			StatusUpdate: func(driver *gomesos.SchedulerDriver, status gomesos.TaskStatus) {
				fmt.Println("Received task status: " + *status.Message)

				if *status.State == gomesos.TaskState_TASK_FINISHED {
					taskLimit--
					if taskLimit <= 0 {
						exit <- true
					}
				}
			},
		},
	}

	driver.Init()
	defer driver.Destroy()

	driver.Start()
	<-exit
	driver.Stop(false)
}

func executorPath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	path := dir + "/example-executor"

	return path, nil
}
