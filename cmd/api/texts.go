package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"dev.theenthusiast.text-bin/internal/data"
	"dev.theenthusiast.text-bin/internal/validator"
)

// createTextHandler will be used to create a text
func (app *application) createTextHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title        string `json:"title"`
		Content      string `json:"content"`
		Format       string `json:"format"`
		ExpiresValue int    `json:"expiresValue"`
		ExpiresUnit  string `json:"expiresUnit"`
		IsPrivate    bool   `json:"is_private"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var expires time.Time
	if input.ExpiresUnit != "" && input.ExpiresValue != 0 {
		expires, err = app.expirationTime(input.ExpiresValue, input.ExpiresUnit)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
	}

	user := app.contextGetUser(r)

	text := &data.Text{
		Title:     input.Title,
		Content:   input.Content,
		Format:    input.Format,
		Expires:   expires,
		IsPrivate: input.IsPrivate,
	}
	if !user.IsAnonymous() {
		text.UserID = &user.ID
	}

	slug, err := app.models.Texts.GenerateUniqueSlug(text.Title)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	text.Slug = slug

	v := validator.New()
	if data.ValidateText(v, text); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Texts.Insert(text)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/texts/%s", text.Slug))

	err = app.writeJSON(w, http.StatusCreated, envelope{"text": text}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showTextHandler will be used to show a text
func (app *application) showTextHandler(w http.ResponseWriter, r *http.Request) {
	slug, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)
	var userID *int64
	if !user.IsAnonymous() {
		userID = &user.ID
	}

	text, err := app.models.Texts.Get(slug, userID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"text": text}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTextHandler(w http.ResponseWriter, r *http.Request) {
	slug, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	text, err := app.models.Texts.Get(slug, &user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title        *string `json:"title"`
		Content      *string `json:"content"`
		Format       *string `json:"format"`
		ExpiresUnit  *string `json:"expiresUnit"`
		ExpiresValue *int    `json:"expiresValue"`
		IsPrivate    *bool   `json:"is_private"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		text.Title = *input.Title
	}
	if input.Content != nil {
		text.Content = *input.Content
	}
	if input.Format != nil {
		text.Format = *input.Format
	}
	if input.ExpiresUnit != nil && input.ExpiresValue != nil {
		text.Expires, err = app.expirationTime(*input.ExpiresValue, *input.ExpiresUnit)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
	}
	if input.IsPrivate != nil {
		text.IsPrivate = *input.IsPrivate
	}

	v := validator.New()
	if data.ValidateText(v, text); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Texts.Update(text, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"text": text}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTextHandler(w http.ResponseWriter, r *http.Request) {
	slug, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	err = app.models.Texts.Delete(slug, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "text successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addLikeHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	textID, err := app.readIntParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Likes.AddLike(user.ID, textID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Like added successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) removeLikeHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	textID, err := app.readIntParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Likes.RemoveLike(user.ID, textID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Like removed successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	textID, err := app.readIntParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Content string `json:"content"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &data.Comment{
		UserID:  user.ID,
		TextID:  textID,
		Content: input.Content,
	}

	err = app.models.Comments.AddComment(comment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"comment": comment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	commentID, err := app.readIntParam(r, "commentID")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Comments.DeleteComment(commentID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Comment deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
