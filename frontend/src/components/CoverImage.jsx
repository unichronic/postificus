import React, { useRef, useState } from 'react';
import { Image as ImageIcon, X } from 'lucide-react';
import { Button } from "@/components/ui/button";

const CoverImage = ({ url, onUpdate }) => {
    const [isHovered, setIsHovered] = useState(false);
    const fileInputRef = useRef(null);

    const handleUpload = () => {
        if (fileInputRef.current) {
            fileInputRef.current.value = '';
            fileInputRef.current.click();
        }
    };

    const handleFileChange = async (event) => {
        const file = event.target.files?.[0];
        if (!file) {
            return;
        }
        if (!file.type.startsWith('image/')) {
            alert('Please select an image file.');
            return;
        }

        // Upload to API
        setIsHovered(false); // Hide overlay during upload

        const formData = new FormData();
        formData.append('file', file);

        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/upload`, {
                method: 'POST',
                body: formData,
            });

            if (!response.ok) {
                throw new Error('Upload failed');
            }

            const data = await response.json();
            if (data.url) {
                onUpdate(data.url);
            }
        } catch (error) {
            console.error("Upload error:", error);
            alert("Failed to upload image.");
        }
    };

    const handleRemove = (e) => {
        e.stopPropagation();
        onUpdate('');
    };

    if (url) {
        return (
            <div
                className="relative w-full h-64 md:h-80 rounded-2xl overflow-hidden group border border-gray-200/60 bg-white/70 backdrop-blur-sm"
                onMouseEnter={() => setIsHovered(true)}
                onMouseLeave={() => setIsHovered(false)}
            >
                <img src={url} alt="Cover" className="w-full h-full object-cover" />
                <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={handleFileChange}
                />

                {isHovered && (
                    <div className="absolute inset-0 bg-black/30 flex items-center justify-center gap-4 transition-opacity">
                        <Button variant="secondary" onClick={handleUpload}>
                            Change Cover
                        </Button>
                        <Button variant="outline" size="icon" onClick={handleRemove} className="border-white text-white hover:bg-white/20">
                            <X className="w-4 h-4" />
                        </Button>
                    </div>
                )}
            </div>
        );
    }

    return (
        <div className="flex items-center gap-4 group">
            <Button
                variant="outline"
                onClick={handleUpload}
                className="text-gray-700 hover:text-brand hover:border-brand/50 transition-all gap-2 text-base px-5 py-2.5"
            >
                <ImageIcon className="w-4 h-4" />
                Add Cover Image
            </Button>
            <input
                ref={fileInputRef}
                type="file"
                accept="image/*"
                className="hidden"
                onChange={handleFileChange}
            />
            <p className="text-sm text-gray-400 opacity-0 group-hover:opacity-100 transition-opacity">
                Use a high-quality image for better engagement
            </p>
        </div>
    );
};

export default CoverImage;
