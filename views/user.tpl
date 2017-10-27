<% template "header.tpl" .%>

<el-row :gutter="5">
<el-col :span="3"><% template "left.tpl" .%></el-col>
<el-col :span="21">
	<el-button @click="formVisible = true" style="margin:5px" type="primary">add</el-button>
	<div class="block" style="margin:5px">
	    <span class="demonstration">username:</span>
	    <el-input v-model="page.search_username" placeholder="username" style="width:180px"></el-input>
	    <el-button type="primary" @click="search">search</el-button>
	</div>
	<el-table :data="tableData" border style="width: 100%">
		<el-table-column prop="Id" label="id" width="60"></el-table-column>
		<el-table-column prop="UserName" label="userName" width="120"></el-table-column>
		<el-table-column prop="Mobile" label="mobile" width="110"></el-table-column>
		<el-table-column prop="Gender" label="gender" width="250"></el-table-column>
		<el-table-column prop="RealName" label="realName" width="250"></el-table-column>
		<el-table-column prop="IdentityCard" label="identityCard" width="250"></el-table-column>
		<el-table-column prop="Avatar" label="avatar" width="150">
			<template scope="scope"><a v-bind:href="scope.row.avatar">
				<img v-bind:src="scope.row.Avatar" style="width:150px;height:80px"/>
			</a></template>
		</el-table-column>
		<el-table-column label="operate" width="300">
	      <template scope="scope">
			<el-row>
		        <el-button @click="handleEdit(scope.$index, scope.row)"  size="small" type="text">edit</el-button>
		        <el-button @click="handleDelete(scope.$index, scope.row)" size="small" type="text">delete</el-button>
			</el-row>
			<el-row>
				<el-button @click="handleDepart(scope.$index, scope.row)" size="small" type="text">allotGroup</el-button>
				<el-button @click="handleRole(scope.$index, scope.row)" size="small" type="text">allotRole</el-button>
			</el-row>
			<el-button @click="handleRes(scope.$index, scope.row)" size="small" type="text">allotRes</el-button>
	      </template>
	    </el-table-column>
	</el-table>
	<el-pagination style="float:right"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="page.currentPage"
      :page-sizes="[10, 15, 50, 100]"
      :page-size="page.pageSize"
      layout="total, sizes, prev, pager, next, jumper"
      :total="page.totalSize">
    </el-pagination>
	<el-dialog title="dialog" :visible.sync="formVisible">
	  <el-form :model="form">
		<el-form-item label="userName" :label-width="formLabelWidth">
	      <el-input v-model="form.UserName" auto-complete="off"></el-input>
	    </el-form-item>
	    <el-form-item label="mobile" :label-width="formLabelWidth">
	      <el-input v-model="form.Mobile" auto-complete="off"></el-input>
	    </el-form-item>
	    <el-form-item label="gender" :label-width="formLabelWidth">
	      <el-input v-model="form.Gender" auto-complete="off"></el-input>
	    </el-form-item>
		<el-form-item label="realName" :label-width="formLabelWidth">
	      <el-input v-model="form.RealName" auto-complete="off"></el-input>
	    </el-form-item>
		<el-form-item label="identityCard" :label-width="formLabelWidth">
	      <el-input v-model="form.IdentityCard" auto-complete="off"></el-input>
	    </el-form-item>
	    
	    <el-form-item label="image" :label-width="formLabelWidth">
	    	<el-upload
			  action="/upload"
			  list-type="picture-card"
			  ref="upload"
			  :file-list="form.fileList_image"
			  :multiple=false
			  :on-success="handleUploadSuccess">
			  <i class="el-icon-plus"></i>
			</el-upload>
			<el-dialog v-model="dialogVisible" size="tiny">
			  <img width="100%" :src="dialogImageUrl" alt="">
			</el-dialog>
	    </el-form-item>
	    
	  </el-form>
	  <div slot="footer" class="dialog-footer">
	    <el-button @click="formVisible = false">取 消</el-button>
	    <el-button type="primary" @click="handleOk">确 定</el-button>
	  </div>
	</el-dialog>
	<el-dialog title="depart" :visible.sync="departDialogVisible" 
	 @open="handleDepartBeforeOpen">
		<el-tree
		  :data="depart"
		  show-checkbox
		  node-key="Id"
		  ref="departTree"
		  highlight-current
		  check-strictly 
		  default-expand-all
		  :props="defaultProps">
		</el-tree>
		<div slot="footer" class="dialog-footer">
		  <el-button @click="departDialogVisible = false">取 消</el-button>
		  <el-button type="primary" @click="handleDepartOK">确 定</el-button>
		</div>
	</el-dialog>
	<el-dialog title="role" :visible.sync="roleDialogVisible" 
	 @open="handleRoleBeforeOpen">
		<el-tree
		  :data="role"
		  show-checkbox
		  node-key="Id"
		  ref="roleTree"
		  highlight-current
		  check-strictly 
		  default-expand-all
		  :props="defaultProps">
		</el-tree>
		<div slot="footer" class="dialog-footer">
		  <el-button @click="roleDialogVisible = false">取 消</el-button>
		  <el-button type="primary" @click="handleRoleOK">确 定</el-button>
		</div>
	</el-dialog>
	
	<el-dialog title="resource" :visible.sync="resDialogVisible" 
	 @open="handleResBeforeOpen">
		<el-tree
		  :data="res"
		  show-checkbox
		  node-key="Id"
		  ref="resTree"
		  highlight-current
		  :props="defaultProps">
		</el-tree>
		<div slot="footer" class="dialog-footer">
		  <el-button @click="resDialogVisible = false">取 消</el-button>
		  <el-button type="primary" @click="handleResOK">确 定</el-button>
		</div>
	</el-dialog>
