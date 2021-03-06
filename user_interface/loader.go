package user_interface

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"io/ioutil"
	"os"
	"strings"
)

func GetFileName(name string) string {
	dot := len(name) - 1
	for dot >= 0 && name[dot] != '.' {
		dot--
	}
	slh := dot - 1
	for slh >= 0 && name[slh] != '/' {
		slh--
	}
	return name[slh+1 : dot]
}

func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(err)
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer CloseFile(file)

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

//func Sprite(path string) {
//	pic, err := loadPicture(path)
//	if err != nil {
//		panic(err)
//	}
//	sprite := pixel.NewSprite(pic, pic.Bounds())
//	name := getFileName(path)
//	if _, ok := Sprites[name]; ok {
//		name = path
//	}
//	Sprites[name] = sprite
//}

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer CloseFile(file)

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font_, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font_, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func LoadFont(path string, size float64) {
	face, err := loadTTF(path, size)
	if err != nil {
		panic(err)
	}

	name := GetFileName(path)
	if _, ok := Fonts[name]; !ok {
		Fonts[name] = map[float64]*text.Atlas{}
	}
	Fonts[name][size] = text.NewAtlas(face, text.ASCII)
}

func LoadAllSprites(directoryPath string) {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		filePath := directoryPath + "/" + file.Name()
		if file.IsDir() {
			LoadAllSprites(filePath)
		} else if strings.HasSuffix(filePath, ".png") {
			LoadComposeSprite(filePath)
		}
	}
}

func LoadAllFonts(directoryPath string, size float64) {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		filePath := directoryPath + "/" + file.Name()
		if file.IsDir() {
			LoadAllFonts(filePath, size)
		} else {
			LoadFont(filePath, size)
		}
	}
}

func LoadComposeSprite(path string) {
	pic, err := loadPicture(path)
	if err != nil {
		panic(err)
	}

	dot := len(path) - 1
	for dot >= 0 && path[dot] != '.' {
		dot--
	}
	infoPath := path[:dot] + ".txt"
	infoStr, err := ioutil.ReadFile(infoPath)
	if err != nil { // sprite is not composed
		sprite := pixel.NewSprite(pic, pic.Bounds())
		name := GetFileName(path)
		Sprites[name] = sprite
		return
	}

	infoList := strings.Split(strings.ReplaceAll(string(infoStr), "\r", ""), "\n")
	bounds := [4]float64{}
	for _, info := range infoList {
		nameList := strings.Split(info, ": ")
		List := strings.Split(nameList[1], ", ")
		for j := 0; j < 4; j++ {
			bounds[j] = float64(StrToInt(List[j]))
		}
		sprite := pixel.NewSprite(pic, pixel.R(bounds[0], bounds[1], bounds[2], bounds[3]))
		name := GetFileName(path)
		SpriteTypes[name] = append(SpriteTypes[name], nameList[0])
		Sprites[name+"_"+nameList[0]] = sprite
	}
}
