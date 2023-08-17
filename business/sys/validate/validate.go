package validate

import (
	"reflect"
	"strings"

	emv "github.com/AfterShip/email-verifier"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
)

var validate *validator.Validate
var translator ut.Translator

func init() {
	validate = validator.New()

	translator, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	en_translations.RegisterDefaultTranslations(validate, translator)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}

func Check(val any) error {
	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(translator),
			}
			fields = append(fields, field)
		}

		return &fields
	}

	return nil
}

func Email(email string) error {
	_, err := emv.NewVerifier().Verify(email)
	if err != nil {
		return err
	}
	return nil
}

func GenerateID() string {
	return uuid.NewString()
}

func CheckID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}
	return nil
}
