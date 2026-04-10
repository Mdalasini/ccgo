package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"unicode/utf8"
)

func usage(toStdErr bool) {
	var output io.Writer
	if toStdErr {
		output = os.Stderr
	} else {
		output = os.Stdout
	}
	fmt.Fprintln(output, "Usage: wc <flags> <file>")
	fmt.Fprintln(output, "Flags:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(output, "  -%s\t%s\n", f.Name, f.Usage)
	})
}

func countAll(r io.Reader, stats *FileStats) error {
	buf := make([]byte, 32*1024)
	inWord := false
	ByteCount, LineCount, WordCount, CharCount := 0, 0, 0, 0

	for {
		n, e := r.Read(buf)
		if n > 0 {
			ByteCount += n

			for _, b := range buf[:n] {
				// lines
				if b == '\n' {
					LineCount++
				}

				// words (simple ASCII definition)
				if isSpace(b) {
					inWord = false
				} else {
					if !inWord {
						WordCount++
						inWord = true
					}
				}
			}

			// chars (runes)
			CharCount += utf8.RuneCount(buf[:n])
		}

		if e == io.EOF {
			stats.ByteCount = ByteCount
			stats.LineCount = LineCount
			stats.WordCount = WordCount
			stats.CharCount = CharCount
			return nil
		}
		if e != nil {
			return e
		}
	}
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\r'
}

type FileStats struct {
	ByteCount int
	LineCount int
	WordCount int
	CharCount int
	FileName  string
}

func (fs *FileStats) print(printLines, printWords, printBytes, printChars, printFilename bool) {
	if printLines {
		fmt.Printf("%d ", fs.LineCount)
	}
	if printWords {
		fmt.Printf("%d ", fs.WordCount)
	}
	if printBytes {
		fmt.Printf("%d ", fs.ByteCount)
	}
	if printChars {
		fmt.Printf("%d ", fs.CharCount)
	}
	if printFilename {
		fmt.Printf("%s\n", fs.FileName)
	} else {
		fmt.Println()
	}
}

func (fs *FileStats) countStdin() error {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		err := countAll(os.Stdin, fs)
		if err != nil {
			return fmt.Errorf("wc: error counting: %v", err)
		}
	}
	return nil
}

func (fs *FileStats) countFile() error {
	file, err := os.Open(fs.FileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("wc: %v: No such file or directory", fs.FileName)
		}
		return fmt.Errorf("wc: error opening %v: %v", fs.FileName, err)
	}
	defer file.Close()

	err = countAll(file, fs)
	if err != nil {
		return fmt.Errorf("wc: error counting %v: %v", fs.FileName, err)
	}

	return nil
}

func main() {
	bytesPtr := flag.Bool("c", false, "print the byte counts")
	linesPtr := flag.Bool("l", false, "print the line counts")
	wordPtr := flag.Bool("w", false, "print the word counts")
	charPtr := flag.Bool("m", false, "print the character counts")
	helpPtr := flag.Bool("h", false, "print this help message")
	flag.Usage = func() { usage(true) }
	flag.Parse()
	args := flag.Args()

	if *helpPtr {
		usage(false)
		return
	}

	if flag.NFlag() == 0 {
		*bytesPtr, *linesPtr, *wordPtr = true, true, true
	}

	if len(args) == 0 {
		stats := FileStats{}
		err := stats.countStdin()
		if err != nil {
			fmt.Println(err)
		}
		stats.print(*linesPtr, *wordPtr, *bytesPtr, *charPtr, false)
		return
	}

	var mu sync.Mutex
	var totalStats FileStats
	totalStats.FileName = "total"

	var wg sync.WaitGroup
	errCh := make(chan error, len(args))
	outputCh := make(chan FileStats, len(args))

	for _, arg := range args {
		wg.Add(1)
		go func(arg string) {
			defer wg.Done()
			stats := FileStats{FileName: arg}
			if err := stats.countFile(); err != nil {
				errCh <- fmt.Errorf("%v", err)
				return
			}
			outputCh <- stats
			mu.Lock()
			totalStats.ByteCount += stats.ByteCount
			totalStats.LineCount += stats.LineCount
			totalStats.WordCount += stats.WordCount
			totalStats.CharCount += stats.CharCount
			mu.Unlock()
		}(arg)
	}

	wg.Wait()
	close(errCh)
	close(outputCh)

	for stats := range outputCh {
		stats.print(*linesPtr, *wordPtr, *bytesPtr, *charPtr, true)
	}
	totalStats.print(*linesPtr, *wordPtr, *bytesPtr, *charPtr, true)
	for err := range errCh {
		fmt.Println(err)
	}
}
