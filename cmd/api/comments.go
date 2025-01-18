package main

import (
	"net/http"

	"github.com/ariefro/threads-server/internal/store"
)

// type Comment struct {
// 	ID        int64  `json:"id"`
// 	PostID    int64  `json:"post_id"`
// 	UserID    int64  `json:"user_id"`
// 	Content   string `json:"content"`
// 	CreatedAt string `json:"created_at"`
// 	User      User   `json:"user"`
// }

type CreateCommentPayload struct {
	Content string `json:"content"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	post, err := getPostFromCtx(r)
	if err != nil {
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

	comment := &store.Comment{
		PostID:  post.ID,
		UserID:  1, // todo: change after auth
		Content: payload.Content,
	}

	ctx := r.Context()

	// Add the new comment to the database
	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Retrieve the updated comments for the post
	comments, err := app.store.Comments.GetByPostID(ctx, post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Include the updated comments in the post response
	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
