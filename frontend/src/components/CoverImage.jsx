import React, { useState } from 'react';
import { Image as ImageIcon, X } from 'lucide-react';
import { Button } from "@/components/ui/button";

const CoverImage = ({ url, onUpdate }) => {
    const [isHovered, setIsHovered] = useState(false);

    const handleUpload = () => {
        // Mock upload for now - in real app would upload to server/S3
        const mockUrl = "https://images.unsplash.com/photo-1519389950473-47ba0277781c?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2070&q=80";
        onUpdate(mockUrl);
    };

    const handleRemove = (e) => {
        e.stopPropagation();
        onUpdate('');
    };

    if (url) {
        return (
            <div
                className="relative w-full h-64 md:h-80 rounded-2xl overflow-hidden group mb-8 shadow-md"
                onMouseEnter={() => setIsHovered(true)}
                onMouseLeave={() => setIsHovered(false)}
            >
                <img src={url} alt="Cover" className="w-full h-full object-cover" />

                {isHovered && (
                    <div className="absolute inset-0 bg-black/30 flex items-center justify-center gap-4 transition-opacity">
                        <Button variant="secondary" onClick={handleUpload}>
                            Change Cover
                        </Button>
                        <Button variant="destructive" size="icon" onClick={handleRemove}>
                            <X className="w-4 h-4" />
                        </Button>
                    </div>
                )}
            </div>
        );
    }

    return (
        <div className="flex items-center gap-4 mb-8 group">
            <Button
                variant="outline"
                onClick={handleUpload}
                className="text-gray-500 hover:text-magical-violet hover:border-magical-violet/50 transition-all gap-2"
            >
                <ImageIcon className="w-4 h-4" />
                Add Cover Image
            </Button>
            <p className="text-sm text-gray-400 opacity-0 group-hover:opacity-100 transition-opacity">
                Use a high-quality image for better engagement
            </p>
        </div>
    );
};

export default CoverImage;
