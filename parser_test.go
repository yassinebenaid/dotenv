package dotenv

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	input := []byte(`abcd`)

	parser := parser{input: input}

	assert.Equal(t, "a", string(parser.next()))
	assert.Equal(t, 0, int(parser.previous))
	assert.Equal(t, "b", string(parser.next()))
	assert.Equal(t, "a", string(parser.previous))
	assert.Equal(t, "c", string(parser.next()))
	assert.Equal(t, "b", string(parser.previous))
	assert.Equal(t, "d", string(parser.next()))
	assert.Equal(t, "c", string(parser.previous))
	assert.Equal(t, uint8(0x0), parser.next())
	assert.Equal(t, "d", string(parser.previous))
	assert.Equal(t, 0, int(parser.next()))
	assert.Equal(t, 0, int(parser.previous))
}

func TestParse(t *testing.T) {
	testCases := []struct {
		input    string
		expected map[string]string
	}{
		{"", map[string]string{}},
		{`

		`, map[string]string{}},
		{"EMPTY_KEY=", map[string]string{"EMPTY_KEY": ""}},
		{"EMPTY_KEY_WITH_SPACE=  ", map[string]string{"EMPTY_KEY_WITH_SPACE": ""}},
		{"KEY=value", map[string]string{"KEY": "value"}},
		{"ARBITRARY_KEY= hello @ world 123 _+-;:, ", map[string]string{"ARBITRARY_KEY": "hello @ world 123 _+-;:,"}},
		{"INLINE_COMMENT= value #comment", map[string]string{"INLINE_COMMENT": "value"}},
		{"EMPTY_KEY=#comment", map[string]string{"EMPTY_KEY": ""}},
		{"#comment", map[string]string{}},
		{`QUOTED="value"`, map[string]string{"QUOTED": "value"}},
		{`QUOTED_COMMENT="#quoted-comment"`, map[string]string{"QUOTED_COMMENT": "#quoted-comment"}},
		{`KEY=txt" #test"`, map[string]string{"KEY": `txt"`}},
		{`WITH_HASH=txt"#test"`, map[string]string{"WITH_HASH": `txt"#test"`}},
		{`ONLY_SPACE=" "`, map[string]string{"ONLY_SPACE": " "}},
		{`WITH_PADDING=" value "`, map[string]string{"WITH_PADDING": " value "}},
		{`WITH_PADDING=" value "`, map[string]string{"WITH_PADDING": " value "}},
		{`ESCAPED_QUOTE=" value\" "`, map[string]string{"ESCAPED_QUOTE": ` value" `}},
		{`ESCAPED_ESCAPE_CHAR=" value\\s "`, map[string]string{"ESCAPED_ESCAPE_CHAR": ` value\s `}},
		{`
			KEY_1=value-1
			KEY_2=value-2
		`, map[string]string{"KEY_1": "value-1", "KEY_2": "value-2"}},
		{`
			KEY_1 = value 1
			KEY_2 = value 2
		`, map[string]string{"KEY_1": "value 1", "KEY_2": "value 2"}},
		{"key1=value-1\nkey2=value-2", map[string]string{"key1": "value-1", "key2": "value-2"}},
		{`
			# some comments go here
			key1 = value-1 # comment
				# comment
			key2=8dD63dBQ3Gf+MQO2ZScyLU6culzas5PeoaYj3Q6DddU=
			# commet
			key3=some-#strange-value
		`, map[string]string{
			"key1": "value-1",
			"key2": "8dD63dBQ3Gf+MQO2ZScyLU6culzas5PeoaYj3Q6DddU=",
			"key3": "some-#strange-value"}},
		{`
			KEY_1=value
			KEY_2=$KEY_1
		`, map[string]string{"KEY_1": "value", "KEY_2": "value"}},
		{`
			KEY_1=value
			KEY_2="$KEY_1-world"
		`, map[string]string{"KEY_1": "value", "KEY_2": "value-world"}},
		{`
			KEY_1=value
			KEY_2="hello-$KEY_1-world"
		`, map[string]string{"KEY_1": "value", "KEY_2": "hello-value-world"}},
		{`KEY="$UNDEFINED"`, map[string]string{"KEY": ""}},
		{`KEY=hello $UNDEFINED world`, map[string]string{"KEY": "hello  world"}},
		{`KEY=hello$UNDEFINEDworld`, map[string]string{"KEY": "helloworld"}},
		{`KEY=\$ESCAPED`, map[string]string{"KEY": "$ESCAPED"}},
		{`KEY=	$UNDEFINED s`, map[string]string{"KEY": "s"}},
		{`KEY=inline-\$ESCAPED-value`, map[string]string{"KEY": "inline-$ESCAPED-value"}},
		{`
			KEY_1=value
			KEY_2=inline-\$KEY_1-value
		`, map[string]string{"KEY_1": "value", "KEY_2": "inline-$KEY_1-value"}},
		{`
			# all at once
			KEY_1=value # inline comment

			# empty lines
			# multi line
			# comment
			KEY_2=inline-$KEY_1-value
			KEY_3=$KEY_2 # comment
			KEY_4	=	$UNDEFINED_KEY	 # comment with tabs
			KEY_5 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_world	 " # comment with tabs
		`, map[string]string{
			"KEY_1": "value",
			"KEY_2": "inline-value-value",
			"KEY_3": "inline-value-value",
			"KEY_4": "",
			"KEY_5": "value_inline-value-value_inline-value-value__world	 ",
		}},

		//
		{`
			KEY_1=value
			KEY_2=${KEY_1}
		`, map[string]string{"KEY_1": "value", "KEY_2": "value"}},
		{`
			KEY_1=value
			KEY_2="${KEY_1}-world"
		`, map[string]string{"KEY_1": "value", "KEY_2": "value-world"}},
		{`
			KEY_1=value
			KEY_2="hello-${KEY_1}-world"
		`, map[string]string{"KEY_1": "value", "KEY_2": "hello-value-world"}},
		{`KEY="${UNDEFINED}"`, map[string]string{"KEY": ""}},
		{`KEY=hello ${UNDEFINED} world`, map[string]string{"KEY": "hello  world"}},
		{`KEY=hello${UNDEFINED}world`, map[string]string{"KEY": "helloworld"}},
		{`KEY=\${ESCAPED}`, map[string]string{"KEY": "${ESCAPED}"}},
		{`KEY=	${UNDEFINED} s`, map[string]string{"KEY": "s"}},
		{`KEY=inline-\${ESCAPED}-value`, map[string]string{"KEY": "inline-${ESCAPED}-value"}},
		{`
			KEY_1=value
			KEY_2=inline-\${KEY_1}-value
		`, map[string]string{"KEY_1": "value", "KEY_2": "inline-${KEY_1}-value"}},
		{`
			# all at once
			KEY_1=value # inline comment

			# empty lines
			# multi line
			# comment
			KEY_2=inline-${KEY_1}-value
			KEY_3=${KEY_2} # comment
			KEY_4	=	${UNDEFINED_KEY}	 # comment with tabs
			KEY_5 = "${KEY_1}\_${KEY_2}\_${KEY_3}\_${KEY_4}\_world	 " # comment with tabs
		`, map[string]string{
			"KEY_1": "value",
			"KEY_2": "inline-value-value",
			"KEY_3": "inline-value-value",
			"KEY_4": "",
			"KEY_5": "value_inline-value-value_inline-value-value__world	 ",
		}},
		{`KEY=${SOME_KEY`, map[string]string{"KEY": "${SOME_KEY"}},
		{`
			KEY_1=value
			KEY_2=${KEY_1
		`, map[string]string{"KEY_1": "value", "KEY_2": "${KEY_1"}},
		{`
			KEY_1=value
			KEY_2=${KEY_1 test
		`, map[string]string{"KEY_1": "value", "KEY_2": "${KEY_1 test"}},
		{`
			KEY_1=value
			KEY_2=${KEY_1test
		`, map[string]string{"KEY_1": "value", "KEY_2": "${KEY_1test"}},
		{`
			KEY_1=value
			KEY_2="${KEY_1test"
		`, map[string]string{"KEY_1": "value", "KEY_2": "${KEY_1test"}},
		{`
			KEY_1=value
			KEY_2="${KEY_1}test"
		`, map[string]string{"KEY_1": "value", "KEY_2": "valuetest"}},
	}

	for _, tc := range testCases {
		p := parser{input: []byte(tc.input), env: make(map[string]string)}

		err := p.parse()
		assert.NoError(t, err)

		assert.NotNil(t, p.env)

		assert.Equal(t, tc.expected, p.env)
	}
}

