import React, { useState } from 'react';
import { useEditor, EditorContent } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import Image from '@tiptap/extension-image';
import Link from '@tiptap/extension-link';
import { Bold, Italic, Heading2, List, Quote, Code, Link as LinkIcon, Image as ImageIcon, Save, Send, ArrowLeft } from 'lucide-react';
import { Link as RouterLink } from 'react-router-dom';
import CoverImage from './CoverImage';
import TagInput from './TagInput';
import { Button } from "@/components/ui/button";

const Editor = () => {
    const [title, setTitle] = useState('');
    const [coverImage, setCoverImage] = useState('');
    const [tags, setTags] = useState([]);
    const [isPublishing, setIsPublishing] = useState(false);

    const editor = useEditor({
        extensions: [
            StarterKit,
            Image,
            Link.configure({
                openOnClick: false,
            }),
        ],
        content: '<p>Start writing your magical story...</p>',
        editorProps: {
            attributes: {
                class: 'prose prose-lg max-w-none focus:outline-none min-h-[50vh] p-4 text-gray-700',
            },
        },
    });

    const handlePublish = async () => {
        if (!editor || !title) {
            alert('Please add a title before publishing');
            return;
        }

        setIsPublishing(true);
        const content = editor.getHTML();
        console.log("Publishing:", { title, coverImage, tags, content });

        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/publish/devto`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    title,
                    content,
                    cover_image: coverImage,
                    tags
                })
            });
            const data = await response.json();
            alert(data.status === 'published' ? 'Published successfully!' : 'Failed to publish: ' + (data.error || 'Unknown error'));
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
            title={title}
            className={`p-2 rounded-lg transition-all duration-200 ${isActive
                ? 'bg-magical-violet/10 text-magical-violet shadow-sm'
                : 'text-gray-500 hover:bg-gray-100 hover:text-gray-700'
                }`}
        >
            <Icon className="w-5 h-5" />
        </button>
    );

    const addImage = () => {
        const url = window.prompt('URL');
        if (url) {
            editor.chain().focus().setImage({ src: url }).run();
        }
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

    return (
        <div className="max-w-4xl mx-auto py-12 px-6">
            <div className="flex justify-between items-center mb-8">
                <div className="flex items-center gap-4">
                    <RouterLink to="/dashboard">
                        <Button variant="ghost" size="icon" className="rounded-full hover:bg-gray-100">
                            <ArrowLeft className="w-5 h-5 text-gray-600" />
                        </Button>
                    </RouterLink>
                    <div className="flex items-center gap-2">
                        <div className="w-3 h-3 rounded-full bg-yellow-400 animate-pulse"></div>
                        <span className="text-sm font-medium text-gray-500">Draft - Unsaved changes</span>
                    </div>
                </div>
                <div className="flex gap-3">
                    <Button variant="ghost" className="text-gray-600">
                        <Save className="w-4 h-4 mr-2" />
                        Save Draft
                    </Button>
                    <Button
                        onClick={handlePublish}
                        disabled={isPublishing}
                        className="bg-gradient-to-r from-magical-violet to-magical-fuchsia text-white shadow-lg hover:shadow-xl hover:scale-105 transition-all"
                    >
                        <Send className="w-4 h-4 mr-2" />
                        {isPublishing ? 'Publishing...' : 'Publish'}
                    </Button>
                </div>
            </div>

            <CoverImage url={coverImage} onUpdate={setCoverImage} />

            <div className="space-y-6">
                <input
                    type="text"
                    placeholder="Article Title..."
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    className="text-5xl font-extrabold w-full bg-transparent border-none outline-none placeholder-gray-300 text-gray-800 leading-tight"
                />

                <TagInput tags={tags} setTags={setTags} />

                <div className="relative group">
                    {/* Floating Toolbar */}
                    <div className="sticky top-24 z-10 mx-auto max-w-fit mb-6 px-4 py-2 bg-white/80 backdrop-blur-md border border-gray-200/50 shadow-lg rounded-full flex gap-1 transition-all">
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
                        <div className="w-px h-6 bg-gray-200 my-auto mx-1"></div>
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
                        <div className="w-px h-6 bg-gray-200 my-auto mx-1"></div>
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
                    </div>

                    <div className="min-h-[60vh]">
                        <EditorContent editor={editor} />
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Editor;
