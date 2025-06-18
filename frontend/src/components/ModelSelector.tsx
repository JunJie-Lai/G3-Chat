import { useState } from 'react';
import { useChatStore } from '../utils/store';
import { MODELS } from '../services/chatService';

export default function ModelSelector() {
  const { modelType, model, setModelType, setModel } = useChatStore();
  const [isOpen, setIsOpen] = useState(false);

  const handleModelTypeChange = (type: 'Google' | 'OpenAI' | 'Anthropic') => {
    setModelType(type);
    // Set default model for the selected type
    setModel(MODELS[type][0]);
    setIsOpen(false);
  };

  const handleModelChange = (selectedModel: string) => {
    setModel(selectedModel);
    setIsOpen(false);
  };

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center justify-between w-full px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-800 dark:text-gray-200 dark:border-gray-600 dark:hover:bg-gray-700"
      >
        <span>{modelType}: {model}</span>
        <svg
          className="w-5 h-5 ml-2 -mr-1"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 20 20"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            fillRule="evenodd"
            d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
            clipRule="evenodd"
          />
        </svg>
      </button>

      {isOpen && (
        <div className="absolute z-10 w-full mt-1 bg-white rounded-md shadow-lg dark:bg-gray-800 dark:text-gray-200 dark:border-gray-600">
          <div className="py-1">
            <div className="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400">
              Provider
            </div>
            {Object.keys(MODELS).map((type) => (
              <button
                key={type}
                onClick={() => handleModelTypeChange(type as 'Google' | 'OpenAI' | 'Anthropic')}
                className={`block w-full text-left px-4 py-2 text-sm ${
                  modelType === type
                    ? 'bg-blue-100 text-blue-900 dark:bg-blue-900 dark:text-blue-100'
                    : 'text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700'
                }`}
              >
                {type}
              </button>
            ))}
          </div>
          <div className="border-t border-gray-200 dark:border-gray-600">
            <div className="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400">
              Models
            </div>
            {MODELS[modelType].map((modelName) => (
              <button
                key={modelName}
                onClick={() => handleModelChange(modelName)}
                className={`block w-full text-left px-4 py-2 text-sm ${
                  model === modelName
                    ? 'bg-blue-100 text-blue-900 dark:bg-blue-900 dark:text-blue-100'
                    : 'text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700'
                }`}
              >
                {modelName}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}