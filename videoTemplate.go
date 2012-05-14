package main

var (
	videoTemplateString = `
	<!DOCTYPE html>
	<html>
	<body>
	{{.BaseName |html}}
	<br />
	<br />
	<video id="myvideo" controls="controls" preload="auto" autoplay="autoplay"> 
	<source src="/files/{{.Name}}" />
	Your browser does not support the video element.
	</video>
	<br />
	<br />
	press enter to fullscreen

	<script>

  var videoElement = document.getElementById("myvideo");
    
  function toggleFullScreen() {
    if (!document.mozFullScreen && !document.webkitFullScreen) {
      if (videoElement.mozRequestFullScreen) {
        videoElement.mozRequestFullScreen();
      } else {
        videoElement.webkitRequestFullScreen(Element.ALLOW_KEYBOARD_INPUT);
      }
    } else {
      if (document.mozCancelFullScreen) {
        document.mozCancelFullScreen();
      } else {
        document.webkitCancelFullScreen();
      }
    }
  }

  
  document.addEventListener("keydown", function(e) {
    if (e.keyCode == 13) {
      toggleFullScreen();
    }
  }, false);
	
	</script>
	</body>
	</html>
`
)
