import React, { useState } from 'react';

const Navbar = () => {
    const [showModal, setShowModal] = useState(false);
    const [token, setToken] = useState('');

    const styles = {
        nav: {
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            padding: '1rem 0',
            marginBottom: '2rem',
            borderBottom: `2px solid var(--color-sky-blue)`,
        },
        logo: {
            fontSize: '1.5rem',
            fontWeight: 'bold',
            color: 'var(--color-fuchsia)',
            textDecoration: 'none',
        },
        links: {
            display: 'flex',
            gap: '1rem',
            alignItems: 'center',
        },
        link: {
            color: 'var(--color-parma-violet)',
            textDecoration: 'none',
            fontWeight: '500',
        },
        modalOverlay: {
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: 'rgba(0,0,0,0.5)',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            zIndex: 1000,
        },
        modal: {
            backgroundColor: 'white',
            padding: '2rem',
            borderRadius: '8px',
            width: '400px',
            display: 'flex',
            flexDirection: 'column',
            gap: '1rem',
        }
    };

    const handleConnect = async () => {
        try {
            const response = await fetch('http://localhost:8080/config/devto', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ token })
            });
            const data = await response.json();
            if (data.status === 'saved') {
                alert('Token saved successfully!');
                setShowModal(false);
            } else {
                alert('Failed to save token: ' + data.error);
            }
        } catch (e) {
            console.error(e);
            alert('Error connecting to backend');
        }
    };

    return (
        <>
            <nav style={styles.nav}>
                <a href="/" style={styles.logo}>Postificus</a>
                <div style={styles.links}>
                    <a href="#" style={styles.link}>New Post</a>
                    <button onClick={() => setShowModal(true)} style={{ fontSize: '0.9rem', padding: '0.4rem 0.8rem' }}>
                        Connect Dev.to
                    </button>
                </div>
            </nav>

            {showModal && (
                <div style={styles.modalOverlay} onClick={() => setShowModal(false)}>
                    <div style={styles.modal} onClick={e => e.stopPropagation()}>
                        <h3>Connect Dev.to</h3>
                        <p>Paste your <code>remember_user_token</code> cookie value here:</p>
                        <input
                            type="text"
                            value={token}
                            onChange={e => setToken(e.target.value)}
                            placeholder="Paste token here..."
                        />
                        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.5rem' }}>
                            <button onClick={() => setShowModal(false)} style={{ backgroundColor: '#ccc' }}>Cancel</button>
                            <button onClick={handleConnect}>Save Token</button>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
};

export default Navbar;
