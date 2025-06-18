import { ChatMessage } from '../services/chatService';

interface MessageDisplayProps {
  messages: ChatMessage[];
}

export default function MessageDisplay({ messages }: MessageDisplayProps) {
  return (
    <div className="flex-1 overflow-y-auto p-4 space-y-4">
      {messages.length === 0 ? (
        <div className="flex items-center justify-center h-full">
          <div className="text-center text-gray-500 dark:text-gray-400">
            <h3 className="text-xl font-semibold mb-2">Start a conversation</h3>
            <p>Send a message to begin chatting with the AI</p>
          </div>
        </div>
      ) : (
        messages.map((message, index) => (
          <div
            key={index}
            className={`flex ${
              message.role === 'human' ? 'justify-end' : 'justify-start'
            }`}
          >
            <div
              className={`max-w-3/4 p-3 rounded-lg ${
                message.role === 'human'
                  ? 'bg-blue-500 text-white rounded-tr-none'
                  : 'bg-gray-200 dark:bg-gray-700 dark:text-white rounded-tl-none'
              }`}
            >
              <div className="whitespace-pre-wrap">{message.text}</div>
            </div>
          </div>
        ))
      )}
    </div>
  );
}