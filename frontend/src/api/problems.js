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

const problemsApi = {
  // 获取所有题目 (带分页)
  getProblems() {
    return apiClient.get(`/problems`).then(response => response.data)
  },
  
  // 获取单个题目详情
  getProblem(id) {
    return apiClient.get(`/problems/${id}`).then(response => response.data)
  },

  // 提交代码到指定接口
  submitCode(submissionData) {
    return apiClient.post('/submit', submissionData).then(response => response.data)
  },
  
  // 创建题目
  createProblem(problemData) {
    return apiClient.post('/problems', problemData).then(response => response.data)
  },
  
  // 更新题目
  updateProblem(id, problemData) {
    return apiClient.put(`/problems/${id}`, problemData).then(response => response.data)
  },
  
  // 删除题目
  deleteProblem(id) {
    return apiClient.delete(`/problems/${id}`).then(response => response.data)
  }
}

export default problemsApi