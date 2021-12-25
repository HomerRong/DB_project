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
  <div class="resetpassword">
    <el-container style="width: 100%; height: 100%">
      <el-main>
        <div style="display: flex; justify-content: center">
          <div v-show="hidden">
            <el-card style="margin-top: 50%; width: 400px; height: auto">
              <el-form
                :model="myForm"
                status-icon
                ref="myForm"
                label-position="left"
                label-width="0px"
                class="demo-ruleForm login-page"
              >
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

                <el-form-item>
                  <el-button
                    type="primary"
                    style="width: 100%; font-size: 14px"
                    @click="goToNext"
                    >下一步</el-button
                  >
                </el-form-item>
              </el-form>
            </el-card>
          </div>

          <div v-show="!hidden">
            <el-card style="margin-top: 50%; width: 400px; height: auto">
              <el-form
                :model="myForm2"
                status-icon
                ref="myForm"
                label-position="left"
                label-width="0px"
                class="demo-ruleForm login-page"
              >
                <el-form-item>
                  <p align="left">密保问题: {{ question }}</p>
                  <el-input
                    type="text"
                    v-model="myForm2.answer"
                    auto-complete="off"
                    placeholder="正确答案"
                  >
                  </el-input>
                </el-form-item>

                <el-form-item prop="password">
                  <el-input
                    type="text"
                    v-model="myForm2.password"
                    auto-complete="off"
                    placeholder="新密码"
                    show-password
                  >
                  </el-input>
                </el-form-item>

                <el-form-item prop="password">
                  <el-input
                    type="text"
                    v-model="myForm2.password2"
                    auto-complete="off"
                    placeholder="确认密码"
                    show-password
                  >
                  </el-input>
                </el-form-item>

                <el-form-item>
                  <el-button
                    type="primary"
                    style="width: 100%; font-size: 14px"
                    @click="handleSubmit"
                    >重置密码
                  </el-button>
                </el-form-item>
              </el-form>
            </el-card>
          </div>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { User } from "@element-plus/icons";
import { ref } from "vue";
import { ElNotification } from "element-plus";
import service from "../utils/request";

export default {
  setup() {
    const activeIndex = ref("/login");
    const handleSelect = (key, keyPath) => {
      console.log(key, keyPath);
    };
    return {
      activeIndex,
      handleSelect,
    };
  },
  name: "resetpassword",
  //使用组件
  components: {
    User,
  },
  data() {
    return {
      hidden: true,
      question: "",
      myForm: {
        username: "",
        //为了登录方便，可以直接在这里写好用户名和密码的值
      },
      myForm2: {
        answer: "",
        password: "",
        password2: "",
      },
    };
  },
  methods: {
    //重置密码按钮
    handleSubmit() {
      if (this.myForm2.password !== this.myForm2.password2) {
        ElNotification.info({
          message: "重新确认密码",
        });
        return;
      }

      if(this.myForm2.answer===''){
        ElNotification.info({
          message:"密保回答不能为空"
        })
        return
      }


      if(this.myForm2.password === ''){
        ElNotification.info({
          message:"密码不能为空"
        })
        return
      }

      service
        .post("/api/resetpassword", {
          new_password: this.myForm2.password,
          username: this.myForm.username,
          answer: this.myForm2.answer,
        })
        .then((response) => {
          ElNotification({
            message: response.message,
          });
          this.$router.push("/login");
        });
    },

    // 下一步按钮
    goToNext() {
      if (this.myForm.username === '') {
        ElNotification.info({
          message: "用户名不能为空",
        });
        return;
      }
      service
        .post("/api/getquestion", {
          username: this.myForm.username,
        })
        .then((response) => {
          this.question = response.message;
          //console.log(response.message)
          this.hidden = !this.hidden;
          //console.log(this.hidden);
        });
    },
  },
};
</script>

<style>
.resetpassword {
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
