package lib

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Album struct {
	BackgroundImage  string `json:"backgroundImage"`
	Title            string `json:"title"`
	SubTitle         string `json:"subTitle"`
	CreatedAt        string `json:"createdAt"`
	RootPath         string `json:"-"`
	Data             *Node  `json:"data"`
	ServeStatically  bool   `json:"serveStatically"`
	NoScaledThumbs   bool   `json:"-"`
	NoScaledPreviews bool   `json:"-"`
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

func ScanDir(root string) (result *Node, err error) {
	log.Print("Reading images in %s", root)
	thumbDir := root + "/thumbs"
	previewDir := root + "/preview"
	parents := make(map[string]*Node)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(info.Name(), ".") {
			// skip dotfiles. log?
			return nil
		}
		if strings.HasPrefix(path, thumbDir) || strings.HasPrefix(path, previewDir) {
			// skip folder created by us
			return nil
		}
		parents[path] = &Node{
			FullPath: strings.TrimPrefix(path, root),
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
	result.FullPath = "/"
	return
}

func renderIndexTemplate(data Album) []byte {
	buf := &bytes.Buffer{}
	templateBinary := _escFSMustByte(false, "/index.html")
	tpl, err := template.New("index").Parse(string(templateBinary))
	if err != nil {
		log.Fatalf("Template parsing error: %v\n", err)
	}
	err = tpl.Execute(buf, data)
	if err != nil {
		log.Printf("Template execution error: %v\n", err)
	}
	return buf.Bytes()
}
