package dotenv

import (
	"bytes"
	"fmt"
)

type parser struct {
	env      map[string]string
	input    []byte
	position int
	current  byte
	previous byte
}

func (p *parser) parse() error {
	p.next()
	p.consumeSpace(true)

	for i := p.current; i != 0x0; i = p.current {
		if p.current == '#' {
			p.consumeComment()
			p.consumeSpace(true)
			continue
		}

		key := p.readKey()
		p.consumeSpace(false)
		if p.current != '=' {
			return p.errorUnuexpectedSymbole("=")
		}
		p.next()
		p.consumeSpace(false)
		value, err := p.readValue()
		if err != nil {
			return err
		}

		p.env[string(key)] = string(value)
		p.consumeSpace(true)
	}

	return nil
}

func (p *parser) next() byte {
	p.previous = p.current
	if p.position >= len(p.input) {
		p.current = 0
	} else {
		p.current = p.input[p.position]
		p.position++
	}
	return p.current
}

func (p *parser) readKey() []byte {
	var key []byte

	var is_valid_key = func(i byte) bool {
		return (i >= 65 && i <= 90) || (i >= 97 && i <= 122) || (i >= 48 && i <= 57) || i == 95
	}

	for i := p.current; is_valid_key(i); i = p.next() {
		key = append(key, i)
	}

	return key
}

func (p *parser) readValue() ([]byte, error) {
	var value []byte

	var quoted_with byte
	if p.current == '"' || p.current == '\'' {
		quoted_with = p.current
		p.next()
	}

	for i := p.current; i != 0 && i != '\n'; {
		if i == '\\' {
			p.next()
		}
		if quoted_with == 0 && i == '#' && (p.previous == ' ' || p.previous == '=') {
			value = bytes.TrimRight(value, " ")
			break
		}
		if quoted_with == i { // if the value is quoted and this is the right quote
			break
		}
		if i == '$' {
			p.next()
			value = append(value, p.readVariable()...)
			i = p.current
			continue
		}
		value = append(value, p.current)
		i = p.next()
	}

	if quoted_with != 0 {
		if p.current == quoted_with {
			p.next()
		} else {
			return nil, p.errorUnterminatedQuotedValue(value)
		}
	} else {
		value = bytes.Trim(value, " ")
		value = bytes.Trim(value, "\t")
	}

	return value, nil
}

func (p *parser) consumeSpace(new_line_too bool) {
	c := p.current
	for {
		switch c {
		case ' ', '\t':
			c = p.next()
		case '\n':
			if new_line_too {
				c = p.next()
			} else {
				return
			}
		default:
			return
		}
	}
}

func (p *parser) consumeComment() {
	for p.current != '\n' && p.current != 0 {
		p.next()
	}
}

func (p *parser) readVariable() []byte {
	var var_name []byte

	var with_braces = p.current == '{'
	if with_braces {
		p.next()
	}

	for i := p.current; (i >= 65 && i <= 90) || (i >= 48 && i <= 57) || i == 95; i = p.next() {
		var_name = append(var_name, i)
	}

	if with_braces {
		if p.current == '}' {
			p.next()
		} else {
			return append([]byte("${"), var_name...)
		}
	}

	return []byte(p.env[string(var_name)])
}

func (p *parser) errorUnuexpectedSymbole(expected string) error {
	var line = 1
	for i := 0; i < p.position; i++ {
		if p.input[i] == '\n' {
			line++
		}
	}

	if p.current == 0 {
		return fmt.Errorf(`expected "%s", found end of file, line %d:%d`, expected, line, p.position)
	}
	return fmt.Errorf(`expected "%s", found "%c", line %d:%d`, expected, p.current, line, p.position)
}

func (p *parser) errorUnterminatedQuotedValue(v []byte) error {
	var line = 1
	for i := 0; i < p.position; i++ {
		if p.input[i] == '\n' {
			line++
		}
	}

	return fmt.Errorf(`unterminated quoted value "%s", line %d:%d`, v, line, p.position)
}

func (p *parser) reset() {
	p.input = nil
	p.position = 0
	p.previous = 0
	p.current = 0
}
