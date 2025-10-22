import axios from 'axios'

const API_URL = 'http://localhost:8080/api'

// 创建axios实例
const apiClient = axios.create({
  baseURL: API_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器，添加token
apiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器，处理错误
apiClient.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      // 未授权，清除token并跳转到登录页
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

const authApi = {
  /**
   * 用户登录
   * @param {Object} credentials - 登录凭证
   * @param {string} credentials.username - 用户名
   * @param {string} credentials.password - 密码
   * @returns {Promise} - 返回包含token和用户信息的Promise
   */
  login(credentials) {
    return apiClient.post('/auth/login', credentials)
      .then(response => {
        localStorage.setItem('id', response.data.user.ID)
        return response.data
      })
  },

  /**
   * 用户注册
   * @param {Object} userData - 用户数据
   * @returns {Promise} - 返回注册结果的Promise
   */
  register(userData) {
    return apiClient.post('/auth/register', userData)
      .then(response => response.data)
  },

  /**
   * 获取当前用户信息
   * @returns {Promise} - 返回用户信息的Promise
   */
  getCurrentUser() {
    return apiClient.get('/auth/me')
      .then(response => response.data)
  },
  
  /**
   * 退出登录
   */
  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }
}

export default authApi