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

  <el-row>

    <el-menu
      :default-active="activeIndex"
      @select="handleSelect"
      active-text-color="#409eff"
      router="true"
      style="dispaly: flex;"
    >
      <el-menu-item @click="dialogVisible = true">
        <template #title>
          <el-icon><circle-plus /></el-icon>
          <span style="padding-left: 10px">新建分享</span>
        </template>
      </el-menu-item>
    </el-menu>


  
  </el-row>


  <el-dialog v-model="comment_visble" title="评论" width="40%">
    <div v-for="(comment_item,index) in comment_items" :key="index">
    <el-row>
      <div style="display:inline; position:relative; top:3px;right:3px">
        <el-avatar :size="40" :src="url_head+comment_item.useravatar"></el-avatar>
      </div>
      <span style="position:relative; left:10px;font-size:18px;font-weight: 700;top:10px">{{comment_item.username}}</span>
      <el-button style="position:absolute;right:20px;" size="small" @click="deletecomment(comment_item.comment_id)"><el-icon><delete /></el-icon></el-button>
      <el-button  v-show="comment_items[index].star_active" style="position:absolute;right:100px;" size="small" @click="addstar(comment_item.comment_id,index)">
        <el-icon><star /></el-icon>
        </el-button>
      <div v-show="!comment_items[index].star_active">
      <el-button  style="position:absolute;right:100px;" size="small" @click="cancelstar(comment_item.comment_id,index)">
        <el-icon><star-filled /></el-icon>
        </el-button>
      </div>
      <span style="position:absolute; right:70px;top:10px;display:block; width:30px">{{comment_item.like_num}}</span>
    </el-row>
    <br>
    <el-row>
      <span style=" position:relative;left:50px; font-size:18px">{{comment_item.content}}</span>
    </el-row>
    <br>
    <el-row>
      <span  style="position:absolute; right:100px">{{comment_item.created_at.slice(0,10)}}</span>
      <span  style="position:absolute; right:20px" >{{comment_item.created_at.slice(11,19)}}</span>
    </el-row>
    <br>
    <hr class="new2"/>
    </div>

    <el-form
      status-icon
      label-position="left"
      label-width="0px"
      class="new_comment_demo"
    >
  <el-form-item>
    <el-input
    v-model="new_comment_content"
    :rows="3"
    type="textarea"
    placeholder="请输入评论内容"
    />
  </el-form-item>

  </el-form>
  <el-button type="primary" @click="gotonewcomment()"
          >发布</el-button>
    
  </el-dialog>

  <div v-for="(share,index) in shares" :key="index">
    <el-card style="width: 60%;
    margin: 10px auto;">
  <el-row>
    <div class = "share_userinfo">  
      <div style="display:inline; position:relative; top:15px;right:3px">
        <el-avatar :size="40" :src="url_head+share.useravatar"></el-avatar>
      </div>
      <span style="font-size: 20px;left:50px;position: absolute;top:22px">{{share.username}}</span> 
      
    </div>
    
    <div class="share_block" >
      
      
      <el-image style="width: 185px; height: 185px"  :src="url_head+share.picture" :preview-src-list="[url_head+share.picture]"></el-image>
      <span style="margin-left:80px;margin-top:80px; font-size:18px">{{share.content}}</span>
      
      </div>
  </el-row>
  <el-row>
    <div class="share_button">
      <el-button @click="getcomment(share.share_id)"><el-icon><comment /></el-icon></el-button>
      <el-button @click="editdialogVisible_list[index] = true" :disabled="!share.authority"><el-icon><edit /></el-icon></el-button>
    <el-button @click="deleteshare(share.share_id)" :disabled="!share.authority"><el-icon><delete /></el-icon></el-button>
    <span style="position:absolute; right:80px">{{share.updated_at.slice(0,10)}}</span>
    <span style="position:absolute; right:10px">{{share.updated_at.slice(11,19)}}</span>
    </div>
  </el-row>

  <el-dialog v-model="editdialogVisible_list[index]" title="修改分享" width="30%">
    <el-form
      status-icon
      label-position="left"
      label-width="100px"
      class="demo-ruleForm login-page"
    >
  <el-form-item label="分享内容">
    <el-input
    v-model="share.content"
    :rows="3"
    type="textarea"
    placeholder="请输入分享内容"
    />
  </el-form-item>

  </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="editdialogVisible_list[index] = false">Cancel</el-button>
        <el-button type="primary" @click="gotoeditshare(share.share_id,share.content, index)"
          >Confirm</el-button>
      </span>
    </template>
  </el-dialog>

    </el-card>
  </div>

  <el-dialog v-model="dialogVisible" title="分享" width="30%">
    <el-form
      :model="myForm"
      status-icon
      ref="myForm"
      label-position="left"
      label-width="100px"
      class="demo-ruleForm login-page"
    >
  <el-form-item label="分享内容">
    <el-input
    v-model="myForm.content"
    :rows="3"
    type="textarea"
    placeholder="请输入分享内容"
    />
  </el-form-item>

  <el-form-item label="表情包">
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

  <el-form-item label="类型">
  <el-select v-model="myForm.category_name" placeholder="Select">
    <el-option
      v-for="item in options"
      :key="item.value"
      :label="item.label"
      :value="item.value"
    >
    </el-option>
  </el-select>
  </el-form-item>

  </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="gotoNewshare"
          >Confirm</el-button
        >
      </span>
    </template>
  </el-dialog>

  <br>
  <el-pagination v-model:current-page="page_num" layout="prev, pager, next" :total="50" 
  @current-change="handleCurrentChange">
  </el-pagination>

  <hr class="new1"/>
  <el-footer style="margin-top:10px">©Copyright 2021 Powered by HomerRong and YoungBest</el-footer>

