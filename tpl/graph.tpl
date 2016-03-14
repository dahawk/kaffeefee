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
			<h1>{{.UserId}}'s Graphs</h1>
			<div class="row" id="dailychart" style="height: 200px;"></div>
			<div class="row" id="weeklychart" style="height: 200px;"></div>
			<div class="row" id="monthlychart" style="height: 200px;"></div>
		</div>
		<script>
			$(function() {
				var chartDaily = Morris.Line({
					element: 'dailychart',
					data: [{datum: '0000-00-00 00:00:00', count: '0'}],
					xkey: 'datum',
					ykeys: ['count'],
					labels: ['Tagesverlauf'],
					xlabels: 'day'
				});
				var chartWeekly = Morris.Line({
					element: 'weeklychart',
					data: [{datum: '0000-00-00 00:00:00', count: '0'}],
					xkey: 'datum',
					ykeys: ['count'],
					labels: ['Wochenverlauf'],
					xlabels: 'day'
				});
				var chartMonthly = Morris.Line({
					element: 'monthlychart',
					data: [{datum: '0000-00-00 00:00:00', count: '0'}],
					xkey: 'datum',
					ykeys: ['count'],
					labels: ['Monatsverlauf'],
					xlabels: 'day'
				});

				$.ajax({
					type: "GET",
					dataType: 'json',
					url: "/json",
					data: { interval: 'daily', user: {{.UserId}}}
				})
				.done(function(data){
					if (data.length > 0) {
					chartDaily.setData(data)
					}
				})
				.fail(function(){
					alert('problem')
				});

				$.ajax({
					type: "GET",
					dataType: 'json',
					url: "/json",
					data: { interval: 'monthly', user: {{.UserId}}}
				})
				.done(function(data){
					if (data.length > 0) {
					chartMonthly.setData(data)
					}
				})
				.fail(function(){
					alert('problem')
				});

				$.ajax({
					type: "GET",
					dataType: 'json',
					url: "/json",
					data: { interval: 'weekly', user: {{.UserId}}}
				})
				.done(function(data){
					if (data.length > 0) {
					chartWeekly.setData(data)
					}
				})
				.fail(function(){
					alert('problem')
				});
			});
		</script>
	</body>
</html>