</el-col>
</el-row>
</div>
<script type="text/javascript">
var app = new Vue({
	el:"#app",
	data:{
		tableData:<%.tableData%>,
		depart:<%.depart%>,
		res:<%.res%>,
		role:<%.role%>,
		defaultProps: {
          children: 'Children',
          label: 'Label',
		  disabled:'Disabled'
        },
		form:{
			Id:0,
			UserName:'',
			Birth:'',
			Gender:0,
			Addr:'',
			Mobile:'',
			Email:'',
			RealName:'',
			Status:0,
			IdentityCard:'',
			index:0,
			depart:[],
			role:[],
			res:[]
		},
		page:{
			currentPage:1,
			pageSize:10,
			totalSize:<%.totalSize%>,
			search_username:''
		},
		formVisible:false,
		dialogVisible:false,
		dialogImageUrl:'',
		formLabelWidth:'120px',
		departDialogVisible:false,
		roleDialogVisible:false,
		resDialogVisible:false,
	},
	methods:{
		handleDepart(index,row){
			this.form.Id = row.Id;
			Vue.http.post('/rpc/getdepartbyuser',this.form).then(response=>{
				if(response.body.Code==1){
					this.form.depart = [];
					if(response.body.Data!=null)
						for(let i=0;i<response.body.Data.length;i++){
							this.form.depart[i] = response.body.Data[i].Id;
						}
					this.departDialogVisible = true;
	  			}else{
	  				this.$message.error('get depart failure');
	  			}
			});
		},
		handleDepartOK:function(){
			let tmpSelected = this.$refs.departTree.getCheckedNodes(false);
			this.form.depart = [];
			for(let i=0;i<tmpSelected.length;i++){
				this.form.depart[i] = tmpSelected[i].Id;
			}
			Vue.http.post('allotdepart',this.form).then(response=>{
				if(response.body.Code==1){
					this.$message({
			          message: 'allot depart success',
			          type: 'success'
			        });
					this.departDialogVisible = false;
				}else{
					this.$message.error('allot depart failure');
				}
			});
			
		},
		handleRole(index,row){
			this.form.Id = row.Id;
			Vue.http.post('/rpc/getrolebyuser',this.form).then(response=>{
				if(response.body.Code==1){
					this.form.role = [];
					if(response.body.Data!=null)
						for(let i=0;i<response.body.Data.length;i++){
							this.form.role[i] = response.body.Data[i].Id;
						}
					this.roleDialogVisible = true;
	  			}else{
	  				this.$message.error('get role failure');
	  			}
			});
		},
		handleRoleOK:function(){
			let tmpSelected = this.$refs.roleTree.getCheckedNodes(false);
			this.form.role = [];
			for(let i=0;i<tmpSelected.length;i++){
				this.form.role[i] = tmpSelected[i].Id;
			}
			Vue.http.post('allotrole',this.form).then(response=>{
				if(response.body.Code==1){
					this.$message({
			          message: 'allot role success',
			          type: 'success'
			        });
					this.roleDialogVisible = false;
				}else{
					this.$message.error('allot role failure');
				}
			});
			
		},
		handleRes:function(index,row){
			this.form.Id = row.Id;
			Vue.http.post('/rpc/getresbyuser',this.form).then(response=>{
				if(response.body.Code==1){
					this.form.res = [];
					if(response.body.Data!=null)
						for(let i=0;i<response.body.Data.length;i++){
							this.form.res[i] = response.body.Data[i].Id;
						}
					this.resDialogVisible = true;
	  			}else{
	  				this.$message.error('get res failure');
	  			}
			});
		},
		handleResOK:function(){
			let tmpSelected = this.$refs.resTree.getCheckedNodes(true);
			this.form.res = [];
			for(let i=0;i<tmpSelected.length;i++){
				this.form.res[i] = tmpSelected[i].Id;
			}
			Vue.http.post('allotres',this.form).then(response=>{
				if(response.body.Code==1){
					this.$message({
			          message: 'allot resource success',
			          type: 'success'
			        });
					this.resDialogVisible = false;
				}else{
					this.$message.error('allot resource failure');
				}
			});
		},
		handleDepartBeforeOpen:function(){
			setTimeout(() => {
		        app.$refs.departTree.setCheckedKeys(this.form.depart);
		    }, 100);
		},
		handleRoleBeforeOpen:function(){
			setTimeout(() => {
				console.log(this.form.role);
		        app.$refs.roleTree.setCheckedKeys(this.form.role);
		    }, 100);
		},
		handleResBeforeOpen:function(){
			setTimeout(() => {
				console.log(this.form.res);
		        app.$refs.resTree.setCheckedKeys(this.form.res);
		    }, 100);
		},
		handleEdit:function(index,row){
	  		this.formVisible = true;
	  		this.form.Id = row.Id;
	  		this.form.UserName = row.UserName;
	  		this.form.Mobile = row.Mobile;
	  		this.form.Gender = row.Gender;
			this.form.RealName = row.RealName;
			this.form.IdentityCard = row.IdentityCard;
			this.form.index = index;
	  		//this.form.upload_url = row.image_path;
	  		//this.form.fileList_image = [{name:'name',url:row.image_path}];
	  	},
	  	handleOk:function(){
	  		if(this.form.Id==0) this.add();
	  		else this.edit();
	  	},
	  	handleDelete:function(index,row){
			console.log(index);
	  		this.$confirm('this operate will delete forever,ready to delete?','confirm',{
	  			confirmButtonText: 'ok',
		        cancelButtonText: 'cancel',
		        type: 'warning'
	  		}).then(()=>{
	  			Vue.http.post('delete',{id:row.Id}).then(response=>{
		  			if(response.body.Code==1){
		  				this.$message({
				          message: 'delete success',
				          type: 'success'
				        });
						this.tableData.splice(index,1);
		  			}else{
		  				this.$message.error('delete failure');
		  			}
		  		});
	  		});
	  		
	  	},
	  	add(){
	  		Vue.http.post('add',this.form).then(response=>{
	  			console.log(response);
	  			if(response.body.Code==1){
	  				this.$message({
			          message: 'add success',
			          type: 'success'
			        });
			        this.formVisible = false;
					this.form.Id = response.body.Data; 
         			var obj=JSON.parse(JSON.stringify(this.templateData));//deep copy
					this.tableData.splice(0,0,obj);
			        this.clearForm();
					
	  			}else{
	  				this.$message.error('add failure');
	  			}
	  		});
	  	},
	  	edit(){
	  		Vue.http.post('edit',this.form).then(response=>{
	  			console.log(response);
	  			if(response.body.Code==1){
	  				this.$message({
			          message: 'edit success',
			          type: 'success'
			        });
			        this.formVisible = false;
         			var obj=JSON.parse(JSON.stringify(this.form));//deep copy
					this.tableData.splice(this.form.index,1,obj);
	  			}else{
	  				this.$message.error('edit failure');
	  			}
	  		});
	  	},
		handleSizeChange(val) {
  			this.page.pageSize = val;
	  			this.$http.post('',this.page).then(response=>{
	    			this.tableData = response.body;
	    		});
	    },
	    handleCurrentChange(val) {
		    	this.page.currentPage = val;
		    	this.$http.post('',this.page).then(response=>{
	    			this.tableData = response.body;
		    	});
	    },
	    search(){
		    	this.page.currentPage = 1;
		    	this.page.pageSize = 10;
		    	this.$http.post('',this.page).then(response=>{
	    			this.tableData = response.body;
		    	});
	    },
		handleUploadSuccess(){
			
		}
	}
})
</script>
<% template "footer.tpl" .%>