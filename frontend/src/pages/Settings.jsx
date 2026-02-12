import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { User, Link as LinkIcon, CheckCircle, Edit2 } from 'lucide-react';
import ConnectPlatformButton from '@/components/ConnectPlatformButton';
import TagInput from '@/components/TagInput';

const Settings = () => {
    const [activeTab, setActiveTab] = useState("profile");
    const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8080';
    const emptyProfile = {
        fullName: '',
        username: '',
        headline: '',
        bio: '',
        location: '',
        website: '',
        publicEmail: '',
        skills: [],
    };
    const [profile, setProfile] = useState(emptyProfile);
    const [profileStatus, setProfileStatus] = useState('idle');
    const [profileError, setProfileError] = useState('');

    useEffect(() => {
        const loadProfile = async () => {
            try {
                const response = await fetch(`${apiBase}/api/profile`);
                if (!response.ok) {
                    throw new Error('Failed to load profile');
                }
                const data = await response.json();
                setProfile({
                    fullName: data.full_name || '',
                    username: data.username || '',
                    headline: data.headline || '',
                    bio: data.bio || '',
                    location: data.location || '',
                    website: data.website || '',
                    publicEmail: data.public_email || '',
                    skills: Array.isArray(data.skills) ? data.skills : [],
                });
            } catch (e) {
                setProfileError('Unable to load profile right now.');
            }
        };

        loadProfile();
    }, [apiBase]);

    const handleProfileSave = async () => {
        setProfileStatus('saving');
        setProfileError('');
        try {
            const response = await fetch(`${apiBase}/api/profile`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    full_name: profile.fullName,
                    username: profile.username,
                    headline: profile.headline,
                    bio: profile.bio,
                    location: profile.location,
                    website: profile.website,
                    public_email: profile.publicEmail,
                    skills: profile.skills,
                })
            });
            if (!response.ok) {
                const data = await response.json().catch(() => ({}));
                throw new Error(data.error || 'Failed to save profile');
            }
            setProfileStatus('saved');
            setTimeout(() => setProfileStatus('idle'), 2000);
        } catch (e) {
            setProfileStatus('error');
            setProfileError(e.message || 'Failed to save profile');
        }
    };

    return (
        <div className="max-w-7xl mx-auto animate-in fade-in slide-in-from-bottom-4 duration-500 space-y-12">
            <div>
                <h1 className="text-5xl font-semibold text-gray-900 font-heading">Settings</h1>
                <p className="text-gray-500 mt-4 text-lg leading-relaxed">Manage your profile and publishing connections.</p>
            </div>

            <div className="flex flex-col md:flex-row gap-10">
                {/* Sidebar */}
                <div className="w-full md:w-72 space-y-3">
                    <button
                        type="button"
                        className={`w-full flex items-center gap-2 rounded-lg border px-5 py-4 text-lg transition-colors ${activeTab === "profile" ? "border-brand/30 text-brand bg-brand/5" : "border-gray-200 text-gray-700 hover:bg-gray-50"}`}
                        onClick={() => setActiveTab("profile")}
                    >
                        <User className="w-4 h-4" />
                        Profile Settings
                    </button>
                    <button
                        type="button"
                        className={`w-full flex items-center gap-2 rounded-lg border px-5 py-4 text-lg transition-colors ${activeTab === "connections" ? "border-brand/30 text-brand bg-brand/5" : "border-gray-200 text-gray-700 hover:bg-gray-50"}`}
                        onClick={() => setActiveTab("connections")}
                    >
                        <LinkIcon className="w-4 h-4" />
                        Connections
                    </button>
                </div>

                {/* Content */}
                <div className="flex-1 space-y-6">
                    {activeTab === "profile" && (
                        <Card className="border-gray-200/60 shadow-none bg-white/70">
                            <CardHeader>
                                <CardTitle className="text-3xl font-semibold text-gray-900 font-heading">Profile Settings</CardTitle>
                                <CardDescription className="text-lg">Update the information shown on your public profile.</CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-6">
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                    <div className="space-y-2">
                                        <label className="text-base font-medium text-gray-700">Full Name</label>
                                        <Input
                                            placeholder="Jane Doe"
                                            className="h-12 text-lg"
                                            value={profile.fullName}
                                            onChange={(e) => setProfile((prev) => ({ ...prev, fullName: e.target.value }))}
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <label className="text-base font-medium text-gray-700">Headline</label>
                                        <Input
                                            placeholder="Senior UI/UX Designer"
                                            className="h-12 text-lg"
                                            value={profile.headline}
                                            onChange={(e) => setProfile((prev) => ({ ...prev, headline: e.target.value }))}
                                        />
                                    </div>
                                </div>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                    <div className="space-y-2">
                                        <label className="text-base font-medium text-gray-700">Username</label>
                                        <Input
                                            placeholder="janedoe"
                                            className="h-12 text-lg"
                                            value={profile.username}
                                            onChange={(e) => setProfile((prev) => ({ ...prev, username: e.target.value }))}
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <label className="text-base font-medium text-gray-700">Public Email</label>
                                        <Input
                                            type="email"
                                            placeholder="jane@postificus.com"
                                            className="h-12 text-lg"
                                            value={profile.publicEmail}
                                            onChange={(e) => setProfile((prev) => ({ ...prev, publicEmail: e.target.value }))}
                                        />
                                    </div>
                                </div>
                                <div className="space-y-2">
                                    <label className="text-base font-medium text-gray-700">Bio</label>
                                    <textarea
                                        rows={3}
                                        placeholder="Short summary about your work..."
                                        className="w-full rounded-md border border-gray-200 px-4 py-3 text-lg text-gray-900 focus:outline-none focus:ring-2 focus:ring-brand/30"
                                        value={profile.bio}
                                        onChange={(e) => setProfile((prev) => ({ ...prev, bio: e.target.value }))}
                                    />
                                </div>
                                <div className="space-y-3">
                                    <label className="text-base font-medium text-gray-700">Skills</label>
                                    <TagInput
                                        tags={profile.skills}
                                        setTags={(nextSkills) => setProfile((prev) => ({ ...prev, skills: nextSkills }))}
                                        maxTags={10}
                                        placeholder="Add a skill..."
                                    />
                                    <p className="text-sm text-gray-500">These appear on your public profile.</p>
                                </div>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                    <div className="space-y-2">
                                        <label className="text-base font-medium text-gray-700">Location</label>
                                        <Input
                                            placeholder="San Francisco, CA"
                                            className="h-12 text-lg"
                                            value={profile.location}
                                            onChange={(e) => setProfile((prev) => ({ ...prev, location: e.target.value }))}
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <label className="text-base font-medium text-gray-700">Website</label>
                                        <Input
                                            placeholder="https://janedoe.design"
                                            className="h-12 text-lg"
                                            value={profile.website}
                                            onChange={(e) => setProfile((prev) => ({ ...prev, website: e.target.value }))}
                                        />
                                    </div>
                                </div>
                                <div className="flex flex-wrap items-center gap-3">
                                    <Button
                                        className="bg-brand hover:bg-brand-dark text-white w-fit px-7 py-3 text-base"
                                        onClick={handleProfileSave}
                                        disabled={profileStatus === 'saving'}
                                    >
                                        {profileStatus === 'saving' ? 'Saving...' : 'Save Changes'}
                                    </Button>
                                    {profileStatus === 'saved' && (
                                        <span className="text-sm text-brand">Saved</span>
                                    )}
                                    {profileStatus === 'error' && profileError && (
                                        <span className="text-sm text-brand-dark">{profileError}</span>
                                    )}
                                </div>
                            </CardContent>
                        </Card>
                    )}

                    {activeTab === "connections" && (
                        <Card className="border-gray-200/60 shadow-none bg-white/70">
                            <CardHeader>
                                <CardTitle className="text-3xl font-semibold text-gray-900 font-heading">Connections</CardTitle>
                                <CardDescription className="text-lg">Store your platform credentials for publishing.</CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-6">
                                {/* Dev.to */}
                                <ConnectionCard
                                    platform="devto"
                                    name="Dev.to"
                                    description="Connect your Dev.to account to publish posts."
                                    icon="DEV"
                                    automated={true}
                                />

                                {/* Medium */}
                                <ConnectionCard
                                    platform="medium"
                                    name="Medium"
                                    description="Connect your Medium account to publish posts."
                                    icon="M"
                                    automated={true}
                                />

                                {/* LinkedIn */}
                                <ConnectionCard
                                    platform="linkedin"
                                    name="LinkedIn"
                                    description="Connect your LinkedIn account (Coming Soon)."
                                    icon="IN"
                                    automated={true}
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
const ConnectionCard = ({ platform, name, description, icon, fields, automated }) => {
    const [isConnected, setIsConnected] = useState(false);
    const [accountName, setAccountName] = useState('');
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
            if (data.connected && data.account) {
                setAccountName(data.account);
            }
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
        <div className="p-5 border border-gray-200/60 rounded-xl bg-white/70">
            <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-4">
                    <div className="w-12 h-12 border border-gray-200 rounded-lg flex items-center justify-center">
                        <span className="text-gray-700 font-semibold text-sm">{icon}</span>
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 text-base">{name}</h3>
                        <p className="text-base text-gray-500">{description}</p>
                    </div>
                </div>
                {isConnected && !isEditing && (
                    <div className="flex flex-col items-end gap-1">
                        <div className="flex items-center gap-2 text-brand border border-brand/30 bg-brand/5 px-3 py-1 rounded-full text-sm font-medium">
                            <CheckCircle className="w-4 h-4" />
                            Connected
                        </div>
                        {accountName && (
                            <p className="text-sm text-gray-500 font-medium">
                                as <span className="text-gray-900">{accountName}</span>
                            </p>
                        )}
                    </div>
                )}
            </div>

            {isConnected && !isEditing ? (
                <div className="flex justify-end gap-2">
                    {automated ? (
                        <ConnectPlatformButton
                            platform={platform}
                            name={name}
                            isConnected={true}
                            onConnected={() => {
                                alert(`${name} re-connected successfully!`);
                                checkConnection(); // Refresh status
                            }}
                            className="bg-brand hover:bg-brand-dark text-white px-4 py-2"
                        />
                    ) : (
                        <Button
                            variant="outline"
                            onClick={() => setIsEditing(true)}
                            className="gap-2 border-gray-200 text-gray-700 hover:bg-gray-50 text-base px-4 py-2.5"
                        >
                            <Edit2 className="w-4 h-4" />
                            Edit Credentials
                        </Button>
                    )}
                </div>
            ) : (
                <div className="space-y-4 animate-in fade-in slide-in-from-top-2">
                    {automated ? (
                        <div className="pt-2">
                            <ConnectPlatformButton
                                platform={platform}
                                name={name}
                                onConnected={() => {
                                    alert(`${name} connected successfully!`);
                                    setIsConnected(true);
                                }}
                                className="bg-brand hover:bg-brand-dark text-white px-5 py-2.5 text-base w-full sm:w-auto"
                            />
                            <p className="text-sm text-gray-500 mt-3">
                                A browser window will open to log you in securely.
                            </p>
                        </div>
                    ) : (
                        <>
                            {fields && fields.map((field) => (
                                <div key={field.key} className="space-y-1">
                                    <label className="text-base font-medium text-gray-700">
                                        {field.label}
                                    </label>
                                    <Input
                                        type={field.type}
                                        value={formData[field.key] || ''}
                                        onChange={(e) => handleChange(field.key, e.target.value)}
                                        placeholder={field.placeholder}
                                        className="h-11 text-base"
                                    />
                                </div>
                            ))}
                            <div className="flex gap-2 pt-2">
                                <Button
                                    onClick={handleSave}
                                    disabled={isSaving}
                                    className="bg-brand hover:bg-brand-dark text-white px-5 py-2.5 text-base"
                                >
                                    {isSaving ? 'Saving...' : 'Save'}
                                </Button>
                                {isConnected && (
                                    <Button
                                        variant="ghost"
                                        onClick={() => setIsEditing(false)}
                                        className="text-base"
                                    >
                                        Cancel
                                    </Button>
                                )}
                            </div>
                            <p className="text-sm text-gray-500">
                                Credentials are stored securely and used only for automation.
                            </p>
                        </>
                    )}
                </div>
            )}
        </div>
    );
};


export default Settings;
