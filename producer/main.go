package main

import (
	"github.com/IBM/sarama"
	"fmt"
	"flag"
	"time"
	"sync"
	"context"
	"os"
	"os/signal"
	"syscall"
	"math/rand"
)

var(
	hostname string
	topic string
)

func init() {
	flag.StringVar(&hostname, "hostname", "localhost:9092", "")
	flag.StringVar(&topic, "topic", "t1", "")
}


func runner(ctx context.Context, id string, onClose []func()) (<-chan []int64, <-chan int, <-chan string) {
	producer, err := sarama.NewSyncProducer([]string{hostname}, nil)
	if err != nil {
		panic(err)
	}	
	messageNo := 0
	cSent := make(chan int, 1024)
	cWin := make(chan []int64, 1024)
	cFail := make(chan string, 1024)
	go func() {
		defer func() {
			if err := producer.Close(); err != nil {
				fmt.Printf("Damn %s wouldn't close!\n", id)
			} else {
				fmt.Printf("%s - closed that crap!\n", id)
			}
			for _, v := range onClose {
				v()
			}			
		}()			
		for {
			//fmt.Printf("Sending message!\n")
			msg := &sarama.ProducerMessage{Key: sarama.StringEncoder(id), Topic: topic, Value:  sarama.StringEncoder(fmt.Sprintf("Hello from %s! Message # %d", id, messageNo))}
			part, offset, err := producer.SendMessage(msg)
			//fmt.Printf("Message sent\n")
			cSent <- 1
			messageNo++
			if err != nil {
				cFail <- fmt.Sprintf("%v", err)
			} else {
				cWin <- []int64{int64(part), offset}
			}
			select {
				case <-ctx.Done():
					fmt.Printf("Producer %s shutting down!\n", id)
					return
				default:
			}
			x := 1 + rand.Intn(300)
			time.Sleep(time.Duration(x) * time.Millisecond)
		}
	}()
	return cWin, cSent, cFail
} 

type FINStat struct {
	ID string
	CurrOffset int64
	CurrPartition int64
	Fails map[string]int64
	Sent int
}

func collector(ctx context.Context, id string, cWin <-chan []int64, cSent <-chan int, cFail <-chan string, back chan<- FINStat, onClose []func())  {
	fmt.Printf("Collecting on ID %s\n", id)
	fs := FINStat{ID: id}
	go func() {
		defer func() {
			for _, v := range onClose {
				v()
			}
		}()
		check := time.After(3 * time.Second)
		for {
			select {
				case <- cSent:
					fs.Sent += 1
				case errorString := <- cFail:
					cnt, ok := fs.Fails[errorString]
					if !ok {
						fs.Fails[errorString] = 0
					}
					fs.Fails[errorString] = cnt + 1
				case win := <- cWin:
					fs.CurrOffset = win[1]
					fs.CurrPartition =win[0]
				case <-check:
					check = time.After(3 * time.Second)
					//fmt.Printf("Collector %s sending stats!\n", id)
					back <- fs
				default:
					select {
						case <-ctx.Done():
							fmt.Printf("Collector %s shutting down!\n", id)
							back <- fs
							return
						default:	
					}				
					time.Sleep(10 * time.Millisecond)
			}
		}
	}()
}

func main() {
	flag.Parse()
	fmt.Printf("Producing for topic %s! Sending to %s\n", topic, hostname)
	var wg sync.WaitGroup
	cancellers := make(map[string][]func())
	stats := make(chan FINStat, 1024)
	for i := 0; i < 8; i++ {
		id := fmt.Sprintf("%d", i)
		ctx, cancel := context.WithCancel(context.Background())
		ctx2, cancel2 := context.WithCancel(context.Background())
		cancellers[id] = []func(){cancel, cancel2}
		
		wg.Add(1)
		cWin, cSent, cFail := runner(ctx, id, []func(){wg.Done})
		
		wg.Add(1)
		collector(ctx2, id, cWin, cSent, cFail, stats, []func(){wg.Done})
	}
	
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	
	loop := true
	wg.Add(1)
	go func() {
		defer wg.Done()
		for loop {
			select  {
				case <-sigterm:
					fmt.Printf("Got Shutdown signal\n")
					loop = false;
				case next := <-stats:
					fmt.Printf("ID:%s Sent:%d CurrOffset:%d CurrPartition:%d Fails:%v\n", next.ID, next.Sent, next.CurrOffset, next.CurrPartition, next.Fails)
					loop = next.Sent <= 100
				default:
					time.Sleep(1)
			}
		}
		for k, v := range cancellers {
			fmt.Printf("Running cancel callbacks for %s...", k)
			for _, x := range v {
				x()
				fmt.Printf(" ok")
			}
			fmt.Printf("\n")
		}		
		
		loop = true
		check := time.After(2 * time.Second)
		report := make(map[string]*FINStat)
		for loop {
			select {
				case <-sigterm:
					fmt.Printf("Got Shutdown signal\n")
					loop = false;				
				case next := <-stats:
					report[next.ID] = &next
				case <-check:
					loop = false
				default:
					time.Sleep(10 * time.Millisecond)
			}
		}		
		fmt.Printf("------------------------------------\n")
		for _, next := range report {
			fmt.Printf("ID:%s Sent:%d CurrOffset:%d CurrPartition:%d Fails:%v\n", next.ID, next.Sent, next.CurrOffset, next.CurrPartition, next.Fails)
		}		
	}()

	fmt.Printf("Waiting...\n")
	wg.Wait()
	fmt.Printf("Done!\n")

}