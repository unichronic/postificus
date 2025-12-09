import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { User, Link as LinkIcon, Save, CheckCircle, Edit2 } from 'lucide-react';

const Settings = () => {
    const [activeTab, setActiveTab] = useState("connections");
    const [devToToken, setDevToToken] = useState('');
    const [isSaving, setIsSaving] = useState(false);
    const [isConnected, setIsConnected] = useState(false);
    const [isEditing, setIsEditing] = useState(false);

    useEffect(() => {
        checkConnection();
    }, []);

    const checkConnection = async () => {
        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/config/devto`);
            const data = await response.json();
            setIsConnected(data.connected);
        } catch (e) {
            console.error("Failed to check connection", e);
        }
    };

    const handleSaveDevTo = async () => {
        setIsSaving(true);
        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/config/devto`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ token: devToToken })
            });
            const data = await response.json();
            if (data.status === 'saved') {
                alert('Token saved successfully!');
                setIsConnected(true);
                setIsEditing(false);
                setDevToToken('');
            } else {
                alert('Failed to save token: ' + data.error);
            }
        } catch (e) {
            console.error(e);
            alert('Error connecting to backend');
        } finally {
            setIsSaving(false);
        }
    };

    return (
        <div className="max-w-5xl mx-auto animate-in fade-in slide-in-from-bottom-4 duration-500">
            <h1 className="text-3xl font-bold mb-8 bg-gradient-to-r from-magical-fuchsia to-magical-violet bg-clip-text text-transparent">
                Settings
            </h1>

            <div className="flex flex-col md:flex-row gap-8">
                {/* Sidebar */}
                <div className="w-full md:w-64 space-y-2">
                    <Button
                        variant={activeTab === "profile" ? "secondary" : "ghost"}
                        className={`w-full justify-start gap-2 ${activeTab === "profile" ? "bg-magical-violet/10 text-magical-violet" : "text-gray-600"}`}
                        onClick={() => setActiveTab("profile")}
                    >
                        <User className="w-4 h-4" />
                        Profile Settings
                    </Button>
                    <Button
                        variant={activeTab === "connections" ? "secondary" : "ghost"}
                        className={`w-full justify-start gap-2 ${activeTab === "connections" ? "bg-magical-violet/10 text-magical-violet" : "text-gray-600"}`}
                        onClick={() => setActiveTab("connections")}
                    >
                        <LinkIcon className="w-4 h-4" />
                        Connections
                    </Button>
                </div>

                {/* Content */}
                <div className="flex-1">
                    {activeTab === "profile" && (
                        <Card className="border-none shadow-md bg-white/50 backdrop-blur-sm">
                            <CardHeader>
                                <CardTitle>Profile Settings</CardTitle>
                                <CardDescription>Manage your public profile information.</CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div className="space-y-2">
                                    <label className="text-sm font-medium text-gray-700">Display Name</label>
                                    <Input placeholder="Jane Doe" />
                                </div>
                                <div className="space-y-2">
                                    <label className="text-sm font-medium text-gray-700">Bio</label>
                                    <Input placeholder="Senior UI/UX Designer" />
                                </div>
                                <Button className="bg-magical-violet hover:bg-magical-fuchsia text-white">
                                    Save Changes
                                </Button>
                            </CardContent>
                        </Card>
                    )}

                    {activeTab === "connections" && (
                        <Card className="border-none shadow-md bg-white/50 backdrop-blur-sm">
                            <CardHeader>
                                <CardTitle>Connections</CardTitle>
                                <CardDescription>Manage your external account connections.</CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-6">
                                <div className="p-4 border border-gray-200 rounded-xl bg-white/50">
                                    <div className="flex items-center justify-between mb-4">
                                        <div className="flex items-center gap-4">
                                            <div className="w-10 h-10 bg-black rounded-lg flex items-center justify-center">
                                                <span className="text-white font-bold text-xs">DEV</span>
                                            </div>
                                            <div>
                                                <h3 className="font-semibold text-gray-800">Dev.to</h3>
                                                <p className="text-sm text-gray-500">Connect your Dev.to account to publish posts.</p>
                                            </div>
                                        </div>
                                        {isConnected && !isEditing && (
                                            <div className="flex items-center gap-2 text-green-600 bg-green-50 px-3 py-1 rounded-full text-sm font-medium">
                                                <CheckCircle className="w-4 h-4" />
                                                Connected
                                            </div>
                                        )}
                                    </div>

                                    {isConnected && !isEditing ? (
                                        <div className="flex justify-end">
                                            <Button
                                                variant="outline"
                                                onClick={() => setIsEditing(true)}
                                                className="gap-2"
                                            >
                                                <Edit2 className="w-4 h-4" />
                                                Edit Token
                                            </Button>
                                        </div>
                                    ) : (
                                        <div className="space-y-3 animate-in fade-in slide-in-from-top-2">
                                            <label className="text-sm font-medium text-gray-700">
                                                API Key / Token
                                                <span className="ml-2 text-xs text-gray-400 font-normal">(remember_user_token)</span>
                                            </label>
                                            <div className="flex gap-2">
                                                <Input
                                                    type="password"
                                                    value={devToToken}
                                                    onChange={(e) => setDevToToken(e.target.value)}
                                                    placeholder="Paste your token here..."
                                                    className="flex-1"
                                                />
                                                <Button
                                                    onClick={handleSaveDevTo}
                                                    disabled={isSaving}
                                                    className="bg-magical-violet hover:bg-magical-fuchsia text-white"
                                                >
                                                    {isSaving ? 'Saving...' : 'Save'}
                                                </Button>
                                                {isConnected && (
                                                    <Button
                                                        variant="ghost"
                                                        onClick={() => setIsEditing(false)}
                                                    >
                                                        Cancel
                                                    </Button>
                                                )}
                                            </div>
                                            <p className="text-xs text-gray-500">
                                                Your token is stored locally and used only to authenticate with Dev.to.
                                            </p>
                                        </div>
                                    )}
                                </div>
                            </CardContent>
                        </Card>
                    )}
                </div>
            </div>
        </div>
    );
};

export default Settings;
