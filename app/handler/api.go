package handler

import (
	"strconv"

	"github.com/dinever/golf"
	"github.com/dingoblog/dingo/app/model"
	"github.com/dingoblog/dingo/app/modules/posts"
	"github.com/jinzhu/gorm"
)

// Generic API holds a temporary way to group API handlers and share dependencies with them
// Even though we're currently holding a pointer to the DB we should abstract the data handling to
// Other submodules taking advantage of the dependency injection
type GenericAPI struct {
	DB    *gorm.DB      `inject:""`
	Posts *posts.Module `inject:""`
}

func (g GenericAPI) Ping(c *golf.Context) {

	c.JSON(map[string]interface{}{
		"status": "pong",
	})
}

// Get all tags
func (group GenericAPI) GetAllTags(c *golf.Context) {

	// Expressive and cleaner
	tags := group.Posts.GetTags()

	c.JSON(tags)
}

// APIDocumentationHandler shows which routes match with what functionality,
// similar to https://api.github.com
func APIDocumentationHandler(ctx *golf.Context) {
	// Go doesn't display maps in the order they appear here, so if the order
	// of these routes is important, it might be better to use a struct
	routes := map[string]interface{}{
		"auth_url":              "/auth/",
		"api_documentation_url": "/api/",
		"comments_url":          "/api/comments",
		"comment_url":           "/api/comments/:id",
		"comment_post_url":      "/api/comments/post/:id",
		"posts_url":             "/api/posts/",
		"post_url":              "/api/posts/:id",
		"post_slug_url":         "/api/posts/slug/:slug",
		"tags_url":              "/api/tags/",
		"tag_url":               "/api/tags/:id",
		"tag_slug_url":          "/api/tags/slug/:slug",
		"users_url":             "/api/users/",
		"user_url":              "/api/users/:id",
		"user_slug_url":         "/api/users/slug/:slug",
		"user_email_url":        "/api/users/email/:email",
	}
	ctx.JSONIndent(routes, "", "  ")
}

// APICommentsHandler retrieves all the comments.
func APICommentsHandler(ctx *golf.Context) {
	ctx.JSONIndent(map[string]interface{}{
		"message": "Not implemented",
	}, "", "  ")
}

// APICommentHandler retrieves a comment with the given comment id.
func APICommentHandler(ctx *golf.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		handleErr(ctx, 500, err)
		return
	}
	comment := &model.Comment{Id: int64(id)}
	err = comment.GetCommentById()
	if err != nil {
		handleErr(ctx, 404, err)
		return
	}
	ctx.JSONIndent(comment, "", "  ")
}

// APICommentPostHandler retrives the tag with the given post id.
func APICommentPostHandler(ctx *golf.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		handleErr(ctx, 500, err)
		return
	}
	comments := new(model.Comments)
	err = comments.GetCommentsByPostId(int64(id))
	if err != nil {
		handleErr(ctx, 404, err)
		return
	}
	ctx.JSONIndent(comments, "", "  ")
}

// APIPostsHandler gets every page, ordered by publication date.
func APIPostsHandler(ctx *golf.Context) {
	posts := new(model.Posts)
	err := posts.GetAllPostList(false, true, "published_at DESC")
	if err != nil {
		handleErr(ctx, 404, err)
		return
	}
	ctx.JSONIndent(posts, "", "  ")
}

// APIPostHandler retrieves the post with the given ID.
func APIPostHandler(ctx *golf.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		handleErr(ctx, 500, err)
		return
	}
	post := &model.Post{Id: int64(id)}
	err = post.GetPostById()
	if err != nil {
		handleErr(ctx, 404, err)
		return
	}
	ctx.JSONIndent(post, "", "  ")
}

// APIPostSlugHandler retrieves the post with the given slug.
func APIPostSlugHandler(ctx *golf.Context) {
	slug := ctx.Param("slug")
	post := new(model.Post)
	err := post.GetPostBySlug(slug)
	if err != nil {
		handleErr(ctx, 404, err)
		return
	}
	ctx.JSONIndent(post, "", "  ")
}

// APITagsHandler retrieves all the tags.
func APITagsHandler(ctx *golf.Context) {
	tags := new(model.Tags)
	err := tags.GetAllTags()
	if err != nil {
		handleErr(ctx, 404, err)
		return
	}
	ctx.JSONIndent(tags, "", "  ")
}

// APITagHandler retrieves the tag with the given id.
func APITagHandler(ctx *golf.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		handleErr(ctx, 500, err)
		return
	}
	tag := &model.Tag{Id: int64(id)}
	err = tag.GetTag()
	if err != nil {
		handleErr(ctx, 404, err)
		return
	}
	ctx.JSONIndent(tag, "", "  ")
}

// APITagSlugHandler retrieves the tag(s) with the given slug.
func APITagSlugHandler(ctx *golf.Context) {
	slug := ctx.Param("slug")
	tags := &model.Tag{Slug: slug}
	err := tags.GetTagBySlug()
	if err != nil {
		handleErr(ctx, 500, err)
		return
	}
	ctx.JSONIndent(tags, "", "  ")
}

// APIUsersHandler retrieves all users.
func APIUsersHandler(ctx *golf.Context) {
	ctx.JSONIndent(map[string]interface{}{
		"message": "Not implemented",
	}, "", "  ")
}

// APIUserHandler retrieves the user with the given id.
func APIUserHandler(ctx *golf.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		handleErr(ctx, 500, err)
		return
	}
	user := &model.User{Id: int64(id)}
	err = user.GetUserById()
	if err != nil {
		handleErr(ctx, 404, err)
		return
	}
	ctx.JSONIndent(user, "", "  ")
}

// APIUserSlugHandler retrives the user with the given slug.
func APIUserSlugHandler(ctx *golf.Context) {
	slug := ctx.Param("slug")
	user := &model.User{Slug: slug}
	err := user.GetUserBySlug()
	if err != nil {
		handleErr(ctx, 404, err)
		return
	}
	ctx.JSONIndent(user, "", "  ")
}

// APIUserEmailHandler retrieves the user with the given email.
func APIUserEmailHandler(ctx *golf.Context) {
	email := ctx.Param("email")
	user := &model.User{Email: email}
	err := user.GetUserByEmail()
	if err != nil {
		handleErr(ctx, 404, err)
		return
	}
	ctx.JSONIndent(user, "", "  ")
}

// handleErr sends the staus code and error message formatted as JSON.
func handleErr(ctx *golf.Context, statusCode int, err error) {
	ctx.JSONIndent(map[string]interface{}{
		"statusCode": statusCode,
		"error":      err.Error(),
	}, "", "  ")
}
