<template>


<el-affix>
    <el-menu
      :default-active="activeIndex"
      class="header-navigation"
      mode="horizontal"
      @select="handleSelect"
      active-text-color="#409eff"
      router="true"
    >
      <el-menu-item index="/">主页</el-menu-item>
      <el-menu-item index="/login" style="float: right">登录</el-menu-item>
    </el-menu>
    <div class="line"></div>
  </el-affix>


  <div class="background"></div>
  <div class="register">
    <el-container style="width: 100%; height: 100%">
      <el-main>
      <div style="display: flex; justify-content: center">
        <el-card style="margin-top: 10%; width: 400px; height: auto;">
          <el-form
            :model="myForm"
            status-icon
            ref="myForm"
            label-position="left"
            label-width="0px"
            class="demo-ruleForm login-page"
          >
            <h2 class="title" style="color: #409eff">注册账号</h2>
            <el-form-item prop="username">
              <el-input
                type="text"
                v-model="myForm.username"
                auto-complete="off"
                placeholder="用户名"
              >
                <template #suffix>
                  <el-icon class="el-input__icon"><user /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item prop="password">
              <el-input
                type="text"
                v-model="myForm.password"
                auto-complete="off"
                placeholder="密码"
                show-password
              >
              </el-input>
            </el-form-item>

            <el-form-item prop="password">
              <el-input
                type="text"
                v-model="myForm.password2"
                auto-complete="off"
                placeholder="确认密码"
                show-password
              >
              </el-input>
            </el-form-item>

            <el-form-item prop="question">
              <el-input
                type="text"
                v-model="myForm.question"
                auto-complete="off"
                placeholder="密保问题"
              >
              
              </el-input>
            </el-form-item>

            <el-form-item prop="answer">
              <el-input
                type="text"
                v-model="myForm.answer"
                auto-complete="off"
                placeholder="密保答案"
              >
              </el-input>
            </el-form-item>


            <el-form-item>
                <el-button
                  type="primary"
                  style="width: 100%;font-size:14px"
                  @click="handleSubmit"
                  :loading="registering"
                  >立即注册</el-button
                >
              </el-form-item>
          </el-form>
        </el-card>
      </div>
      </el-main>
    </el-container>
  </div>
  <!-- <HelloWorld msg="Welcome to Your Vue.js App"/> -->
</template>

<script>
import { User} from "@element-plus/icons"
import {ref } from 'vue'
import { ElNotification } from 'element-plus'
import  service from "../utils/request"


let registering = false;
export default {


  setup(){
    const activeIndex = ref('/login')
    const handleSelect = (key, keyPath) => {
      console.log(key, keyPath)
    }
    return{
      activeIndex,
      handleSelect,
    }
  },
  name: "register",
  //使用组件
  components: {
    User,
  },
  data() {
    return {
      registering,
      myForm: {
        username: "",
        password: "",
        password2:"",
        question:"",
        answer:"",
        //为了登录方便，可以直接在这里写好用户名和密码的值
      },
    };
  },
  methods: {
    handleSubmit(){
      
      if(this.myForm.password!==this.myForm.password2){
        ElNotification.info({
          message:"重新确认密码"
        })
        return
      }

      if(this.myForm.username === ''){
        ElNotification.info({
          message:"用户名不能为空"
        })
        return
      }


      if(this.myForm.question === '' || this.myForm.answer===''){
        ElNotification.info({
          message:"密保项不能为空"
        })
        return
      }


      if(this.myForm.password === ''){
        ElNotification.info({
          message:"密码不能为空"
        })
        return
      }
      this.registering=true;

      service.post(
            '/api/register',
          {
            username: this.myForm.username,
            password: this.myForm.password,
            question:this.myForm.question,
            answer:this.myForm.answer,
          }
      ).then(response =>{
         ElNotification({
          message: response.message
        })
        this.registering=false;
        this.$router.push('/login');
      }
      )
    }
  },
};
</script>

<style>
.register {
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

.background {
  width: 100%;
  height: 100%;
  z-index: -2;
  position: absolute;
  background: url("../assets/login_backgroud.png");
  background-size: cover;
}

body {
  margin-top: 0px;
}
</style>
