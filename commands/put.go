package commands

import (
	"strconv"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// TODO Inline syntax
const putUsage string = `usage: put <newline>
Create or update parameters. Enter one option per line, ending with a blank line.
Example below. All fields except name, value,a nd type are optionalll fields except name
and value are optional. Default type is String.
/>put
Input options. End with a blank line.
... name=/foo/bar
... value=baz
... type=string
... description=foobar
... key=arn:aws:kms:us-west-2:012345678901:key/321ec4ec-ed00-427f-9729-748ba2254794
... overwrite=true
... pattern=[A-z]+
...
`

var putParamInput ssm.PutParameterInput

var validTypes = []string{"String", "StringList", "SecureString"}

// Add or update parameters
func put(c *ishell.Context) {
	var err error
	defer reset()
	// Set the prompt explicitly rather than use SetMultiPrompt due to the unexpected 2nd line behavior
	shell.SetPrompt("... ")
	c.Println("Input options. End with a blank line.")
	c.ReadMultiLinesFunc(putOptions)
	if putParamInput.Type == nil {
		putParamInput.Type = aws.String("String")
	}
	if putParamInput.Name == nil || putParamInput.Value == nil {
		shell.Println("Error: name and value are required fields.")
	} else {
		err = ps.Put(&putParamInput)
		if err != nil {
			shell.Println("Error: ", err)
		}
	}
}

func reset() {
	putParamInput = ssm.PutParameterInput{}
	setPrompt(ps.Cwd)
}

func putOptions(s string) bool {
	//var err error
	if s == "" {
		return false
	}
	paramOption := strings.Split(s, "=")
	if len(paramOption) != 2 {
		shell.Println("Invalid input.")
		shell.Println(putUsage)
		return false
	}
	field := strings.ToLower(paramOption[0])
	val := paramOption[1]
	if validate(field, val) {
		return true
	}
	return false
}

func validate(f string, v string) bool {
	m := map[string]func(string) bool{
		"type":        validateType,
		"name":        validateName,
		"value":       validateValue,
		"description": validateDescription,
		"key":         validateKey,
		"pattern":     validatePattern,
		"overwrite":   validateOverwrite,
	}
	if validator, ok := m[strings.ToLower(f)]; ok {
		if validator(v) {
			return true
		}
		return false
	}
	shell.Println("Invalid input " + f + "=" + v)
	shell.Println(putUsage)
	return false
}

func validateType(s string) bool {
	for i := 0; i < len(validTypes); i++ {
		if strings.EqualFold(s, validTypes[i]) {
			putParamInput.Type = aws.String(validTypes[i])
			return true
		}
	}
	shell.Println("Invalid type " + s)
	return false
}

func validateValue(s string) bool {
	putParamInput.Value = aws.String(s)
	return true
}

func validateName(s string) bool {
	putParamInput.Name = aws.String(s)
	return true
}

func validateDescription(s string) bool {
	putParamInput.Description = aws.String(s)
	return true
}

func validateKey(s string) bool {
	putParamInput.KeyId = aws.String(s)
	return true
}

func validatePattern(s string) bool {
	putParamInput.AllowedPattern = aws.String(s)
	return true
}

func validateOverwrite(s string) bool {
	overwrite, err := strconv.ParseBool(s)
	if err != nil {
		shell.Println("overwrite must be true or false")
		return false
	}
	putParamInput.Overwrite = aws.Bool(overwrite)
	return true
}
