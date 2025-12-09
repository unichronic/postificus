import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Eye, ThumbsUp, MessageSquare, TrendingUp } from 'lucide-react';

const Dashboard = () => {
    // Mock data
    const analytics = [
        { title: "Total Views", value: "12,345", change: "+12%", icon: Eye },
        { title: "Total Likes", value: "1,234", change: "+5%", icon: ThumbsUp },
        { title: "Comments", value: "456", change: "+8%", icon: MessageSquare },
        { title: "Engagement", value: "8.5%", change: "+2%", icon: TrendingUp },
    ];

    const posts = [
        { id: 1, title: "How to build a magical UI", status: "Published", date: "2023-10-25", views: 1200 },
        { id: 2, title: "Understanding React Hooks", status: "Draft", date: "2023-10-28", views: 0 },
        { id: 3, title: "The future of AI in coding", status: "Published", date: "2023-11-01", views: 850 },
    ];

    return (
        <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <div>
                <h1 className="text-3xl font-bold bg-gradient-to-r from-magical-fuchsia to-magical-violet bg-clip-text text-transparent">
                    Dashboard
                </h1>
                <p className="text-gray-500 mt-2">Welcome back! Here's how your content is performing.</p>
            </div>

            {/* Analytics Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {analytics.map((item, index) => (
                    <Card key={index} className="border-none shadow-md hover:shadow-lg transition-shadow bg-white/50 backdrop-blur-sm">
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium text-gray-600">
                                {item.title}
                            </CardTitle>
                            <item.icon className="h-4 w-4 text-magical-violet" />
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold text-gray-800">{item.value}</div>
                            <p className="text-xs text-green-500 font-medium mt-1">
                                {item.change} from last month
                            </p>
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* Posts List */}
            <div className="bg-white/50 backdrop-blur-sm rounded-xl shadow-md border border-white/20 overflow-hidden">
                <div className="p-6 border-b border-gray-100">
                    <h2 className="text-xl font-semibold text-gray-800">Recent Posts</h2>
                </div>
                <div className="divide-y divide-gray-100">
                    {posts.map((post) => (
                        <div key={post.id} className="p-6 flex items-center justify-between hover:bg-white/60 transition-colors">
                            <div>
                                <h3 className="font-medium text-gray-800 hover:text-magical-violet cursor-pointer transition-colors">
                                    {post.title}
                                </h3>
                                <p className="text-sm text-gray-500 mt-1">
                                    Last edited on {post.date}
                                </p>
                            </div>
                            <div className="flex items-center gap-4">
                                <div className="text-sm text-gray-500 flex items-center gap-1">
                                    <Eye className="w-4 h-4" /> {post.views}
                                </div>
                                <Badge variant={post.status === 'Published' ? 'default' : 'secondary'}
                                    className={`${post.status === 'Published' ? 'bg-green-100 text-green-700 hover:bg-green-200' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}`}>
                                    {post.status}
                                </Badge>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
