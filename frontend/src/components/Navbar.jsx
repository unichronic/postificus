import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';

const Navbar = () => {
    const location = useLocation();
    const [isOpen, setIsOpen] = useState(false);

    const isActive = (path) => {
        return location.pathname === path ? 'text-brand font-semibold' : 'text-gray-700 hover:text-brand font-medium';
    };

    return (
        <nav className="nav-soft fixed top-0 left-0 right-0 z-50 bg-gradient-to-b from-[#f6efff]/90 via-white/60 to-transparent backdrop-blur-2xl shadow-[0_10px_24px_rgba(17,17,17,0.06)]">
            <div className="absolute inset-x-0 top-0 h-[2px] bg-gradient-to-r from-brand/40 via-brand/10 to-transparent" />
            <div className="relative mx-auto w-full max-w-7xl px-6 md:px-8">
                <div className="flex items-center justify-between py-5 md:grid md:grid-cols-[1fr_auto_1fr] md:items-center">
                <Link to="/" className="flex items-center gap-3 text-2xl font-semibold text-gray-900 hover:text-brand transition-colors font-heading">
                    <img
                        src="/quill-drawing-a-line.png"
                        alt="Postificus quill"
                        className="landing-quill-logo h-6 w-6 object-contain"
                    />
                    Postificus
                </Link>
                    <div className="hidden md:flex gap-8 items-center text-base justify-self-center">
                        <Link to="/dashboard" className={`${isActive('/dashboard')} transition-colors`}>
                            Dashboard
                        </Link>
                        <Link to="/profile" className={`${isActive('/profile')} transition-colors`}>
                            Profile
                        </Link>
                        <Link to="/settings" className={`${isActive('/settings')} transition-colors`}>
                            Settings
                        </Link>
                    </div>
                    <div className="hidden md:flex justify-self-end">
                        <Link to="/editor" className="px-5 py-2.5 text-base font-medium text-white bg-brand rounded-full hover:bg-brand-dark transition-colors">
                            New Post
                        </Link>
                    </div>
                    <button
                        type="button"
                        className="md:hidden inline-flex items-center justify-center rounded-lg border border-gray-200 px-4 py-2 text-base text-gray-700 hover:bg-gray-50 transition-colors"
                        onClick={() => setIsOpen((prev) => !prev)}
                        aria-expanded={isOpen}
                        aria-controls="mobile-nav"
                        aria-label="Toggle navigation"
                    >
                        Menu
                    </button>
                </div>
            </div>
            {isOpen && (
                <div id="mobile-nav" className="md:hidden px-6 pb-4">
                    <div className="flex flex-col gap-3 rounded-xl bg-white border border-gray-200 p-5 shadow-sm text-base">
                        <Link to="/dashboard" className={`${isActive('/dashboard')} transition-colors`} onClick={() => setIsOpen(false)}>
                            Dashboard
                        </Link>
                        <Link to="/profile" className={`${isActive('/profile')} transition-colors`} onClick={() => setIsOpen(false)}>
                            Profile
                        </Link>
                        <Link to="/settings" className={`${isActive('/settings')} transition-colors`} onClick={() => setIsOpen(false)}>
                            Settings
                        </Link>
                        <Link to="/editor" className="px-5 py-2.5 text-base font-medium text-white bg-brand rounded-full hover:bg-brand-dark transition-colors" onClick={() => setIsOpen(false)}>
                            New Post
                        </Link>
                    </div>
                </div>
            )}
        </nav>
    );
};

export default Navbar;
