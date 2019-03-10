package lib

import (
	"bytes"
	"encoding/json"
	"github.com/nfnt/resize"
	"html/template"
	"image"
	"image/jpeg"
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
	Parent   *Node    `json:"-"`
	Children []*Node  `json:"children"`
	Album    *Album   `json:"-"`
	FullPath string   `json:"-"`
	WebPath  string   `json:"path"`
	Name     string   `json:"name"`
	Size     int64    `json:"size"`
	IsDir    bool     `json:"is_dir"`
	IsImage  bool     `json:"is_image"`
	ExifData ExifData `json:"exifdata,omitempty"`
}

func BuildAlbum(a Album) {
	log.Printf("Building Album '%s'", a.Title)

	if !a.NoScaledPreviews || !a.NoScaledThumbs {
		a.buildThumbsAndPreviews()
	}

	indexHTML := renderIndexTemplate(a)
	indexFile := a.RootPath + string(filepath.Separator) + "index.html"
	log.Printf("Rendering index.html into %s", indexFile)
	err := ioutil.WriteFile(indexFile, indexHTML, 0644)
	if err != nil {
		log.Printf("write %s ERROR: %s", indexFile, err)
	}

	albumFile := a.RootPath + string(filepath.Separator) + "albumdata.js"
	writeAlbumDataJS(a, err, albumFile)
	// add zipped version?

	assetsPath := filepath.FromSlash(a.RootPath + "/assets")
	log.Printf("Copying assets into %s", assetsPath)
	a.copyAssets(assetsPath)
	log.Printf("Nice! All done. Now open %s", indexFile)
}

func writeAlbumDataJS(a Album, err error, outFile string) error {
	albumData, _ := json.Marshal(a)
	log.Printf("Rendering albumdata.js into %s", outFile)
	f, err := os.Create(outFile)
	if err != nil {
		log.Printf("write %s ERROR: %s", outFile, err)
		return err
	}
	// FIXME handleErr
	f.Write([]byte("albumData = "))
	f.Write(albumData)
	f.Write([]byte(";\n"))
	f.Close()
	return nil
}

