import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
// Assuming you have a toast hook, if not we'll use window.alert fallback or standard hook
// In Settings.jsx it wasn't explicit, but checking imports might help. 
// Settings.jsx didn't import useToast. I'll stick to simple logic or check imports.
// Settings.jsx uses manual alerts. I will try to use a better UI if possible, 
// but consistent with Settings.jsx is safer. I'll use a prop for toast if passed, else alert.

const ConnectPlatformButton = ({ platform, name, onConnected, className, isConnected }) => {
    const [isLoading, setIsLoading] = useState(false);
    // Use env defined in Settings.jsx context or global. Vite exposes import.meta.env.
    const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8080';

    const handleConnect = async () => {
        setIsLoading(true);
        // Show immediate feedback
        // If we had a toast system: toast({ title: "Launching...", description: "Check the popup window" })

        try {
            const response = await fetch(`${apiBase}/api/connect/${platform}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    platform: platform,
                    user_id: 1 // TODO: Get from auth context
                })
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.message || data.error || 'Connection failed');
            }

            // Success
            if (onConnected) onConnected(data);

        } catch (error) {
            console.error("Connection error:", error);
            alert(`Failed to connect to ${name}: ${error.message}`);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Button
            onClick={handleConnect}
            disabled={isLoading}
            className={`gap-2 ${className}`}
        >
            {isLoading ? (
                <>
                    <span className="animate-spin">⏳</span>
                    Waiting for Login...
                </>
            ) : (
                <>
                    {isConnected ? 'Update' : 'Connect'} {name}
                </>
            )}
        </Button>
    );
};

export default ConnectPlatformButton;
