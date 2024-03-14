# TCell Blit 


This Go package draws an image into a tcell screen. 
See the Example for a full use including handling resize. 

Here's a demo, but Asciinema doesn't show resizes right:
![Demo](show.gif)
To run the demo:
```
  cd example
  go run main.go
```

Function: Draw(Screen, Image, Fill)
   Fill: false: see whole image. true: no black borders. 