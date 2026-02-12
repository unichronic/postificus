import React from 'react';
import { Link } from 'react-router-dom';
const integrations = ['Medium', 'Dev.to', 'LinkedIn', 'Postificus'];

const Landing = () => {
    const backgroundStyle = {
        backgroundColor: '#f5f5f3',
        backgroundImage: `
            radial-gradient(circle at 15% 50%, rgba(148, 0, 255, 0.08), transparent 45%),
            radial-gradient(circle at 85% 20%, rgba(148, 0, 255, 0.05), transparent 35%),
            radial-gradient(rgba(17, 17, 17, 0.04) 1px, transparent 1px)
        `,
        backgroundSize: 'auto, auto, 28px 28px',
        backgroundPosition: 'center, center, 0 0',
    };

    return (
        <div className="h-screen overflow-hidden text-gray-900" style={backgroundStyle}>
            <nav className="nav-soft fixed left-1/2 top-8 z-50 w-[min(92%,960px)] -translate-x-1/2 rounded-full border border-white/60 bg-gradient-to-b from-[#f6efff]/90 via-white/60 to-transparent shadow-[0_12px_24px_rgba(17,17,17,0.06)] backdrop-blur-2xl">
                <div className="flex h-16 items-center justify-between px-6">
                    <Link to="/" className="flex items-center gap-3 text-2xl font-semibold text-gray-900 transition-colors hover:text-brand font-heading">
                        <img
                            src="/quill-drawing-a-line.png"
                            alt="Postificus quill"
                            className="landing-quill-logo h-6 w-6 object-contain"
                        />
                        Postificus
                    </Link>

                    <div className="flex items-center gap-4">
                        <Link to="/dashboard" className="text-sm font-medium text-gray-500 hover:text-gray-900 transition-colors">
                            Log in
                        </Link>
                        <Link
                            to="/editor"
                            className="rounded-full bg-brand px-5 py-2 text-sm font-medium text-white shadow-[0_6px_18px_rgba(148,0,255,0.28)] transition-colors hover:bg-brand-dark"
                        >
                            Get Started
                        </Link>
                    </div>
                </div>
            </nav>

            <main className="flex h-full flex-col px-6 pb-8 pt-24 md:px-8 md:pb-10">
                <section className="flex flex-1 items-center justify-center text-center">
                    <div className="mx-auto flex max-w-4xl flex-col items-center landing-hero-scale">
                        <h1 className="text-5xl font-bold leading-tight text-gray-900 md:text-6xl font-heading">
                            Write Once.<br />
                            <em className="font-heading text-brand italic font-bold">Publish everywhere.</em>
                        </h1>
                        <p className="mt-6 max-w-2xl text-lg leading-relaxed text-gray-500 md:text-xl">
                            Draft with focus, refine with clarity, and publish to Medium, Dev.to, and LinkedIn without the noise.
                        </p>
                        <div className="mt-7 flex flex-col items-center gap-4 sm:flex-row">
                            <Link
                                to="/editor"
                                className="rounded-full bg-brand px-6 py-3 text-sm font-medium text-white shadow-[0_8px_20px_rgba(148,0,255,0.28)] transition-colors hover:bg-brand-dark"
                            >
                                Start Writing for free
                            </Link>

                        </div>
                    </div>
                </section>

                <section id="integrations" className="text-center">
                    <div className="mx-auto max-w-3xl">
                        <span className="mb-2 inline-block text-base font-semibold uppercase tracking-[0.45em] text-brand md:text-lg">
                            Integrations
                        </span>
                        <div className="relative mt-6 overflow-hidden">
                            <div className="pointer-events-none absolute left-0 top-0 h-full w-24 bg-gradient-to-r from-[#f5f5f3] to-transparent" />
                            <div className="pointer-events-none absolute right-0 top-0 h-full w-24 bg-gradient-to-l from-[#f5f5f3] to-transparent" />
                            <div className="landing-marquee">
                                <div className="landing-marquee-track text-lg font-semibold text-gray-500/80 font-heading">
                                    {[...integrations, ...integrations].map((platform, index) => (
                                        <span
                                            key={`${platform}-${index}`}
                                            className="px-6 opacity-70 grayscale transition-opacity hover:opacity-100"
                                        >
                                            {platform}
                                        </span>
                                    ))}
                                </div>
                            </div>
                        </div>
                    </div>
                </section>
            </main>
        </div>
    );
};

export default Landing;
