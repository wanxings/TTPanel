package terminal

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

type RecType string

const (
	InputType  RecType = "i"
	OutPutType RecType = "o"
)

type RecHeader struct {
	Version   int   `json:"version"`
	Width     int   `json:"width"`
	Height    int   `json:"height"`
	Timestamp int64 `json:"timestamp"`
	Env       struct {
		Shell string `json:"SHELL"`
		Term  string `json:"TERM"`
	} `json:"env"`
}

func defaultRecHeader() *RecHeader {
	recHeader := new(RecHeader)
	recHeader.Version = 2
	recHeader.Env.Shell = "/bin/bash"
	recHeader.Env.Term = "xterm-256color"
	return recHeader
}

type Recorder struct {
	StartTime time.Time
	Writer    io.Writer
	sync.Mutex
}

func NewRecorder(writer io.Writer) *Recorder {
	return &Recorder{
		StartTime: time.Now(),
		Writer:    writer,
	}
}

func (rec *Recorder) WriteHeader(height, width int) {
	header := defaultRecHeader()
	header.Timestamp = rec.StartTime.Unix()
	header.Height = height
	header.Width = width
	b, _ := json.Marshal(header)
	_, _ = rec.Writer.Write(b)
	_, _ = rec.Writer.Write([]byte("\r\n"))
}

func (rec *Recorder) WriteData(recType RecType, data string) {
	recData := make([]interface{}, 3)
	recData[0] = float64(time.Since(rec.StartTime).Microseconds()) / float64(1000000)
	recData[1] = recType
	recData[2] = data
	b, _ := json.Marshal(recData)
	_, _ = rec.Writer.Write(b)
	_, _ = rec.Writer.Write([]byte("\r\n"))
}
