package main

import (
	"errors"
	"io"
	"os"

	"github.com/michalkurzeja/proton"
	"github.com/spf13/cobra"
)

func init() {
	convertCmd.Flags().StringVarP(&inputEncoding, flagInputEncoding, "i", "", "input messages encoding [binary, json, text] (required)")
	convertCmd.Flags().StringVarP(&outputEncoding, flagOutputEncoding, "o", "", "output messages encoding [json, proto, text] (required)")
	convertCmd.Flags().StringVarP(&messageName, flagMessageName, "m", "", "proto message name")
	convertCmd.Flags().StringVarP(&delimiter, flagDelimiter, "d", "", "input messages delimiter (0 or 1 characters; if unset, the input is treated as a single message)")
	convertCmd.Flags().StringVarP(&filter, flagFilter, "f", "", "output message filter")

	convertCmd.MarkFlagRequired(flagInputEncoding)
	convertCmd.MarkFlagRequired(flagOutputEncoding)
	convertCmd.MarkFlagRequired(flagMessageName)
}

// Flag names.
const (
	flagInputEncoding  = "input-encoding"
	flagOutputEncoding = "output-encoding"
	flagMessageName    = "message"
	flagDelimiter      = "delimiter"
	flagFilter         = "filter"
)

// Config values.
var (
	inputEncoding  string
	outputEncoding string
	messageName    string
	delimiter      string
	filter         string
)

var convertCmd = &cobra.Command{
	Use: "convert",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(delimiter) > 1 {
			return errors.New("convert: delimiter cannot be longer than 1 character")
		}
		return nil
	},
	RunE: runConvertCmd,
}

func runConvertCmd(cmd *cobra.Command, args []string) error {
	var delim byte
	if len(delimiter) > 0 {
		delim = delimiter[0]
	}

	cfg := proton.ConverterConfig{
		GlobalConfig: proton.GlobalConfig{
			ImportPaths: importPaths,
			ProtoFile:   protoFile,
		},
		InputEncoding:  proton.Encoding(inputEncoding),
		OutputEncoding: proton.Encoding(outputEncoding),
		MessageName:    messageName,
		Delimiter:      delim,
		Filter:         proton.FilterType(filter),
	}

	converter := proton.NewConverter(cfg)

	in, err := getInputPipe()
	if err != nil {
		return err
	}

	return converter.Convert(in, os.Stdout)
}

// getInputPipe returns an io.Reader or error if there is no open input pipe.
func getInputPipe() (io.Reader, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		return nil, errors.New("convert: no input")
	}

	return os.Stdin, nil
}
