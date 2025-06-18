import api from './api';
import axios from 'axios';

// Helper function to check if code is running in browser
const isBrowser = typeof window !== 'undefined';

export interface User {
  id: string;
  name: string;
  email: string;
  picture: string;
}

export interface SessionToken {
  token: string;
  expiry: number;
}

export interface AuthResponse {
  session_token: SessionToken;
  user: User;
}

export const authService = {
  /**
   * Get Google OAuth URL
   * @returns Promise with the Google OAuth URL
   */
  getGoogleAuthUrl: async (): Promise<string> => {
    try {
      const response = await api.get('/auth/google/login');
      return response.data.auth_url;
    } catch (error) {
      console.error('Error getting Google auth URL:', error);
      throw error;
    }
  },

  /**
   * Handle Google OAuth callback
   * @returns Promise with the user and session token
   */
  handleGoogleCallback: async (): Promise<AuthResponse> => {
    try {
      // Parse token from URL query parameters
      const urlParams = new URLSearchParams(window.location.search);
      const token = urlParams.get('token');

      if (!token) {
        throw new Error('No token found in URL');
      }

      // Store token in localStorage
      if (isBrowser) {
        localStorage.setItem('session_token', token);
      }

      // Make request to get user data
      // Use axios directly to call localhost:8080/user
      const response = await axios.get('http://localhost:8080/user', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      const user = response.data;

      // Store user data in localStorage
      if (isBrowser) {
        localStorage.setItem('user', JSON.stringify(user));
      }

      return {
        session_token: {
          token,
          expiry: 2592000000000000 // Default expiry (30 days)
        },
        user
      };
    } catch (error) {
      console.error('Error handling Google callback:', error);
      throw error;
    }
  },

  /**
   * Check if user is authenticated
   * @returns Boolean indicating if user is authenticated
   */
  isAuthenticated: (): boolean => {
    if (!isBrowser) return false;
    return !!localStorage.getItem('session_token');
  },

  /**
   * Get current user
   * @returns User object or null if not authenticated
   */
  getCurrentUser: (): User | null => {
    if (!isBrowser) return null;
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  },

  /**
   * Logout user
   */
  logout: (): void => {
    if (!isBrowser) return;
    localStorage.removeItem('session_token');
    localStorage.removeItem('user');
  }
};
