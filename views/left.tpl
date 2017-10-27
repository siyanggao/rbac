<el-menu default-active=<%.menu_index%> class="el-menu-vertical-demo"  theme="dark">
  <el-submenu index="1">
    <template slot="title">导航一</template>
	<a href="/user/toview"><el-menu-item index="1-1">user</el-menu-item></a>
    <a href="/depart/toview"><el-menu-item index="1-2">depart</el-menu-item></a>
    <a href="/role/toview"><el-menu-item index="1-3">role</el-menu-item></a>
    <a href="/resource/toview"><el-menu-item index="1-4">resource</el-menu-item></a>
    
    
    <el-submenu index="1-4">
      <template slot="title">选项4</template>
      <el-menu-item index="1-4-1">选项1</el-menu-item>
    </el-submenu>
  </el-submenu>
  <el-submenu index="2">
    <template slot="title">system</template>
    <el-menu-item index="2-1">user</el-menu-item>
    <el-menu-item index="3">导航三</el-menu-item>
  </el-submenu>
  
</el-menu>
