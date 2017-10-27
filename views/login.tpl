<% template "header.tpl" .%>
<el-row style="margin-top:100px">
<el-col :span="8" :offset="8">
<el-form :model="form" :rules="rules" ref="ruleForm" label-width="100px" class="demo-ruleForm">
  <el-form-item label="username:" prop="username" label-width="100px">
    <el-input v-model="form.username"></el-input>
  </el-form-item>
  <el-form-item label="password:" prop="pwd" >
    <el-input v-model="form.pwd" type="password"></el-input>
  </el-form-item>
  
  <el-form-item>
    <el-button type="primary" @click="submitForm('ruleForm')">login</el-button>
    <el-button @click="resetForm('ruleForm')">reset</el-button>
  </el-form-item>
</el-form>
</el-col>
</el-row>
</div>
<script>
var app = new Vue({
	el:"#app",
	data:{
		form:{
			username:'',
			pwd:''
		},
		rules: {
          username: [
            { required: true, message: 'please input username', trigger: 'blur' },
            { min: 3, max: 10, message: '长度在 3 到 10 个字符', trigger: 'blur' }
          ],
          pwd: [
            { required: true, message: 'please input password', trigger: 'blur' },
			{ min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
          ]
        }
	},
	methods: {
      submitForm(formName) {
        this.$refs[formName].validate((valid) => {
          if (valid) {
            Vue.http.post('/user/login',this.form).then(response=>{
				if(response.body.Code==1){
					this.$message({
			          message: 'login success',
			          type: 'success'
			        });
					window.location.href="/user/toview";
				}else{
					this.$message.error('login failure');
				}
			});
          } else {
            console.log('error submit!!');
            return false;
          }
        });
      },
	  resetForm(formName) {
        this.$refs[formName].resetFields();
      }
    }
});
</script>
<% template "footer.tpl" .%>