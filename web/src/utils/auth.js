// src/utils/auth.js

// 简单的 token 获取/存储示例
export function getToken() {
  return localStorage.getItem('token') // 或者你项目里存 token 的 key
}

export function setToken(token) {
  localStorage.setItem('token', token)
}

export function removeToken() {
  localStorage.removeItem('token')
}
