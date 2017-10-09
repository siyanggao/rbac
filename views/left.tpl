<el-menu class="el-menu-vertical-demo"  theme="dark">
  <el-submenu index="1">
    <template slot="title">导航一</template>
    <a href="/webSentence/sentence"><el-menu-item index="1-1">sentence</el-menu-item></a>
    <a href="/webWord/word"><el-menu-item index="1-2">word</el-menu-item></a>
    <a href="/webArticle/article"><el-menu-item index="1-3">article</el-menu-item></a>
    <a href="/webMessage/message"><el-menu-item index="1-4">message</el-menu-item></a>
    
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