func BenchmarkParser(b *testing.B) {
	input := `
	# all at once	
	KEY_1=value # inline comment	

	# empty lines	
	# multi line	
	# comment	
	KEY_2=inline-$KEY_1-value	
	KEY_3=$KEY_2 # comment	
	KEY_4	=	$UNDEFINED_KEY	 # comment with tabs
	KEY_5 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_world	 " # comment with tabs
	KEY_6 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_$KEY_5\_world	 " # comment with tabs
	KEY_7 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_$KEY_5\_$KEY_6\_world	 " # comment with tabs
	KEY_8 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_$KEY_5\_$KEY_6\_$KEY_7\_world	 " # comment with tabs
	KEY_9 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_$KEY_5\_$KEY_6\_$KEY_7\_$KEY_8\_world	 " # comment with tabs
	KEY_10 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_$KEY_5\_$KEY_6\_$KEY_7\_$KEY_8\_$KEY_9\_world	 " # comment with tabs
	`

	for i := 0; i < b.N; i++ {
		p := parser{input: []byte(input), env: make(map[string]string)}
		if err := p.parse(); err != nil {
			b.Fail()
		}
	}
}

func BenchmarkGodotenvPackage(b *testing.B) {
	input := `
	# all at once	
	KEY_1=value # inline comment	

	# empty lines	
	# multi line	
	# comment	
	KEY_2=inline-$KEY_1-value	
	KEY_3=$KEY_2 # comment	
	KEY_4	=	$UNDEFINED_KEY	 # comment with tabs
	KEY_5 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_world	 " # comment with tabs
	KEY_6 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_$KEY_5\_world	 " # comment with tabs
	KEY_7 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_$KEY_5\_$KEY_6\_world	 " # comment with tabs
	KEY_8 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_$KEY_5\_$KEY_6\_$KEY_7\_world	 " # comment with tabs
	KEY_9 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_$KEY_5\_$KEY_6\_$KEY_7\_$KEY_8\_world	 " # comment with tabs
	KEY_10 = "$KEY_1\_$KEY_2\_$KEY_3\_$KEY_4\_$KEY_5\_$KEY_6\_$KEY_7\_$KEY_8\_$KEY_9\_world	 " # comment with tabs
	`

	for i := 0; i < b.N; i++ {
		_, err := godotenv.Unmarshal(input)
		if err != nil {
			b.Fail()
		}
	}
}
