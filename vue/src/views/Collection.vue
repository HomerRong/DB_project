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

  <div v-show="hidden">请登录后查看收藏</div>
  <el-row style="margin-left:300px; margin-right:300px;margin-bottom:10px">
      <div v-for="(collection,index) in collections" :key="index">
          <el-image style="width: 200px; height: 200px ;margin-left:50px; margin-right:50px; margin-top: 120px"  :src="url_head+collection.picture" :preview-src-list="[url_head+collection.picture]"></el-image>
          <br>
          <el-button @click="cancelcollection(collections[index].collection_id,index)"><el-icon><delete /></el-icon></el-button>
      </div>
  </el-row>

  <el-pagination v-model:current-page="page_num" layout="prev, pager, next" :total="50" 
  @current-change="handleCurrentChange">
  </el-pagination>

  <hr />
  <el-footer style="margin-top:10px">©Copyright 2021 Powered by HomerRong and YoungBest</el-footer>

</template>


<script>
import { ref } from "vue";
import { useRoute } from "vue-router";
import { Star,User,Collection,StarFilled,CollectionTag,Delete,ArrowDown} from "@element-plus/icons";
import service from "../utils/request";
import { ElNotification } from "element-plus";
import axios from 'axios';

export default {
  setup() {
    const activeIndex = ref("/collection");
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
      Delete,
      ArrowDown,
  },

  data() {
    return {
      hidden: true,
      collections: ref([]),
      page_num: 1,
      url_head:"/api/getimg?imgname=",
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
    refresh() {
        service.post("/api/getcollection",{
      "page_num" : this.page_num,
      "session_id": sessionStorage.getItem("session_id"),
    })
      .then((response)=>{
        this.collections = response.collection_items;
        for(let i = 0;i < 9;i++){
            if(this.collections[i].picture == ''){
                this.collections = this.collections.splice(0,i);
                break;
            }
        }
       
        //console.log(this.categories);
       
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
    

    handlebutton(value){
        this.category_name = value;
        this.refresh();
    },

    cancelcollection(value,index){
      service.post("api/deletecollection",{
        "collection_id":value,
        "session_id": sessionStorage.getItem("session_id"),
      }).then((response) => {
          this.refresh();
      })
      return;
    },
    gotoUploadAvatar(){
          this.$router.push("/uploadavatar");
    },
  },
};
</script>

<style>


hr {
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



</style>
