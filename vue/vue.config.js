// 跨域路由
module.exports = {
    publicPath: './',
    devServer: {
        proxy: {
            '/api': {
                // target: 'http://172.18.52.124:9000/api', //接口域名
                target: 'http://127.0.0.1:9000', //接口域名
                changeOrigin: true,             //是否跨域
                ws: true,                       //是否代理 websockets
                secure: false,                   //是否https接口
                pathRewrite: {                  //路径重置
                    '^/api': ''
                }
            }
        }
    }
};
