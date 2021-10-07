<div class="container">
	<div class="well">
		<div class="media">
			<div class="media-body">
				<h4 class="media-heading">
					<strong>设备序列号：{{.Camera.SerialNumber}}</strong>
				</h4>
				<p class="text-right">By
					<strong>{{.Camera.CreateBy}}</strong>
				</p>
				<p>设备验证码：{{.Camera.ValidateCode}}</p>
				<p>设备IP：{{.Camera.Ip}}</p>
				<p>设备流协议：{{.Camera.Protocol}}</p>
				<p>设备流地址：{{.Camera.Url}}</p>
				<p>设备绑定信息：{{.Camera.BindInfo}}</p>
				<p>设备状态：{{.Camera.State}}</p>
				<ul class="list-inline list-unstyled">
					<li>
						<span>
							<i class="glyphicon glyphicon-calendar"></i>
              {{.Camera.CreateTime.Format "02-01-2006 15:04:05"}}
						</span>
					</li>
					<li>
						<span>
							<i class="glyphicon glyphicon-time"></i>        {{.Camera.UpdateTime.Format "02-01-2006 15:04:05"}}
						</span>
					</li>
				</ul>
			</div>
		</div>
	</div>
	<a href="/device" class="btn btn-danger"><span aria-hidden="true"></span><span><strong>Back</strong></span></a>
</div>