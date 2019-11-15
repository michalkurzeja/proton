package proton

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

const (
	FormatJSON  Encoding = "json"
	FormatProto Encoding = "proto"
	FormatText  Encoding = "text"
)

// Encoding
type Encoding string

// Unmarshal decodes the bytes using the given encoding.
func Unmarshal(b []byte, msg *dynamic.Message, e Encoding) error {
	switch e {
	case FormatJSON:
		return msg.UnmarshalJSON(b)
	case FormatProto:
		return msg.Unmarshal(b)
	case FormatText:
		return msg.UnmarshalText(b)
	}

	return fmt.Errorf("unmarshal: unknown encoding: %s", e)
}

// Marshal encodes the message with the given encoding.
func Marshal(msg *dynamic.Message, e Encoding) ([]byte, error) {
	switch e {
	case FormatJSON:
		return msg.MarshalJSON()
	case FormatProto:
		return msg.Marshal()
	case FormatText:
		return msg.MarshalText()
	}

	return nil, fmt.Errorf("marshal: unknown encoding: %s", e)
}

// ConverterConfig represents the configuration of the converter.
type ConverterConfig struct {
	GlobalConfig

	InputEncoding  Encoding
	OutputEncoding Encoding
	MessageName    string
	Delimiter      byte
	Filter         FilterType
}

// Converter is a service responsible for converting proto message encoding formats.
type Converter struct {
	cfg    ConverterConfig
	parser protoparse.Parser
}

// NewConverter returns a new Converter instance.
func NewConverter(cfg ConverterConfig) Converter {
	return Converter{
		cfg:    cfg,
		parser: protoparse.Parser{ImportPaths: cfg.ImportPaths},
	}
}

// Convert converts the input from the reader and writes it to the writer .
func (c Converter) Convert(in io.Reader, out io.Writer) error {
	resolved, err := protoparse.ResolveFilenames(c.cfg.ImportPaths, c.cfg.ProtoFile) // Always returns 1 resolved filename.
	if err != nil {
		return err
	}

	descriptors, err := c.parser.ParseFiles(resolved[0]) // Always returns 1 descriptor.
	if err != nil {
		return err
	}

	reader := bufio.NewReader(in)

	b, eof, err := c.readNext(reader)

	for {
		err = c.convert(descriptors[0], b, out)
		if err != nil {
			return err
		}

		if eof {
			break
		}

		b, eof, err = c.readNext(reader)
	}

	return nil
}

// readNext returns a slice of bytes containing  the next input message.
// The reading ends upon encountering the delimiter or EOF.
func (c Converter) readNext(r *bufio.Reader) (b []byte, eof bool, err error) {
	b, err = r.ReadBytes(c.cfg.Delimiter)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, true, err
		}

		return b, true, nil
	}

	return b, false, nil
}

// convert performs the encoding conversion and writes the output to the given writer.
func (c Converter) convert(d *desc.FileDescriptor, b []byte, out io.Writer) error {
	msg := dynamic.NewMessage(d.FindMessage(c.cfg.MessageName))

	err := Unmarshal(b, msg, c.cfg.InputEncoding)
	if err != nil {
		return err
	}

	b, err = Marshal(msg, c.cfg.OutputEncoding)
	if err != nil {
		return err
	}

	if c.cfg.Filter != "" {
		b, err = Filter(b, c.cfg.Filter)
		if err != nil {
			return err
		}
	}

	_, err = out.Write(b)
	if err != nil {
		return err
	}

	if c.cfg.Delimiter != 0x0 {
		_, err = out.Write([]byte{c.cfg.Delimiter})
		if err != nil {
			return err
		}
	}

	return nil
}
