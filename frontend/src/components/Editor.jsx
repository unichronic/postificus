import React, { useState, useEffect, useCallback, useRef } from 'react';
import { useEditor, EditorContent } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import Image from '@tiptap/extension-image';
import Link from '@tiptap/extension-link';
import { Bold, Italic, Heading2, List, Quote, Code, Link as LinkIcon, Image as ImageIcon, Save, Send, CheckCircle, Loader2 } from 'lucide-react';
import CoverImage from './CoverImage';
import TagInput from './TagInput';
import { Button } from "@/components/ui/button";

// Simple debounce hook
const useDebounce = (value, delay) => {
    const [debouncedValue, setDebouncedValue] = useState(value);
    useEffect(() => {
        const handler = setTimeout(() => {
            setDebouncedValue(value);
        }, delay);
        return () => {
            clearTimeout(handler);
        };
    }, [value, delay]);
    return debouncedValue;
};

const Editor = ({ draftId, isExistingDraft = false }) => {
    const [title, setTitle] = useState('');
    const [coverImage, setCoverImage] = useState('');
    const [tags, setTags] = useState([]);
    const [isPublishing, setIsPublishing] = useState(false);
    const [isPublishOpen, setIsPublishOpen] = useState(false);
    const [publishError, setPublishError] = useState('');
    const [isEditorEmpty, setIsEditorEmpty] = useState(true);
    const [selectedPlatforms, setSelectedPlatforms] = useState({
        linkedin: false,
        medium: false,
        devto: false,
    });
    const [saveStatus, setSaveStatus] = useState('saved'); // 'saved', 'saving', 'unsaved'
    const [draftReady, setDraftReady] = useState(!isExistingDraft);
    const imageInputRef = useRef(null);
    const hasLoadedDraft = useRef(false);

    const editor = useEditor({
        extensions: [
            StarterKit,
            Image,
            Link.configure({
                openOnClick: false,
            }),
        ],
        content: '',
        editorProps: {
            attributes: {
                class: 'prose prose-base prose-slate max-w-none focus:outline-none min-h-[50vh] p-5 text-gray-800 leading-relaxed',
            },
        },
        onCreate: ({ editor }) => {
            setIsEditorEmpty(editor.isEmpty);
        },
        onUpdate: ({ editor }) => {
            setSaveStatus('unsaved');
            setIsEditorEmpty(editor.isEmpty);
        },
    });

    // Auto-Save Logic
    const debouncedTitle = useDebounce(title, 1000);
    const debouncedContent = useDebounce(editor?.getHTML(), 1000);
    const debouncedTags = useDebounce(tags, 1000);

    const saveDraft = useCallback(async () => {
        if (!draftId || !editor) return;

        setSaveStatus('saving');
        try {
            const content = editor.getHTML();
            const publishTargets = Object.entries(selectedPlatforms)
                .filter(([, enabled]) => enabled)
                .map(([key]) => key);
            const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/drafts/${draftId}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    title,
                    content,
                    cover_image: coverImage,
                    tags,
                    publish_targets: publishTargets
                })
            });
            if (!response.ok) {
                throw new Error(`Save failed with status ${response.status}`);
            }
            setSaveStatus('saved');
        } catch (error) {
            console.error("Auto-save failed:", error);
            setSaveStatus('unsaved');
        }
    }, [draftId, editor, title, coverImage, tags, selectedPlatforms]);

    useEffect(() => {
        if (!draftReady) return;
        if (debouncedTitle || debouncedContent || debouncedTags) {
            saveDraft();
        }
    }, [debouncedTitle, debouncedContent, debouncedTags, saveDraft, draftReady]);

    useEffect(() => {
        if (!isExistingDraft || !draftId || !editor || hasLoadedDraft.current) return;

        const loadDraft = async () => {
            try {
                const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/drafts/${draftId}`);
                if (!response.ok) {
                    setDraftReady(true);
                    return;
                }
                const data = await response.json();
                if (data.title) {
                    setTitle(data.title);
                }
                if (data.content) {
                    editor.commands.setContent(data.content, false);
                }
                if (data.cover_image) {
                    setCoverImage(data.cover_image);
                }
                if (Array.isArray(data.publish_targets)) {
                    setSelectedPlatforms({
                        linkedin: data.publish_targets.includes('linkedin'),
                        medium: data.publish_targets.includes('medium'),
                        devto: data.publish_targets.includes('devto'),
                    });
                }
                setIsEditorEmpty(editor.isEmpty);
                setSaveStatus('saved');
                hasLoadedDraft.current = true;
                setDraftReady(true);
            } catch (error) {
                console.error("Failed to load draft:", error);
                setDraftReady(true);
            }
        };

        loadDraft();
    }, [draftId, editor, isExistingDraft]);


    const handlePublish = async (platforms) => {
        if (!editor || !title) {
            alert('Please add a title before publishing');
            return;
        }

        if (!platforms.length) {
            setPublishError('Select at least one platform.');
            return;
        }

        setPublishError('');
        setIsPublishing(true);
        const content = editor.getHTML();

        const payload = {
            title,
            content,
            cover_image: coverImage,
            blog_url: '',
            blog_content: content,
        };

        const platformLabels = {
            linkedin: 'LinkedIn',
            medium: 'Medium',
            devto: 'Dev.to',
        };

        const publishToPlatform = async (platformKey) => {
            let endpoint = '/api/publish/linkedin';
            if (platformKey === 'medium') {
                endpoint = '/api/publish/medium';
            } else if (platformKey === 'devto') {
                endpoint = '/api/publish/devto';
            }

            const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}${endpoint}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (!response.ok) {
                let errorMessage = 'Unknown error';
                try {
                    const data = await response.json();
                    errorMessage = data.error || errorMessage;
                } catch (err) {
                    errorMessage = response.statusText || errorMessage;
                }
                throw new Error(errorMessage);
            }
        };

        try {
            const results = await Promise.allSettled(platforms.map((platformKey) => publishToPlatform(platformKey)));
            const failed = results
                .map((result, index) => (result.status === 'rejected' ? platforms[index] : null))
                .filter(Boolean);

            if (failed.length === 0) {
                alert(`Successfully queued for ${platforms.map((p) => platformLabels[p]).join(', ')}.`);
                setIsPublishOpen(false);
            } else {
                const succeeded = platforms.filter((p) => !failed.includes(p));
                const successLabel = succeeded.length ? `Queued for ${succeeded.map((p) => platformLabels[p]).join(', ')}.` : '';
                const failedLabel = `Failed for ${failed.map((p) => platformLabels[p]).join(', ')}.`;
                alert([successLabel, failedLabel].filter(Boolean).join(' '));
            }
        } catch (e) {
            console.error(e);
            alert('Error publishing');
        } finally {
            setIsPublishing(false);
        }
    };

    if (!editor) {
        return null;
    }

    const ToolbarButton = ({ onClick, isActive, icon: Icon, title }) => (
        <button
            onClick={onClick}
            type="button"
            title={title}
            aria-label={title}
            className={`p-2 rounded-md transition-all duration-200 ${isActive
                ? 'bg-brand/15 text-brand'
                : 'text-gray-600 hover:bg-gray-100/80 hover:text-gray-800'
                }`}
        >
            <Icon className="w-4 h-4" />
        </button>
    );

    const addImage = () => {
        if (imageInputRef.current) {
            imageInputRef.current.value = '';
            imageInputRef.current.click();
        }
    };

    const handleInlineImage = (event) => {
        const file = event.target.files?.[0];
        if (!file || !file.type.startsWith('image/')) {
            return;
        }

        const reader = new FileReader();
        reader.onload = () => {
            const src = typeof reader.result === 'string' ? reader.result : '';
            if (src) {
                editor.chain().focus().setImage({ src }).run();
            }
        };
        reader.readAsDataURL(file);
    };

    const setLink = () => {
        const previousUrl = editor.getAttributes('link').href;
        const url = window.prompt('URL', previousUrl);

        if (url === null) {
            return;
        }

        if (url === '') {
            editor.chain().focus().extendMarkRange('link').unsetLink().run();
            return;
        }

        editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run();
    };

    const togglePlatform = (key) => {
        setSelectedPlatforms((prev) => ({ ...prev, [key]: !prev[key] }));
        setPublishError('');
    };

    const confirmPublish = async () => {
        const platforms = Object.keys(selectedPlatforms).filter((key) => selectedPlatforms[key]);
        await handlePublish(platforms);
    };

    const selectedCount = Object.values(selectedPlatforms).filter(Boolean).length;

    return (
        <>
            <div className="max-w-6xl mx-auto py-10 px-6 space-y-6">
                <div className="relative rounded-2xl border border-gray-200/50 bg-white/85 p-6 md:p-8">
                    <div className="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
                        <div>
                            <p className="text-sm uppercase tracking-[0.3em] text-brand/80">Editor</p>
                            <h1 className="mt-3 text-3xl font-semibold text-gray-900 font-heading">Shape your next post</h1>
                            <p className="text-gray-500 mt-2 text-base leading-relaxed max-w-2xl">
                                Draft with focus, refine with clarity, and publish when you are ready.
                            </p>
                        </div>
                        <div className="flex flex-col sm:flex-row items-start sm:items-center gap-3">
                            <div className="flex items-center gap-2 rounded-full border border-gray-200/60 bg-white/85 px-3 py-1.5 text-xs text-gray-600">
                                {saveStatus === 'saving' && <Loader2 className="w-3 h-3 animate-spin text-gray-400" />}
                                {saveStatus === 'saved' && <CheckCircle className="w-3 h-3 text-brand" />}
                                {saveStatus === 'unsaved' && <div className="w-3 h-3 rounded-full bg-brand/40 animate-pulse"></div>}
                                <span className="text-xs font-medium">
                                    {saveStatus === 'saving' ? 'Saving...' : saveStatus === 'saved' ? 'Saved' : 'Unsaved changes'}
                                </span>
                            </div>
                            <div className="flex gap-3">
                                <Button variant="ghost" className="text-gray-700 text-sm px-3 py-2" onClick={saveDraft}>
                                    <Save className="w-4 h-4 mr-2" />
                                    Save Draft
                                </Button>
                                <Button
                                    onClick={() => {
                                        setPublishError('');
                                        setIsPublishOpen(true);
                                    }}
                                    disabled={isPublishing}
                                    className="bg-brand hover:bg-brand-dark text-white transition-colors px-4 py-2 text-sm"
                                >
                                    <Send className="w-4 h-4 mr-2" />
                                    {isPublishing ? 'Publishing...' : 'Publish'}
                                </Button>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="rounded-2xl border border-gray-200/50 bg-white/85 px-6 py-4 md:px-8 md:py-5 space-y-4">
                    <CoverImage url={coverImage} onUpdate={setCoverImage} />

                    <div className="space-y-2">
                        <p className="text-sm uppercase tracking-[0.25em] text-brand/70">Title</p>
                        <input
                            type="text"
                            placeholder="Article Title..."
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            aria-label="Article title"
                            className="text-[2.75rem] font-semibold w-full bg-transparent border-none outline-none placeholder-gray-300 text-gray-900 leading-tight font-heading"
                        />
                    </div>

                    <div className="space-y-2">
                        <p className="text-sm uppercase tracking-[0.25em] text-brand/70">Tags</p>
                        <TagInput tags={tags} setTags={setTags} />
                    </div>
                </div>

                <div className="rounded-2xl border border-gray-200/50 bg-white/85 px-6 py-4 md:px-8 md:py-5">
                    <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                        <div>
                            <p className="text-sm uppercase tracking-[0.25em] text-brand/70">Content</p>
                        </div>
                    </div>

                    <div className="mt-5 relative group">
                        <div className="sticky top-24 z-10 mb-4 flex flex-wrap gap-2 rounded-full border border-gray-200/60 bg-white/85 px-3 py-1.5">
                            <div className="flex items-center gap-1 rounded-full bg-gray-50/80 px-2 py-1">
                                <ToolbarButton
                                    onClick={() => editor.chain().focus().toggleBold().run()}
                                    isActive={editor.isActive('bold')}
                                    icon={Bold}
                                    title="Bold"
                                />
                                <ToolbarButton
                                    onClick={() => editor.chain().focus().toggleItalic().run()}
                                    isActive={editor.isActive('italic')}
                                    icon={Italic}
                                    title="Italic"
                                />
                            </div>
                            <div className="flex items-center gap-1 rounded-full bg-gray-50/80 px-2 py-1">
                                <ToolbarButton
                                    onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
                                    isActive={editor.isActive('heading', { level: 2 })}
                                    icon={Heading2}
                                    title="Heading 2"
                                />
                                <ToolbarButton
                                    onClick={() => editor.chain().focus().toggleBulletList().run()}
                                    isActive={editor.isActive('bulletList')}
                                    icon={List}
                                    title="Bullet List"
                                />
                                <ToolbarButton
                                    onClick={() => editor.chain().focus().toggleBlockquote().run()}
                                    isActive={editor.isActive('blockquote')}
                                    icon={Quote}
                                    title="Quote"
                                />
                                <ToolbarButton
                                    onClick={() => editor.chain().focus().toggleCodeBlock().run()}
                                    isActive={editor.isActive('codeBlock')}
                                    icon={Code}
                                    title="Code Block"
                                />
                            </div>
                            <div className="flex items-center gap-1 rounded-full bg-gray-50/80 px-2 py-1">
                                <ToolbarButton
                                    onClick={setLink}
                                    isActive={editor.isActive('link')}
                                    icon={LinkIcon}
                                    title="Link"
                                />
                                <ToolbarButton
                                    onClick={addImage}
                                    isActive={false}
                                    icon={ImageIcon}
                                    title="Image"
                                />
                                <input
                                    ref={imageInputRef}
                                    type="file"
                                    accept="image/*"
                                    className="hidden"
                                    onChange={handleInlineImage}
                                />
                            </div>
                        </div>

                        <div className="relative min-h-[50vh]">
                            {isEditorEmpty && (
                                <div className="pointer-events-none absolute top-0 left-0 select-none p-5 text-gray-400">
                                    Start writing...
                                </div>
                            )}
                            <EditorContent editor={editor} />
                        </div>
                    </div>
                </div>
            </div>

            {isPublishOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 p-4">
                    <div className="w-full max-w-md rounded-xl border border-gray-200 bg-white p-5 shadow-lg">
                        <div className="flex items-start justify-between gap-4">
                            <div>
                                <h2 className="text-lg font-semibold text-gray-900 font-heading">Publish to platforms</h2>
                                <p className="text-sm text-gray-500 mt-1">Choose where you want to publish this post.</p>
                            </div>
                            <button
                                type="button"
                                className="text-gray-400 hover:text-gray-600"
                                onClick={() => {
                                    setIsPublishOpen(false);
                                    setPublishError('');
                                }}
                                aria-label="Close publish dialog"
                            >
                                x
                            </button>
                        </div>

                        <div className="mt-5 space-y-3">
                            {[
                                { key: 'linkedin', label: 'LinkedIn' },
                                { key: 'medium', label: 'Medium' },
                                { key: 'devto', label: 'Dev.to' },
                            ].map((platform) => (
                                <label
                                    key={platform.key}
                                    className="flex items-center justify-between gap-3 rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
                                >
                                    <span>{platform.label}</span>
                                    <input
                                        type="checkbox"
                                        checked={selectedPlatforms[platform.key]}
                                        onChange={() => togglePlatform(platform.key)}
                                        className="h-4 w-4 accent-brand"
                                    />
                                </label>
                            ))}
                        </div>

                        {publishError && (
                            <p className="mt-3 text-sm text-brand-dark">{publishError}</p>
                        )}

                        <div className="mt-6 flex items-center justify-between">
                            <span className="text-xs text-gray-500">{selectedCount} selected</span>
                            <div className="flex gap-2">
                                <Button
                                    variant="ghost"
                                    onClick={() => {
                                        setIsPublishOpen(false);
                                        setPublishError('');
                                    }}
                                    className="text-gray-700 text-sm px-3 py-2"
                                >
                                    Cancel
                                </Button>
                                <Button
                                    onClick={confirmPublish}
                                    disabled={isPublishing || selectedCount === 0}
                                    className="bg-brand hover:bg-brand-dark text-white text-sm px-4 py-2"
                                >
                                    {isPublishing ? 'Publishing...' : 'Publish'}
                                </Button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
};

export default Editor;
