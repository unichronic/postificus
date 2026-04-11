# Frontend Integration Guide: Connect Medium

I have implemented the **Backend-Managed Login Workflow** on the server. Here is how to integrate it into your React frontend.

## 1. Backend Status
✅ **Endpoint:** `POST /api/connect/medium`
✅ **Logic:** Launches a browser window, waits for login, updates DB.
✅ **Response:**
```json
{
  "status": "connected",
  "platform": "medium",
  "message": "Using account: <uid>"
}
```

## 2. React Component (`src/components/ConnectMediumButton.jsx`)

Create this component to handle the interaction. It handles the "Connecting..." state while the backend waits for the user.

```jsx
import React, { useState } from 'react';
import { Button } from '@/components/ui/button'; // Assuming shadcn/ui
import { useToast } from '@/components/ui/use-toast';

export default function ConnectMediumButton({ onConnected }) {
    const [isLoading, setIsLoading] = useState(false);
    const { toast } = useToast();

    const handleConnect = async () => {
        setIsLoading(true);
        toast({
            title: "Launching Login Window...",
            description: "Please log in to Medium in the popup window."
        });

        try {
            // Call our new backend endpoint
            const response = await fetch('/api/connect/medium', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    platform: 'medium',
                    user_id: 1 // TODO: Replace with actual logged-in user ID
                })
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.message || 'Connection failed');
            }

            // Success!
            toast({
                title: "✅ Connected!",
                description: "Medium account linked successfully.",
                variant: "success"
            });
            
            if (onConnected) onConnected(data);

        } catch (error) {
            toast({
                title: "❌ Connection Failed",
                description: error.message,
                variant: "destructive"
            });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Button 
            onClick={handleConnect} 
            disabled={isLoading}
            className="gap-2"
        >
            {isLoading ? (
                <>
                    <span className="animate-spin">⏳</span> 
                    Waiting for Login...
                </>
            ) : (
                <>Connect Medium</>
            )}
        </Button>
    );
}
```

## 3. Integrating into Connections Page (`src/pages/Platforms.jsx`)

```jsx
import ConnectMediumButton from '../components/ConnectMediumButton';

export default function PlatformsPage() {
    const [mediumConnected, setMediumConnected] = useState(false);

    return (
        <div className="p-6">
            <h1 className="text-2xl font-bold mb-6">Connected Platforms</h1>
            
            <div className="border p-4 rounded-lg flex items-center justify-between">
                <div className="flex items-center gap-4">
                    <img src="/medium-logo.png" className="w-10 h-10" />
                    <div>
                        <h3 className="font-semibold">Medium</h3>
                        <p className="text-sm text-gray-500">
                            {mediumConnected ? "✅ Connected" : "Not connected"}
                        </p>
                    </div>
                </div>
                
                {mediumConnected ? (
                    <Button variant="outline" onClick={() => setMediumConnected(false)}>
                        Disconnect
                    </Button>
                ) : (
                    <ConnectMediumButton onConnected={() => setMediumConnected(true)} />
                )}
            </div>
        </div>
    );
}
```

## 4. How to Test
1. Make sure the backend API is running (`./bin/api`).
2. Run your frontend (`npm run dev`).
3. Click "Connect Medium".
4. A separate Chrome window will pop up.
5. Log in to Medium.
6. The window will close, and your button will turn to "Connected"!
