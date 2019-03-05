package lib

import (
	"github.com/xor-gate/goexif2/exif"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Album struct {
	BackgroundImage string `json:"backgroundImage"`
	Title           string `json:"title"`
	SubTitle        string `json:"subTitle"`
	CreatedAt       string `json:"createdAt"`
	RootPath        string `json:"rootPath"`
	Data            *Node  `json:"data"`
}

type ExifData struct {
	DateTime              string `json:"dateTime,omitempty"`
	Height                int    `json:"height,omitempty"`
	Width                 int    `json:"width,omitempty"`
	ExposureBiasValue     int    `json:"exposureBiasValue,omitempty"`
	ExposureMode          string `json:"exposureMode,omitempty"`
	ExposureTime          string `json:"exposureTime,omitempty"`
	FNum                  string `json:"fNum,omitempty"`
	FNumber               string `json:"fNumber,omitempty"`
	FileSize              int    `json:"fileSize,omitempty"`
	Flash                 string `json:"flash,omitempty"`
	FocalLength           string `json:"focalLength,omitempty"`
	FocalLengthIn35mmFilm int    `json:"focalLengthIn35mmFilm,omitempty"`
	ISOSpeedRatings       int    `json:"iSOSpeedRatings,omitempty"`
	Model                 string `json:"model,omitempty"`
	WhiteBalance          string `json:"whiteBalance,omitempty"`
}

// Node represents a node in a directory tree.
type Node struct {
	FullPath string   `json:"path"`
	Name     string   `json:"name"`
	Size     int64    `json:"size"`
	IsDir    bool     `json:"is_dir"`
	IsImage  bool     `json:"is_image"`
	Children []*Node  `json:"children"`
	Parent   *Node    `json:"-"`
	ExifData ExifData `json:"exifdata,omitempty"`
}

func scanDir(root string) (result *Node, err error) {
	// FIXME: re-add abspath...?
	parents := make(map[string]*Node)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		parents[path] = &Node{
			FullPath: path,
			Name:     info.Name(),
			IsDir:    info.IsDir(),
			Size:     info.Size(),
			Children: make([]*Node, 0),
		}
		if !info.IsDir() && isImage(path) {
			parents[path].IsImage = true
			ed, err := getExif(path)
			if err == nil {
				parents[path].ExifData = ed
			}
		}
		return nil
	}
	if err = filepath.Walk(root, walkFunc); err != nil {
		return
	}
	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists {
			result = node
		} else {
			node.Parent = parent
			parent.Children = append(parent.Children, node)
		}
	}
	return
}

func isImage(path string) bool {
	extension := strings.ToLower(filepath.Ext(path))
	if extension == ".jpg" || extension == ".jpeg" {
		return true
	}
	return false
}

func getExif(path string) (data ExifData, err error) {
	f, err := _getExifData(path)
	if err == nil {
		if tag, err := f.Get(exif.Model); err == nil {
			if val, err := tag.StringVal(); err == nil {
				data.Model = val
			}
		}
		if tag, err := f.Get(exif.ExposureTime); err == nil {
			if val, err := tag.Rat(0); err == nil {
				data.ExposureTime = val.String()
			}
		}
		if tag, err := f.Get(exif.FNumber); err == nil {
			if val, err := tag.Rat(0); err == nil {
				data.FNumber = val.String()
			}
		}
		if tag, err := f.Get(exif.PixelXDimension); err == nil {
			if val, err := tag.Int(0); err == nil {
				data.Width = val
			}
		}
		if tag, err := f.Get(exif.PixelYDimension); err == nil {
			if val, err := tag.Int(0); err == nil {
				data.Height = val
			}
		}
		if tag, err := f.Get(exif.DateTime); err == nil {
			if val, err := tag.StringVal(); err == nil {
				data.DateTime = val
			}
		}
	}
	return
}

func _getExifData(fname string) (x *exif.Exif, err error) {
	log.Printf("getExif: %s", fname)
	f, err := os.Open(fname)
	if err != nil {
		log.Printf("Error reading EXIF from %s: %s", fname, err)
	}

	//FIXME: fails on my images...
	//exif.RegisterParsers(mknote.All...)

	x, err = exif.Decode(f)
	if err != nil {
		log.Printf("Decode NOEXIF in %s: %s", fname, err)
		return nil, err
	}
	return x, nil
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
