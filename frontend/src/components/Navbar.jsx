import React from 'react';

const Navbar = () => {
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
        },
        link: {
            color: 'var(--color-parma-violet)',
            textDecoration: 'none',
            fontWeight: '500',
        }
    };

    return (
        <nav style={styles.nav}>
            <a href="/" style={styles.logo}>Postificus</a>
            <div style={styles.links}>
                <a href="#" style={styles.link}>New Post</a>
                <a href="#" style={styles.link}>Settings</a>
            </div>
        </nav>
    );
};

export default Navbar;
