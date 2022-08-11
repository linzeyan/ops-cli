/*
Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bytes"
	"log"
	"os"
	"strings"

	hashtag "github.com/abhinav/goldmark-hashtag"
	mermaid "github.com/abhinav/goldmark-mermaid"
	toc "github.com/abhinav/goldmark-toc"
	"github.com/spf13/cobra"
	pdf "github.com/stephenafamo/goldmark-pdf"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/storage"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert data format",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

/* CSV. */
var convertSubCmdCSV2JSON = &cobra.Command{
	Use:   "csv2json",
	Short: "Convert csv to json format",
	Run:   cf.Run,
}
var convertSubCmdCSV2TOML = &cobra.Command{
	Use:   "csv2toml",
	Short: "Convert csv to toml format",
	Run:   cf.Run,
}
var convertSubCmdCSV2XML = &cobra.Command{
	Use:   "csv2xml",
	Short: "Convert csv to xml format",
	Run:   cf.Run,
}
var convertSubCmdCSV2YAML = &cobra.Command{
	Use:   "csv2yaml",
	Short: "Convert csv to yaml format",
	Run:   cf.Run,
}

/* JSON. */
var convertSubCmdJSON2CSV = &cobra.Command{
	Use:   "json2csv",
	Short: "Convert json to csv format",
	Run:   cf.Run,
}
var convertSubCmdJSON2TOML = &cobra.Command{
	Use:   "json2toml",
	Short: "Convert json to toml format",
	Run:   cf.Run,
}
var convertSubCmdJSON2XML = &cobra.Command{
	Use:   "json2xml",
	Short: "Convert json to xml format",
	Run:   cf.Run,
}
var convertSubCmdJSON2YAML = &cobra.Command{
	Use:   "json2yaml",
	Short: "Convert json to yaml format",
	Run:   cf.Run,
}

/* Markdown. */
var convertSubCmdMarkdown2HTML = &cobra.Command{
	Use:   "markdown2html",
	Short: "Convert markdown to html format",
	Run:   cf.Run,
}
var convertSubCmdMarkdown2PDF = &cobra.Command{
	Use:   "markdown2pdf",
	Short: "Convert markdown to pdf format",
	Run:   cf.Run,
}

/* TOML. */
var convertSubCmdTOML2CSV = &cobra.Command{
	Use:   "toml2csv",
	Short: "Convert toml to csv format",
	Run:   cf.Run,
}
var convertSubCmdTOML2JSON = &cobra.Command{
	Use:   "toml2json",
	Short: "Convert toml to json format",
	Run:   cf.Run,
}
var convertSubCmdTOML2XML = &cobra.Command{
	Use:   "toml2xml",
	Short: "Convert toml to xml format",
	Run:   cf.Run,
}
var convertSubCmdTOML2YAML = &cobra.Command{
	Use:   "toml2yaml",
	Short: "Convert toml to yaml format",
	Run:   cf.Run,
}

/* XML. */
var convertSubCmdXML2CSV = &cobra.Command{
	Use:   "xml2csv",
	Short: "Convert xml to csv format",
	Run:   cf.Run,
}
var convertSubCmdXML2JSON = &cobra.Command{
	Use:   "xml2json",
	Short: "Convert xml to json format",
	Run:   cf.Run,
}
var convertSubCmdXML2TOML = &cobra.Command{
	Use:   "xml2toml",
	Short: "Convert xml to toml format",
	Run:   cf.Run,
}
var convertSubCmdXML2YAML = &cobra.Command{
	Use:   "xml2yaml",
	Short: "Convert xml to yaml format",
	Run:   cf.Run,
}

/* YAML. */
var convertSubCmdYAML2CSV = &cobra.Command{
	Use:   "yaml2csv",
	Short: "Convert yaml to csv format",
	Run:   cf.Run,
}
var convertSubCmdYAML2JSON = &cobra.Command{
	Use:   "yaml2json",
	Short: "Convert yaml to json format",
	Run:   cf.Run,
	Example: Examples(`# Convert yaml to json
ops-cli convert yaml2json -i input.yaml -o output.json`),
}
var convertSubCmdYAML2TOML = &cobra.Command{
	Use:   "yaml2toml",
	Short: "Convert yaml to toml format",
	Run:   cf.Run,
}
var convertSubCmdYAML2XML = &cobra.Command{
	Use:   "yaml2xml",
	Short: "Convert yaml to xml format",
	Run:   cf.Run,
}

