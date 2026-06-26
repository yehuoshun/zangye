// HTTP 请求封装
// 统一错误处理、JSON 解析
// 类比 Axios，但使用 fetch API

const BASE_URL = '/api'

interface RequestOptions {
  method?: string
  body?: any
  headers?: Record<string, string>
}

class ApiError extends Error {
  status: number
  constructor(message: string, status: number) {
    super(message)
    this.status = status
  }
}

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, headers = {} } = options

  const config: RequestInit = {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...headers,
    },
  }

  if (body && method !== 'GET') {
    config.body = JSON.stringify(body)
  }

  const response = await fetch(`${BASE_URL}${path}`, config)

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: '请求失败' }))
    throw new ApiError(errorData.error || '请求失败', response.status)
  }

  // 处理空响应
  const text = await response.text()
  if (!text) return {} as T

  return JSON.parse(text) as T
}

// 获取列表响应（含 total/page/size）
async function requestList<T>(path: string, options: RequestOptions = {}): Promise<{ data: T[]; total: number; page: number; size: number }> {
  const result = await request<{ data: T[]; total: number; page: number; size: number }>(path, options)
  return result
}

export { request, requestList, ApiError }
export default { request, requestList }
