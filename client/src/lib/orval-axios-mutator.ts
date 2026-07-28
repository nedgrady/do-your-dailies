import type { AxiosRequestConfig } from 'axios'
import axios from './axios'

export const customInstance = async <T>(
  config: AxiosRequestConfig,
): Promise<T> => {
  const response = await axios(config)
  return response.data
}
