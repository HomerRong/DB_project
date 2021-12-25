<template>
  <el-affix>
    <el-menu
      :default-active="activeIndex"
      class="header-navigation"
      mode="horizontal"
      @select="handleSelect"
      active-text-color="#409eff"
      router="true"
      style="dispaly: flex; justify-content: flex-end"
    >
      <el-menu-item index="/">主页</el-menu-item>
      <el-menu-item index="/category">分类</el-menu-item>
      <el-menu-item index="/collection">收藏</el-menu-item>
      <div v-show="hidden">
        <el-menu-item index="/login" style="float: right">登录</el-menu-item>
      </div>

      <div
        v-show="!hidden"
        style="line-height: 4; float: right; padding-right: 20px"
      >
      <div style="display:inline; position:relative; top:15px;right:10px">
        <el-avatar :size="35" :src="url_head+useravatar"></el-avatar>
      </div>
        <el-dropdown>
          <span class="el-dropdown-link">
            {{ username }}
            <el-icon class="el-icon--right">
              <arrow-down />
            </el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="gotoLogout">登出</el-dropdown-item>
              <el-dropdown-item @click="gotoUploadAvatar">修改头像</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-menu>
    <div class="line"></div>
  </el-affix>

  <br>

  <div style="display:flex; justify-content:center">

  <el-form
      :model="myForm"
      status-icon
      ref="myForm"
      label-position="left"
      label-width="100px"
      class="demo-ruleForm login-page"
    >

  <el-form-item label="个人头像">
    <el-upload
    class="avatar-uploader"
    action=""
    :show-file-list="false"
    :http-request="myUpload"
    :auto-upload="true"
  >
    <img v-if="imageUrl" :src="imageUrl" class="avatar" />
    <el-icon v-else class="avatar-uploader-icon"><plus /></el-icon>
  </el-upload>
  </el-form-item>


  </el-form>
    </div>
  <div>

  <el-button type="primary" @click="UploadAvatar"
          >确认修改</el-button>
  </div>
  

</template>


<script>
import { ref } from "vue";
import { useRoute } from "vue-router";
import { Plus,ArrowDown} from "@element-plus/icons";
import service from "../utils/request";
import { ElNotification } from "element-plus";
import axios from 'axios';

export default {
  setup() {
    const activeIndex = ref("/");
    const handleSelect = (key, keyPath) => {
      console.log(key, keyPath);
    };
    
    const route = useRoute();
    const getParams = () => {
      return route.params;
    };
    const dialogVisible = ref(false);
    return {
      activeIndex,
      handleSelect,
      getParams,
      dialogVisible,
    };
  },

  components: {
      ArrowDown,
        Plus,
  },

  data() {
    return {
      hidden: true,
      file: new Object(),
      url_head:"/api/getimg?imgname=",
      imageUrl:"",
      useravatar:"",
    };
  },

  mounted: function () {// 保持登录状态
    this.imageUrl = this.url_head + sessionStorage.getItem("useravatar");
    this.useravatar = sessionStorage.getItem("useravatar");
    if (sessionStorage.getItem("is_login") == "true") {
      this.username = sessionStorage.getItem("username");
      this.hidden = false;
      return;
    }
    let temp = this.getParams().username;
    //console.log("mounted" + temp);
    this.username = temp;
    sessionStorage.setItem("username", temp);
    if (typeof this.username == "undefined") {
      // undefined 判断
      this.hidden = true;
    } else {
      sessionStorage.setItem("is_login", "true");
      this.hidden = !this.hidden;
    }

    
  },



  methods: {
    gotoLogout() {
      console.log("logout");
      let session_id = sessionStorage.getItem("session_id");
      console.log(session_id);
      service
        .post("/api/logout", {
          session_id: session_id,
        })
        .then((response) => {
          sessionStorage.setItem("is_login", "false");
          ElNotification({
            message: response.message,
          });

          this.$router.push("/login");
        });
    },

    gotoUploadAvatar(){
          this.$router.push("/uploadavatar");
    },
    
    myUpload(param){// 自己定义了上传请求，原来的上传成功失败的hook函数都失效了
      this.file = param.file;// 先保存,等到提交时再上传
      //console.log(this.file);
      this.imageUrl = URL.createObjectURL(this.file);
      
      //console.log(this.imageUrl);
      return;
    },

    UploadAvatar(){
        this.dialogVisible = false;
        let formdata = new FormData();
        formdata.append("file",this.file);
        formdata.append("session_id", sessionStorage.getItem("session_id"));
        service.post("/api/uploadavatar",formdata)
        .then(response=>{
            sessionStorage.setItem("useravatar", response.useravatar)
            ElNotification({
            message: response.message
            })
        });
        return;
    },

  },
};
</script>

<style>

.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 50%;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}
.avatar-uploader .el-upload:hover {
  border-color: #409eff;
}
.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
}
.avatar-uploader-icon svg {
  margin-top: 74px; /* (178px - 28px) / 2 - 1px */
}
.avatar {
  width: 178px;
  height: 178px;
  display: block;
}



</style>
