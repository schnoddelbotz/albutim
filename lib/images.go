package lib

import (
	"bytes"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"github.com/xor-gate/goexif2/exif"
)

func isImage(path string) bool {
	extension := strings.ToLower(filepath.Ext(path))
	if extension == ".jpg" || extension == ".jpeg" {
		return true
	}
	return false
}

func getScaled(path string, width uint, height uint) (image []byte, err error) {
	buf := &bytes.Buffer{}
	file, err := os.Open(path)
	if err == nil {
		img, err := jpeg.Decode(file)
		if err == nil {
			err = file.Close()
			if err != nil {
				return nil, err
			}
			m := resize.Resize(width, height, img, resize.Bicubic)
			err := jpeg.Encode(buf, m, nil /* FIXME add quality config option */)
			if err != nil {
				return nil, err
			}
			image = buf.Bytes()
			return image, nil
		}
	}
	return
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
	log.Printf("Getting EXIF data: %s", fname)
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
