import api from './api';

export interface ChatMessage {
  role: 'human' | 'ai';
  text: string;
}

export interface ChatTitle {
  id: number;
  title: string;
}

export interface ChatHistory {
  chatHistory: ChatMessage[];
}

export interface ChatRequest {
  id: number;
  model_type: 'Google' | 'OpenAI' | 'Anthropic';
  model: string;
  prompt: string;
}

export const MODELS = {
  OpenAI: [
    'gpt-4.1-nano',
    'gpt-4.1-mini',
    'gpt-4.1',
    'gpt-4o',
    'gpt-4o-mini',
    'o4-mini',
    'o3',
    'o3-mini',
    'o3-pro',
    'gpt-4.5-preview',
  ],
  Google: [
    'gemini-2.5-flash-preview-05-20',
    'gemini-2.5-pro-preview-06-05',
    'gemini-2.0-flash',
    'gemini-2.0-flash-lite',
  ],
  Anthropic: [
    'claude-sonnet-4-0',
    'claude-opus-4-0',
    'claude-3-7-sonnet-latest',
    'claude-3-5-sonnet-latest',
  ],
};

export const chatService = {
  /**
   * Send a chat message
   * @param request The chat request
   * @returns Promise with the chat response
   */
  sendMessage: async (request: ChatRequest): Promise<any> => {
    try {
      const response = await api.post('/chat', request);
      return response.data;
    } catch (error) {
      console.error('Error sending message:', error);
      throw error;
    }
  },

  /**
   * Get chat titles
   * @returns Promise with the chat titles
   */
  getChatTitles: async (): Promise<ChatTitle[]> => {
    try {
      const response = await api.get('/chat');
      return response.data.titles;
    } catch (error) {
      console.error('Error getting chat titles:', error);
      throw error;
    }
  },

  /**
   * Get chat history
   * @param id The chat ID
   * @returns Promise with the chat history
   */
  getChatHistory: async (id: number): Promise<ChatHistory> => {
    try {
      const response = await api.get(`/chat/${id}`);
      return response.data;
    } catch (error) {
      console.error('Error getting chat history:', error);
      throw error;
    }
  },
};