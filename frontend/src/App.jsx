import React from 'react';
import Navbar from './components/Navbar';
import Editor from './components/Editor';

function App() {
    return (
        <>
            <Navbar />
            <main>
                <Editor />
            </main>
        </>
    );
}

export default App;
