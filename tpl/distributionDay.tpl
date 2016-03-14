<html>
	<head>
		<title>Kaffeefee</title>
		<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/morris.js/0.5.1/morris.css">
		<script src="//ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/raphael/2.1.0/raphael-min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/morris.js/0.5.1/morris.js"></script>
		<script src="http://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/js/bootstrap.min.js"></script>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" href="http://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
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
			<h1>Day Distribution</h1>
			<div class="row" id="dailychart" style="height: 300px;"></div>
    </div>
		<script>
			$(function() {
				var chart = Morris.Bar({
					element: "dailychart",
					xkey: "Hour",
					ykeys: ["Cnt"],
					labels: ["Kaffees"]
				});

				$.ajax({
					type: "GET",
					dataType: 'json',
					url: "/jsonDaily"
				})
				.done(function(data){
					if (data.length > 0) {
					chart.setData(data)
					}
				})
				.fail(function(){
					alert('problem')
				});
			});
		</script>
  </body>
</html>