var cf convertFlag

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.PersistentFlags().StringVarP(&cf.inFile, "in", "i", "", "Input file")
	convertCmd.PersistentFlags().StringVarP(&cf.outFile, "out", "o", "", "Output file")
	/* Markdown */
	convertCmd.AddCommand(convertSubCmdMarkdown2HTML, convertSubCmdMarkdown2PDF)
	/* CSV */
	convertCmd.AddCommand(convertSubCmdCSV2JSON, convertSubCmdCSV2TOML, convertSubCmdCSV2XML, convertSubCmdCSV2YAML)
	/* JSON */
	convertCmd.AddCommand(convertSubCmdJSON2CSV, convertSubCmdJSON2TOML, convertSubCmdJSON2XML, convertSubCmdJSON2YAML)
	/* TOML */
	convertCmd.AddCommand(convertSubCmdTOML2CSV, convertSubCmdTOML2JSON, convertSubCmdTOML2XML, convertSubCmdTOML2YAML)
	/* XML */
	convertCmd.AddCommand(convertSubCmdXML2CSV, convertSubCmdXML2JSON, convertSubCmdXML2TOML, convertSubCmdXML2YAML)
	/* YAML */
	convertCmd.AddCommand(convertSubCmdYAML2CSV, convertSubCmdYAML2JSON, convertSubCmdYAML2TOML, convertSubCmdYAML2XML)
}

type convertFlag struct {
	inFile  string
	inType  string
	outFile string
	outType string

	readFile []byte
}

func (c *convertFlag) Run(cmd *cobra.Command, _ []string) {
	if !ValidFile(c.inFile) || c.outFile == "" {
		_ = cmd.Help()
		return
	}
	switch cmd.Name() {
	case FileTypeMarkdown + "2" + FileTypeHTML:
		c.ConvertMarkdown2HTML()
		return
	case FileTypeMarkdown + "2" + FileTypePDF:
		c.ConvertMarkdown2PDF()
		return
	}
	slice := strings.Split(cmd.Name(), "2")
	if len(slice) != 2 {
		_ = cmd.Help()
		return
	}
	c.inType = slice[0]
	c.outType = slice[1]
	if err := c.Convert(); err != nil {
		log.Println(err)
	}
}

func (c *convertFlag) Convert() error {
	node, err := dasel.NewFromFile(c.inFile, c.inType)
	if err != nil {
		return err
	}
	if err := node.WriteToFile(c.outFile, c.outType, []storage.ReadWriteOption{
		{Key: storage.OptionIndent},
	}); err != nil {
		return err
	}
	return nil
}

func (c *convertFlag) Load() {
	var err error
	if c.readFile, err = os.ReadFile(c.inFile); err != nil {
		log.Println(err)
		os.Exit(0)
	}
}

func (c *convertFlag) ConvertMarkdown2HTML() {
	c.Load()
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Linkify,
			extension.Strikethrough,
			extension.Table,
			extension.TaskList,
			extension.Typographer,
			emoji.Emoji,
			&hashtag.Extender{},
			highlighting.Highlighting,
			&mermaid.Extender{},
			&toc.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithBlockParsers(),
			parser.WithInlineParsers(),
			parser.WithParagraphTransformers(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(c.readFile, &buf); err != nil {
		log.Println(err)
		os.Exit(0)
	}
	c.WriteFile(buf.Bytes())
}

func (c *convertFlag) ConvertMarkdown2PDF() {
	c.Load()
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Linkify,
			extension.Strikethrough,
			extension.Table,
			extension.TaskList,
			extension.Typographer,
			emoji.Emoji,
			&hashtag.Extender{},
			highlighting.Highlighting,
			&mermaid.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithBlockParsers(),
			parser.WithInlineParsers(),
			parser.WithParagraphTransformers(),
		),
		goldmark.WithRenderer(
			pdf.New(
				pdf.WithContext(rootContext),
				pdf.WithHeadingFont(pdf.GetTextFont("IBM Plex Serif", pdf.FontLora)),
				pdf.WithBodyFont(pdf.GetTextFont("Open Sans", pdf.FontRoboto)),
				pdf.WithCodeFont(pdf.GetCodeFont("Inconsolata", pdf.FontRobotoMono)),
			),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(c.readFile, &buf); err != nil {
		log.Println(err)
		os.Exit(0)
	}
	c.WriteFile(buf.Bytes())
}

func (c *convertFlag) WriteFile(content []byte) {
	if err := os.WriteFile(c.outFile, content, os.ModePerm); err != nil {
		log.Println(err)
	}
}
