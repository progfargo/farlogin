package util

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

func IsValidEmail(str string) error {
	re := regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@" +
		"(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
	res := re.MatchString(str)
	if !res {
		return errors.New("Invalid e-mail address.")
	}

	return nil
}

func IsValidPassword(str string) error {
	if len(str) < 6 {
		return errors.New("Password is too short. It must contain at least 6 characters.")
	}

	re := regexp.MustCompile("[a-z0-9!#$%&]+")
	res := re.MatchString(str)
	if !res {
		return errors.New("Invalid character in password. Allowed characters:|| 'a-z0-9!#$%&'")
	}

	return nil
}

func IsValidIdentifier(str string) bool {
	rv, err := regexp.MatchString("^[a-zA-Z]{1}[a-zA-Z0-9]*", str)
	if err != nil {
		panic(err)
	}

	return rv
}

func IsValidHexColor(str string) bool {
	rv, err := regexp.MatchString("^#[0-9a-fA-F]{6}$", str)
	if err != nil {
		panic(err)
	}

	return rv
}

var (
	htmlCommentRegex = regexp.MustCompile("(?i)<!--([\\s\\S]*?)-->")
	svgRegex         = regexp.MustCompile(`(?i)^\s*(?:<\?xml[^>]*>\s*)?(?:<!doctype svg[^>]*>\s*)?<svg[^>]*>[^*]*<\/svg>\s*$`)
)

// isBinary checks if the given buffer is a binary file.
func isBinary(buf []byte) bool {
	if len(buf) < 24 {
		return false
	}
	for i := 0; i < 24; i++ {
		charCode, _ := utf8.DecodeRuneInString(string(buf[i]))
		if charCode == 65533 || charCode <= 8 {
			return true
		}
	}
	return false
}

// returns true if the given buffer is a valid SVG image.
func IsValidSvg(buf []byte) bool {
	return !isBinary(buf) && svgRegex.Match(htmlCommentRegex.ReplaceAll(buf, []byte{}))
}
