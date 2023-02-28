package main

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	//CutVideo("test.mp4", "test_cut.mp4", "00:00:40", "00:00:50")
	StartStream("https://demiroren.daioncdn.net/kanald/kanald_480p.m3u8?&sid=56t16au2wka1&app=da2109ea-5dfe-4107-89ab-23593336ed61&ce=3", "test_%Y%m%d%H%M%S.mp4", "30")
}

func CutVideo(sorucePath string, destPath string, startTime string, endTime string) {
	ffmpeg.Input(sorucePath, ffmpeg.KwArgs{"ss": startTime, "to": endTime, "async": "1", "strict": "-2"}).Output(destPath).OverWriteOutput().ErrorToStdOut().Run()

}

func StartStream(source string, dest string, segmentTime string) {
	ffmpeg.Input(source).Output(dest, ffmpeg.KwArgs{"b:v": "1M", "f": "segment", "segment_time": segmentTime, "reset_timestamps": "1", "strftime": "1"}).ErrorToStdOut().Run()
}
