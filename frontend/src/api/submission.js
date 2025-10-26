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
const submissionApi = {
  // 获取所有提交记录 (带分页)
  getSubmissions(page = 1, pageSize = 10) {
    return apiClient.get(`/submissions?page=${page}&page_size=${pageSize}`).then(response => response.data)
  },
  
  // 获取单个提交详情
  getSubmission(id) {
    return apiClient.get(`/submissions/${id}`).then(response => response.data)
  },
  
  // 提交代码
  submitCode(data) {
    return apiClient.post('/submissions', data).then(response => response.data)
  },
  
  // 获取提交结果
  getSubmissionResult(id) {
    return apiClient.get(`/submissions/${id}`).then(response => response.data)
  }
}

export default submissionApi