<template>
  <div id="app">
  

  <!-- 路由出口 -->
  <!-- 路由匹配到的组件将渲染在这里 -->
  <router-view></router-view>
</div>
  <!-- <HelloWorld msg="Welcome to Your Vue.js App"/> -->
</template>

<script>
// import HelloWorld from './components/HelloWorld.vue'
//注册icon组件

import {ref} from 'vue'
import {useRoute} from 'vue-router'
// let imgsrc = "http://127.0.0.1:9000/api/getimg?imgname=a.png";

export default {
  setup(){
    const activeIndex = ref('/')
    const handleSelect = (key, keyPath) => {
      console.log(key, keyPath)
    }
    
    return{
      activeIndex,
      handleSelect,

    }
  },
  name: "App",
  
  //使用组件
  components: {
  },
  data() {
    return {
      
    };
  },

  methods: {
    
  },

  
  
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  /* margin-top: 60px; */
  z-index: 2;
  position: absolute;
  width: 100%;
  height: 100%; /**宽高100%是为了图片铺满屏幕 */
}

</style>
