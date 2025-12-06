package code

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func MAIN() {
	input := UserInput{
		Name:  "zhou",
		Email: "invalid@email.",
		Age:   150,
	}
	inputV := UserInputV{
		Name:  "joey",
		Email: "weirong.zhou@outlook.com",
		Age:   100,
	}

	err := Validate(input)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(input)
	}

	fmt.Println("====separarte====")

	errV := inputV.ValidateV()
	if errV != nil {
		fmt.Println(errV)
	} else {
		fmt.Println(inputV)
	}
}

type UserInput struct {
	Name  string
	Email string
	Age   int
}

type UserInputV struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"eq=100"`
}

func Validate(input UserInput) []string {
	var errors []string
	if input.Name == "" {
		errors = append(errors, "Name 不能为空")
	}

	if input.Email == "" {
		errors = append(errors, "Email不能为空")
	} else if !strings.Contains(input.Email, "@") || !strings.Contains(input.Email, ".") {
		errors = append(errors, "Email格式不正确")
	}

	if input.Age < 0 {
		errors = append(errors, "年龄不能小于0")
	} else if input.Age > 100 {
		errors = append(errors, "年龄不能大于100")
	}
	return errors
}

func (v UserInputV) ValidateV() error {
	validate := validator.New()
	return validate.Struct(v)
}
