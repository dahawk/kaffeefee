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
			{{range $k, $v := .Total}}
			<h3><a href="graph?user={{$k}}">{{$k}}</a></h3>
			<div class="row">
				<div class="col-xs-1"><img src="{{ index $.Images $k}}"></div>
				<div class="col-xs-11">
					<table class="table table-condensed">
						<thead>
							<tr>
								<td>Heute<br>Tagesschnitt</td>
								<td>Diese Woche<br>Wochenschnitt</td>
								<td>Diesen Monat<br>Monatsschnitt</td>
								<td>Summe</td>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td>{{index $.Daily $k}}</td>
								<td>{{index $.Weekly $k}}</td>
								<td>{{index $.Monthly $k}}</td>
								<td>{{$v}}</td>
							</tr>
							<tr>
								<td>{{index $.DailyAvgs $k}}</td>
								<td>{{index $.WeeklyAvgs $k}}</td>
								<td>{{index $.MonthlyAvgs $k}}</td>
								<td>-</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
			{{end}}
		</div>
	</body>
</html>
