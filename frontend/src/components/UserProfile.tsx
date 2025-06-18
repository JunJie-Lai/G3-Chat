import { useUserStore } from '../utils/store';
import { authService } from '../services/authService';
import { useRouter } from 'next/navigation';

export default function UserProfile() {
  const { user, isAuthenticated, logout } = useUserStore();
  const router = useRouter();

  const handleLogout = () => {
    authService.logout();
    logout();
    router.push('/');
  };

  if (!isAuthenticated || !user) {
    return (
      <div className="flex items-center space-x-2">
        <button
          onClick={() => router.push('/auth')}
          className="px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-blue-600 dark:hover:bg-blue-700"
        >
          Sign In
        </button>
      </div>
    );
  }

  return (
    <div className="relative group">
      <button className="flex items-center space-x-2 focus:outline-none">
        <div className="w-8 h-8 overflow-hidden rounded-full">
          {user.picture ? (
            <img
              src={user.picture}
              alt={user.name}
              className="object-cover w-full h-full"
            />
          ) : (
            <div className="flex items-center justify-center w-full h-full text-white bg-blue-500">
              {user.name}
            </div>
          )}
        </div>
        <span className="text-sm font-medium text-gray-700 dark:text-gray-200">
          {user.name}
        </span>
      </button>

      <div className="absolute right-0 invisible mt-2 w-48 bg-white rounded-md shadow-lg dark:bg-gray-800 group-hover:visible">
        <div className="py-1">
          <div className="px-4 py-2 text-sm text-gray-700 dark:text-gray-200">
            <div className="font-medium">{user.name}</div>
            <div className="text-xs text-gray-500 dark:text-gray-400 truncate">
              {user.email}
            </div>
          </div>
          <hr className="border-gray-200 dark:border-gray-700" />
          <button
            onClick={() => router.push('/settings')}
            className="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700"
          >
            Settings
          </button>
          <button
            onClick={handleLogout}
            className="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-100 dark:text-red-400 dark:hover:bg-gray-700"
          >
            Sign out
          </button>
        </div>
      </div>
    </div>
  );
}