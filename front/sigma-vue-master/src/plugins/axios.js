import Vue from 'vue'
import axios from 'axios'

// axios.defaults.baseURL = "https://curso-vue-50e2a.firebaseio.com/"


Vue.use({
  install(Vue) {
    // Vue.prototype.$http = axios  
    // Vue.prototype.$httpBaseURL = process.env.NODE_ENV === 'production' ? 'http://localhost:8081' : 'http://localhost:8081'

    Vue.prototype.$http = axios.create({
      baseURL: process.env.NODE_ENV === 'production' ? 'http://localhost:8081' : 'http://localhost:8081',
      
    })   

    // Vue.prototype.$http.interceptors.request.use(config => {
    //   return config
    // }, error => Promise.reject(error))

    // Vue.prototype.$http.interceptors.response.use(resp => {
    //   return resp
    // }, error => Promise.reject(error))

  }
})