<% template "header.tpl" .%>

<el-row :gutter="5">
<el-col :span="3"><% template "left.tpl" .%></el-col>
<el-col :span="21">
	<el-input
	  placeholder="输入关键字进行过滤"
	  v-model="filterText">
	</el-input>
	<el-tree
	  :data="tree"
	  :props="defaultProps"
	  node-key="Id"
	  default-expand-all
	  :filter-node-method="filterNode"
	  class="filter-tree"
	  ref="tree"
	  :render-content="renderContent">
	</el-tree>
	<el-dialog title="dialog" :visible.sync="dialogFormVisible">
	  <el-form :model="form">
	    <el-form-item label="depart name:" :label-width="formLabelWidth"
		:rules="{required:true}">
	      <el-input v-model="form.name" auto-complete="off"></el-input>
	    </el-form-item>
	    
	  </el-form>
	  <div slot="footer" class="dialog-footer">
	    <el-button @click="dialogFormVisible = false">取 消</el-button>
	    <el-button type="primary" @click="handleOK">确 定</el-button>
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
		tree:<%.tree%>,
		res:<%.resTree%>,
		role:<%.roleTree%>,
		defaultProps: {
          children: 'Children',
          label: 'Label',
		  disabled:'Disabled'
        },
		form:{
			id:0,
			name:'',
			pid:0,
			role:[],
			res:[]
		},
		selectedRole:[],
		selectedRes:[],
		formLabelWidth:'120px',
		dialogFormVisible:false,
		roleDialogVisible:false,
		resDialogVisible:false,
		currentData:'',
		currentStore:'',
		filterText: ''
	},
	methods:{
		add:function(){
			Vue.http.post('add',this.form).then(response=>{
	  			if(response.body.Code==1){
	  				this.$message({
			          message: 'add success',
			          type: 'success'
			        });
					this.currentData.Children = this.currentData.Children || [];
					let obj = {
						Id:response.body.Data,
						Label:this.form.name,
						Role:{
							Id:response.body.Data,
							Name:this.form.name,
							PId:this.currentData.Id
						}
					}
					console.log(obj);
					this.currentStore.append(obj,this.currentData);
					this.currentData.Children.push(obj);
					this.dialogFormVisible = false;
					this.clearForm();
	  			}else{
	  				this.$message.error('add failure');
	  			}
	  		});
		},
		edit:function(){
			Vue.http.post('edit',this.form).then(response=>{
	  			if(response.body.Code==1){
	  				this.$message({
			          message: 'edit success',
			          type: 'success'
			        });
					this.currentData.Label = this.form.name
					this.currentData.Depart.Name = this.form.name
			        this.dialogFormVisible = false;
	  			}else{
	  				this.$message.error('edit failure');
	  			}
	  		});
		},
		handleOK:function(){
			if(this.form.id==0) this.add();
			else this.edit();
		},
		handleEdit:function(store,node,data){
			this.form.id = data.Id;
			this.form.name = data.Depart.Name;
			this.form.pid = data.Pid;
			
			this.currentData = data;
			this.currentStore = store;
			
			this.dialogFormVisible = true;
			
		},
		handleAppend(store,data){
			var pid = data.Id;
			this.clearForm();
			this.form.pid = pid;
			this.dialogFormVisible = true;
			this.currentData = data;
			this.currentStore = store;
			
		},
		handleRemove(store,node,data){
			this.currentData = data;
			this.$confirm('此操作将永久删除, 是否继续?', '提示', {
	          confirmButtonText: '确定',
	          cancelButtonText: '取消',
	          type: 'warning'
	        }).then(() => {
				Vue.http.post('delete',{id:data.Id}).then(response=>{
		  			if(response.body.Code==1){
		  				this.$message({
				          message: 'delete success',
				          type: 'success'
				        });
						const parent = node.parent;
				        const index = parent.data.Children.findIndex(d => d.Id === data.Id);
				        parent.data.Children.splice(index, 1);
						store.remove(data);
		  			}else{
		  				this.$message.error('delete failure');
		  			}
		  		});
	        });  
		},
		handleRole:function(store,node,data){
			this.form.id = data.Id;
			Vue.http.post('/rpc/getrolebydepart',this.form).then(response=>{
				if(response.body.Code==1){
					this.selectedRole = [];
					this.form.role = [];
					if(response.body.Data!=null)
						for(let i=0;i<response.body.Data.length;i++){
							this.selectedRole[i] = {};
							this.selectedRole[i].value = response.body.Data[i].Id;
							this.selectedRole[i].label = response.body.Data[i].Name;
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
			this.selectedRole = [];
			this.form.role = [];
			for(let i=0;i<tmpSelected.length;i++){
				this.selectedRole[i] = {};
				this.selectedRole[i].value = tmpSelected[i].Id;
				this.selectedRole[i].label = tmpSelected[i].Label;
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
		handleRes:function(store,node,data){
			this.form.id = data.Id;
			Vue.http.post('/rpc/getresbydepart',this.form).then(response=>{
				if(response.body.Code==1){
					this.selectedRes = [];
					this.form.res = [];
					if(response.body.Data!=null)
						for(let i=0;i<response.body.Data.length;i++){
							this.selectedRes[i] = {};
							this.selectedRes[i].value = response.body.Data[i].Id;
							this.selectedRes[i].label = response.body.Data[i].RoleName;
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
			this.selectedRes = [];
			this.form.res = [];
			for(let i=0;i<tmpSelected.length;i++){
				this.selectedRes[i] = {};
				this.selectedRes[i].value = tmpSelected[i].Id;
				this.selectedRes[i].label = tmpSelected[i].Label;
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
		clearForm:function(){
			this.form.id = 0,
			this.form.name = '',
			this.form.pid = 0;
			this.form.role = [];
			this.form.res = [];
			this.selectedRole = [];
			this.selectedRes = [];
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
		filterNode(value, data) {
	        if (!value) return true;
	        return data.Label.indexOf(value) !== -1;
	    },
		renderContent:function(createElement, { node, data, store }) {  
            var self = this;  
            return createElement('span',{attrs:{
				style:"flex: 1; display: flex; align-items: center; justify-content: space-between; font-size: 14px; padding-right: 8px;"
			}}, [  
                createElement('span', node.label),  
                createElement('span', [  
                    createElement('el-button',{attrs:{  
                        style:"font-size:12px;",
					   type:"text"
						
                    },on:{  
                        click:function() { 
							self.handleAppend(store,data); 
                            //console.info("点击了节点" + data.id + "的添加按钮");  
                            //store.append({ id: self.baseId++, label: 'testtest', children: [] }, data);  
                        }  
                    }},"append"), 
					createElement('el-button',{attrs:{  
                        style:"font-size:12px;",
					   type:"text"
                    },on:{  
                        click:function() {
							self.handleEdit(store,node,data); 
                        }  
                    }},"edit"),   
                    createElement('el-button',{attrs:{  
                        style:"font-size:12px;",
					   type:"text"
                    },on:{  
                        click:function() {
							self.handleRemove(store,node,data); 
                        }  
                    }},"delete"),  
					createElement('el-button',{attrs:{  
                        style:"font-size:12px;",
					   type:"text"
                    },on:{  
                        click:function() {
							self.handleRole(store,node,data);
                        }  
                    }},"allotRole"),  
					createElement('el-button',{attrs:{  
                        style:"font-size:12px;",
					   type:"text"
                    },on:{  
                        click:function() {
							self.handleRes(store,node,data);
                        }  
                    }},"allotRes"),  
                ]),  
            ]); 
		}
	},
	watch: {
      filterText:function(val) {
        this.$refs.tree.filter(val);
		
      }
    }
});
</script>
<% template "footer.tpl" .%>