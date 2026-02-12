import React, { useState } from 'react';
import { X, Hash } from 'lucide-react';
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";

const TagInput = ({ tags, setTags, maxTags = 4, placeholder }) => {
    const [input, setInput] = useState('');

    const handleKeyDown = (e) => {
        if (e.key === 'Enter' && input.trim()) {
            e.preventDefault();
            if (!tags.includes(input.trim()) && tags.length < maxTags) {
                setTags([...tags, input.trim()]);
                setInput('');
            }
        } else if (e.key === 'Backspace' && !input && tags.length > 0) {
            setTags(tags.slice(0, -1));
        }
    };

    const removeTag = (tagToRemove) => {
        setTags(tags.filter(tag => tag !== tagToRemove));
    };

    return (
        <div className="flex flex-wrap items-center gap-2 mb-6">
            {tags.map((tag) => (
                <Badge
                    key={tag}
                    variant="secondary"
                    className="pl-3 pr-2 py-1.5 bg-white/70 text-gray-700 border border-gray-200/50 gap-1 text-base"
                >
                    <Hash className="w-4 h-4 opacity-50" />
                    {tag}
                    <button
                        type="button"
                        onClick={() => removeTag(tag)}
                        className="ml-1 hover:bg-gray-200 rounded-full p-0.5 transition-colors"
                    >
                        <X className="w-4 h-4" />
                    </button>
                </Badge>
            ))}

            {tags.length < maxTags && (
                <Input
                    type="text"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyDown={handleKeyDown}
                    placeholder={tags.length === 0 ? (placeholder || `Add up to ${maxTags} tags...`) : "Add another..."}
                    className="w-48 border-none shadow-none focus-visible:ring-0 px-0 placeholder:text-gray-400 h-9 text-base"
                />
            )}
        </div>
    );
};

export default TagInput;
