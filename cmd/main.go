package main

import (
	"blog-restapi/internal/Blog"
	"log"
	"net/http"
)

func main() {
	storage := &Blog.JsonStorage{}
	storage = storage.NewStorage("jsonStorage")
	handler := Blog.BlogHandler{Storage: storage}

	http.HandleFunc("POST /blogs", handler.CreateBlog)
	http.HandleFunc("GET /blogs/{id}", handler.GetBlogByID)
	http.HandleFunc("GET /blogs", handler.GetBlogs)
	http.HandleFunc("PUT /blogs/{id}", handler.UpdateBlog)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
