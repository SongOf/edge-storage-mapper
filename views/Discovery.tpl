<div class="container">
	{{if .Cameras}}
    	<div class="row">
    		<div class="panel widget">
    			<div class="panel-body">
    				<ul class="list-group">
                        {{range $camera := .Cameras}}
    					<li class="list-group-item">
    						<div class="row">
    							<div class="col-xs-2 col-md-1">
    								<img src="http://placehold.it/80" class="img-circle img-responsive" alt="" />
    							</div>
    							<div class="col-xs-10 col-md-11">
    							    <div>绑定状态：{{$camera.BindInfo}}</div>
    							    <div class="text-info mt-3">设备IP: {{ $camera.Ip }}</div>
    							    <div class="text-info mt-3">流协议: {{ $camera.Protocol }}</div>
    							    <div class="text-info mt-3">流地址: {{ $camera.Url }}</div>
    							    <div class="text-info mt-3">设备在线状态: {{ $camera.State }}</div>
    								<div class="action mt-3">
    									<a href="device/edit/{{$camera.Id}}" class="btn btn-primary"><span class="glyphicon glyphicon-edit" aria-hidden="true"></span><span><strong>Edit</strong></span>
    									</a>
    									<a href="device/view/{{$camera.Id}}"  class="btn btn-success"><span class="glyphicon glyphicon-eye-open" aria-hidden="true"></span><span> <strong>View</strong></span>
    									</a>
    									<a href="javascript:void(0);" onclick="deleteCamera({{$camera.Id}})"class="btn btn-danger"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span><span> <strong>Delete</strong></span>
    									</a>
    								</div>
    							</div>
    						</div>
    					</li>
                        {{end}}
    				</ul>
    			</div>
    		</div>
    	</div>
    {{end}}
</div>