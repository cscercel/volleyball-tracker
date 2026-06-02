let accessToken = $state(localStorage.getItem('access_token') ?? '')

export const auth = {
  get token() { return accessToken },
  get isAuthenticated() { return accessToken !== '' },

  login(token: string) {
    accessToken = token
    localStorage.setItem('access_token', token)
  },

  logout() {
    accessToken = ''
    localStorage.removeItem('access_token')
  }
}
