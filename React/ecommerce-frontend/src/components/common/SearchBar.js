import React from 'react';
import { Input } from 'antd';

function SearchBar({ onSearch }) {
    return (
        <Input.Search
            placeholder="Search products"
            onSearch={value => onSearch(value)}
            enterButton
        />
    );
}

export default SearchBar;
