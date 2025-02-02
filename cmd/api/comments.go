package main

import (
	"net/http"

	"github.com/getsentry/sentry-go"

	"github.com/ariefro/threads-server/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	post, err := getPostFromCtx(r)
	if err != nil {
		sentry.CaptureException(err)
		app.internalServerError(w, r, err)
		return
	}

	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserFromContext(r)

	comment := &store.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: payload.Content,
	}

	ctx := r.Context()

	// Add the new comment to the database
	if err := app.store.Comments.Create(ctx, comment); err != nil {
		sentry.CaptureException(err)
		app.internalServerError(w, r, err)
		return
	}

	// Retrieve the updated comments for the post
	comments, err := app.store.Comments.GetByPostID(ctx, post.ID)
	if err != nil {
		sentry.CaptureException(err)
		app.internalServerError(w, r, err)
		return
	}

	// Include the updated comments in the post response
	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		sentry.CaptureException(err)
		app.internalServerError(w, r, err)
		return
	}
}
