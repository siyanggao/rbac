<% template "header.tpl" .%>

<el-row :gutter="5">
<el-col :span="3"><% template "left.tpl" .%></el-col>
<el-col :span="21">
	<el-tree
	  :data="tree"
	  :props="defaultProps"
	  node-key="Id"
	  default-expand-all
	  :render-content="renderContent">
	</el-tree>
	<el-dialog title="dialog" :visible.sync="dialogFormVisible">
	  <el-form :model="form">
	    <el-form-item label="role name:" :label-width="formLabelWidth"
		:rules="{required:true}">
	      <el-input v-model="form.role_name" auto-complete="off"></el-input>
	    </el-form-item>
	    <el-form-item label="description:" :label-width="formLabelWidth"
		:rules="{required:true}">
	      <el-input v-model="form.description" auto-complete="off"></el-input>
	    </el-form-item>
		<el-form-item label="resource" :label-width="formLabelWidth">
			<el-button @click="resDialogVisible = true">resource</el-button>
			<el-row style="margin-top:10px">
			<el-select v-model="form.res" multiple style="width:100%">
			  <el-option
			    v-for="item in selectedRes"
			    :key="item.value"
			    :label="item.label"
			    :value="item.value">
			  </el-option>
			</el-select>
			</el-row>
			
		</el-form-item>
	  </el-form>
	  <div slot="footer" class="dialog-footer">
	    <el-button @click="dialogFormVisible = false">取 消</el-button>
	    <el-button type="primary" @click="handleOK">确 定</el-button>
	  </div>
	</el-dialog>
	
	<el-dialog title="dialog" :visible.sync="resDialogVisible" 
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
		defaultProps: {
          children: 'Children',
          label: 'Label'
        },
		form:{
			id:0,
			role_name:'',
			description:'',
			pid:0,
			res:[]
		},
		selectedRes:[],
		formLabelWidth:'120px',
		dialogFormVisible:false,
		resDialogVisible:false,
		currentData:'',
		currentStore:''
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
						Label:this.form.role_name,
						Role:{
							Id:response.body.Data,
							RoleName:this.form.role_name,
							Description:this.form.description,
							PId:this.currentData.Id
						}
					}
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
					this.currentData.Label = this.form.role_name
					this.currentData.Role.ResName = this.form.role_name
					this.currentData.Role.Description = this.form.Description
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
			this.form.role_name = data.Role.RoleName;
			this.form.description = data.Role.Description;
			this.form.pid = data.Pid;
			
			this.currentData = data;
			this.currentStore = store;
			Vue.http.post('/rpc/getresbyrole',this.form).then(response=>{
				if(response.body.Code==1){
					if(response.body.Data==null){
						this.selectedRes = [];
						this.form.res = [];
					}else{
						for(let i=0;i<response.body.Data.length;i++){
							this.selectedRes[i] = {};
							this.selectedRes[i].value = response.body.Data[i].Id;
							this.selectedRes[i].label = response.body.Data[i].ResName;
							this.form.res[i] = response.body.Data[i].Id;
						}
					}
					
					this.dialogFormVisible = true;
	  			}else{
	  				this.$message.error('get res failure');
	  			}
			});
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
			this.resDialogVisible = false;
		},
		clearForm:function(){
			this.form.id = 0,
			this.form.role_name = '',
			this.form.description = '',
			this.form.pid = 0;
			this.form.res = [];
			this.selectedRes = [];
		},
		handleResBeforeOpen:function(){
			setTimeout(() => {
				console.log(this.form.res);
		        app.$refs.resTree.setCheckedKeys(this.form.res);
		    }, 100);
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
                            //console.info("点击了节点" + data.id + "的删除按钮");  
                            //store.remove(data);  
                        }  
                    }},"edit"),   
                    createElement('el-button',{attrs:{  
                        style:"font-size:12px;",
					   type:"text"
                    },on:{  
                        click:function() {
							self.handleRemove(store,node,data);
                            //console.info("点击了节点" + data.id + "的删除按钮");  
                            //store.remove(data);  
                        }  
                    }},"delete"),  
                ]),  
            ]); 
		}
	}
});
</script>
<% template "footer.tpl" .%>