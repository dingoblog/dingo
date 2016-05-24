package handler

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dinever/golf"
	"github.com/dingoblog/dingo/app/model"
	"github.com/dingoblog/dingo/app/utils"
	"github.com/facebookgo/inject"
)

// Handler module will be used to concentrate all sub-modules & provide them with their dependencies
type Module struct {
	Generic GenericAPI // Temporal name for API handlers
}

var app *golf.Application

func (m *Module) Populate(g inject.Graph) {

	// Share dependencies with structs that holds handlers
	err := g.Provide(
		&inject.Object{Value: &m.Generic},
	)

	if err != nil {
		err = fmt.Errorf("failed to provide di on handler module: %v", err)
		panic(err)
	}

	if err := g.Populate(); err != nil {
		err = fmt.Errorf("failed to populate di: %v", err)
		panic(err)
	}
}

func (m *Module) Initialize() *golf.Application {

	app = golf.New()

	app.Config.Set("app/static_dir", "static")
	app.Config.Set("app.log_dir", "tmp/log")
	app.Config.Set("app/upload_dir", "upload")
	upload_dir, _ := app.Config.GetString("app/upload_dir", "upload")
	m.registerMiddlewares()
	m.registerFuncMap()
	RegisterFunctions(app)
	theme := model.GetSettingValue("theme")
	app.View.SetTemplateLoader("base", "view")
	app.View.SetTemplateLoader("admin", filepath.Join("view", "admin"))
	app.View.SetTemplateLoader("theme", filepath.Join("view", theme))
	//      static_dir, _ := app.Config.GetString("app/static_dir", "static")
	app.Static("/upload/", upload_dir)
	app.Static("/admin/", filepath.Join("view", "admin", "assets", "dist"))
	app.Static("/", filepath.Join("view", theme, "assets", "dist"))

	app.SessionManager = golf.NewMemorySessionManager()
	app.Error(404, NotFoundHandler)

	m.registerAdminURLHandlers()
	m.registerHomeHandler()
	m.registerAPIHandler()

	return app
}

func (m *Module) registerFuncMap() {
	app.View.FuncMap["DateFormat"] = utils.DateFormat
	app.View.FuncMap["Now"] = utils.Now
	app.View.FuncMap["Html2Str"] = utils.Html2Str
	app.View.FuncMap["FileSize"] = utils.FileSize
	app.View.FuncMap["Setting"] = model.GetSettingValue
	app.View.FuncMap["Navigator"] = model.GetNavigators
	app.View.FuncMap["Md2html"] = utils.Markdown2HtmlTemplate
}

func (m *Module) registerMiddlewares() {
	app.Use(
		golf.LoggingMiddleware(os.Stdout),
		golf.RecoverMiddleware,
		golf.SessionMiddleware,
	)
}

func (m *Module) registerAdminURLHandlers() {
	authChain := golf.NewChain(AuthMiddleware)
	app.Get("/login/", AuthLoginPageHandler)
	app.Post("/login/", AuthLoginHandler)

	app.Get("/signup/", AuthSignUpPageHandler)
	app.Post("/signup/", AuthSignUpHandler)

	app.Get("/logout/", AuthLogoutHandler)

	app.Get("/admin/", authChain.Final(AdminHandler))

	app.Get("/admin/profile/", authChain.Final(ProfileHandler))
	app.Post("/admin/profile/", authChain.Final(ProfileChangeHandler))

	app.Get("/admin/editor/post/", authChain.Final(PostCreateHandler))
	app.Post("/admin/editor/post/", authChain.Final(PostSaveHandler))

	app.Get("/admin/editor/page/", authChain.Final(PageCreateHandler))
	app.Post("/admin/editor/page/", authChain.Final(PageSaveHandler))

	app.Get("/admin/posts/", authChain.Final(AdminPostHandler))
	app.Get("/admin/pages/", authChain.Final(AdminPageHandler))

	app.Get("/admin/editor/:id/", authChain.Final(ContentEditHandler))
	app.Post("/admin/editor/:id/", authChain.Final(ContentSaveHandler))
	app.Delete("/admin/editor/:id/", authChain.Final(ContentRemoveHandler))

	app.Get("/admin/comments/", authChain.Final(CommentViewHandler))
	app.Post("/admin/comments/", authChain.Final(CommentAddHandler))
	app.Put("/admin/comments/", authChain.Final(CommentUpdateHandler))
	app.Delete("/admin/comments/", authChain.Final(CommentRemoveHandler))

	app.Get("/admin/setting/", authChain.Final(SettingViewHandler))
	app.Post("/admin/setting/", authChain.Final(SettingUpdateHandler))
	app.Post("/admin/setting/custom/", authChain.Final(SettingCustomHandler))
	app.Post("/admin/setting/nav/", authChain.Final(SettingNavHandler))
	//
	app.Get("/admin/files/", authChain.Final(FileViewHandler))
	app.Delete("/admin/files/", authChain.Final(FileRemoveHandler))
	app.Post("/admin/files/upload/", authChain.Final(FileUploadHandler))

	app.Get("/admin/password/", authChain.Final(AdminPasswordPage))
	app.Post("/admin/password/", authChain.Final(AdminPasswordChange))

	app.Get("/admin/monitor/", authChain.Final(AdminMonitorPage))
}

func (m *Module) registerHomeHandler() {
	statsChain := golf.NewChain()
	app.Get("/", statsChain.Final(HomeHandler))
	app.Get("/page/:page/", HomeHandler)
	app.Post("/comment/:id/", CommentHandler)
	app.Get("/tag/:tag/", TagHandler)
	app.Get("/tag/:tag/page/:page/", TagHandler)
	app.Get("/feed/", RssHandler)
	app.Get("/sitemap.xml", SiteMapHandler)
	app.Get("/:slug/", statsChain.Final(ContentHandler))
}

func (m *Module) registerAPIHandler() {
	// Auth
	app.Post("/auth", JWTAuthLoginHandler)
	app.Get("/auth", JWTAuthValidateHandler)

	// register the API handler
	app.Get("/api", APIDocumentationHandler)
	app.Get("/api/ping", m.Generic.Ping)

	// Posts
	app.Get("/api/posts", APIPostsHandler)
	app.Get("/api/posts/:id", APIPostHandler)
	app.Get("/api/posts/slug/:slug", APIPostSlugHandler)

	// Tags
	app.Get("/api/tags", m.Generic.GetAllTags)
	app.Get("/api/tags/:id", APITagHandler)
	app.Get("/api/tags/slug/:slug", APITagSlugHandler)

	// Users
	app.Get("/api/users", APIUsersHandler)
	app.Get("/api/users/:id", APIUserHandler)
	app.Get("/api/users/slug/:slug", APIUserSlugHandler)
	app.Get("/api/users/email/:email", APIUserEmailHandler)
}
