package handler

import (
	"github.com/dinever/golf"
	"github.com/dingoblog/dingo/app/model"
	"github.com/dingoblog/dingo/app/utils"
	"html/template"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func RegisterFunctions(app *golf.Application) {
	app.View.FuncMap["Tags"] = getAllTags
	app.View.FuncMap["RecentPosts"] = getRecentPosts
	app.View.FuncMap["RecentComments"] = getRecentComments
}

func HomeHandler(ctx *golf.Context) {
	p := ctx.Param("page")

	var page int
	if p == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(p)
	}
	posts := new(model.Posts)
	pager, err := posts.GetPostList(int64(page), 5, false, true, "published_at DESC")
	if err != nil {
		ctx.Abort(404)
		return
	}
	// theme := model.GetSetting("site_theme")
	data := map[string]interface{}{
		"Title": "Home",
		"Posts": posts,
		"Pager": pager,
	}
	//	updateSidebarData(data)
	ctx.Loader("theme").Render("index.html", data)
}

func ContentHandler(ctx *golf.Context) {
	slug := ctx.Param("slug")
	post := new(model.Post)
	err := post.GetPostBySlug(slug)
	if err != nil || !post.IsPublished {
		ctx.Abort(404)
		return
	}
	post.Hits++
	data := map[string]interface{}{
		"Title":    post.Title,
		"Post":     post,
		"Content":  post,
		"Comments": post.Comments,
	}
	if post.IsPage {
		ctx.Loader("theme").Render("page.html", data)
	} else {
		ctx.Loader("theme").Render("article.html", data)
	}
}

func CommentHandler(ctx *golf.Context) {
	id := ctx.Param("id")
	cid, _ := strconv.Atoi(id)
	post := new(model.Post)
	post.Id = int64(cid)
	err := post.GetPostById()
	if cid < 1 || err != nil {
		ctx.JSON(map[string]interface{}{
			"status": "error",
		})
	}
	c := model.NewComment()
	c.Author = ctx.Request.FormValue("author")
	c.Email = ctx.Request.FormValue("email")
	c.Website = ctx.Request.FormValue("website")
	c.Content = strings.Replace(utils.Html2Str(template.HTMLEscapeString(ctx.Request.FormValue("comment"))), "\n", "<br/>", -1)
	c.Avatar = utils.Gravatar(c.Email, "50")
	c.PostId = post.Id
	pid, _ := strconv.Atoi(ctx.Request.FormValue("pid"))
	c.Parent = int64(pid)
	c.Ip = ctx.Request.RemoteAddr
	c.UserAgent = ctx.Request.UserAgent()
	c.UserId = 0
	msg := c.ValidateComment()
	if msg == "" {
		if err := c.Save(); err != nil {
			ctx.JSON(map[string]interface{}{
				"status": "error",
				"msg":    "Can not comment on this post.",
			})
		}
		post.CommentNum++
        err = post.Save(post.Tags()...)
		if err != nil {
			log.Printf("[Error]: Can not increase comment count for post %v: %v", post.Id, err.Error())
		}
		ctx.JSON(map[string]interface{}{
			"res":     true,
			"comment": c.ToJson(),
		})
		if err = model.NewMessage("comment", c).Insert(); err != nil {
			panic(err)
		}
	} else {
		ctx.JSON(map[string]interface{}{
			"status": "error",
			"msg":    msg,
		})
	}
}

func TagHandler(ctx *golf.Context) {
	p := ctx.Param("page")

	var page int
	if p == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(p)
	}

	t := ctx.Param("tag")
	tagSlug, _ := url.QueryUnescape(t)
	tag := &model.Tag{Slug: tagSlug}
	err := tag.GetTagBySlug()
	if err != nil {
		NotFoundHandler(ctx)
		return
	}
	posts := new(model.Posts)
	pager, err := posts.GetPostsByTag(tag.Id, int64(page), 5, true)
	data := map[string]interface{}{
		"Posts": posts,
		"Pager": pager,
		"Tag":   tag,
		"Title": tag.Name,
	}
	ctx.Loader("theme").Render("tag.html", data)
}

func SiteMapHandler(ctx *golf.Context) {
	baseUrl := model.GetSettingValue("site_url")
	posts := new(model.Posts)
	_, _ = posts.GetPostList(1, 50, false, true, "published_at DESC")
	navigators := model.GetNavigators()
	now := utils.Now().Format(time.RFC3339)

	articleMap := make([]map[string]string, posts.Len())
	for i := 0; i < posts.Len(); i++ {
		m := make(map[string]string)
		m["Link"] = strings.Replace(baseUrl+posts.Get(i).Url(), baseUrl+"/", baseUrl, -1)
		m["Created"] = posts.Get(i).PublishedAt.Format(time.RFC3339)
		articleMap[i] = m
	}

	navMap := make([]map[string]string, 0)
	for _, n := range navigators {
		m := make(map[string]string)
		if n.Url == "/" {
			continue
		}
		if strings.HasPrefix(n.Url, "/") {
			m["Link"] = strings.Replace(baseUrl+n.Url, baseUrl+"/", baseUrl, -1)
		} else {
			m["Link"] = n.Url
		}
		m["Created"] = now
		navMap = append(navMap, m)
	}

	ctx.SetHeader("Content-Type", "application/rss+xml;charset=UTF-8")
	ctx.Loader("base").Render("sitemap.xml", map[string]interface{}{
		"Title":      model.GetSettingValue("site_title"),
		"Link":       baseUrl,
		"Created":    now,
		"Posts":      articleMap,
		"Navigators": navMap,
	})
}

func RssHandler(ctx *golf.Context) {
	baseUrl := model.GetSettingValue("site_url")
	posts := new(model.Posts)
	_, _ = posts.GetPostList(1, 20, false, true, "published_at DESC")
	articleMap := make([]map[string]string, posts.Len())
	for i := 0; i < posts.Len(); i++ {
		m := make(map[string]string)
		m["Title"] = posts.Get(i).Title
		m["Link"] = posts.Get(i).Url()
		m["Author"] = posts.Get(i).Author().Name
		m["Desc"] = posts.Get(i).Excerpt()
		m["Created"] = posts.Get(i).CreatedAt.Format(time.RFC822)
		articleMap[i] = m
	}
	ctx.SetHeader("Content-Type", "text/xml; charset=utf-8")

	ctx.Loader("base").Loader("base").Render("rss.xml", map[string]interface{}{
		"Title":   model.GetSettingValue("site_title"),
		"Link":    baseUrl,
		"Desc":    model.GetSettingValue("site_description"),
		"Created": utils.Now().Format(time.RFC822),
		"Posts":   articleMap,
	})
}
