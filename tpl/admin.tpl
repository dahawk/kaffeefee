<html>
	<head>
		<title>Kaffeefee</title>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" href="http://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
		<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js"></script>
		<script src="http://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/js/bootstrap.min.js"></script>
	</head>
	<body>
		<nav class="navbar navbar-default">
			<div class="container-fluid">
				<div class="navbar-header">
					<a class="navbar-brand" href="/">Kaffeefee</a>
				</div>
      </div>
    </nav>
    <div class="container">
      {{range $u := .}}
      <div class="row">
        <div class="col-xs-1">
          <img src="/static/{{$u.Name}}.png">
        </div>
        <div class="col-xs-2">
          {{$u.Name}}
        </div>
        <div class="col-xs-9">
          <div class="btn-group">
            <a href="?user={{$u.Name}}&enable=1" class="btn btn-default">{{if $u.Active}}Disable{{else}}Enable{{end}}</a>
            <a href="/editUser?user={{$u.Name}}&edit=1" class="btn btn-default">Edit</a>
            <a href="?user={{$u.Name}}&delete=1" class="btn btn-danger">Delete</a>
          </div>
        </div>
      </div>
      {{end}}
			<div class="row">
				<a href="/createUser" class="btn btn-default">Add user</a>
			</div>
    </div>
  </body>
</html>
