<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8"/>
<title>Test</title>
<link href="https://cdn.bootcss.com/element-ui/2.0.0-alpha.2/theme-chalk/index.css" rel="stylesheet">
<link href="/static/css/style.css" rel="stylesheet">
<script src="/static/js/vue.js"></script>
<script src="/static/js/vue-resource@1.3.4.js"></script>
<script src="https://cdn.bootcss.com/element-ui/2.0.0-alpha.2/index.js"></script>
<script>
	Vue.http.options.emulateJSON = true;
	Vue.http.options.headers = {
	  'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
	};
</script>
</head>
<body>
<div id="app">
<header>
<div class="container" style="margin-bottom: 5px">
	<ul class="el-menu el-menu-demo el-menu--horizontal el-menu--dark">
		<li class="el-menu-item is-active"><a href="/">home</a></li>
		<li class="el-menu-item"><a href="/login">login</a></li>
		<li class="el-menu-item"><a href="/reg">register</a></li>
	</ul>
</div>
</header>
