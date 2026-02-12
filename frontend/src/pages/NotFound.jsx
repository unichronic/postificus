import React from 'react';
import { Link } from 'react-router-dom';

const NotFound = () => {
    return (
        <div className="mx-auto max-w-2xl text-center py-24">
            <p className="text-base font-semibold text-brand">404</p>
            <h1 className="mt-4 text-4xl font-semibold text-gray-900 font-heading">Page not found</h1>
            <p className="mt-4 text-gray-500 text-base">
                The page you are looking for does not exist or has been moved.
            </p>
            <div className="mt-8">
                <Link
                    to="/dashboard"
                    className="inline-flex items-center justify-center rounded-full px-6 py-2.5 text-base font-medium text-white bg-brand hover:bg-brand-dark transition-colors"
                >
                    Go to Dashboard
                </Link>
            </div>
        </div>
    );
};

export default NotFound;
