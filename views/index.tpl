<div class="container">
    {{if .ErrorMsg}}
        <div class="alert alert-danger">
                {{.ErrorMsg}}
        </div>
    {{end}}
	<form method="post">
		<div class="form-group">
			<label for="username">User Name</label>
			<div>
				<input type="text" class="form-control" name="username" id="username" maxlength="48" placeholder="UserName" />
			</div>
		</div>
		<div class="form-group">
			<label for="password">Password</label>
			<div>
				<input type="password" class="form-control" name="password" id="password" maxlength="48" placeholder="Password" />
			</div>
		</div>
		<input type="submit" class="btn btn-primary" value="Login" />
	</form>
</div>