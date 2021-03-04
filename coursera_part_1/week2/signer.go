package main

import (
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func SingleHash(in, out chan interface{}) {
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			chMd := make(chan string)
			chCrc := make(chan string)
			chCrcMd := make(chan string)

			go func(ch chan string, data string) {
				chCrc <- DataSignerCrc32(data)
			}(chCrc, data)

			mu.Lock()
			go func(ch chan string, data string) {
				ch <- DataSignerMd5(data)
			}(chMd, data)
			time.Sleep(11 * time.Millisecond)
			mu.Unlock()

			go func(chMd chan string) {
				chCrcMd <- DataSignerCrc32(<-chMd)
			}(chMd)

			out <- <-chCrc + "~" + <-chCrcMd

		}(strconv.Itoa(data.(int)))
		runtime.Gosched()
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for data := range in {
		stringedData := data.(string)
		wg.Add(1)
		go func() {
			var chs = []chan string{
				make(chan string),
				make(chan string),
				make(chan string),
				make(chan string),
				make(chan string),
				make(chan string),
			}
			var hashSum string
			for i, ch := range chs {
				go func(ch chan string, data string) {
					ch <- DataSignerCrc32(data)
				}(ch, strconv.Itoa(i)+stringedData)
				runtime.Gosched()
			}
			mu.Lock()
			for _, ch := range chs {
				hashSum += <-ch
			}
			out <- hashSum
			mu.Unlock()

			wg.Done()
		}()
		//runtime.Gosched()
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {

	var sorted_res []string

	for data := range in {
		sorted_res = append(sorted_res, data.(string))
	}
	sort.Strings(sorted_res)
	out <- strings.Join(sorted_res[:], "_")

}

func jobMaker(job job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	job(in, out)
}

func ExecutePipeline(jobs ...job) {

	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	for _, singleJob := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		go jobMaker(singleJob, in, out, wg)
		in = out
	}
	wg.Wait()
}

func main() {

}

// сюда писать код
