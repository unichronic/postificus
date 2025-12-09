import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { MapPin, Link as LinkIcon, Calendar } from 'lucide-react';

const Profile = () => {
    return (
        <div className="max-w-4xl mx-auto space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
            {/* Profile Header */}
            <div className="relative mb-20">
                <div className="h-48 bg-gradient-to-r from-magical-sky via-magical-violet to-magical-fuchsia rounded-2xl shadow-md"></div>
                <div className="absolute -bottom-16 left-8 flex items-end gap-6">
                    <div className="w-32 h-32 rounded-full border-4 border-white shadow-lg bg-white overflow-hidden">
                        <img src="https://github.com/shadcn.png" alt="Profile" className="w-full h-full object-cover" />
                    </div>
                    <div className="mb-2">
                        <h1 className="text-3xl font-bold text-gray-800">Jane Doe</h1>
                        <p className="text-gray-500">Senior UI/UX Designer</p>
                    </div>
                </div>
                <div className="absolute -bottom-12 right-8">
                    <Button className="bg-magical-violet hover:bg-magical-fuchsia text-white rounded-full shadow-md transition-all hover:scale-105">
                        Edit Profile
                    </Button>
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-8 pt-8">
                {/* Left Column */}
                <div className="space-y-6">
                    <Card className="border-none shadow-md bg-white/50 backdrop-blur-sm">
                        <CardContent className="pt-6 space-y-4">
                            <div className="flex items-center gap-3 text-gray-600">
                                <MapPin className="w-5 h-5 text-magical-violet" />
                                <span>San Francisco, CA</span>
                            </div>
                            <div className="flex items-center gap-3 text-gray-600">
                                <LinkIcon className="w-5 h-5 text-magical-violet" />
                                <a href="#" className="hover:text-magical-fuchsia transition-colors">janedoe.design</a>
                            </div>
                            <div className="flex items-center gap-3 text-gray-600">
                                <Calendar className="w-5 h-5 text-magical-violet" />
                                <span>Joined October 2023</span>
                            </div>
                        </CardContent>
                    </Card>
                </div>

                {/* Right Column */}
                <div className="md:col-span-2 space-y-6">
                    <Card className="border-none shadow-md bg-white/50 backdrop-blur-sm">
                        <CardHeader>
                            <CardTitle>About</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <p className="text-gray-600 leading-relaxed">
                                Passionate about creating beautiful and functional user experiences.
                                With 7 years of experience in the field, I love blending aesthetics with usability.
                                Currently working on the next generation of content creation tools.
                            </p>
                        </CardContent>
                    </Card>

                    <Card className="border-none shadow-md bg-white/50 backdrop-blur-sm">
                        <CardHeader>
                            <CardTitle>Skills</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="flex flex-wrap gap-2">
                                {['UI Design', 'UX Research', 'React', 'Tailwind CSS', 'Figma', 'Prototyping'].map((skill) => (
                                    <span key={skill} className="px-3 py-1 bg-magical-violet/10 text-magical-violet rounded-full text-sm font-medium">
                                        {skill}
                                    </span>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                </div>
            </div>
        </div>
    );
};

export default Profile;
