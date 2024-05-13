package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 功能测试
func TestTaskGroup(t *testing.T) {
	var (
		tg TaskGroup

		tasks = []*Task{
			NewTask(1, true, task1),
			NewTask(2, true, task2),
			NewTask(3, true, task3),
		}
	)

	taskResult, err := tg.SetWorkerNums(4).AddTask(tasks...).Run()
	fmt.Printf("**************TaskGroup************\n%+v, %+v\n", taskResult, err)
	for fno, result := range taskResult {
		fmt.Printf("FNO: %d, RESULT: %v , STATUS: %v\n", fno, result.Result(), result.Error())
	}
}

func TestTaskGroupBoundary(t *testing.T) {
	var (
		tg TaskGroup

		tasks = []*Task{}
	)

	taskResult, err := tg.SetWorkerNums(4).AddTask(tasks...).Run()
	fmt.Printf("**************TaskGroup************\n%+v, %+v\n", taskResult, err)
	for fno, result := range taskResult {
		fmt.Printf("FNO: %d, RESULT: %v , STATUS: %v\n", fno, result.Result(), result.Error())
	}
}

func TestTaskGroupAbnormal(t *testing.T) {
	var (
		tg TaskGroup

		tasks = []*Task{
			NewTask(1, true, task1),
			nil,
			NewTask(2, false, task2),
			NewTask(2, true, task1),
		}
	)

	taskResult, err := tg.SetWorkerNums(4).AddTask(tasks...).Run()
	fmt.Printf("**************TaskGroup************\n%+v, %+v\n", taskResult, err)
	for fno, result := range taskResult {
		fmt.Printf("FNO: %d, RESULT: %v , STATUS: %v\n", fno, result.Result(), result.Error())
	}
}

func getRandomNum(mod int) int {
	return rand.Int() % mod
}

func task1() (interface{}, error) {
	const taskFlag = "running TASK1"
	fmt.Println(taskFlag)
	simulateMemIO(taskFlag)
	return getRandomNum(2e3), fmt.Errorf("%s err", taskFlag)
}

type task2Struct struct {
	a int
	b string
}

func task2() (interface{}, error) {
	const taskFlag = "running TASK2"
	fmt.Println(taskFlag)
	simulateMemIO(taskFlag)
	return task2Struct{
		a: getRandomNum(1e1),
		b: "mlee",
	}, nil
}

func task3() (interface{}, error) {
	const taskFlag = "running TASK3"
	fmt.Println(taskFlag)
	simulateMemIO(taskFlag)
	return fmt.Sprintf("%s: The data is %d", taskFlag, getRandomNum(12)), fmt.Errorf("%s err", taskFlag)
}

type task4Struct struct {
	field0 task2Struct
	field1 uint32
	field2 string
	field3 []task2Struct
	field4 map[string]*task2Struct
	field5 *string
}

func task4() (interface{}, error) {
	const taskFlag = "TASK4"
	// fmt.Println(taskFlag)
	simulateMemIO(taskFlag)
	var field5 = "mleeeeeeeeeeeeeee"
	return task4Struct{
		field0: task2Struct{10, "ok"},
		field1: 1024,
		field2: "12",
		field3: []task2Struct{{12, "mlee1"}, {122, "mlee2"}, {1222, "mlee3"}},
		field4: map[string]*task2Struct{"a1": {10, "@@"}, "a2": {111, "##"}},
		field5: &field5,
	}, fmt.Errorf("%s err", taskFlag)
}

func simulateMemIO(s string) {
	var buf bytes.Buffer
	_, _ = buf.WriteString(s)
	_ = buf.String()
}

// -----------性能测试-----------
func BenchmarkTaskGroupZero(b *testing.B) {
	tasks := buildTestCaseData(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var tg TaskGroup
		_, _ = tg.SetWorkerNums(2).AddTask(tasks...).Run()
	}
}

func BenchmarkTaskGroupLow(b *testing.B) {
	tasks := buildTestCaseData(3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var tg TaskGroup
		_, _ = tg.SetWorkerNums(2).AddTask(tasks...).Run()
	}
}

func BenchmarkTaskGroupNormal(b *testing.B) {
	tasks := buildTestCaseData(8)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var tg TaskGroup
		_, _ = tg.SetWorkerNums(2).AddTask(tasks...).Run()
	}
}

func BenchmarkTaskGroupMedium(b *testing.B) {
	tasks := buildTestCaseData(15)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var tg TaskGroup
		_, _ = tg.SetWorkerNums(2).AddTask(tasks...).Run()
	}
}

func BenchmarkTaskGroupHigh(b *testing.B) {
	tasks := buildTestCaseData(40)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var tg TaskGroup
		_, _ = tg.SetWorkerNums(2).AddTask(tasks...).Run()
	}
}

func BenchmarkTaskGroupExtremelyHigh(b *testing.B) {
	tasks := buildTestCaseData(200)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var tg TaskGroup
		_, _ = tg.SetWorkerNums(2).AddTask(tasks...).Run()
	}
}

var taskSet = []func() (interface{}, error){task1, task2, task3, task4}

func buildTestCaseData(taskNums uint32) []*Task {
	if taskNums == 0 {
		return nil
	}
	tasks := make([]*Task, 0, taskNums)
	for i := 1; i <= int(taskNums); i++ {
		tasks = append(tasks, NewTask(uint32(getRandomNum(1e10)), true, taskSet[getRandomNum(len(taskSet))]))
	}
	return tasks
}
