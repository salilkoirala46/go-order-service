import http from './axios'

export const notificationsApi = {
  list: () =>
    http.get('/notifications'),

  get: (id) =>
    http.get(`/notifications/${id}`),

  markRead: (id) =>
    http.patch(`/notifications/${id}/read`)
}
