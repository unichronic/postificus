import React from 'react';
import Navbar from '../components/Navbar';
import { Outlet } from 'react-router-dom';

const MainLayout = () => {
    return (
        <div className="min-h-screen bg-transparent text-gray-900">
            <Navbar />
            <main className="container mx-auto px-6 py-12 pt-24">
                <Outlet />
            </main>
        </div>
    );
};

export default MainLayout;
