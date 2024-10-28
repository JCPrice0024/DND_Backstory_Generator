import React from 'react';

const Footer = () => {
    console.log("Footer rendered");
    return (
        <footer>
            <p>My App Footer</p>
        </footer>
    );
};

export default React.memo(Footer);