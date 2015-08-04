package jsont

import (
	"strings"
	"testing"
)

func TestTokenizer_UsageExample(t *testing.T) {

	jsonInput := `{ "Hello" : "World" }`
	tokenizer := NewTokenizer(strings.NewReader(jsonInput))

	jsonResult := ""
	for tokenizer.Next() {
		token := tokenizer.Token()
		jsonResult += token.Value
	}

	if jsonResult != jsonInput {
		t.Error("Error. All tokens concatenated should be equal to json input.", jsonInput, jsonResult)
	}

	if tokenizer.Error() != nil {
		t.Errorf("An unexpected error occured: '%s'", tokenizer.Error())
	}
}

func TestTokenizer_JsonDocument_WithWhitespace(t *testing.T) {

	jsonInput := `{
    "key1" : "value",
    "key2" : 123,
    "key3" : 123.99,
    "key4" : true,
    "key5" : false,
    "key6" : null,
    "key7" : {},
    "key8" : [ "value", 123, 123.99, true, false, null, {}, [] ]
    }`

	expected := []Token{
		Token{ObjectStart, "{"}, Token{Whitespace, "\n    "},
		Token{FieldName, `"key1"`}, Token{Whitespace, " "}, Token{Colon, ":"}, Token{Whitespace, " "}, Token{String, `"value"`}, Token{Comma, ","}, Token{Whitespace, "\n    "},
		Token{FieldName, `"key2"`}, Token{Whitespace, " "}, Token{Colon, ":"}, Token{Whitespace, " "}, Token{Integer, "123"}, Token{Comma, ","}, Token{Whitespace, "\n    "},
		Token{FieldName, `"key3"`}, Token{Whitespace, " "}, Token{Colon, ":"}, Token{Whitespace, " "}, Token{Float, "123.99"}, Token{Comma, ","}, Token{Whitespace, "\n    "},
		Token{FieldName, `"key4"`}, Token{Whitespace, " "}, Token{Colon, ":"}, Token{Whitespace, " "}, Token{True, "true"}, Token{Comma, ","}, Token{Whitespace, "\n    "},
		Token{FieldName, `"key5"`}, Token{Whitespace, " "}, Token{Colon, ":"}, Token{Whitespace, " "}, Token{False, "false"}, Token{Comma, ","}, Token{Whitespace, "\n    "},
		Token{FieldName, `"key6"`}, Token{Whitespace, " "}, Token{Colon, ":"}, Token{Whitespace, " "}, Token{Null, "null"}, Token{Comma, ","}, Token{Whitespace, "\n    "},
		Token{FieldName, `"key7"`}, Token{Whitespace, " "}, Token{Colon, ":"}, Token{Whitespace, " "}, Token{ObjectStart, "{"}, Token{ObjectEnd, "}"}, Token{Comma, ","}, Token{Whitespace, "\n    "},
		Token{FieldName, `"key8"`}, Token{Whitespace, " "}, Token{Colon, ":"}, Token{Whitespace, " "},
		Token{ArrayStart, "["}, Token{Whitespace, " "},
		Token{String, `"value"`}, Token{Comma, ","}, Token{Whitespace, " "},
		Token{Integer, "123"}, Token{Comma, ","}, Token{Whitespace, " "},
		Token{Float, "123.99"}, Token{Comma, ","}, Token{Whitespace, " "},
		Token{True, "true"}, Token{Comma, ","}, Token{Whitespace, " "},
		Token{False, "false"}, Token{Comma, ","}, Token{Whitespace, " "},
		Token{Null, "null"}, Token{Comma, ","}, Token{Whitespace, " "},
		Token{ObjectStart, "{"}, Token{ObjectEnd, "}"}, Token{Comma, ","}, Token{Whitespace, " "},
		Token{ArrayStart, "["}, Token{ArrayEnd, "]"}, Token{Whitespace, " "},
		Token{ArrayEnd, "]"}, Token{Whitespace, "\n    "},
		Token{ObjectEnd, "}"}, Token{Empty, ""}, Token{Empty, ""},
	}

	tokenizer := NewTokenizer(strings.NewReader(jsonInput))
	for _, expectedToken := range expected {
		ValidateNextToken(t, tokenizer, expectedToken)
	}
}

func TestTokenizer_JsonDocument_WithoutWhitespace(t *testing.T) {

	jsonInput := `{"key1":"value","key2":123,"key3":123.99,"key4":true,"key5":false,"key6":null,"key7":{},"key8":["value",123,123.99,true,false,null,{},[]]}`

	expected := []Token{
		Token{ObjectStart, "{"},
		Token{FieldName, `"key1"`}, Token{Colon, ":"}, Token{String, `"value"`}, Token{Comma, ","},
		Token{FieldName, `"key2"`}, Token{Colon, ":"}, Token{Integer, "123"}, Token{Comma, ","},
		Token{FieldName, `"key3"`}, Token{Colon, ":"}, Token{Float, "123.99"}, Token{Comma, ","},
		Token{FieldName, `"key4"`}, Token{Colon, ":"}, Token{True, "true"}, Token{Comma, ","},
		Token{FieldName, `"key5"`}, Token{Colon, ":"}, Token{False, "false"}, Token{Comma, ","},
		Token{FieldName, `"key6"`}, Token{Colon, ":"}, Token{Null, "null"}, Token{Comma, ","},
		Token{FieldName, `"key7"`}, Token{Colon, ":"}, Token{ObjectStart, "{"}, Token{ObjectEnd, "}"}, Token{Comma, ","},
		Token{FieldName, `"key8"`}, Token{Colon, ":"},
		Token{ArrayStart, "["},
		Token{String, `"value"`}, Token{Comma, ","},
		Token{Integer, "123"}, Token{Comma, ","},
		Token{Float, "123.99"}, Token{Comma, ","},
		Token{True, "true"}, Token{Comma, ","},
		Token{False, "false"}, Token{Comma, ","},
		Token{Null, "null"}, Token{Comma, ","},
		Token{ObjectStart, "{"}, Token{ObjectEnd, "}"}, Token{Comma, ","},
		Token{ArrayStart, "["}, Token{ArrayEnd, "]"},
		Token{ArrayEnd, "]"},
		Token{ObjectEnd, "}"}, Token{Empty, ""}, Token{Empty, ""},
	}

	tokenizer := NewTokenizer(strings.NewReader(jsonInput))
	for _, expectedToken := range expected {
		ValidateNextToken(t, tokenizer, expectedToken)
	}
}

func ValidateNextToken(t *testing.T, tokenizer Tokenizer, expectedToken Token) {

	if expectedToken.Type == Empty {
		if tokenizer.Next() {
			t.Error("tokenizer.Next() == true")
		}
	} else {
		if !tokenizer.Next() {
			t.Error("tokenizer.Next() != true")
		}
	}

	if tokenizer.Token() != expectedToken {
		t.Error("tokenizer.Token() != expectedToken", tokenizer.Token(), expectedToken)
	}

	if tokenizer.Error() != nil {
		t.Error("tokenizer.Error() != nil")
	}
}
