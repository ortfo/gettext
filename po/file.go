// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package po

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// File represents an PO File.
//
// See http://www.gnu.org/software/gettext/manual/html_node/PO-Files.html
type File struct {
	MimeHeader     Header
	Messages       []Message
	SourceLanguage language.Tag
}

// Load loads po file format data.
func Load(data []byte) (*File, error) {
	return loadData(data)
}

// LoadFile loads a named po file.
func LoadFile(path string) (*File, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return loadData(data)
}

func loadData(data []byte) (*File, error) {
	r := newLineReader(string(data))
	var file File
	file.SetSourceLanguage(language.English)
	for {
		var msg Message
		if err := msg.readPoEntry(r); err != nil {
			if err == io.EOF {
				return &file, nil
			}
			return nil, err
		}
		if msg.MsgId == "" {
			file.MimeHeader.parseHeader(&msg)
			continue
		}
		file.Messages = append(file.Messages, msg)
	}
}

// Save saves a po file.
func (f *File) Save(name string) error {
	return ioutil.WriteFile(name, []byte(f.String()), 0666)
}

func (f *File) SetSourceLanguage(language language.Tag) {
	f.SourceLanguage = language
}

// Save returns a po file format data.
func (f *File) Data() []byte {
	// sort the massge as ReferenceFile/ReferenceLine field
	var messages []Message
	messages = append(messages, f.Messages...)
	collator := collate.New(f.SourceLanguage, collate.Loose)
	sort.SliceStable(messages, func(i, j int) bool {
		msgid := collator.CompareString(messages[i].MsgId, messages[j].MsgId)
		msgctx := collator.CompareString(messages[i].MsgContext, messages[j].MsgContext)
		return msgid < 0 || (msgid == 0 && msgctx < 0)
	})

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n", f.MimeHeader.String())
	for i := 0; i < len(messages); i++ {
		fmt.Fprintf(&buf, "%s\n", messages[i].String())
	}
	return buf.Bytes()
}

// String returns the po format file string.
func (f *File) String() string {
	return string(f.Data())
}
