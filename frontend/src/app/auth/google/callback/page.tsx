'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { authService } from '@/services/authService';
import { useUserStore } from '@/utils/store';

export default function GoogleCallbackPage() {
  const router = useRouter();
  const { setUser, setAuthenticated } = useUserStore();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // Redirect to auth page with token
    // This page is no longer needed as the callback is handled directly in the auth page
    router.push('/auth');
  }, [router]);

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="w-full max-w-md p-8 space-y-8 bg-white rounded-lg shadow-md dark:bg-gray-800">
        <div className="text-center">
          <h1 className="text-3xl font-extrabold text-gray-900 dark:text-white">G3.Chat</h1>
          {error ? (
            <div className="mt-4 p-4 text-sm text-red-700 bg-red-100 rounded-lg dark:bg-red-900 dark:text-red-100">
              {error}
              <p className="mt-2">Redirecting to login page...</p>
            </div>
          ) : (
            <p className="mt-4 text-gray-600 dark:text-gray-400">
              Authenticating with Google...
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
