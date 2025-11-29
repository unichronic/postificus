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
                class: 'prose prose-sm sm:prose lg:prose-lg xl:prose-2xl m-5 focus:outline-none',
            },
        },
    });

    const styles = {
        container: {
            display: 'flex',
            flexDirection: 'column',
            gap: '1rem',
            textAlign: 'left',
            maxWidth: '800px',
            margin: '0 auto',
        },
        input: {
            fontSize: '2rem',
            fontWeight: 'bold',
            width: '100%',
            boxSizing: 'border-box',
            border: 'none',
            borderBottom: `2px solid var(--color-sky-blue)`,
            outline: 'none',
            padding: '0.5rem 0',
            backgroundColor: 'transparent',
        },
        editorContainer: {
            minHeight: '60vh',
            border: `1px solid var(--color-parma-violet)`,
            borderRadius: '8px',
            padding: '1rem',
            backgroundColor: '#fff',
            overflowY: 'auto',
        },
        toolbar: {
            display: 'flex',
            gap: '0.5rem',
            marginBottom: '1rem',
            borderBottom: '1px solid #eee',
            paddingBottom: '0.5rem',
        },
        toolbarButton: (isActive) => ({
            padding: '0.4rem 0.8rem',
            borderRadius: '4px',
            border: 'none',
            cursor: 'pointer',
            backgroundColor: isActive ? 'var(--color-parma-violet)' : '#f0f0f0',
            color: isActive ? 'white' : '#333',
            fontWeight: '500',
        }),
        actions: {
            display: 'flex',
            justifyContent: 'flex-end',
            gap: '1rem',
            marginTop: '1rem',
        }
    };

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

    return (
        <div style={styles.container}>
            <input
                type="text"
                placeholder="Post Title"
                style={styles.input}
                value={title}
                onChange={(e) => setTitle(e.target.value)}
            />

            <div style={styles.editorContainer}>
                <div style={styles.toolbar}>
                    <button
                        onClick={() => editor.chain().focus().toggleBold().run()}
                        style={styles.toolbarButton(editor.isActive('bold'))}
                    >
                        Bold
                    </button>
                    <button
                        onClick={() => editor.chain().focus().toggleItalic().run()}
                        style={styles.toolbarButton(editor.isActive('italic'))}
                    >
                        Italic
                    </button>
                    <button
                        onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
                        style={styles.toolbarButton(editor.isActive('heading', { level: 2 }))}
                    >
                        H2
                    </button>
                    <button
                        onClick={() => editor.chain().focus().toggleBulletList().run()}
                        style={styles.toolbarButton(editor.isActive('bulletList'))}
                    >
                        Bullet List
                    </button>
                </div>
                <EditorContent editor={editor} />
            </div>

            <div style={styles.actions}>
                <button style={{ backgroundColor: 'var(--color-pink)', color: '#333' }}>Save Draft</button>
                <button onClick={handlePublish}>Publish to Dev.to</button>
            </div>
        </div>
    );
};

export default Editor;
