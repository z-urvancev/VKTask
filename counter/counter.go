package counter

import (
	"TestVK/analyze"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"sync"
	"sync/atomic"
)

type Counter struct {
	ch     chan struct{}
	wg     *sync.WaitGroup
	writer io.Writer
}

func NewCounter(k uint, writer io.Writer) *Counter {
	return &Counter{
		ch:     make(chan struct{}, k), //инициализируем канал с буфером размера k
		wg:     &sync.WaitGroup{},
		writer: writer,
	}
}

func (c *Counter) Execute(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	var total uint64
	for scanner.Scan() {
		str := scanner.Text()
		c.ch <- struct{}{} //можем писать в канал, пока его буфер не заполнится (не запустится k горутин)
		c.wg.Add(1)
		go func(str string, c *Counter) {
			defer func() {
				<-c.ch //освобождаем место для следующей горутины
				c.wg.Done()
			}()
			count, err := counting(str)
			if err != nil {
				log.Fatalln("counting error: ", err)
			}
			atomic.AddUint64(&total, count)
			_, err = fmt.Fprintf(c.writer, "Count for %s: %d\n", str, count)
			if err != nil {
				log.Fatalln("output error: ", err)
			}
		}(str, c)
	}

	if scanner.Err() != nil {
		return fmt.Errorf("scanner error: %w", scanner.Err())
	}

	c.wg.Wait()
	_, err := fmt.Fprintln(c.writer, "Total: ", total)
	if err != nil {
		return fmt.Errorf("output error: %w", err)
	}

	return scanner.Err()
}

func counting(str string) (uint64, error) {
	_, isUrl := url.ParseRequestURI(str)
	if isUrl == nil {
		return analyze.CountingURL(str)
	}
	return analyze.CountingFile(str)
}
