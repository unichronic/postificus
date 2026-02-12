import React, { useMemo } from 'react';
import { v4 as uuidv4 } from 'uuid';
import Editor from '../components/Editor';
import Navbar from '../components/Navbar';
import { useSearchParams } from 'react-router-dom';

const EditorPage = () => {
    const [searchParams] = useSearchParams();
    const draftParam = searchParams.get('draft');
    // Generate a new Draft ID if one isn't provided in the URL (future proofing)
    // For now, we just generate one for the session.
    const draftId = useMemo(() => draftParam || uuidv4(), [draftParam]);
    const isExistingDraft = Boolean(draftParam);

    return (
        <div className="min-h-screen editor-surface pt-24">
            <Navbar />
            <Editor draftId={draftId} isExistingDraft={isExistingDraft} />
        </div>
    );
};

export default EditorPage;
