# Go Modulation and Demodulation

This Go program reads a file, modulates its data into symbols, adds noise to the symbols, demodulates the noisy symbols back into data, and writes the restored data to a new file.

## How it Works

1. The program reads a file in chunks of a specified size.
2. Each chunk of data is converted into bits.
3. The bits are modulated into symbols and written to a CSV file.
4. The symbols in the CSV file are read, noise is added to them, and they are demodulated back into bits.
5. The restored bits are converted back into bytes and written to a new file.

## Running the Program

To run the program, you need to have Go installed on your machine. You can download Go from the [official website](https://golang.org/dl/).

Once you have Go installed, you can run the program with the following command:

```bash
go run main.go

