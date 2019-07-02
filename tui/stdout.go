package tui

import (
	"encoding/json"
	"fmt"
	"github.com/nuweba/faasbenchmark/report/output/graph"
	"io"
	"os"
)

type StdoutStream struct {
	*os.File
	ch chan []byte
}

func (s *StdoutStream) Write(p []byte) (n int, err error) {
	s.ch <- p

	return len(p), nil
}

func NewsStdoutStream(f *os.File) *StdoutStream {
	s := &StdoutStream{f, make(chan []byte, 100)}
	go func() {
		io.Copy(s, f)
	}()
	return s
}

func hookStdout() (*StdoutStream, error) {
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	output := NewsStdoutStream(pr)
	os.Stdout = pw

	return output, nil
}


type graphStream struct {
	ch chan *graph.Result
}

func (gs *graphStream) parse(p []byte) *graph.Result {
	result := &graph.Result{}
	err := json.Unmarshal(p, result)
	if err != nil {
		panic(fmt.Sprintf("cant parse graph result, %s", p))
	}
	//fields := strings.Fields(string(p))
	//
	//f, err := strconv.ParseFloat(fields[1], 64)
	//if err != nil {
	//	fmt.Println(err)
	//}
	return result

}

func (gs *graphStream) Write(p []byte) (n int, err error) {
	gs.ch <- gs.parse(p)
	return len(p), nil
}