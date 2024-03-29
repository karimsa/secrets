package main

import (
	"fmt"
	"os"

	"github.com/karimsa/secrets"
	"github.com/urfave/cli/v2"
)

var cmdEncrypt = &cli.Command{
	Name:    "encrypt",
	Aliases: []string{"enc"},
	Usage:   "Encrypt values in a given file",
	Flags: []cli.Flag{
		inFlag,
		outFlag,
		formatFlag,
		strategyFlag,
		passphraseFlag,
		keyFlag,
		keyFileFlag,
		flagLogLevel,
	},
	Action: func(ctx *cli.Context) error {
		format := ctx.String("format")
		inPath := ctx.String("in")

		securePaths, err := getInputPaths(ctx)
		if err != nil {
			return err
		}

		if format == "" {
			format = getFormatFromPath(inPath)
		}

		inFile, err := os.OpenFile(inPath, os.O_RDONLY, 0)
		if err != nil {
			return err
		}

		cipher, err := getCipher(ctx)
		if err != nil {
			return err
		}

		logLevel, err := getLogLevel(ctx)
		if err != nil {
			return err
		}

		envFile, err := secrets.New(secrets.NewEnvOptions{
			Format:      format,
			Reader:      inFile,
			Cipher:      cipher,
			LogLevel:    logLevel,
			SecurePaths: securePaths,
		})
		if err != nil {
			return err
		}

		buff, err := envFile.Export(format)
		if err != nil {
			return err
		}

		outPath := ctx.String("out")
		switch outPath {
		case "/dev/stdout":
			fmt.Printf(string(buff))
			return nil
		case "/dev/stderr":
			fmt.Fprintf(os.Stderr, string(buff))
			return nil
		}

		// For in-place edits, overwrite the file
		outFileMode := os.O_EXCL
		if outPath == inPath {
			outFileMode = os.O_WRONLY | os.O_TRUNC
		}
		return envFile.ExportFile(format, outPath, outFileMode)
	},
}