</template>


<script>
import { ref } from "vue";
import { useRoute } from "vue-router";
import { ArrowDown, CirclePlus,Plus,Delete,Edit,Comment,User,Star,StarFilled} from "@element-plus/icons";
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
    CirclePlus,
    Plus,
    Delete,
    Edit,
    Comment,
    //User,
    Star,
    StarFilled,
  },

  data() {
    return {
      hidden: true,
      username: "",
      myForm:{
        content:"this is a test!!!",
        category_name: "",
      },
      imageUrl: '',
      file: new Object(),
      options: ref([
        {
          value: '猫猫',
          label: '猫猫',
        },
        {
          value: '狗狗',
          label: '狗狗',
        },
        {
          value: '机器猫',
          label: '机器猫',
        },
        {
          value: '沙雕',
          label: '沙雕',
        },
        {
          value: '美女',
          label: '美女',
        },
        {
          value: '其他',
          label: '其他',
        },
      ]),
      url_head:"/api/getimg?imgname=",
      shares: ref([]),
      comment_items :ref([]),
      editdialogVisible_list : ref([]),
      page_num: 1,
      comment_visble: false,
      new_comment_content:"",
      comment2shareid:0,
      useravatar:"",
    };
  },

  mounted: function () {// 保持登录状态
    this.refresh();
    this.useravatar = sessionStorage.getItem("useravatar");
    if (sessionStorage.getItem("is_login") == "true") {
      this.username = sessionStorage.getItem("username");
      this.hidden = false;
      return;
    }
    let temp = this.getParams().username;
    console.log("mounted" + temp);
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
    refresh() {
        service.post("/api/getshare",{
      "session_id" : sessionStorage.getItem("session_id"),
      "page_num" : this.page_num
    })
      .then((response)=>{
        this.shares = response.shares;
        for(var i = 0;i < this.shares.length;i++){
            if(this.shares[i].content==''){
              this.shares = this.shares.splice(0,i);
              break;
            }
            this.editdialogVisible_list.push(false);
        }
       
      });
    },
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
      console.log(this.file);
      this.imageUrl = URL.createObjectURL(this.file);
      
      console.log(this.imageUrl);
      return;
    },
    gotoNewshare(){
      this.dialogVisible = false;
      let formdata = new FormData();
      formdata.append("file",this.file);
      formdata.append("session_id", sessionStorage.getItem("session_id"));
      formdata.append("content",this.myForm.content);
      formdata.append("category_name",this.myForm.category_name);
      console.log(this.myForm.category_name);
      service.post("/api/newshare",formdata)
      .then(response=>{
        this.refresh();
        ElNotification({
          message: response.message
        })
      });
      
    },

    deleteshare(value){
      console.log("delete" + value);
      service.post("/api/deleteshare",{
        "session_id": sessionStorage.getItem("session_id"),
        "share_id": value,
      })
      .then(response=>{
        this.refresh();
        ElNotification({
          message: response.message
        })
      });
      
      return;
    },

    gotoeditshare(share_id,content,index){
      this.editdialogVisible_list[index] = false;
      let formdata = new FormData();
      formdata.append("session_id", sessionStorage.getItem("session_id"));
      formdata.append("content",content);
      formdata.append("share_id",share_id);
      service.post("/api/editshare",formdata)
      .then(response=>{
        ElNotification({
          message: response.message
        })
      });
    },

    getcomment(value){
      service.post("/api/getcomment",{
        "share_id": value,
      })
      .then(response=>{
        this.comment_items = response.comment_items;
        if(this.comment_items != null){
          for(let i = 0;i < this.comment_items.length;i++){
              this.comment_items[i]["star_active"] = true;
          }
        }
        
        console.log(response.comment_items);
        this.comment_visble = true;
      });
      
      this.comment2shareid = value;
      return;
    },

    gotonewcomment(){
      service.post("/api/newcomment",{
        "share_id": this.comment2shareid,
        "session_id": sessionStorage.getItem("session_id"),
        "content": this.new_comment_content,
      })
      .then(response=>{
        console.log(response);
        ElNotification({
          message: response.message
        })
      });
      this.new_comment_content = "";
      this.comment_visble = false;
      return;
    },

    deletecomment(value){
      service.post("/api/deletecomment",{
        "comment_id":value,
        "session_id": sessionStorage.getItem("session_id"),
      })
      .then(response=>{
        console.log(response);
        this.comment_visble = false;
        ElNotification({
          message: response.message
        })
      });
      
    },

    handleCurrentChange(){
      this.refresh();
    },

    addstar(value,index){
      service.post("/api/addcommentlike",{
        "comment_id":value,
      })
      this.comment_items[index].like_num += 1;
      let temp = this.comment_items[index].star_active;
      this.comment_items[index].star_active = !temp;
      //console.log(this.comment_items[index].star_active);
      return;
    },
    cancelstar(value,index){
      service.post("/api/reducecommentlike",{
        "comment_id":value,
      })
      this.comment_items[index].like_num -= 1;
      let temp = this.comment_items[index].star_active;
      this.comment_items[index].star_active = !temp;
      //console.log(this.comment_items[index].star_active);
      return;
    },
  },
};
</script>

<style>
.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
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

.share_block{
  height:230px;
  display: flex;
  margin-top: 10px;
}

.share_button{
  margin-left: 18%;
  display: flex;
}


hr.new1{
    width: 80%;
    margin: 10px auto;
    height: 1px;
    background-image: linear-gradient(
      to right,
      rgba(255, 255, 255, 0),
      cyan,
      rgba(0, 0, 0, 0)
    );
  }

hr.new2{
    width: 95%;
    margin: 10px auto;
    height: 1px;
    background-image: linear-gradient(
      to right,
      rgba(255, 255, 255, 0),
      rgba(0, 0, 0, 0)
    );
  }


.share_userinfo{
  display: flex;
  width: 18%;
}


</style>
