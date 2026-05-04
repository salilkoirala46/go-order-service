import http from './axios'

export const authApi = {
  register: (name, email, password) =>
    http.post('/auth/register', { name, email, password }),

  login: (email, password) =>
    http.post('/auth/login', { email, password }),

  me: () =>
    http.get('/users/me')
}
