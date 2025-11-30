import React, { useState } from 'react';
import { useEditor, EditorContent } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';

const Editor = () => {
    const [title, setTitle] = useState('');

    const editor = useEditor({
        extensions: [
            StarterKit,
        ],
        content: '<p>Start writing your post here...</p>',
        editorProps: {
            attributes: {
                class: 'prose prose-lg max-w-none focus:outline-none min-h-[50vh] p-4',
            },
        },
    });

    const handlePublish = async () => {
        if (!editor) return;
        const content = editor.getHTML();
        console.log("Publishing:", { title, content });

        try {
            const response = await fetch('http://localhost:8080/publish/devto', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    title,
                    content
                })
            });
            const data = await response.json();
            alert(data.status === 'published' ? 'Published successfully!' : 'Failed to publish: ' + (data.error || 'Unknown error'));
        } catch (e) {
            console.error(e);
            alert('Error publishing');
        }
    };

    if (!editor) {
        return null;
    }

    const ToolbarButton = ({ onClick, isActive, children }) => (
        <button
            onClick={onClick}
            className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${isActive
                    ? 'bg-magical-violet text-white shadow-sm'
                    : 'text-gray-600 hover:bg-gray-100'
                }`}
        >
            {children}
        </button>
    );

    return (
        <div className="max-w-4xl mx-auto flex flex-col gap-6 p-6">
            <input
                type="text"
                placeholder="Post Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="text-4xl font-bold w-full bg-transparent border-b-2 border-transparent focus:border-magical-sky/30 outline-none placeholder-gray-300 transition-colors pb-2"
            />

            <div className="bg-white rounded-2xl shadow-xl border border-gray-100 overflow-hidden">
                {/* Toolbar */}
                <div className="flex gap-2 p-3 border-b border-gray-100 bg-gray-50/50 backdrop-blur-sm">
                    <ToolbarButton
                        onClick={() => editor.chain().focus().toggleBold().run()}
                        isActive={editor.isActive('bold')}
                    >
                        Bold
                    </ToolbarButton>
                    <ToolbarButton
                        onClick={() => editor.chain().focus().toggleItalic().run()}
                        isActive={editor.isActive('italic')}
                    >
                        Italic
                    </ToolbarButton>
                    <ToolbarButton
                        onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
                        isActive={editor.isActive('heading', { level: 2 })}
                    >
                        H2
                    </ToolbarButton>
                    <ToolbarButton
                        onClick={() => editor.chain().focus().toggleBulletList().run()}
                        isActive={editor.isActive('bulletList')}
                    >
                        Bullet List
                    </ToolbarButton>
                </div>

                {/* Editor Content */}
                <div className="bg-white">
                    <EditorContent editor={editor} />
                </div>
            </div>

            <div className="flex justify-end gap-4 mt-4">
                <button className="px-6 py-2.5 text-gray-700 font-medium bg-magical-pink/20 hover:bg-magical-pink/40 rounded-xl transition-colors">
                    Save Draft
                </button>
                <button
                    onClick={handlePublish}
                    className="px-6 py-2.5 text-white font-medium bg-gradient-to-r from-magical-violet to-magical-fuchsia rounded-xl shadow-lg hover:shadow-xl hover:scale-[1.02] transition-all duration-200"
                >
                    Publish to Dev.to
                </button>
            </div>
        </div>
    );
};

export default Editor;
