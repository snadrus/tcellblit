package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/gdamore/tcell/v2"
)

func runLoop(s tcell.Screen, image string) {
	img, err := load(image)
	if err != nil {
		panic(err)
	}

	//draw(s, img)

	redrawTimer := time.NewTimer(1 * time.Millisecond)
	go func() {
		for range redrawTimer.C {
			blit.Draw(s, img)
		}
	}()

	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			key, ch := ev.Key(), ev.Rune()
			if key == tcell.KeyEscape || ch == 'q' {
				return
			}
		case *tcell.EventResize:
			redrawTimer.Reset(time.Millisecond * 100)
		default:
			time.Sleep(10 * time.Millisecond)
		}

	}
}

func main() {
	imgName := "ai-hill-2.png"
	if len(os.Args) == 2 {
		imgName = os.Args[1]
	} else {
		fmt.Printf("Usage: %s <filename>...\n\n", os.Args[0])
		fmt.Println("Close the image with <ESC> or by pressing 'q'.")
		os.Exit(1)
	}

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default text style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	runLoop(s, imgName)
	s.Fini()
}

// load an image stored in the given path
func load(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}
