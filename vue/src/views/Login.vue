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
  <div class="login">
    <el-container style="width: 100%; height: 100%">
      <!-- <el-aside width="50%">
        <div class="title" style="margin-top: 40%; font-size: 20px"></div>
      </el-aside>
      <div class="square"></div> -->
      <el-main>
        <div style="display: flex; justify-content: center">
          <el-card style="margin-top: 10%; width: 350px; height: auto">
            <el-form
              :model="myForm"
              status-icon
              ref="myForm"
              label-position="left"
              label-width="0px"
              class="demo-ruleForm login-page"
            >
              <h2 class="title" style="color: #409eff">账号密码登录</h2>
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
                  type="password"
                  v-model="myForm.password"
                  auto-complete="off"
                  placeholder="密码"
                  show-password
                />
              </el-form-item>
              <div
                style="
                  display: flex;
                  padding-left: 0px;
                  justify-content: space-between;
                "
              >
                <el-checkbox
                  v-model="checked"
                  class="rememberme"
                  label="自动登录"
                ></el-checkbox>
                <el-link
                  style="color: #409eff; font-size: 14px"
                  @click="goToresetpassword"
                  >忘记密码</el-link
                >
              </div>
              <p></p>
              <el-form-item>
                <el-button
                  type="primary"
                  style="width: 100%; font-size: 14px"
                  @click="handleSubmit"
                  >立即登录</el-button
                >
              </el-form-item>
              <div style="display: flex; justify-content: center">
                <el-link
                  style="color: #409eff; font-size: 14px"
                  @click="goToRegister"
                  >注册账号</el-link
                >
              </div>
            </el-form>
          </el-card>
        </div>
      </el-main>
    </el-container>

    <!-- <div class="demo-image__preview">
        <el-image
          style="width: 500px; height: 300px"
          :src="imgsrc"
        >
        </el-image>
      </div> -->
  </div>
  <!-- <HelloWorld msg="Welcome to Your Vue.js App"/> -->
</template>

<script>
// import HelloWorld from './components/HelloWorld.vue'
//注册icon组件
import { User } from "@element-plus/icons";
import { ref } from "vue";
import service from "../utils/request";
import { ElNotification } from "element-plus";
import { useRouter } from 'vue-router';
// let imgsrc = "http://127.0.0.1:9000/api/getimg?imgname=a.png";

let checked = false;
export default {
  setup() {
    const activeIndex = ref('/login')
    const handleSelect = (key, keyPath) => {
      console.log(key, keyPath)
    }
    return{
      activeIndex,
      handleSelect,
    }
  },
  name: "login",
  //使用组件
  components: {
    User,
  },
  data() {
    return {
      checked,
      myForm: {
        username: "john",
        password: "321",
        //为了登录方便，可以直接在这里写好用户名和密码的值
      },
    };
  },
  mounted() {
    //this.getData();
    // this.getImg();
  },
  methods: {
    handleSubmit() {
      //console.log(this.myForm.username + this.myForm.password);

      //console.log(this.logining);
      if (this.myForm.username === '' || this.myForm.password === '') {
        ElNotification.info({
          message: "用户名和密码不能为空",
        });
        return;
      }


      service
        .post("/api/login", this.myForm)
        .then((response) => {
          sessionStorage.setItem("session_id",response.session_id)
          sessionStorage.setItem("useravatar", response.useravatar)
          this.$router.push({
            name: 'home',
            params:{
              username: this.myForm.username,
            },
          });
        });
    },
    goToRegister() {
      this.$router.push("/register");
    },
    goToresetpassword() {
      this.$router.push("/resetpassword");
    },
  },
};
</script>

<style>
.login {
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

.square {
  width: 50%;
  height: 100%;
  z-index: -1;
  position: absolute;
  background: white;
  opacity: 0.6;
  right: 0px;
}

body {
  margin-top: 0px;
}
</style>
