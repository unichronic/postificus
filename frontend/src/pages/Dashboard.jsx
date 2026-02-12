import React, { useState, useEffect, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { RefreshCw, ChevronLeft, ChevronRight } from "lucide-react";

const Dashboard = () => {
    const [posts, setPosts] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [isSyncing, setIsSyncing] = useState(false);
    const [currentPage, setCurrentPage] = useState(1);
    const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8080';
    const itemsPerPage = 10;

    useEffect(() => {
        fetchActivity();
    }, []);

    const fetchActivity = async () => {
        try {
            const res = await fetch(`${apiBase}/api/dashboard/activity?limit=100`);
            const data = await res.json();
            if (data.posts) {
                setPosts(data.posts);
            }
        } catch (err) {
            console.error("Failed to fetch activity", err);
        } finally {
            setIsLoading(false);
        }
    };

    const handleSync = async () => {
        setIsSyncing(true);
        try {
            await fetch(`${apiBase}/api/dashboard/sync`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ platform: 'all' })
            });
            setTimeout(fetchActivity, 4000);
        } catch (err) {
            console.error("Sync failed", err);
        } finally {
            setTimeout(() => setIsSyncing(false), 1200);
        }
    };

    const formatDate = (value) => {
        if (!value) return 'Unknown date';
        const raw = typeof value === 'string' ? value : String(value);
        if (raw.startsWith('0001-01-01')) return 'Unknown date';
        const parsed = new Date(raw.includes('T') ? raw : `${raw}T00:00:00`);
        if (Number.isNaN(parsed.getTime())) return 'Unknown date';
        return parsed.toLocaleDateString('en-US', {
            month: 'long',
            day: 'numeric',
            year: 'numeric'
        });
    };

    const platformLabels = {
        medium: 'Medium',
        devto: 'Dev.to',
        linkedin: 'LinkedIn',
        postificus: 'Postificus',
    };

    const platformClasses = {
        LinkedIn: 'bg-brand/15 text-brand-dark border-brand/25',
        Medium: 'bg-brand/8 text-brand border-brand/20',
        'Dev.to': 'bg-brand/22 text-brand-dark border-brand/30',
        Postificus: 'bg-brand/12 text-brand border-brand/25',
    };

    const displayPosts = useMemo(() => {
        const groups = new Map();
        const normalizeTitle = (title) => (title || '').trim().toLowerCase();
        const toTimestamp = (value) => {
            if (!value) return null;
            const raw = typeof value === 'string' ? value : String(value);
            if (raw.startsWith('0001-01-01')) return null;
            const parsed = new Date(raw);
            const time = parsed.getTime();
            return Number.isNaN(time) ? null : time;
        };

        posts.forEach((post, index) => {
            const title = post.title || 'Untitled post';
            const key = normalizeTitle(title) || `untitled-${post.remote_id || index}`;
            const platformKey = (post.platform || '').toLowerCase();
            const platformLabel = platformLabels[platformKey] || post.platform || 'Unknown';
            const targetLabels = Array.isArray(post.publish_targets)
                ? post.publish_targets.map((target) => platformLabels[target] || target)
                : [];
            const labels = targetLabels.length > 0 ? targetLabels : [platformLabel];
            const timestamp = toTimestamp(post.published_at);
            const existing = groups.get(key);

            if (!existing) {
                groups.set(key, {
                    id: key,
                    title,
                    date: post.published_at,
                    dateValue: timestamp,
                    platforms: labels,
                    orderIndex: index,
                    editDraftId: platformKey === 'postificus' ? post.remote_id : null,
                });
                return;
            }

            labels.forEach((label) => {
                if (!existing.platforms.includes(label)) {
                    existing.platforms.push(label);
                }
            });

            if (timestamp !== null && (existing.dateValue === null || timestamp > existing.dateValue)) {
                existing.dateValue = timestamp;
                existing.date = post.published_at;
            }

            if (platformKey === 'postificus' && post.remote_id) {
                existing.editDraftId = post.remote_id;
            }
        });

        return Array.from(groups.values()).sort((a, b) => {
            if (a.dateValue !== null && b.dateValue !== null && a.dateValue !== b.dateValue) {
                return b.dateValue - a.dateValue;
            }
            return a.orderIndex - b.orderIndex;
        });
    }, [posts]);

    const totalPages = Math.ceil(displayPosts.length / itemsPerPage);

    return (
        <div className="space-y-12 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <div className="relative overflow-hidden rounded-2xl border border-gray-200/60 bg-white/70 p-8 md:p-10 backdrop-blur-sm">
                <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_18%_20%,rgba(148,0,255,0.12),transparent_55%)]" />
                <div className="pointer-events-none absolute right-0 top-0 h-48 w-48 -translate-y-16 translate-x-12 rounded-full bg-brand/10 blur-2xl" />
                <div className="relative">
                    <p className="text-sm uppercase tracking-[0.3em] text-brand/80">Dashboard</p>
                    <h1 className="mt-4 text-5xl font-semibold text-gray-900 font-heading">
                        Your latest posts, everywhere.
                    </h1>
                    <p className="text-gray-500 mt-4 text-lg leading-relaxed max-w-2xl">
                        Track the posts you have published and see where each one is live.
                    </p>
                </div>
            </div>

            <Card className="border-gray-200/60 bg-white/70 shadow-none">
                <CardHeader className="pb-4 flex flex-row items-center justify-between gap-4">
                    <CardTitle className="text-3xl font-semibold text-gray-900 font-heading">Latest Posts</CardTitle>
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={handleSync}
                        disabled={isSyncing}
                        className="border-gray-200/60 text-gray-600 hover:text-gray-800"
                    >
                        <RefreshCw className={`w-3.5 h-3.5 ${isSyncing ? 'animate-spin' : ''}`} />
                        {isSyncing ? 'Syncing' : 'Sync'}
                    </Button>
                </CardHeader>
                <CardContent className="space-y-8">
                    {isLoading ? (
                        <div className="py-8 text-center text-gray-500">Loading latest posts...</div>
                    ) : displayPosts.length === 0 ? (
                        <div className="py-8 text-center text-gray-500">No posts found yet.</div>
                    ) : (
                        displayPosts
                            .slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)
                            .map((post) => (
                                <div
                                    key={post.id}
                                    className="flex flex-col gap-4 border-b border-gray-100 pb-8 last:border-b-0 last:pb-0"
                                >
                                    <div className="flex items-center justify-between gap-4">
                                        <div>
                                            <h3 className="text-xl font-medium text-gray-900 font-heading">
                                                {post.title}
                                            </h3>
                                            <p className="text-lg text-gray-500 mt-2 leading-relaxed">
                                                Last edited on {formatDate(post.date)}
                                                {post.editDraftId && (
                                                    <>
                                                        <span className="mx-2 text-gray-300">•</span>
                                                        <Link
                                                            to={`/editor?draft=${post.editDraftId}`}
                                                            className="text-brand font-medium hover:text-brand-dark transition-colors"
                                                        >
                                                            Edit
                                                        </Link>
                                                    </>
                                                )}
                                            </p>
                                        </div>
                                    </div>
                                    <div className="flex flex-wrap gap-2">
                                        {post.editDraftId && (
                                            <span className="inline-flex items-center rounded-full border px-4 py-1.5 text-base font-medium bg-gray-100 text-gray-600 border-gray-200">
                                                Draft
                                            </span>
                                        )}
                                        {post.platforms
                                            .filter(p => p !== 'Postificus')
                                            .map((platform) => (
                                                <span
                                                    key={platform}
                                                    className={`inline-flex items-center rounded-full border px-4 py-1.5 text-base font-medium ${platformClasses[platform] || 'bg-brand/10 text-brand-dark border-brand/30'}`}
                                                >
                                                    {platform}
                                                </span>
                                            ))}
                                    </div>
                                </div>
                            ))
                    )}
                </CardContent>
            </Card>

            {totalPages > 1 && (
                <div className="flex items-center justify-center gap-2 pb-8">
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                        disabled={currentPage === 1}
                        className="h-9 w-9 p-0 border-transparent bg-brand/10 text-brand hover:bg-brand/20 hover:text-brand-dark disabled:opacity-50"
                    >
                        <ChevronLeft className="h-4 w-4" />
                        <span className="sr-only">Previous</span>
                    </Button>

                    <div className="flex items-center gap-1">
                        {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
                            <button
                                key={page}
                                onClick={() => setCurrentPage(page)}
                                className={`h-9 w-9 rounded-md text-sm font-medium transition-colors ${currentPage === page
                                    ? 'bg-brand text-white shadow-sm'
                                    : 'text-gray-600 hover:bg-gray-100'
                                    }`}
                            >
                                {page}
                            </button>
                        ))}
                    </div>

                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                        disabled={currentPage === totalPages}
                        className="h-9 w-9 p-0 border-transparent bg-brand/10 text-brand hover:bg-brand/20 hover:text-brand-dark disabled:opacity-50"
                    >
                        <ChevronRight className="h-4 w-4" />
                        <span className="sr-only">Next</span>
                    </Button>
                </div>
            )}
        </div>
    );
};

export default Dashboard;
