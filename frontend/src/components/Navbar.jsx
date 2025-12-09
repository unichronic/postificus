import React from 'react';
import { Link, useLocation } from 'react-router-dom';

const Navbar = () => {
    const location = useLocation();

    const isActive = (path) => {
        return location.pathname === path ? 'text-magical-violet font-semibold' : 'text-gray-600 hover:text-magical-violet font-medium';
    };

    return (
        <nav className="fixed top-0 left-0 right-0 z-50 flex justify-between items-center px-8 py-4 bg-white/70 backdrop-blur-md border-b border-white/20 shadow-sm">
            <Link to="/dashboard" className="text-2xl font-bold bg-gradient-to-r from-magical-fuchsia to-magical-violet bg-clip-text text-transparent hover:opacity-80 transition-opacity">
                Postificus
            </Link>
            <div className="flex gap-6 items-center">
                <Link to="/dashboard" className={`${isActive('/dashboard')} transition-colors`}>
                    Dashboard
                </Link>
                <Link to="/profile" className={`${isActive('/profile')} transition-colors`}>
                    Profile
                </Link>
                <Link to="/settings" className={`${isActive('/settings')} transition-colors`}>
                    Settings
                </Link>
                <Link to="/editor" className="px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-magical-violet to-magical-fuchsia rounded-full shadow-md hover:shadow-lg hover:scale-105 transition-all duration-200">
                    New Post
                </Link>
            </div>
        </nav>
    );
};

export default Navbar;
