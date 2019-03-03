package lib

// Album doc
type Image struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" db:"name"`
	Path        string  `json:"path" db:"path"`
	//Tracks      []Track `json:"tracks"`
	ArtistName  string  `json:"artist_name" db:"artist_name"`
	ArtistCount int     `json:"artist_count" db:"artist_count"`
}

// Albums doc
type Album struct {
	Error  string  `json:"error"`
	Albums []Image `json:"albums"`
}

/*

{
  "background": null,
  "created": "just now dude",
  "images": {
    "/subfolder1/": {
      "trump burgers.jpg": {
        "FileSize": "0.21 MB"
      },
      "vitamine.jpg": {
        "FileSize": "1.17 MB"
      }
    },
    "/subfolder2/": {
      "10143d473b7702fb.jpg": {
        "FileSize": "0.07 MB"
      },
      "2017-01-03_14.47.49_720.jpg": {
        "DateTime": "2017:01:03 14:47:49",
        "ExifImageLength": "3024",
        "ExifImageWidth": "4032",
        "ExposureBiasValue": "0",
        "ExposureMode": "Auto Exposure",
        "ExposureTime": "1250119/62500000",
        "FNum": "2",
        "FNumber": "2",
        "FileSize": "0.07 MB",
        "Flash": "Flash did not fire",
        "FocalLength": "467/100",
        "FocalLengthIn35mmFilm": "0",
        "ISOSpeedRatings": "208",
        "Model": "Nexus 5X",
        "WhiteBalance": "Auto"
      },
  },
  "subtitle": "Yet another photo blog",
  "title": "test"
}

 */