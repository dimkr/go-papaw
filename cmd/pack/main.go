// This file is part of go-papaw.
//
// Copyright (c) 2020 Dima Krasner
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"bytes"
	"compress/flate"
	"debug/elf"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/itchio/lzma"
)

const (
	repoURL         = "https://github.com/dimkr/papaw"
	downloadTimeout = 5 * time.Second
)

func getInputArchitecture(exe *elf.File) string {
	switch exe.FileHeader.Machine {
	case elf.EM_386, elf.EM_X86_64:
		return "i386"

	case elf.EM_ARM, elf.EM_AARCH64:
		if exe.FileHeader.ByteOrder == binary.BigEndian {
			return "armeb"
		}
		return "arm"

	case elf.EM_MIPS:
		if exe.FileHeader.ByteOrder == binary.LittleEndian {
			return "mipsel"
		}

		return "mips"
	}

	return ""
}

func getStub(algo, arch string) ([]byte, error) {
	client := http.Client{Timeout: downloadTimeout}
	response, err := client.Get(repoURL + fmt.Sprintf("/releases/latest/download/papaw-%s-%s", algo, arch))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to download the stub")
	}

	stub, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if _, err = elf.NewFile(bytes.NewReader(stub)); err != nil {
		return nil, err
	}

	return stub, nil
}

func main() {
	var inputPath, outputPath, algo string
	flag.StringVar(&inputPath, "input", "", "executable path")
	flag.StringVar(&outputPath, "output", "", "packed executable path")
	flag.StringVar(&algo, "algo", "deflate", "compression algorithm")
	flag.Parse()
	if inputPath == "" || outputPath == "" {
		flag.Usage()
		return
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	inputData, err := ioutil.ReadAll(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	inputSize := len(inputData)
	if inputSize == 0 || inputSize > math.MaxUint32 {
		log.Fatal("Failed to read the executable")
	}

	exe, err := elf.NewFile(bytes.NewReader(inputData))
	if err != nil {
		log.Fatal(err)
	}

	if exe.FileHeader.Type != elf.ET_EXEC && exe.FileHeader.Type != elf.ET_DYN {
		log.Fatal("Not an executable")
	}

	arch := getInputArchitecture(exe)
	if arch == "" {
		log.Fatal("Unsupported architecture")
	}

	stub, err := getStub(algo, arch)
	if err != nil {
		log.Fatal(err)
	}

	stubSize := len(stub)
	packed := bytes.NewBuffer(stub)

	switch algo {
	case "lzma":
		compressedBuffer := bytes.Buffer{}

		compressor := lzma.NewWriterSizeLevel(&compressedBuffer, int64(inputSize), lzma.BestCompression)

		compressed, err := compressor.Write(inputData)
		if err != nil {
			log.Fatal(err)
		}

		if compressed != inputSize {
			log.Fatal("Failed to compress the executable")
		}

		if err := compressor.Close(); err != nil {
			panic(err)
		}

		packed.Write(compressedBuffer.Bytes()[:5])
		// obfuscation
		packed.WriteByte(8)
		packed.Write(compressedBuffer.Bytes()[5:])

	case "deflate":
		compressor, err := flate.NewWriter(packed, flate.BestCompression)
		if err != nil {
			log.Fatal(err)
		}

		compressed, err := compressor.Write(inputData)
		if err != nil {
			log.Fatal(err)
		}

		if compressed != inputSize {
			log.Fatal("Failed to compress the executable")
		}

		if err := compressor.Flush(); err != nil {
			log.Fatal(err)
		}

		if err := compressor.Close(); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("Unsupported compresion algorithm")
	}

	compressedLength := len(packed.Bytes()) - stubSize

	binary.Write(packed, binary.BigEndian, uint32(inputSize))
	binary.Write(packed, binary.BigEndian, uint32(compressedLength))

	if err := ioutil.WriteFile(outputPath, packed.Bytes(), 0755); err != nil {
		log.Fatal(err)
	}

	log.Printf("Done: output is %.2f%% smaller.", 100.0-float64(len(packed.Bytes()))*100.0/float64(inputSize))
}
