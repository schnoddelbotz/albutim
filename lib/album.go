package lib

import (
	"bytes"
	"html/template"
	"io/ioutil"
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
	NoCacheScaled    bool   `json:"-"`
	NumThreads       int    `json:"-"`
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

func BuildAlbum(a Album) {
	log.Printf("Building Album '%s'", a.Title)

	if !a.NoScaledPreviews || !a.NoScaledThumbs {
		a.buildThumbsAndPreviews()
	}

	log.Print("Copying index.html and template files to album...")
}

func ScanDir(root string) (result *Node, err error) {
	log.Printf("Reading images in %s ...", root)
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
	log.Print("Reading images: completed")
	result.FullPath = "/"
	return
}

func (n *Node) getAllImagePaths() []string {
	var toVisit []*Node
	var images []string
	toVisit = append(toVisit, n)
	for len(toVisit) > 0 {
		c := toVisit[0]
		if c.IsImage {
			images = append(images, c.FullPath)
		} else {
			toVisit = append(toVisit, c.Children...)
		}
		toVisit = toVisit[1:]
	}
	return images
}

func (a *Album) buildThumbsAndPreviews() {
	log.Printf("Building previews and thumbnails using %d threads ...", a.NumThreads)
	jobs := make(chan string, 100)
	results := make(chan int, 100)
	for w := 0; w <= a.NumThreads-1; w++ {
		go a.imageScalingWorker(w, jobs, results)
	}
	imagePaths := a.Data.getAllImagePaths()
	for _, path := range imagePaths {
		jobs <- path
	}
	close(jobs)
	for range imagePaths {
		<-results
	}
	log.Print("Building previews and thumbnails: done.")
}

func (a *Album) imageScalingWorker(id int, jobs <-chan string, results chan<- int) {
	for relativeImagePath := range jobs {
		originalPath := a.RootPath + relativeImagePath
		previewPath := a.RootPath + "/preview" + relativeImagePath
		thumbPath := a.RootPath + "/thumbs" + relativeImagePath
		var todo []string

		if _, err := os.Stat(thumbPath); os.IsNotExist(err) && !a.NoScaledThumbs {
			todo = append(todo, "thumb")
			thumb, err := getScaled(originalPath, 0, 105 /* FIXME config value */)
			if err != nil {
				log.Printf("scalingError for %s: %s", relativeImagePath, err)
				return
			}
			a.addCache(thumbPath, thumb)
		}

		if _, err := os.Stat(previewPath); os.IsNotExist(err) && !a.NoScaledPreviews {
			todo = append(todo, "preview")
			preview, err := getScaled(originalPath, 0, 700 /* FIXME config value */)
			if err != nil {
				log.Printf("scalingError for %s: %s", relativeImagePath, err)
				return
			}
			a.addCache(previewPath, preview)
		}

		if len(todo) > 0 {
			log.Printf("[thread-%d] %s created for %s", id, strings.Join(todo, "+"), relativeImagePath)
		}

		results <- 1
	}
}

func (a *Album) addCache(file string, data []byte) {
	if a.NoCacheScaled {
		return
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		//log.Printf("Add to cache: %s", file)
		err = os.MkdirAll(filepath.Dir(file), os.ModePerm)
		if err != nil {
			log.Printf("mkdir %s error: %s", filepath.Dir(file), err)
			return
		}
		err = ioutil.WriteFile(file, data, 0644)
		if err != nil {
			log.Printf("write %s error: %s", file, err)
		}
	}
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
