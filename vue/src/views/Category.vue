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

  <el-row style="display:flex;justify-content:center">
    <el-tooltip
        class="item"
        effect="dark"
        content="猫猫表情包"
        placement="top-start"
      >
    <el-button @click="handlebutton(category_list[0])"> {{category_list[0]}}</el-button>
    </el-tooltip>

    <el-tooltip
        class="item"
        effect="dark"
        content="狗狗表情包"
        placement="top-start"
      >
    <el-button @click="handlebutton(category_list[1])"> {{category_list[1]}}</el-button>
    </el-tooltip>

    <el-tooltip
        class="item"
        effect="dark"
        content="(⊙o⊙)哆啦A梦表情包"
        placement="top-start"
      >
    <el-button @click="handlebutton(category_list[2])"> {{category_list[2]}}</el-button>
    </el-tooltip>

    <el-tooltip
        class="item"
        effect="dark"
        content="建议上传同学的表情包"
        placement="top-start"
      >
    <el-button @click="handlebutton(category_list[3])"> {{category_list[3]}}</el-button>
    </el-tooltip>

    <el-tooltip
        class="item"
        effect="dark"
        content="oyb喜欢的美女的表情包(⊙o⊙) 只能有美女!! 圆梦美女数据库"
        placement="top-start"
      >
    <el-button @click="handlebutton(category_list[4])"> {{category_list[4]}}</el-button>
    </el-tooltip>

    <el-tooltip
        class="item"
        effect="dark"
        content="分不了类就是说"
        placement="top-start"
      >
    <el-button @click="handlebutton(category_list[5])"> {{category_list[5]}}</el-button>
    </el-tooltip>
  </el-row>

  <br>

  <el-row style="margin-left:300px; margin-right:300px;margin-bottom:10px">
      <div v-for="(category,index) in categories" :key="index">
          <el-image style="width: 200px; height: 200px ;margin-left:50px; margin-right:50px; margin-top: 120px"  :src="url_head+category.picture" :preview-src-list="[url_head+category.picture]"></el-image>
      <div style="font-size: 20px;padding-left:5px;margin-top:10px">
      <div style="display:inline; position:relative; top:10px;right:10px">
        <el-avatar :size="30" :src="url_head+category.useravatar"></el-avatar>
      </div>
        {{category.username}}
      </div>
      <div style="display:flex;margin-left:50px;margin-top:20px">
        
        <el-button  v-show="categories[index].star_active" style="" size="small" @click="addstar(category.sticker_id,index)">
        <el-icon><star /></el-icon>
        </el-button>
      <div v-show="!categories[index].star_active">
      <el-button  style="" size="small" @click="cancelstar(category.sticker_id,index)">
        <el-icon><star-filled /></el-icon>
        </el-button>
      </div>
        <span style="font-size: 20px ;padding-left:10px;padding-top:5px;padding-right:40px">{{category.like_num}}</span>
        <el-button  v-show="categories[index].collection_active" style="" size="small" @click="addcollection(category.sticker_id,index)">
        <el-icon><collection /></el-icon>
        collect
        </el-button>
      <div v-show="!categories[index].collection_active">
      <el-button  style="" size="small" @click="cancelcollection(category.collection_id,index)">
        <el-icon><collection-tag /></el-icon>
        collected
        </el-button>
      </div>
        <span style="font-size: 20px ;padding-left:10px;padding-top:5px">{{category.collection_num}}</span>
      </div>
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
import { Star,User,Collection,StarFilled,CollectionTag,ArrowDown} from "@element-plus/icons";
import service from "../utils/request";
import { ElNotification } from "element-plus";
import axios from 'axios';

export default {
  setup() {
    const activeIndex = ref("/category");
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
    Star,
    //User,
    Collection,
    StarFilled,
    CollectionTag,
    ArrowDown,
  },

  data() {
    return {
      hidden: true,
      categories: ref([]),
      page_num: 1,
      url_head:"/api/getimg?imgname=",
      category_name: "美女",
      category_list:['猫猫','狗狗','机器猫','沙雕','美女','其他'],
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
        service.post("/api/getcategory",{
      "category_name" : this.category_name,
      "page_num" : this.page_num,
      "session_id": sessionStorage.getItem("session_id"),
    })
      .then((response)=>{
        this.categories = response.category_items;
        for(let i = 0;i < 9;i++){
            if(this.categories[i].picture == ''){
                this.categories = this.categories.splice(0,i);
                break;
            }
        }
        if(this.categories != null){
          for(let i = 0;i < this.categories.length;i++){
              this.categories[i]["star_active"] = true;
              if(this.categories[i].collection_id == 0) this.categories[i]['collection_active'] = true;
              else this.categories[i]['collection_active'] = false;
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
    
    getcomment(value){
      return;
    },

    handleCurrentChange(){
      this.refresh();
    },

    handlebutton(value){
        this.category_name = value;
        this.refresh();
    },

    addstar(value,index){
      service.post("/api/addstickerlike",{
        "sticker_id":value,
      })
      this.categories[index].like_num += 1;
      let temp = this.categories[index].star_active;
      this.categories[index].star_active = !temp;
      return;
    },
    cancelstar(value,index){
      service.post("/api/reducestickerlike",{
        "sticker_id":value,
      })
      this.categories[index].like_num -= 1;
      let temp = this.categories[index].star_active;
      this.categories[index].star_active = !temp;
      return;
    },

    addcollection(value,index){
      service.post("/api/newcollection",{
        "sticker_id":value,
        "session_id": sessionStorage.getItem("session_id"),
      }).then((response) => {
          this.categories[index].collection_num += 1;
          let temp = this.categories[index].collection_active;
          this.categories[index].collection_active = !temp;
      })
      
      return;
    },
    cancelcollection(value,index){
      service.post("/api/deletecollection",{
        "collection_id":value,
        "session_id": sessionStorage.getItem("session_id"),
      })
      this.categories[index].collection_num -= 1;
      let temp = this.categories[index].collection_active;
      this.categories[index].collection_active = !temp;
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
