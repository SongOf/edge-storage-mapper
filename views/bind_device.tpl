<div class="container">
    {{if .Errors}}
	    <div class="alert alert-danger">
            {{range $value := .Errors}}
		        <li> {{$value}} </li>
            {{end}}
	    </div>
    {{end}}
	<form method="post" action="/device/bind/{{.Ip}}">
		<div class="form-group">
		    <label for="serialNumber">设备序列号</label>
        	<input type="text" id="serialNumber" name="serialNumber" class="form-control" placeholder="请输入设备序列号" aria-describedby="basic-addon1">
        </div>
        <div class="form-group">
            <label for="validateCode">设备验证码</label>
            <input type="text" id="validateCode" name="validateCode" class="form-control" placeholder="请输入设备验证码" aria-describedby="basic-addon1">
        </div>
        <div class="alert alert-success" role="alert">设备序列号信息和设备验证码信息在设备底部标签可见</div>
		<input type="submit" class="btn btn-primary" value="Bind"/>
		    <a href="/discovery" class="btn btn-danger"><span aria-hidden="true"></span><span><strong>Back</strong></span></a>
	</form>
</div>