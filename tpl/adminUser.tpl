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
      <form method="POST" action="{{if .UserID}}{{else}}/createUser{{end}}">
        <input type="hidden" name="id" value="{{ .UserID }}"/>
        <div class="form-group row">
          <label for="user" class="col-md-2">Username</label>
          <input type="text" name="user" id="user" value="{{ .Name }}" class="col-md-4"/>
        </div>
        <div class="form-group row">
          <label for="email" class="col-md-2">Mail</label>
          <input type="email" name="email" id="email" value="{{ .Mail }}" class="col-md-4"/>
        </div>
         <button type="submit" class="btn btn-default">Submit</button>
      </form>
    </div>
  </body>
</html>
