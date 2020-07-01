import Vue from 'vue'
import axios from 'axios'

// axios.defaults.baseURL = "https://curso-vue-50e2a.firebaseio.com/"


Vue.use({
  install(Vue) {
    // Vue.prototype.$http = axios  
    // Vue.prototype.$httpBaseURL = process.env.NODE_ENV === 'production' ? 'http://localhost:8081' : 'http://localhost:8081'

    Vue.prototype.$acl = axios.create({
      baseURL: process.env.NODE_ENV === 'production' ? 'http://192.168.1.111:20100' : 'http://localhost:20100',
      // timeout:	10000,
    })   

  }
})


// this.$acl.interceptors.request.use((config) => {
//   const token = this.$cookies.get("token");

//   if (token) {
//     config.headers.Authorization = `Bearer ${token}`
//   }

//   return config
// }, (err) => {
//   return Promise.reject(err)
// })
