import React, { useState } from 'react';
import { X, Hash } from 'lucide-react';
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";

const TagInput = ({ tags, setTags }) => {
    const [input, setInput] = useState('');

    const handleKeyDown = (e) => {
        if (e.key === 'Enter' && input.trim()) {
            e.preventDefault();
            if (!tags.includes(input.trim()) && tags.length < 4) {
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
                    className="pl-2 pr-1 py-1 bg-magical-violet/10 text-magical-violet hover:bg-magical-violet/20 border-none gap-1"
                >
                    <Hash className="w-3 h-3 opacity-50" />
                    {tag}
                    <button
                        onClick={() => removeTag(tag)}
                        className="ml-1 hover:bg-magical-violet/20 rounded-full p-0.5 transition-colors"
                    >
                        <X className="w-3 h-3" />
                    </button>
                </Badge>
            ))}

            {tags.length < 4 && (
                <Input
                    type="text"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyDown={handleKeyDown}
                    placeholder={tags.length === 0 ? "Add up to 4 tags..." : "Add another..."}
                    className="w-40 border-none shadow-none focus-visible:ring-0 px-0 placeholder:text-gray-400 h-8"
                />
            )}
        </div>
    );
};

export default TagInput;
