package main

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"github.com/nenadl/gomesos"
)

func main() {
	driver := gomesos.ExecutorDriver{
		Executor: &gomesos.Executor{
			Registered: func(
				driver *gomesos.ExecutorDriver,
				executor gomesos.ExecutorInfo,
				framework gomesos.FrameworkInfo,
				slave gomesos.SlaveInfo) {
				fmt.Println("Executor registered!")
			},

			LaunchTask: func(driver *gomesos.ExecutorDriver, taskInfo gomesos.TaskInfo) {
				fmt.Println("Launch task!")
				driver.SendStatusUpdate(&gomesos.TaskStatus{
					TaskId:  taskInfo.TaskId,
					State:   gomesos.NewTaskState(gomesos.TaskState_TASK_RUNNING),
					Message: proto.String("Go task is running!"),
				})

				driver.SendStatusUpdate(&gomesos.TaskStatus{
					TaskId:  taskInfo.TaskId,
					State:   gomesos.NewTaskState(gomesos.TaskState_TASK_FINISHED),
					Message: proto.String("Go task is done!"),
				})
			},
		},
	}

	driver.Init()
	defer driver.Destroy()

	driver.Run()
}
