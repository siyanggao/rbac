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
	    <el-form-item label="resource name:" :label-width="formLabelWidth"
		:rules="{required:true}">
	      <el-input v-model="form.res_name" auto-complete="off"></el-input>
	    </el-form-item>
	    <el-form-item label="resource code:" :label-width="formLabelWidth"
		:rules="{required:true}">
	      <el-input v-model="form.res_code" auto-complete="off"></el-input>
	    </el-form-item>
	  </el-form>
	  <div slot="footer" class="dialog-footer">
	    <el-button @click="dialogFormVisible = false">取 消</el-button>
	    <el-button type="primary" @click="handleOK">确 定</el-button>
	  </div>
	</el-dialog>
</el-col>
</el-row>
</div>
<script type="text/javascript">
var app = new Vue({
	el:'#app',
	data:{
		tree:<%.tree%>,
		defaultProps: {
          children: 'Children',
          label: 'Label'
        },
		form:{
			id:0,
			res_name:'',
			res_code:'',
			pid:0
		},
		formLabelWidth:'120px',
		dialogFormVisible:false,
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
						Label:this.form.res_name,
						Res:{
							Id:response.body.data,
							ResName:this.form.res_name,
							ResCode:this.form.res_code,
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
	  			console.log(response);
	  			if(response.body.Code==1){
	  				this.$message({
			          message: 'edit success',
			          type: 'success'
			        });
					this.currentData.Label = this.form.res_name
					this.currentData.Res.ResName = this.form.res_name
					this.currentData.Res.ResCode = this.form.res_code
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
			console.log(data);
			this.form.id = data.Id;
			this.form.res_name = data.Res.ResName;
			this.form.res_code = data.Res.ResCode;
			this.form.pid = data.Pid;
			this.dialogFormVisible = true;
			this.currentData = data;
			this.currentStore = store;
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
		clearForm:function(){
			this.form.id = 0,
			this.form.res_name = '',
			this.form.res_code = '',
			this.form.pid = 0;
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
