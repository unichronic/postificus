import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { User, Link as LinkIcon, Save, CheckCircle, Edit2 } from 'lucide-react';

const Settings = () => {
    const [activeTab, setActiveTab] = useState("connections");



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
                                {/* Dev.to */}
                                <ConnectionCard
                                    platform="devto"
                                    name="Dev.to"
                                    description="Connect your Dev.to account to publish posts."
                                    icon="DEV"
                                    fields={[
                                        { key: "token", label: "API Key / Token", type: "password", placeholder: "Paste your token here..." }
                                    ]}
                                />

                                {/* Medium */}
                                <ConnectionCard
                                    platform="medium"
                                    name="Medium"
                                    description="Connect your Medium account (uid, sid, xsrf)."
                                    icon="M"
                                    fields={[
                                        { key: "uid", label: "UID", type: "password", placeholder: "Medium UID cookie" },
                                        { key: "sid", label: "SID", type: "password", placeholder: "Medium SID cookie" },
                                        { key: "xsrf", label: "XSRF", type: "password", placeholder: "Medium XSRF cookie" }
                                    ]}
                                />

                                {/* LinkedIn */}
                                <ConnectionCard
                                    platform="linkedin"
                                    name="LinkedIn"
                                    description="Connect your LinkedIn account (li_at cookie)."
                                    icon="IN"
                                    fields={[
                                        { key: "li_at", label: "li_at Cookie", type: "password", placeholder: "Paste your li_at cookie..." }
                                    ]}
                                />
                            </CardContent>
                        </Card>
                    )}
                </div>
            </div>
        </div>
    );
};

// Reusable Connection Component
const ConnectionCard = ({ platform, name, description, icon, fields }) => {
    const [isConnected, setIsConnected] = useState(false);
    const [isEditing, setIsEditing] = useState(false);
    const [isSaving, setIsSaving] = useState(false);
    const [formData, setFormData] = useState({});

    useEffect(() => {
        checkConnection();
    }, []);

    const checkConnection = async () => {
        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/settings/credentials/${platform}`);
            const data = await response.json();
            setIsConnected(data.connected);
        } catch (e) {
            console.error(`Failed to check ${platform} connection`, e);
        }
    };

    const handleSave = async () => {
        setIsSaving(true);
        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/settings/credentials`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ platform, credentials: formData })
            });
            const data = await response.json();
            if (data.status === 'saved') {
                alert(`${name} connected successfully!`);
                setIsConnected(true);
                setIsEditing(false);
                setFormData({});
            } else {
                alert('Failed to save: ' + data.error);
            }
        } catch (e) {
            console.error(e);
            alert('Error connecting to backend');
        } finally {
            setIsSaving(false);
        }
    };

    const handleChange = (key, value) => {
        setFormData(prev => ({ ...prev, [key]: value }));
    };

    return (
        <div className="p-4 border border-gray-200 rounded-xl bg-white/50">
            <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-4">
                    <div className="w-10 h-10 bg-black rounded-lg flex items-center justify-center">
                        <span className="text-white font-bold text-xs">{icon}</span>
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-800">{name}</h3>
                        <p className="text-sm text-gray-500">{description}</p>
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
                        Edit Credentials
                    </Button>
                </div>
            ) : (
                <div className="space-y-3 animate-in fade-in slide-in-from-top-2">
                    {fields.map((field) => (
                        <div key={field.key} className="space-y-1">
                            <label className="text-sm font-medium text-gray-700">
                                {field.label}
                            </label>
                            <Input
                                type={field.type}
                                value={formData[field.key] || ''}
                                onChange={(e) => handleChange(field.key, e.target.value)}
                                placeholder={field.placeholder}
                            />
                        </div>
                    ))}
                    <div className="flex gap-2 pt-2">
                        <Button
                            onClick={handleSave}
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
                        Credentials are stored securely and used only for automation.
                    </p>
                </div>
            )}
        </div>
    );
};


export default Settings;
