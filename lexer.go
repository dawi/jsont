package jsont

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
)

// return value can be one of those:
//  * valid token
//  * partly filled token + error
type readTokenFunction func(reader *bufio.Reader) (Token, error)

// array of functions that are tried in order to read the next token
var readTokenFunctions = []readTokenFunction{readCharToken, readStringToken, readWhitespaceToken, readNonStringToken}

// reads the next json token from reader
func readToken(reader *bufio.Reader) (Token, error) {
	for _, readTokenFunction := range readTokenFunctions {
		if token, err := readTokenFunction(reader); err != nil || token.Value != "" {
			return token, err
		}
	}
	return Token{}, errors.New("Error while parsing tokens. This error should never happen.")
}

// reads a single character token like { } [ ] : ,
func readCharToken(reader *bufio.Reader) (Token, error) {

	b, err := reader.ReadByte()
	if err != nil {
		return Token{}, err
	}

	switch b {
	case '{':
		return Token{ObjectStart, string(b)}, nil
	case '}':
		return Token{ObjectEnd, string(b)}, nil
	case '[':
		return Token{ArrayStart, string(b)}, nil
	case ']':
		return Token{ArrayEnd, string(b)}, nil
	case ':':
		return Token{Colon, string(b)}, nil
	case ',':
		return Token{Comma, string(b)}, nil
	}

	reader.UnreadByte()
	return Token{}, nil
}

// reads a string token like "abc"
func readStringToken(reader *bufio.Reader) (Token, error) {

	var escape = false

	var buffer bytes.Buffer

	if b, err := reader.ReadByte(); err != nil {
		return Token{}, err
	} else if b != '"' {
		reader.UnreadByte()
		return Token{}, nil
	} else {
		buffer.WriteByte(b)
	}

	var err error

	for {

		b, e := reader.ReadByte()
		if e != nil {
			err = e
			break
		}

		if escape {
			escape = false
			buffer.WriteByte(b)
		} else if b == '\\' {
			escape = true
			buffer.WriteByte(b)
		} else if b == '"' {
			escape = false
			buffer.WriteByte(b)
			break
		} else {
			escape = false
			buffer.WriteByte(b)
		}
	}

	return Token{String, buffer.String()}, err
}

// reads a whitespace token like " \r \n \t "
func readWhitespaceToken(reader *bufio.Reader) (Token, error) {

	var err error
	var buffer bytes.Buffer

	for {
		if b, e := reader.ReadByte(); e != nil {
			err = e
			break
		} else if isWhitespace(b) {
			buffer.WriteByte(b)
		} else {
			reader.UnreadByte()
			break
		}
	}

	if buffer.Len() == 0 {
		return Token{}, err
	}

	return Token{Whitespace, buffer.String()}, err
}

// reads a non string token like true, false, 123.45
func readNonStringToken(reader *bufio.Reader) (Token, error) {

	var err error
	var buffer bytes.Buffer

	for {
		if b, e := reader.ReadByte(); e != nil {
			err = e
			break
		} else if isWhitespace(b) || isJsonChar(b) {
			reader.UnreadByte()
			break
		} else {
			buffer.WriteByte(b)
		}
	}

	if buffer.Len() == 0 {
		return Token{}, err
	}

	token := buffer.String()

	switch token {
	case "true":
		return Token{True, token}, nil
	case "false":
		return Token{False, token}, nil
	case "null":
		return Token{Null, token}, nil
	}

	switch {
	case isInteger(token):
		return Token{Integer, token}, nil
	case isFloat(token):
		return Token{Float, token}, nil
	}

	return Token{Unknown, token}, err
}

// true if byte represents a whitespace
func isWhitespace(b byte) bool {
	return b == ' ' || b == '\r' || b == '\n' || b == '\t'
}

// true if byte represents a json character
func isJsonChar(b byte) bool {
	return b == ':' || b == ',' || b == '{' || b == '}' || b == '[' || b == ']'
}

func isInteger(token string) bool {
	_, err := strconv.Atoi(token)
	return err == nil
}

func isFloat(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}
