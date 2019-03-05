package lib

type Album struct {
	BackgroundImage string   `json:"backgroundImage"`
	Title           string   `json:"title"`
	SubTitle        string   `json:"subTitle"`
	CreatedAt       string   `json:"createdAt"`
	Data            []Folder `json:"data"`
}

type Folder struct {
	Parent *Folder `json:"parent"`
	Name   string  `json:"name"`
	Images []Image `json:"images"`
}

// Album doc
type Image struct {
	Filename string `json:"filename"`
	ExifData string `json:"exif"`
}

type ExifData struct {
	DateTime              string
	ExifImageLength       int
	ExifImageWidth        int
	ExposureBiasValue     int
	ExposureMode          string
	ExposureTime          string
	FNum                  string
	FNumber               string
	FileSize              int
	Flash                 string
	FocalLength           string
	FocalLengthIn35mmFilm int
	ISOSpeedRatings       int
	Model                 string
	WhiteBalance          string
}

/*

https://golang.org/pkg/path/filepath/#Walk
https://stackoverflow.com/questions/12657365/extracting-directory-hierarchy-using-go-language

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
