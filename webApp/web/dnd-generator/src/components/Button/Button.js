import React from 'react';

const Button = ({ text, onClick, value }) => {
    return (
        <button onClick={() => onClick(value)}>
            {text}
        </button>
    );
};

export default Button;