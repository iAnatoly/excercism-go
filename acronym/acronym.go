// This is a "stub" file.  It's a little start on your solution.
// It's not a complete solution though; you have to write some code.

// Package acronym should have a package comment that summarizes what it's about.
// https://golang.org/doc/effective_go.html#commentary
package acronym

import (
    str "strings"
    "unicode"
)

// Abbreviate should have a comment documenting it.
func Abbreviate(s string) string {
	// Write some code here to pass the test suite.
	// Then remove all the stock comments.
	// They're here to help you get started but they only clutter a finished solution.
	// If you leave them in, reviewers may protest!
    var acronym str.Builder
    repl := str.NewReplacer(
        ","," ",
        "-", " ",
        "_"," ")
    for _, word := range str.Fields(repl.Replace(s)) {
        for  _, c := range word {
            if unicode.IsPunct(c) {
                continue
            }
            acronym.WriteRune(c)
            break
        }
    }
	return str.ToUpper(acronym.String())
}