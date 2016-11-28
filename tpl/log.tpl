<html>
	<head>
		<title>Kaffeefee</title>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<script src="https://code.jquery.com/jquery-1.12.0.min.js"></script>
		<link rel="stylesheet" href="http://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
		<script src="http://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/js/bootstrap.min.js"></script>
		<link href="/static/stylesheets/tally.css" media="all" rel="stylesheet" />
		<script src="/static/javascripts/tally.js"></script>

		<script>
			$(document).ready(function(){
			    $(".flex-item").click(function(){
			        var id = $(this).data("id");
					$.get(
						"/logcoffee?id="+id
					);
					var cnt = parseInt($("#id-"+id).text())
					cnt += 1
					$("#id-"+id).text(cnt)
			    });
			});
		</script>
		<style>
			.flex-container {
				display: -webkit-flex;
    			display: flex;
				flex-wrap: wrap;
			}

			.flex-item {
				background: lightgrey;
				width: 350px;
				height: 75px;
				margin: 15px;
				padding: 5px;
			}

			.list {
				font-size: 175%;
			}
		</style>
	</head>
	<body>
		<nav class="navbar navbar-default">
			<div class="container-fluid">
				<div class="navbar-header">
					<a class="navbar-brand" href="#">Kaffeefee</a>
				</div>
				<div>
					<ul class="nav navbar-nav">
						<li><a href="/">Kaffee eintragen</a></li>
						<li><a href="stats">Stats</a></li>
						<li><a href="dailyChart">Daily Breakdown</a></li>
						<li><a href="weeklyChart">Weekly Breakdown</a></li>
					</ul>
				</div>
			</div>
		</nav>
		<div class="container">
			<h1>{{.From}} - {{.To}}</h1>
			{{if .Error}}
			<div class="alert alert-danger">
				<strong>Danger!</strong> {{.Error}}
			</div>
			{{else if .Store}}
			<div class="alert alert-success">
				<strong>Success!</strong> Daten gespeichert.
			</div>
			{{end}}
			<div class="flex-container container">
				{{range .Users}}<div class="flex-item row" data-id="{{.UserID}}"><div class="col-md-3"><img src="{{.Image}}" class="userimg"/></div><div style="float: right;" class="col-md-9">{{.Name}}<br/><span id="id-{{.UserID}}" class="tally list">{{.Today}}</span></div></div>{{end}}
			</div>
		</div>
	</body>
</html>
