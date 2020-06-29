import Vue from 'vue'
import Vuex from 'vuex'
import credencial from './credencial'
import logins from './logins'
import getts from './getters'

Vue.use(Vuex)

export default new Vuex.Store({
  state: credencial,
  mutations: logins,
  getters: getts,
  actions: {},
  modules: {
  }
})
