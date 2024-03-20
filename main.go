package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	defer timeTrack(time.Now(), "Modulation and Demodulation")
	pathFile := "input_video.mp4"
	level := 64
	noise := 0.20
	chunkSize := 8192 // size of each chunk in bytes

	M := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(level)), nil)
	fmt.Println("Modulation level M: ", M)

	file, err := os.Open(pathFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	data := make([]byte, chunkSize)

	for {
		_, err := io.ReadFull(file, data)
		if err == io.EOF {
			break
		} else if err != nil && err != io.ErrUnexpectedEOF {
			fmt.Println("Error reading file:", err)
			return
		}

		bits := bytesToBits(data)
		pathModulatePoints := modulate(bits, M, pathFile+"-modulation.csv")
		bitsRestore := demodulate(pathModulatePoints, M, noise)
		originalBytes := bitsToBytes(bitsRestore)
		writeRestoreFile("restore_"+pathFile, originalBytes)
	}
	fmt.Println("Successfully demodulated the data!")
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func writeRestoreFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}

func bytesToBits(data []byte) []int64 {
	var bits []int64
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bit := (b >> i) & 1
			bits = append(bits, int64(bit))
		}
	}
	return bits
}

func bitsToBytes(bits []int64) []byte {
	var bytes []byte
	for i := 0; i < len(bits); i += 8 {
		b := 0
		for j := 0; j < 8; j++ {
			b = b << 1
			if bits[i+j] == 1 {
				b = b | 1
			}
		}
		bytes = append(bytes, byte(b))
	}
	return bytes
}

func modulate(bits []int64, M *big.Int, filePath string) string {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	bitsPerSymbol := M.BitLen() - 1
	for i := 0; i < len(bits); i += bitsPerSymbol {
		xBits := bits[i : i+bitsPerSymbol/2]
		yBits := bits[i+bitsPerSymbol/2 : i+bitsPerSymbol]

		x := bitsToInt(xBits)
		y := bitsToInt(yBits)

		err := writer.Write([]string{strconv.Itoa(int(x)), strconv.Itoa(int(y))})
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}
	return filePath
}

func bitsToInt(bits []int64) int64 {
	val := int64(0)
	for _, bit := range bits {
		val = (val << 1) | bit
	}
	return val
}

func demodulate(filePath string, M *big.Int, noiseStdDev float64) []int64 {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var bits []int64
	bitsPerSymbol := M.BitLen() - 1

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	for _, record := range records {
		x, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			fmt.Println("Error parsing x coordinate:", err)
			return nil
		}

		y, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			fmt.Println("Error parsing y coordinate:", err)
			return nil
		}

		x += int64(rand.NormFloat64() * noiseStdDev)
		y += int64(rand.NormFloat64() * noiseStdDev)

		xBits := intToBits(x, bitsPerSymbol/2)
		yBits := intToBits(y, bitsPerSymbol/2)

		bits = append(bits, xBits...)
		bits = append(bits, yBits...)
	}

	return bits
}

func intToBits(n int64, bitCount int) []int64 {
	bits := make([]int64, bitCount)
	for i := range bits {
		bits[bitCount-i-1] = n & 1
		n = n >> 1
	}
	return bits
}
