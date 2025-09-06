package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png"{
		respondWithError(w, http.StatusBadRequest, "Invalid file type", nil)
		return
	}
	fileExt := mediaTypeToExt(contentType)


	metadata, err := cfg.db.GetVideo(videoID)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "no video found", err)
		return
	}
	if metadata.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Auhorization required", nil)
		return
	}

	videoKey := make([]byte, 32)
	rand.Read(videoKey)
	videoRef := base64.RawURLEncoding.EncodeToString(videoKey)

	fileName := fmt.Sprintf("%v%v", videoRef, fileExt)
	path := filepath.Join(cfg.assetsRoot, fileName)
	newImageFile, err := os.Create(path)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create file", err)
		return
	}

	_, err = io.Copy(newImageFile, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not write to file", err)
	}

	thumbnailUrl := fmt.Sprintf("http://localhost:%v/assets/%v%v", cfg.port, videoRef, fileExt)
	metadata.ThumbnailURL = &thumbnailUrl

	err = cfg.db.UpdateVideo(metadata)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update video", err)
	}

	respondWithJSON(w, http.StatusOK, metadata)
}
