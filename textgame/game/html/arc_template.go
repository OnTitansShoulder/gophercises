package html

// bind to source for now
// interesting lib to see https://github.com/go-bindata/go-bindata
const ArcTemplateStr = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>TextGame</title>
	</head>
	<body>
        <h2>{{.Title}}</h2>
        <br/>
        {{range .Stories}}
            <p>{{- . -}}</p><br/>
        {{end}}
        {{if gt (len .Options) 0}}
            <p>What would you choose?</p>
        {{else}}
            <p><a href="/intro">Restart Game</a></p>
        {{end}}
        <br/>
        <ol type="1">
            {{range .Options}}
                <li><a href="/{{.ArcName}}">{{ .Text }}</a></li>
            {{end}}
        </ol>
	</body>
</html>`