func ScanDir(root string, a *Album) (result *Node, err error) {
	log.Printf("Reading images in %s ...", root)
	parents := make(map[string]*Node)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(info.Name(), ".") {
			// skip dotfiles. log?
			return nil
		}
		if strings.HasPrefix(path, filepath.FromSlash(a.getPathOfThumbnails())) ||
			strings.HasPrefix(path, filepath.FromSlash(a.RootPath+"/assets")) ||
			strings.HasPrefix(path, filepath.FromSlash(a.getPathOfPreviews())) {
			// skip folder created by us
			return nil
		}
		parents[path] = &Node{
			FullPath: strings.TrimPrefix(path, root),
			WebPath:  filepath.ToSlash(strings.TrimPrefix(path, root)),
			Name:     info.Name(),
			IsDir:    info.IsDir(),
			Size:     info.Size(),
			Album:    a,
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
	result.WebPath = "/"
	return
}

func (n *Node) getAllImages() []Node {
	var toVisit []*Node
	var images []Node
	toVisit = append(toVisit, n)
	for len(toVisit) > 0 {
		c := toVisit[0]
		if c.IsImage {
			images = append(images, *c)
		} else {
			toVisit = append(toVisit, c.Children...)
		}
		toVisit = toVisit[1:]
	}
	return images
}

func (a *Album) getPathOriginals() string {
	return a.RootPath
}

func (a *Album) getPathOfThumbnails() string {
	return a.RootPath + "/thumbs"
}

func (a *Album) getPathOfPreviews() string {
	return a.RootPath + "/preview"
}

func (n *Node) getPathOfOriginal() string {
	return n.Album.RootPath + n.FullPath
}

func (n *Node) getPathOfThumbnail() string {
	return n.Album.RootPath + "/thumbs" + n.FullPath
}

func (n *Node) getPathOfPreview() string {
	return n.Album.RootPath + "/preview" + n.FullPath
}

func (a *Album) buildThumbsAndPreviews() {
	log.Printf("Building previews and thumbnails using %d threads ...", a.NumThreads)
	// FIXME limits max amounts of images...
	// removing 8192 will ...
	// https://stackoverflow.com/questions/26927479/go-language-fatal-error-all-goroutines-are-asleep-deadlock
	jobs := make(chan Node, 8192)
	results := make(chan int, 8192)
	for w := 0; w <= a.NumThreads-1; w++ {
		go a.imageScalingWorker(w, jobs, results)
	}
	imageNodes := a.Data.getAllImages()
	for _, path := range imageNodes {
		jobs <- path
	}
	close(jobs)
	for range imageNodes {
		<-results
	}
	log.Print("Building previews and thumbnails: done.")
}

func (a *Album) imageScalingWorker(id int, jobs <-chan Node, results chan<- int) {
	for imageNode := range jobs {
		originalPath := imageNode.getPathOfOriginal()
		previewPath := imageNode.getPathOfPreview()
		thumbPath := imageNode.getPathOfThumbnail()
		buildThumb := false
		buildPreview := false
		var todo []string
		if _, err := os.Stat(thumbPath); os.IsNotExist(err) && !a.NoScaledThumbs {
			todo = append(todo, "thumb")
			buildThumb = true
		}
		if _, err := os.Stat(previewPath); os.IsNotExist(err) && !a.NoScaledPreviews {
			todo = append(todo, "preview")
			buildPreview = true
		}
		if len(todo) == 0 {
			results <- 1
			continue
		}
		if len(todo) > 0 {
			bufP := &bytes.Buffer{}
			bufT := &bytes.Buffer{}
			file, err := os.Open(originalPath)
			var originalImage image.Image
			var previewImage image.Image
			if err == nil {
				originalImage, err = jpeg.Decode(file)
				if err == nil {
					err = file.Close()
					if err != nil {
						log.Printf("Decoding %s failed: %s", originalPath, err)
						results <- 1
						continue
					}
					if buildPreview {
						previewImage = a.buildView(imageNode, 700, originalImage, nil, bufP, imageNode.getPathOfPreview())
					}
					if buildThumb {
						a.buildView(imageNode, 105, originalImage, previewImage, bufT, imageNode.getPathOfThumbnail())
					}
				}
			}
			log.Printf("[thread-%d] %s created for %s", id, strings.Join(todo, "+"), originalPath)
		}
		results <- 1
	}
}

func (a *Album) buildView(albumImage Node, height uint, original image.Image, preview image.Image, output *bytes.Buffer, outoutPath string) image.Image {
	var m image.Image

	//log.Printf("buildView h %d %s", height, outoutPath)
	if preview == nil {
		m = resize.Resize(0, height, original, resize.Bicubic)
	} else {
		m = resize.Resize(0, height, preview, resize.Bicubic)
	}

	err := jpeg.Encode(output, m, nil /* FIXME add quality config option */)
	if err != nil {
		log.Printf("Resizing %s error: %s", albumImage.getPathOfOriginal(), err)
	} else {
		err = a.addCache(outoutPath, output.Bytes())
		if err != nil {
			log.Printf("addCache %s error: %s", outoutPath, err)
		}
	}
	// FIXME return err here!

	return m
}

func (a *Album) addCache(file string, data []byte) (err error) {
	if a.NoCacheScaled {
		return nil
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		//log.Printf("Add to cache: %s", file)
		err = os.MkdirAll(filepath.Dir(file), os.ModePerm)
		if err != nil {
			log.Printf("mkdir %s error: %s", filepath.Dir(file), err)
			return err
		}
		err = ioutil.WriteFile(file, data, 0644)
		if err != nil {
			log.Printf("write %s error: %s", file, err)
		}
	}
	return
}

func (a *Album) copyAssets(targetDir string) {
	assets := []string{"albutim.css", "albutim.js", "folder-up.svg", "jquery-2.2.2.min.js"}
	for _, asset := range assets {
		log.Printf("  copying %s -> %s", asset, targetDir)
		fileData := _escFSMustByte(false, "/"+asset)
		filename := targetDir + string(filepath.Separator) + asset
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(filename), os.ModePerm)
			if err != nil {
				log.Printf("mkdir %s error: %s", filepath.Dir(filename), err)
				return
			}
			err = ioutil.WriteFile(filename, fileData, 0644)
			if err != nil {
				log.Printf("write %s error: %s", filename, err)
			}
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
