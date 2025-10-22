import { createStore } from 'vuex'
import authApi from '../api/auth'

export default createStore({
  state: {
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user')) || null
  },
  mutations: {
    SET_TOKEN(state, token) {
      state.token = token
      localStorage.setItem('token', token)
    },
    SET_USER(state, user) {
      state.user = user
      localStorage.setItem('user', JSON.stringify(user))
    },
    LOGOUT(state) {
      state.token = ''
      state.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  },
  actions: {
    login({ commit }, credentials) {
      return new Promise((resolve, reject) => {
        authApi.login(credentials)
          .then(response => {
            const { token, user } = response
            commit('SET_TOKEN', token)
            commit('SET_USER', user)
            resolve(response)
          })
          .catch(error => {
            reject(error)
          })
      })
    },
    register({ commit }, userData) {
      return new Promise((resolve, reject) => {
        authApi.register(userData)
          .then(response => {
            const { token, user } = response
            commit('SET_TOKEN', token)
            commit('SET_USER', user)
            resolve(response)
          })
          .catch(error => {
            reject(error)
          })
      })
    },
    logout({ commit }) {
      authApi.logout()
      commit('LOGOUT')
    },
    getCurrentUser({ commit }) {
      return new Promise((resolve, reject) => {
        authApi.getCurrentUser()
          .then(response => {
            const { user } = response
            commit('SET_USER', user)
            resolve(user)
          })
          .catch(error => {
            reject(error)
          })
      })
    }
  },
  getters: {
    isAuthenticated: state => !!state.token,
    currentUser: state => state.user
  }
})