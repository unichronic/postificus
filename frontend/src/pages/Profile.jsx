import React, { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { MapPin, Link as LinkIcon, Calendar } from 'lucide-react';
import { Link } from 'react-router-dom';

const Profile = () => {
    const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8080';
    const [profile, setProfile] = useState({
        fullName: '',
        username: '',
        headline: '',
        bio: '',
        location: '',
        website: '',
        publicEmail: '',
        skills: [],
    });

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
                // Keep defaults on failure
            }
        };

        loadProfile();
    }, [apiBase]);

    return (
        <div className="max-w-7xl mx-auto space-y-14 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-6 border-b border-gray-200 pb-10">
                <div className="flex items-center gap-4">
                    <div className="w-28 h-28 rounded-full border border-gray-200 bg-white overflow-hidden">
                        <img src="https://github.com/shadcn.png" alt="Profile" loading="lazy" className="w-full h-full object-cover" />
                    </div>
                    <div>
                        <h1 className="text-5xl font-semibold text-gray-900 font-heading">
                            {profile.fullName || 'Your Name'}
                        </h1>
                        <p className="text-gray-500 text-lg leading-relaxed">
                            {profile.headline || 'Your headline'}
                        </p>
                    </div>
                </div>
                <Link to="/settings">
                    <Button className="bg-brand hover:bg-brand-dark text-white rounded-full px-7 py-3 text-base">
                        Edit Profile
                    </Button>
                </Link>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-10">
                <div className="space-y-6">
                    <Card className="border-gray-200/60 shadow-none bg-white/70">
                        <CardContent className="pt-6 space-y-5 text-lg text-gray-700">
                            {profile.location && (
                                <div className="flex items-center gap-3">
                                    <MapPin className="w-5 h-5 text-brand" />
                                    <span>{profile.location}</span>
                                </div>
                            )}
                            {profile.website && (
                                <div className="flex items-center gap-3">
                                    <LinkIcon className="w-5 h-5 text-brand" />
                                    <a href={profile.website} className="hover:text-brand transition-colors" target="_blank" rel="noreferrer">
                                        {profile.website}
                                    </a>
                                </div>
                            )}
                            {profile.publicEmail && (
                                <div className="flex items-center gap-3">
                                    <LinkIcon className="w-5 h-5 text-brand" />
                                    <a href={`mailto:${profile.publicEmail}`} className="hover:text-brand transition-colors">
                                        {profile.publicEmail}
                                    </a>
                                </div>
                            )}
                            <div className="flex items-center gap-3">
                                <Calendar className="w-5 h-5 text-brand" />
                                <span>Joined October 2023</span>
                            </div>
                        </CardContent>
                    </Card>
                </div>

                <div className="md:col-span-2 space-y-6">
                    <Card className="border-gray-200/60 shadow-none bg-white/70">
                        <CardHeader>
                            <CardTitle className="text-3xl font-semibold text-gray-900 font-heading">About</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <p className="text-gray-600 leading-relaxed text-lg">
                                {profile.bio || 'Add a short bio in Settings to introduce yourself.'}
                            </p>
                        </CardContent>
                    </Card>

                    <Card className="border-gray-200/60 shadow-none bg-white/70">
                        <CardHeader>
                            <CardTitle className="text-3xl font-semibold text-gray-900 font-heading">Skills</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="flex flex-wrap gap-2">
                                {(profile.skills && profile.skills.length > 0) ? (
                                    profile.skills.map((skill) => (
                                        <span key={skill} className="px-5 py-2 rounded-full border border-gray-200 text-lg text-gray-700 bg-gray-50">
                                            {skill}
                                        </span>
                                    ))
                                ) : (
                                    <span className="text-gray-500 text-base">Add skills in Settings to show them here.</span>
                                )}
                            </div>
                        </CardContent>
                    </Card>
                </div>
            </div>
        </div>
    );
};

export default Profile;
