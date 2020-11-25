package function

//GiphyResponse - represents (part of) JSON response sent by Giphy API
type GiphyResponse struct {
	Data Data `json:"data"`
}

//Data - high level attribute
type Data struct {
	Title  string `json:"title"`
	Images Images `json:"images"`
}

//Images - Contains downsized format GIF info
type Images struct {
	Downsized Downsized `json:"downsized"`
}

//Downsized - Giphy URL for downsized format GIF
type Downsized struct {
	URL string `json:"url"`
}
