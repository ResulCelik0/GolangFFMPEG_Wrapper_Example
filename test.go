package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {

	CutVideo("test.mp4", "test_cut.mp4", "00:00:40", "00:00:50")
	StartStream("https://demiroren.daioncdn.net/kanald/kanald_480p.m3u8?&sid=56t16au2wka1&app=da2109ea-5dfe-4107-89ab-23593336ed61&ce=3", "test_%Y%m%d%H%M%S.mp4", "30")
}

func CutVideo(sorucePath string, destPath string, startTime string, endTime string) {

	writer := &logWriter{
		buffer:    &bytes.Buffer{},
		lastLines: make([][]byte, 0),
	}
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Dosya oluşturma hatası:", err)
		return
	}
	defer file.Close()

	ffmpeg.Input(sorucePath, ffmpeg.KwArgs{"ss": startTime, "to": endTime, "async": "1", "strict": "-2"}).Output(destPath).OverWriteOutput().WithErrorOutput(io.MultiWriter(writer)).Run()

	for _, line := range writer.LastLines() {
		file.Write(line)
	}
}

func StartStream(source string, dest string, segmentTime string) {
	writer := &logWriter{
		buffer:    bytes.NewBuffer(make([]byte, 1024)),
		lastLines: make([][]byte, 0),
	}
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Dosya oluşturma hatası:", err)
		return
	}
	defer file.Close()
	ffmpeg.Input(source).Output(dest, ffmpeg.KwArgs{"b:v": "1M", "f": "segment", "segment_time": segmentTime, "reset_timestamps": "1", "strftime": "1"}).WithErrorOutput(io.MultiWriter(writer)).Run()

}

var ErrorWords = []string{"error", "Error", "ERROR", "failed", "Failed", "FAILED", "fatal", "Fatal", "FATAL", "unable", "Unable", "UNABLE", "invalid", "Invalid", "INVALID", "not", "Not", "NOT", "cannot", "Cannot", "CANNOT", "could not", "Could not", "COULD NOT", "no such", "No such", "NO SUCH", "no input", "No input", "NO INPUT", "no output", "No output", "NO OUTPUT", "no stream", "No stream", "NO STREAM", "no file", "No file", "NO FILE", "no such file", "No such file", "NO SUCH FILE", "no such stream", "No such stream", "NO SUCH STREAM", "no such device", "No such device", "NO SUCH DEVICE", "no such filter", "No such filter", "NO SUCH FILTER", "no such option", "No such option", "NO SUCH OPTION", "no such channel", "No such channel", "NO SUCH CHANNEL", "no such field", "No such field", "NO SUCH FIELD", "no such property", "No such property", "NO SUCH PROPERTY", "no such method", "No such method", "NO SUCH METHOD", "no such class", "No such class", "NO SUCH CLASS", "no such function", "No such function", "NO SUCH FUNCTION", "no such frame", "No such frame", "NO SUCH FRAME", "no such sample", "No such sample", "NO SUCH SAMPLE", "no such device"}

type logWriter struct {
	buffer *bytes.Buffer
	Writer *io.Writer
}

func (w *logWriter) Write(p []byte) (int, error) {
	n, err := w.buffer.Write(p)
	if err != nil {
		return n, err
	}
	if w.buffer.Len() > 1024 {
		w.buffer.Next(w.buffer.Len() - 1024)
	}
	for _, key := range ErrorWords {
		if strings.Contains(string(p), key) {
			break
		}
	}
	return n, nil
}
