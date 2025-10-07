package dto

type ImageUrlsRequest struct {
	Images []ImageUrlRequest `json:"images"`
}

type ImageUrlRequest struct {
	Name string `json:"name"`
}
