package main

var (
	indexTemplateString = `
<!DOCTYPE html>
<html>
<body>
<h1> Videos </h1>
<br />
<br />

{{range .Videos}}

<a href="/watch/{{.Name}}"> {{.BaseName}} </a>
<br />
<br />
{{end}}

<h1> Musics </h1>


{{range .Audios}}

<a href="/listen/{{.Name}}"> {{.BaseName}} </a>
<br />
<br />
{{end}}

</body>

</html>
	`
)
