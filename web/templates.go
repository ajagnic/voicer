package web

type page struct {
	Title  string
	Body   string
	Script string
}

func indexPage() *page {
	return &page{
		Title: "Home",
		Body: `
		<audio controls>
			<source src="media/output.mp3" type="audio/mpeg">
			<source src="media/output.ogg" type="audio/ogg">
			<source src="media/output.wav" type="audio/wav">
		</audio>
		<form action="/post" method="post">
			<textarea name="input" id="input" cols="30" rows="10"></textarea>
			<button type="submit">Convert</button>
		</form>
		`,
		Script: ``,
	}
}

func htmlLayout() string {
	return `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{.Body}}
		{{.Script}}
	</body>
	</html>
	`
}
