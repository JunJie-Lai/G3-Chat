'use client';

import {useEffect} from 'react';
import {useRouter} from 'next/navigation';
import {useUserStore} from '@/utils/store';
import {initializeStores} from '@/utils/store';
import {authService} from '@/services/authService';

export default function Home() {
    const router = useRouter();
    const {isAuthenticated} = useUserStore();

    useEffect(() => {
        // Initialize stores from localStorage
        initializeStores();

        // Use a timeout to ensure the stores are initialized before redirecting
        const timer = setTimeout(() => {
            // Get the latest authentication state
            const isAuth = authService.isAuthenticated();
            if (isAuth) {
                router.push('/chat');
            } else {
                router.push('/auth');
            }
        }, 100);

        return () => clearTimeout(timer);
    }, [router]);

    return (
        <div className="flex items-center justify-center min-h-screen">
            <div className="text-center">
                <h1 className="text-2xl font-bold mb-4">G3.Chat</h1>
                <p className="text-gray-500 dark:text-gray-400">Redirecting...</p>
            </div>
        </div>
    );
}
