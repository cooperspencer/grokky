//
// Copyright (c) 2016 Konstanin Ivanov <kostyarin.ivanov@gmail.com>.
// All rights reserved. This program is free software. It comes without
// any warranty, to the extent permitted by applicable law. You can
// redistribute it and/or modify it under the terms of the Do What
// The Fuck You Want To Public License, Version 2, as published by
// Sam Hocevar. See LICENSE.md file for more details or see below.
//

//
//        DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
//                    Version 2, December 2004
//
// Copyright (C) 2004 Sam Hocevar <sam@hocevar.net>
//
// Everyone is permitted to copy and distribute verbatim or modified
// copies of this license document, and changing it is allowed as long
// as the name is changed.
//
//            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
//   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION
//
//  0. You just DO WHAT THE FUCK YOU WANT TO.
//

//
package grokky

//go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html

import (
	"testing"
)

const patternsTest = "patterns.test"

func TestNew(t *testing.T) {
	h := New()
	if len(h) != 0 {
		t.Error("New returns non-empty host")
	}
	if h == nil {
		t.Error("New returns nil")
	}
}

func TestHost_Add(t *testing.T) {
	h := New()
	if err := h.Add("", "expr"); err == nil {
		t.Error("(Host).Add is missing ErrEmptyName")
	} else if err != ErrEmptyName {
		t.Error("(Host).Add returns non-ErrEmptyName error")
	}
	if len(h) != 0 {
		t.Error("added bad patterns")
	}
	if err := h.Add("zorro", ""); err == nil {
		t.Error("(Host).Add is missing ErrEmptyExpression")
	} else if err != ErrEmptyExpression {
		t.Error("(Host).Add returns non-ErrEmptyExpression error")
	}
	if len(h) != 0 {
		t.Error("added bad patterns")
	}
	if err := h.Add("DIGIT", `\d`); err != nil {
		t.Errorf("(Host).Add returns non-nil error: %v", err)
	}
	if len(h) != 1 {
		t.Error("wrong patterns count")
	}
	if err := h.Add("DIGIT", `[+-](0x)?\d`); err == nil {
		t.Error("(Host).Add is missing ErrAlreadyExist")
	} else if err != ErrAlreadyExist {
		t.Error("(Host).Add returns non-ErrAlreadyExist error")
	}
	if len(h) != 1 {
		t.Error("wrong patterns count")
	}
	if err := h.Add("BAD", `(?![0-5])`); err == nil {
		t.Error("(Host).Add is missing any bad-regexp error")
	}
	if len(h) != 1 {
		t.Error("wrong patterns count")
	}
	if err := h.Add("TWODIG", `%{DIGIT}-%{DIGIT}`); err != nil {
		t.Errorf("(Host).Add returns non-nil error: %v", err)
	}
	if len(h) != 2 {
		t.Error("wrong patterns count")
	}
	if err := h.Add("THREE", `%{NOT}-%{EXIST}`); err == nil {
		t.Errorf("(Host).Add is missing the-pattern-not-exist error")
	}
	if len(h) != 2 {
		t.Error("wrong patterns count")
	}
	if err := h.Add("FOUR", `%{DIGIT:one}-%{DIGIT:two}`); err != nil {
		t.Errorf("(Host).Add returns non-nil error: %v", err)
	}
	if len(h) != 3 {
		t.Error("wrong patterns count")
	}
	if err := h.Add("FIVE", `(?!\d)%{DIGIT}(?!\d)`); err == nil {
		t.Errorf("(Host).Add is missing an error of regexp")
	}
	if len(h) != 3 {
		t.Error("wrong patterns count")
	}
	if err := h.Add("SIX", `%{FOUR:four}-%{DIGIT:six}`); err != nil {
		t.Errorf("(Host).Add returns non-nil error")
	}
	if len(h) != 4 {
		t.Error("wrong patterns count")
	}
}

func TestHost_Get(t *testing.T) {
	h := New()
	if err := h.Add("DIG", `\d`); err != nil {
		t.Error(err)
	}
	if p, err := h.Get("DIG"); err != nil {
		t.Error(err)
	} else if p == nil {
		t.Error("(Host).Get returns nil (and nil-error)")
	}
	if p, err := h.Get("SEVEN"); err == nil {
		t.Error("(Host).Get is missing ErrNotExist")
	} else if p != nil {
		t.Error("(Host).Get returns non-nil not-exsted-pattern")
	}
}

func TestHost_AddFromFile(t *testing.T) {
	h := New()
	if err := h.AddFromFile(patternsTest); err != nil {
		t.Error(err)
	}
	if len(h) != 3 {
		t.Error("wrong patterns count")
	}
	if _, err := h.Get("ONE"); err != nil {
		t.Error(err)
	}
	if _, err := h.Get("TWO"); err != nil {
		t.Error(err)
	}
	if _, err := h.Get("THREE"); err != nil {
		t.Error(err)
	}
}

func TestHost_inject(t *testing.T) {
	h := New()
	h["TWO"] = `(?!\d)`
	if err := h.Add("ONE", `%{TWO:one}`); err == nil {
		t.Error("bad injection returns nil error")
	}
}

func TestHost_badPath(t *testing.T) {
	h := New()
	if err := h.AddFromFile("unexisted-file-without-patterns"); err == nil {
		t.Error("bad path with nil error")
	}
}

func TestHost_addFromLine(t *testing.T) {
	h := New()
	if err := h.addFromLine("ONE (?!\\d)"); err == nil {
		t.Error("bad line with nil error")
	}
}
