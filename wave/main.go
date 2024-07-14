package main

import (
	"encoding/binary"
	"log"
	"math"
	"os"
)

type sampleT int16
type FourCC [4]byte

const (
	SampleMax = 32767
	Duration  = 5
	SR        = 44100
	Nchannels = 1
	Nsamples  = Nchannels * Duration * SR
)

type riffHdr struct {
	id       FourCC
	size     uint32
	riffType FourCC
}

type fmtCk struct {
	id            FourCC
	size          uint32
	fmtTag        uint16
	channels      uint16
	samplesPerSec uint32
	bytesPerSec   uint32
	blockAlign    uint16
	bitsPerSample uint16
}

type dataHdr struct {
	id   FourCC
	size uint32
}

type wavHdr struct {
	riff riffHdr
	fmt  fmtCk
	data dataHdr
}

var buf = make([]sampleT, Nsamples)

func main() {
	const SampleSize = 2
	hdr := wavHdr{}

	file, err := os.Create("output.wav")
	if err != nil {
		log.Fatal(err)
	}
	//https://www.joeshaw.org/dont-defer-close-on-writable-files/
	defer file.Close()

	copy(hdr.riff.id[:], "RIFF")
	hdr.riff.size = 36 + Nsamples*SampleSize
	copy(hdr.riff.riffType[:], "WAVE")

	copy(hdr.fmt.id[:], "fmt ")
	hdr.fmt.size = 16
	hdr.fmt.fmtTag = 1
	hdr.fmt.channels = Nchannels
	hdr.fmt.samplesPerSec = SR
	hdr.fmt.bytesPerSec = Nchannels * SR * SampleSize
	hdr.fmt.blockAlign = Nchannels * SampleSize
	hdr.fmt.bitsPerSample = 16

	copy(hdr.data.id[:], "data")
	hdr.data.size = Nsamples * SampleSize

	frequency := 440.0
	angularFrequency := 2 * math.Pi * frequency / float64(SR)
	for i := 0; i < Nsamples; i++ {
		buf[i] = sampleT(SampleMax * math.Sin(angularFrequency*float64(i)))
	}
	err = binary.Write(file, binary.LittleEndian, hdr)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(file, binary.LittleEndian, buf)
	if err != nil {
		log.Fatal(err)
	}
	if len(buf)%2 == 1 {
		nilByte := byte(0)
		_, err = file.Write([]byte{nilByte})
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("File written")
}
