import React, { useState } from 'react';

const Navbar = () => {
    const [showModal, setShowModal] = useState(false);
    const [token, setToken] = useState('');

    const handleConnect = async () => {
        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/config/devto`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ token })
            });
            const data = await response.json();
            if (data.status === 'saved') {
                alert('Token saved successfully!');
                setShowModal(false);
            } else {
                alert('Failed to save token: ' + data.error);
            }
        } catch (e) {
            console.error(e);
            alert('Error connecting to backend');
        }
    };

    return (
        <>
            <nav className="fixed top-0 left-0 right-0 z-50 flex justify-between items-center px-8 py-4 bg-white/70 backdrop-blur-md border-b border-white/20 shadow-sm">
                <a href="/" className="text-2xl font-bold bg-gradient-to-r from-magical-fuchsia to-magical-violet bg-clip-text text-transparent hover:opacity-80 transition-opacity">
                    Postificus
                </a>
                <div className="flex gap-6 items-center">
                    <a href="#" className="text-gray-600 hover:text-magical-violet font-medium transition-colors">
                        New Post
                    </a>
                    <button
                        onClick={() => setShowModal(true)}
                        className="px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-magical-violet to-magical-fuchsia rounded-full shadow-md hover:shadow-lg hover:scale-105 transition-all duration-200"
                    >
                        Connect Dev.to
                    </button>
                </div>
            </nav>

            {/* Spacer for fixed navbar */}
            <div className="h-20"></div>

            {showModal && (
                <div className="fixed inset-0 z-[100] flex justify-center items-center bg-black/40 backdrop-blur-sm" onClick={() => setShowModal(false)}>
                    <div className="bg-white p-8 rounded-2xl shadow-2xl w-96 flex flex-col gap-4 animate-in fade-in zoom-in duration-200" onClick={e => e.stopPropagation()}>
                        <h3 className="text-xl font-bold text-gray-800">Connect Dev.to</h3>
                        <p className="text-sm text-gray-600">Paste your <code className="bg-gray-100 px-1 py-0.5 rounded text-magical-violet">remember_user_token</code> cookie value here:</p>
                        <input
                            type="text"
                            value={token}
                            onChange={e => setToken(e.target.value)}
                            placeholder="Paste token here..."
                            className="w-full p-3 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-magical-sky/50 transition-all"
                        />
                        <div className="flex justify-end gap-3 mt-2">
                            <button
                                onClick={() => setShowModal(false)}
                                className="px-4 py-2 text-gray-500 hover:bg-gray-100 rounded-lg transition-colors"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={handleConnect}
                                className="px-4 py-2 text-white bg-magical-violet hover:bg-magical-fuchsia rounded-lg shadow-md transition-colors"
                            >
                                Save Token
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
};

export default Navbar;
