package dotenv

import (
	"fmt"
	"strings"
)

type errUnexpectedSymbole struct {
	ch byte
}

func (e errUnexpectedSymbole) Error() string {
	return fmt.Sprintf(`unexpected symbole "%c"`, e.ch)
}

func parse(s []byte) (map[string]string, error) {
	env := make(map[string]string)

	p := parser{input: s}

	p.next()
	p.consumeSpace(true)

	for i := p.current; i != 0x0; i = p.current {
		if p.current == '#' {
			p.consumeComment()
			continue
		}

		key := p.readKey()
		p.consumeSpace(false)
		if p.current != '=' {
			return nil, errUnexpectedSymbole{p.current}
		}
		p.next()
		p.consumeSpace(false)
		value, err := p.readValue()
		if err != nil {
			return nil, err
		}

		env[string(key)] = strings.TrimSpace(string(value))
		p.consumeSpace(true)
	}

	return env, nil
}

type parser struct {
	input    []byte
	position int
	current  byte
}

func (p *parser) next() byte {
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

	var quoted bool
	if p.current == '"' {
		quoted = true
		p.next()
	}

	for i := p.current; i != 0 && i != '\n'; i = p.next() {
		if !quoted && i == '#' {
			break
		}
		if quoted && i == '"' {
			break
		}
		value = append(value, p.current)
	}

	if quoted {
		if p.current == '"' {
			p.next()
		} else {
			return nil, fmt.Errorf(`unterminated quoted value "%s"`, value)
		}
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
