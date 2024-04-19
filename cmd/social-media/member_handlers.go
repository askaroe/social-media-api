package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/askaroe/social-media-api/pkg/social-media/model"
	"github.com/askaroe/social-media-api/pkg/social-media/validator"
)

func (app *application) registerMemberHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	member := &model.Member{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = member.Password.Set(input.Password)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if model.ValidateMember(v, member); !v.Valid() {
		fmt.Println("ERROR HERE?")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Members.Insert(member)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrDuplicateEmail):
			v.AddError("email", "a member with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.models.Tokens.New(member.ID, 3*24*time.Hour, model.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var res struct {
		Token  *string       `json:"token"`
		Member *model.Member `json:"member"`
	}

	res.Token = &token.Plaintext
	res.Member = member

	app.writeJSON(w, http.StatusCreated, envelope{"member": res}, nil)
}
