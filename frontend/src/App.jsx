import React, { Suspense, lazy } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import MainLayout from './layouts/MainLayout';

const Dashboard = lazy(() => import('./pages/Dashboard'));
const Profile = lazy(() => import('./pages/Profile'));
const Settings = lazy(() => import('./pages/Settings'));
const EditorPage = lazy(() => import('./pages/EditorPage'));
const Landing = lazy(() => import('./pages/Landing'));
const NotFound = lazy(() => import('./pages/NotFound'));

function App() {
    return (
        <BrowserRouter>
            <Suspense fallback={<div className="p-8 text-gray-500">Loading...</div>}>
                <Routes>
                    <Route path="/" element={<Landing />} />
                    <Route element={<MainLayout />}>
                        <Route path="dashboard" element={<Dashboard />} />
                        <Route path="profile" element={<Profile />} />
                        <Route path="settings" element={<Settings />} />
                        <Route path="*" element={<NotFound />} />
                    </Route>
                    <Route path="/editor" element={<EditorPage />} />
                </Routes>
            </Suspense>
        </BrowserRouter>
    );
}

export default App;
