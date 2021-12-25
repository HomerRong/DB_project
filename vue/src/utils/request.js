import { ElNotification } from 'element-plus'

import axios from 'axios'


// 创建axios实例
const service = axios.create({
    baseURL: process.env.BASE_API,
    timeout: 50000 // 请求超时时间
})


// response 拦截器
service.interceptors.response.use(
    response => {
        const res = response.data
        if (res.code !== 0) { // 后台返回码，0为成功
            ElNotification.error({
                title: '错误',
                message: res.message, // 错误描述信息
                duration: 0
            })
            return Promise.reject('error')
        } else {
            return response.data
        }
    },
    error => {
        console.log('err' + error) // for debug
        Notification.error({
            title: '错误',
            message: error,
            duration: 0
        })
        return Promise.reject(error)
    }
)


export default service