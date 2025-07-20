package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dokhoyan/go-messenger-microservices/subscriber_service/models"
	"github.com/go-openapi/strfmt"
)

func main() {
	http.HandleFunc("/v1/friends/requests", friendsRequestHandler)//сделать и получить запросы
	http.HandleFunc("/v1/friends/requests/", friendRequestRouter)//отклонить и принять 
	http.HandleFunc("/v1/users/", userFriendsRouter)//список друзей

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}



func friendsRequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleMakeFriendRequest(w, r) //запрос при пост 
	case http.MethodGet:
		getFriendRequests(w, r) //получаем заявки при гет
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func friendRequestRouter(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/confirm") {
		confirmFriendRequest(w, r) // принять
		return
	}
	if strings.HasSuffix(r.URL.Path, "/decline") {
		declineFriendRequest(w, r)   //отклонить 
		return
	}
	http.NotFound(w, r)
}

func userFriendsRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || !strings.HasSuffix(r.URL.Path, "/friends") {
		http.NotFound(w, r)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	idStr := parts[3]
	_, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	getFriends(w, r)
}

// === Handlers ===

func handleMakeFriendRequest(w http.ResponseWriter, r *http.Request) {
	var req models.MakeFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.FromUserID == nil || req.ToUserID == nil {
		http.Error(w, `{"message":"Missing from_user_id or to_user_id"}`, http.StatusBadRequest)
		return
	}

	if *req.FromUserID <= 0 || *req.ToUserID <=0 {
		http.Error(w, `{"message":"user_id must be positive"}`, http.StatusBadRequest)
		return
	}

	if *req.FromUserID == *req.ToUserID  {
		http.Error(w, `{"Cannot send friend request to yourself"}`, http.StatusBadRequest)
		return
	}

	// TODO: логика сохранения в базу
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "request_sent"})
}

func getFriendRequests(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, `{"message":"Missing user_id"}`, http.StatusBadRequest)
		return
	}
	_, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"message":"Invalid user_id"}`, http.StatusBadRequest)
		return
	}

	requests := []models.FriendRequest{
		{RequestID: 1, FromUserID: 2, ToUserID: 123, CreatedAt: strfmt.DateTime{}},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

func confirmFriendRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	requestID, err := getRequestIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, `{"message":"Invalid request ID"}`, http.StatusBadRequest)
		return
	}

	// TODO: логика подтверждения
	_ = requestID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "confirmed"})
}

func declineFriendRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	requestID, err := getRequestIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, `{"message":"Invalid request ID"}`, http.StatusBadRequest)
		return
	}

	// TODO: логика отклонения
	_ = requestID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "declined"})
}

func getFriends(w http.ResponseWriter, r *http.Request) {
	friends := []models.Friend{
		{UserID: 42, Name: "Alice", AvatarPhotoURL: "https://example.com/avatar.jpg"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}



func getRequestIDFromPath(path string) (int64, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		return 0, errors.New("invalid path")
	}
	
	return strconv.ParseInt(parts[4], 10, 64)
}
