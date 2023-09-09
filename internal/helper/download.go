package helper

import (
	"TTPanel/internal/global"
	"TTPanel/pkg/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type WriteCounter struct {
	Total   uint64
	Written uint64
	Key     string
	Name    string
}

type Process struct {
	Total   uint64  `json:"total"`
	Written uint64  `json:"written"`
	Percent float64 `json:"percent"`
	Name    string  `json:"name"`
}

func (w *WriteCounter) Write(p []byte) (n int, err error) {
	n = len(p)
	w.Written += uint64(n)
	w.SaveProcess()
	return n, nil
}

func (w *WriteCounter) SaveProcess() {
	percentValue := 0.0
	if w.Total > 0 {
		percent := float64(w.Written) / float64(w.Total) * 100
		percentValue, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", percent), 64)
	}
	process := Process{
		Total:   w.Total,
		Written: w.Written,
		Percent: percentValue,
		Name:    w.Name,
	}
	by, _ := json.Marshal(process)
	if percentValue < 100 {
		global.GoCache.Set(w.Key, string(by), -1)
	} else {
		global.GoCache.Set(w.Key, string(by), time.Second*time.Duration(60))
	}
}
func AsyncDownloadFile(url, savePath, key string) error {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	request.Header.Set("Accept-Encoding", "identity")
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	go func() {
		counter := &WriteCounter{}
		counter.Key = key
		if resp.ContentLength > 0 {
			counter.Total = uint64(resp.ContentLength)
		}
		counter.Name = filepath.Base(savePath)
		if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
			global.Log.Errorf("AsyncDownloadFile->io.Copy Eroor:%s \n", err.Error())
		}
		_ = out.Close()
		_ = resp.Body.Close()

		value, ok := global.GoCache.Get(counter.Key)
		if !ok {
			global.Log.Errorf("AsyncDownloadFile -> global.GoCache.Get nil")
			return
		}
		process := &Process{}
		_ = util.JsonStrToStruct(value.(string), process)
		process.Percent = 100
		process.Name = counter.Name
		process.Total = process.Written
		by, _ := json.Marshal(process)
		global.GoCache.Set(counter.Key, string(by), time.Second*time.Duration(60))
	}()
	return nil
}
