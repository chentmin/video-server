package main

var (
	audioTemplateString = `
<!DOCTYPE html>
	<html>
	<body>
	{{.BaseName |html}}

	<br />
	<br />
<audio  controls="controls" preload="auto" autoplay="autoplay" loop="loop"> 

<source src="/files/{{.Name}}" />
	Your browser does not support the audio element.
	</audio>
	<br />
	<br />
	`
)
