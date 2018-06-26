// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package api

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContent(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	name := randName(`page`)
	assert.NoError(t, postTx(`NewPage`, &url.Values{
		"ApplicationId": {`1`},
		"Name":          {name},
		"Value":         {`If(true){Div(){Span(My text)Address()}}.Else{Div(Body: Hidden text)}`},
		"Menu":          {`default_menu`},
		"Conditions":    {"true"},
	}))

	cases := []struct {
		url      string
		form     url.Values
		expected string
	}{
		{
			"content/source/" + name,
			url.Values{},
			`[{"tag":"if","attr":{"condition":"true"},"children":[{"tag":"div","children":[{"tag":"span","children":[{"tag":"text","text":"My text"}]},{"tag":"address"}]}],"tail":[{"tag":"else","children":[{"tag":"div","children":[{"tag":"text","text":"Hidden text"}]}]}]}]`,
		},
		{
			"content",
			url.Values{
				"template": {"input Div(myclass, #mytest# Div(mypar) the Div)"},
				"mytest":   {"test value"},
			},
			`[{"tag":"text","text":"input "},{"tag":"div","attr":{"class":"myclass"},"children":[{"tag":"text","text":"test value "},{"tag":"div","attr":{"class":"mypar"}},{"tag":"text","text":" the Div"}]}]`,
		},
		{
			"content",
			url.Values{
				"template":  {"#test_page# input Div(myclass, #test_page# ok) #test_page#"},
				"test_page": {"7"},
			},
			`[{"tag":"text","text":"7 input "},{"tag":"div","attr":{"class":"myclass"},"children":[{"tag":"text","text":"7 ok"}]},{"tag":"text","text":" 7"}]`,
		},
		{
			"content",
			url.Values{
				"template": {"SetVar(mytest, myvar)Div(myclass, Span(#mytest#) Div(mypar){Span(test)}#mytest#)"},
				"source":   {"true"},
			},
			`[{"tag":"setvar","attr":{"name":"mytest","value":"myvar"}},{"tag":"div","attr":{"class":"myclass"},"children":[{"tag":"span","children":[{"tag":"text","text":"#mytest#"}]},{"tag":"div","attr":{"class":"mypar"},"children":[{"tag":"span","children":[{"tag":"text","text":"test"}]}]},{"tag":"text","text":"#mytest#"}]}]`,
		},
		{
			"content",
			url.Values{
				"template": {`DBFind(Name: pages, Source: src).Custom(custom_col){
				Span(Body: "test")
			}`},
				"lang":   {"ru"},
				"source": {"true"},
			},
			`[{"tag":"dbfind","attr":{"name":"pages","source":"src"},"tail":[{"tag":"custom","attr":{"column":"custom_col"},"children":[{"tag":"span","children":[{"tag":"text","text":"test"}]}]}]}]`,
		},
		{
			"content",
			url.Values{
				"template": {`Data(Source: src).Custom(custom_col){
				Span(Body: "test")
			}`},
				"lang":   {"ru"},
				"source": {"true"},
			},
			`[{"tag":"data","attr":{"source":"src"},"tail":[{"tag":"custom","attr":{"column":"custom_col"},"children":[{"tag":"span","children":[{"tag":"text","text":"test"}]}]}]}]`,
		},
		{
			"content",
			url.Values{
				"template": {`Data(myforlist,"id,name",
					"1",Test message 1
					2,"Test message 2"
					3,"Test message 3"
					)`},
				"source": {"true"},
			},
			`[{"tag":"data","attr":{"columns":"id,name","data":"1,Test message 1\n\t\t\t\t\t2,\"Test message 2\"\n\t\t\t\t\t3,\"Test message 3\"","source":"myforlist"}}]`,
		},
		{
			"content",
			url.Values{
				"template": {`
					Data(src_test,"type"){
						text
					}
					ForList(src_test){
					If(#type#==text){
						Span(:#type#)
					}
				}`},
			},
			`[{"tag":"data","attr":{"columns":["type"],"data":[["text"]],"source":"src_test","types":["text"]}},{"tag":"forlist","attr":{"source":"src_test"},"children":[{"tag":"span","attr":{"":"text"}}]}]`,
		},
	}

	var ret contentResult
	for _, v := range cases {
		assert.NoError(t, sendPost(v.url, &v.form, &ret))
		assert.Equal(t, v.expected, string(ret.Tree))
	}
}
