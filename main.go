package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

type Post struct {
	Title   string
	Date    time.Time
	Content template.HTML
	Slug    string
}

type BlogData struct {
	Posts []Post
	Title string
}

func main() {
	posts, err := loadPosts("posts")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveIndex(w, posts)
	})

	http.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/post/")
		servePost(w, posts, slug)
	})

	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadPosts(dir string) ([]Post, error) {
	var posts []Post

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".md") {
			post, err := parseMarkdownFile(path)
			if err != nil {
				return err
			}
			posts = append(posts, post)
		}
		return nil
	})

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, err
}

func parseMarkdownFile(filename string) (Post, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return Post{}, err
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(content, &buf); err != nil {
		return Post{}, err
	}

	base := filepath.Base(filename)
	slug := strings.TrimSuffix(base, ".md")
	
	title := extractTitle(string(content))
	if title == "" {
		title = slug
	}

	date := extractDate(filename)

	return Post{
		Title:   title,
		Date:    date,
		Content: template.HTML(buf.String()),
		Slug:    slug,
	}, nil
}

func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}

func extractDate(filename string) time.Time {
	base := filepath.Base(filename)
	parts := strings.Split(base, "-")
	if len(parts) >= 3 {
		dateStr := strings.Join(parts[:3], "-")
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			return date
		}
	}
	
	info, err := os.Stat(filename)
	if err != nil {
		return time.Now()
	}
	return info.ModTime()
}

func serveIndex(w http.ResponseWriter, posts []Post) {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="/styles/main.css">
</head>
<body>
    <header class="main-header">
        <h1>{{.Title}}</h1>
        <nav class="main-nav">
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="#about">About</a></li>
            </ul>
        </nav>
    </header>
    
    <main class="content">
        <section class="intro">
            <h2>Latest Posts</h2>
            {{range .Posts}}
            <article class="post-preview">
                <h3><a href="/post/{{.Slug}}">{{.Title}}</a></h3>
                <p class="post-date">{{.Date.Format "January 2, 2006"}}</p>
            </article>
            {{end}}
        </section>
    </main>

    <footer class="main-footer">
        <p>&copy; 2024 My Blog. All rights reserved.</p>
    </footer>
</body>
</html>`

	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := BlogData{
		Posts: posts,
		Title: "My Blog",
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func servePost(w http.ResponseWriter, posts []Post, slug string) {
	var post Post
	found := false
	for _, p := range posts {
		if p.Slug == slug {
			post = p
			found = true
			break
		}
	}

	if !found {
		http.NotFound(w, nil)
		return
	}

	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - My Blog</title>
    <link rel="stylesheet" href="/styles/main.css">
</head>
<body>
    <header class="main-header">
        <h1><a href="/">My Blog</a></h1>
        <nav class="main-nav">
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="#about">About</a></li>
            </ul>
        </nav>
    </header>
    
    <main class="content">
        <article class="post">
            <header class="post-header">
                <h1>{{.Title}}</h1>
                <p class="post-date">{{.Date.Format "January 2, 2006"}}</p>
            </header>
            <div class="post-content">
                {{.Content}}
            </div>
        </article>
        <nav class="post-nav">
            <a href="/">&larr; Back to Home</a>
        </nav>
    </main>

    <footer class="main-footer">
        <p>&copy; 2024 My Blog. All rights reserved.</p>
    </footer>
</body>
</html>`

	t, err := template.New("post").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}