import React from 'react';

const Header = () => {
    console.log("Header rendered");
    return (
        <header>
            <h1>My App Header</h1>
        </header>
    );
};

export default React.memo(Header);