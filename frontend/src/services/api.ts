import axios from 'axios';

// Helper function to check if code is running in browser
const isBrowser = typeof window !== 'undefined';

const API_BASE_URL = 'http://localhost:8080/v1';

// Create axios instance with base URL
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
});

// Request interceptor to add authorization header if session token exists
api.interceptors.request.use((config) => {
  // Only access localStorage in browser environment
  if (isBrowser) {
    const sessionToken = localStorage.getItem('session_token');
    if (sessionToken) {
      config.headers.Authorization = `Bearer ${sessionToken}`;
    }

    const apiKey = localStorage.getItem('api_key');
    if (apiKey) {
      config.headers['Api-Key'] = apiKey;
    }
  }

  return config;
});

export default api;
