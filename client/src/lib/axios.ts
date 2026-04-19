import axios from 'axios'

export const apiBaseUrl = import.meta.env.VITE_API_BASE_URL

if (!apiBaseUrl) {
  throw new Error('VITE_API_BASE_URL is not set')
}

export default axios.create({
  baseURL: apiBaseUrl,
})
