import http from './axios'

export const ordersApi = {
  list: () =>
    http.get('/orders'),

  get: (id) =>
    http.get(`/orders/${id}`),

  create: (product, quantity, price) =>
    http.post('/orders', { product, quantity, price }),

  updateStatus: (id, status) =>
    http.patch(`/orders/${id}/status`, { status })
}
