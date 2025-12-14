import React, { useMemo } from 'react';
import { v4 as uuidv4 } from 'uuid';
import Editor from '../components/Editor';

const EditorPage = () => {
    // Generate a new Draft ID if one isn't provided in the URL (future proofing)
    // For now, we just generate one for the session.
    const draftId = useMemo(() => uuidv4(), []);

    return (
        <div className="min-h-screen bg-white">
            <Editor draftId={draftId} />
        </div>
    );
};

export default EditorPage;
