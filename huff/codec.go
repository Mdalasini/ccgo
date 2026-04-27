package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"slices"
)

// BitWriter accumulates individual bits and writes them to an underlying
// io.Writer in whole-byte chunks. Call Flush after writing all bits to pad the
// final byte and return the number of padding bits.
type BitWriter struct {
	w     io.Writer
	buf   byte
	count uint8
}

func NewBitWriter(w io.Writer) *BitWriter {
	return &BitWriter{w: w}
}

// WriteBit appends a single bit to the buffer. bit must be 0 or 1.
func (bw *BitWriter) WriteBit(bit uint8) error {
	bw.buf = (bw.buf << 1) | (bit & 1)
	bw.count++
	if bw.count == 8 {
		if _, err := bw.w.Write([]byte{bw.buf}); err != nil {
			return err
		}
		bw.buf = 0
		bw.count = 0
	}
	return nil
}

// WriteBits writes every character of the string bits (which must consist of
// '0' and '1') to the stream.
func (bw *BitWriter) WriteBits(bits string) error {
	for _, b := range bits {
		if err := bw.WriteBit(uint8(b - '0')); err != nil {
			return err
		}
	}
	return nil
}

// Flush pads the current partial byte with zeros and writes it to the
// underlying writer. It returns the number of padding bits that were added.
func (bw *BitWriter) Flush() (uint8, error) {
	if bw.count == 0 {
		return 0, nil
	}
	bw.buf <<= (8 - bw.count)
	if _, err := bw.w.Write([]byte{bw.buf}); err != nil {
		return 0, err
	}
	padding := 8 - bw.count
	bw.buf = 0
	bw.count = 0
	return padding, nil
}

// BitReader reads individual bits from an underlying io.Reader.
type BitReader struct {
	r     io.Reader
	buf   byte
	count uint8
	err   error
}

func NewBitReader(r io.Reader) *BitReader {
	return &BitReader{r: r}
}

// ReadBit returns the next bit (0 or 1) from the stream. On EOF the returned
// error is io.EOF; other errors are propagated from the underlying reader.
func (br *BitReader) ReadBit() (uint8, error) {
	if br.err != nil {
		return 0, br.err
	}
	if br.count == 0 {
		var b [1]byte
		n, err := br.r.Read(b[:])
		if n == 0 {
			if err == nil {
				err = io.EOF
			}
			br.err = err
			return 0, err
		}
		br.buf = b[0]
		br.count = 8
	}
	br.count--
	bit := (br.buf >> br.count) & 1
	return bit, nil
}

// BuildCodeLengths traverses the Huffman tree and returns the bit-length of
// each symbol's prefix code.
func BuildCodeLengths(root *HuffNode) map[rune]int {
	lengths := make(map[rune]int)
	if root == nil {
		return lengths
	}
	if root.isLeaf() {
		lengths[root.Char] = 1
		return lengths
	}
	var walk func(n *HuffNode, depth int)
	walk = func(n *HuffNode, depth int) {
		if n.isLeaf() {
			lengths[n.Char] = depth
			return
		}
		if n.Left != nil {
			walk(n.Left, depth+1)
		}
		if n.Right != nil {
			walk(n.Right, depth+1)
		}
	}
	walk(root, 0)
	return lengths
}

// BuildCanonicalCodes computes canonical Huffman codes from a set of code
// lengths. Symbols with the same length are ordered by their rune value.
func BuildCanonicalCodes(lengths map[rune]int) map[rune]string {
	if len(lengths) == 0 {
		return nil
	}

	maxLen := 0
	for _, l := range lengths {
		if l > maxLen {
			maxLen = l
		}
	}

	symbolsByLen := make([][]rune, maxLen+1)
	for r, l := range lengths {
		symbolsByLen[l] = append(symbolsByLen[l], r)
	}
	for _, symbols := range symbolsByLen {
		slices.Sort(symbols)
	}

	codes := make(map[rune]string)
	code := uint32(0)
	for l := 1; l <= maxLen; l++ {
		for _, r := range symbolsByLen[l] {
			codes[r] = formatBits(code, l)
			code++
		}
		code <<= 1
	}
	return codes
}

func formatBits(code uint32, length int) string {
	bits := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		bits[i] = '0' + byte(code&1)
		code >>= 1
	}
	return string(bits)
}

// Encode reads uncompressed data from r, compresses it with Huffman coding,
// and writes the result to w.
//
// Binary format:
//
//	uint16  number of unique symbols
//	for each symbol:
//	  int32   rune
//	  uint8   code length (in bits)
//	uint32  uncompressed rune count
//	uint8   padding bits in the final byte
//	[...]   compressed bitstream
func Encode(r io.Reader, w io.Writer) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	runes := []rune(string(data))

	freq := make(map[rune]int)
	for _, ch := range runes {
		freq[ch]++
	}

	root := BuildHuffTree(freq)
	lengths := BuildCodeLengths(root)
	codes := BuildCanonicalCodes(lengths)

	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, uint16(len(codes)))

	for r, code := range codes {
		binary.Write(&buf, binary.LittleEndian, r)
		binary.Write(&buf, binary.LittleEndian, uint8(len(code)))
	}

	binary.Write(&buf, binary.LittleEndian, uint32(len(runes)))

	paddingIdx := buf.Len()
	buf.WriteByte(0)

	bw := NewBitWriter(&buf)
	for _, r := range runes {
		if err := bw.WriteBits(codes[r]); err != nil {
			return err
		}
	}

	padding, err := bw.Flush()
	if err != nil {
		return err
	}

	buf.Bytes()[paddingIdx] = padding

	_, err = w.Write(buf.Bytes())
	return err
}

// Decode reads Huffman-compressed data from r and writes the decompressed
// output to w. The format is the same as produced by Encode.
func Decode(r io.Reader, w io.Writer) error {
	var numCodes uint16
	if err := binary.Read(r, binary.LittleEndian, &numCodes); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	if numCodes == 0 {
		return nil
	}

	lengths := make(map[rune]int, numCodes)
	for i := uint16(0); i < numCodes; i++ {
		var ch rune
		var length uint8
		if err := binary.Read(r, binary.LittleEndian, &ch); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
			return err
		}
		lengths[ch] = int(length)
	}

	var runeCount uint32
	if err := binary.Read(r, binary.LittleEndian, &runeCount); err != nil {
		return err
	}

	var padding uint8
	if err := binary.Read(r, binary.LittleEndian, &padding); err != nil {
		return err
	}

	codes := BuildCanonicalCodes(lengths)
	decodeTable := make(map[string]rune, len(codes))
	for r, code := range codes {
		decodeTable[code] = r
	}

	br := NewBitReader(r)
	var prefix []byte
	var decoded uint32

	for decoded < runeCount {
		bit, err := br.ReadBit()
		if err != nil {
			if err == io.EOF && prefix == nil {
				break
			}
			return err
		}
		prefix = append(prefix, '0'+bit)

		if ch, ok := decodeTable[string(prefix)]; ok {
			if _, err := io.WriteString(w, string(ch)); err != nil {
				return err
			}
			prefix = prefix[:0]
			decoded++
		}
	}

	return nil
}
