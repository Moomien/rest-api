package Blog

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type BlogHandler struct {
	Storage Storage
}

// GET /blogs
func (b *BlogHandler) GetBlogs(w http.ResponseWriter, r *http.Request) {
	data, err := b.Storage.LoadAll()
	if err != nil {
		http.Error(w, "Failed loading data from db", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	w.WriteHeader(http.StatusOK)

}

// GET /blogs/{id}
// just load data from map
func (b *BlogHandler) GetBlogByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Err convert string to int:", http.StatusInternalServerError)
		return
	}

	blog, err := b.Storage.LoadById(id)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

// POST /blogs/Blog{}
func (b *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var newBlog Blog
	if err := json.NewDecoder(r.Body).Decode(&newBlog); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := b.Storage.Save(newBlog); err != nil {
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// PUT(full update)
// PUT /blogs/{id}
func (b *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Wrong numbers", http.StatusBadRequest)
		return
	}

	var newBlog Blog
	if err := json.NewDecoder(r.Body).Decode(&newBlog); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	newBlog.ID = id
	if err := b.Storage.SaveById(id, newBlog); err != nil {
		if err.Error() == "Not found data" {
			http.Error(w, "Blog not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DELETE /blogs/{id}
func (b *BlogHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Wrong ID", http.StatusBadRequest)
		return
	}
	if err := b.Storage.Delete(id); err != nil {
		http.Error(w, "Error delete", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	response := map[string]string{"message": "Blog deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
