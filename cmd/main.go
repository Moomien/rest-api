package main

import (
	"blog-restapi/internal/Blog"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("err: ", err)
	}
	storage := Blog.NewJsonStorage("jsonStorage")
	handler := Blog.BlogHandler{Storage: storage}

	http.HandleFunc("POST /blogs", handler.CreateBlog)
	http.HandleFunc("GET /blogs/{id}", handler.GetBlogByID)
	http.HandleFunc("GET /blogs", handler.GetBlogs)
	http.HandleFunc("PUT /blogs/{id}", handler.UpdateBlog)
	http.HandleFunc("DELETE /blogs/{id}", handler.DeleteBlog)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
